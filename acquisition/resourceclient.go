package acquisition

type ResourceClient interface {
	Get(address string) (string, error)
}
