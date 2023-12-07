package metrics

const serviceName = "wormscan-tx-tracker"

type Metrics interface {
	IncVaaConsumedQueue(chainID uint16)
	IncVaaUnfiltered(chainID uint16)
	IncOriginTxInserted(chainID uint16)
	IncVaaWithoutTxHash(chainID uint16)
	IncVaaWithTxHashFixed(chainID uint16)
	AddVaaProcessedDuration(chainID uint16, duration float64)
}
