package queue

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/wormhole-foundation/wormhole-explorer/parser/internal/metrics"
	"github.com/wormhole-foundation/wormhole-explorer/parser/internal/sqs"
	"go.uber.org/zap"
)

// SQSOption represents a VAA queue in SQS option function.
type SQSOption func(*SQS)

// SQS represents a VAA queue in SQS.
type SQS struct {
	consumer      *sqs.Consumer
	ch            chan ConsumerMessage
	chSize        int
	wg            sync.WaitGroup
	filterConsume FilterConsumeFunc
	metrics       metrics.Metrics
	logger        *zap.Logger
}

// FilterConsumeFunc filter vaaa func definition.
type FilterConsumeFunc func(vaaEvent *VaaEvent) bool

// NewVAASQS creates a VAA queue in SQS instances.
func NewVAASQS(consumer *sqs.Consumer, filterConsume FilterConsumeFunc, metrics metrics.Metrics, logger *zap.Logger, opts ...SQSOption) *SQS {
	s := &SQS{
		consumer:      consumer,
		chSize:        10,
		filterConsume: filterConsume,
		metrics:       metrics,
		logger:        logger}
	for _, opt := range opts {
		opt(s)
	}
	s.ch = make(chan ConsumerMessage, s.chSize)
	return s
}

// WithChannelSize allows to specify an channel size when setting a value.
func WithChannelSize(size int) SQSOption {
	return func(d *SQS) {
		d.chSize = size
	}
}

// Consume returns the channel with the received messages from SQS queue.
func (q *SQS) Consume(ctx context.Context) <-chan ConsumerMessage {
	go func() {
		for {
			messages, err := q.consumer.GetMessages(ctx)
			if err != nil {
				q.logger.Error("Error getting messages from SQS", zap.Error(err))
				continue
			}
			expiredAt := time.Now().Add(q.consumer.GetVisibilityTimeout())
			for _, msg := range messages {

				// unmarshal body to sqsEvent
				var sqsEvent sqsEvent
				err := json.Unmarshal([]byte(*msg.Body), &sqsEvent)
				if err != nil {
					q.logger.Error("Error decoding message from SQS", zap.Error(err))
					continue
				}

				// unmarshal message to vaaEvent
				var vaaEvent VaaEvent
				err = json.Unmarshal([]byte(sqsEvent.Message), &vaaEvent)
				if err != nil {
					q.logger.Error("Error decoding vaaEvent message from SQSEvent", zap.Error(err))
					continue
				}
				q.metrics.IncVaaConsumedQueue(vaaEvent.ChainID)

				// filter vaaEvent by p2p net.
				if q.filterConsume(&vaaEvent) {
					if err := q.consumer.DeleteMessage(ctx, msg.ReceiptHandle); err != nil {
						q.logger.Error("Error deleting message from SQS", zap.Error(err))
					}
					continue
				}
				q.metrics.IncVaaUnfiltered(vaaEvent.ChainID)

				q.wg.Add(1)
				q.ch <- &sqsConsumerMessage{
					id:        msg.ReceiptHandle,
					data:      &vaaEvent,
					wg:        &q.wg,
					logger:    q.logger,
					consumer:  q.consumer,
					expiredAt: expiredAt,
					ctx:       ctx,
				}
			}
			q.wg.Wait()
		}

	}()
	return q.ch
}

// Close closes all consumer resources.
func (q *SQS) Close() {
	close(q.ch)
}

type sqsConsumerMessage struct {
	data      *VaaEvent
	consumer  *sqs.Consumer
	wg        *sync.WaitGroup
	id        *string
	logger    *zap.Logger
	expiredAt time.Time
	ctx       context.Context
}

func (m *sqsConsumerMessage) Data() *VaaEvent {
	return m.data
}

func (m *sqsConsumerMessage) Done() {
	if err := m.consumer.DeleteMessage(m.ctx, m.id); err != nil {
		m.logger.Error("Error deleting message from SQS", zap.Error(err))
	}
	m.wg.Done()
}

func (m *sqsConsumerMessage) Failed() {
	m.wg.Done()
}

func (m *sqsConsumerMessage) IsExpired() bool {
	return m.expiredAt.Before(time.Now())
}
