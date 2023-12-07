package metrics

// DummyMetrics is a dummy implementation of Metric interface.
type DummyMetrics struct{}

// NewDummyMetrics returns a new instance of DummyMetrics.
func NewDummyMetrics() *DummyMetrics {
	return &DummyMetrics{}
}

// IncVaaConsumedQueue is a dummy implementation of IncVaaConsumedQueue.
func (d *DummyMetrics) IncVaaConsumedQueue(chainID uint16) {}

// IncVaaUnfiltered is a dummy implementation of IncVaaUnfiltered.
func (d *DummyMetrics) IncVaaUnfiltered(chainID uint16) {}

// IncOriginTxInserted is a dummy implementation of IncOriginTxInserted.
func (d *DummyMetrics) IncOriginTxInserted(chainID uint16) {}

// IncVaaWithoutTxHash is a dummy implementation of IncVaaWithoutTxHash.
func (d *DummyMetrics) IncVaaWithoutTxHash(chainID uint16) {}

// IncVaaWithTxHashFixed is a dummy implementation of IncVaaWithTxHashFixed.
func (d *DummyMetrics) IncVaaWithTxHashFixed(chainID uint16) {}

// AddVaaProcessedDuration is a dummy implementation of AddVaaProcessedDuration.
func (d *DummyMetrics) AddVaaProcessedDuration(chainID uint16, duration float64) {}
