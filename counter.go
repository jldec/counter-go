package counter

type Counter interface {
	Get() uint64
	Inc()
}
