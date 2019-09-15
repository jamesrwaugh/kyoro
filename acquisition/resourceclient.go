package acquisition

// ResourceClient is an interface for retreiving some internet-based resouce
// for sentences and vocab.
type ResourceClient interface {
	Get(address string) (string, error)
}
