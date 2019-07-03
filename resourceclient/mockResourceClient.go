package resourceclient

import "github.com/stretchr/testify/mock"

type MockResourceClient struct {
	mock.Mock
}

func (this *MockResourceClient) Get(address string) (string, error) {
	args := this.Called(address)
	return args.String(0), nil
}
