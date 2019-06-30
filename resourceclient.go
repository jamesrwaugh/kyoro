package kyoro

type ResourceClient interface {
	Get(address string) (string, error)
}
