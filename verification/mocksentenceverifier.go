package verification

import "github.com/stretchr/testify/mock"

// MockSentenceVerifier is exactly what it says.
type MockSentenceVerifier struct {
	mock.Mock
}

// UserConfirmSentence is the main MockSentenceVerifier function
func (mrc *MockSentenceVerifier) UserConfirmSentence(info *SentenceVerificationInfo) bool {
	return true
}
