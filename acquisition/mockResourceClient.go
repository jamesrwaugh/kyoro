package acquisition

import "github.com/stretchr/testify/mock"

// MockResourceClient is exactly what it says.
type MockResourceClient struct {
	mock.Mock
}

// Get is the main ResourceClient function
func (mrc *MockResourceClient) Get(address string) (string, error) {
	args := mrc.Called(address)
	return args.String(0), nil
}
