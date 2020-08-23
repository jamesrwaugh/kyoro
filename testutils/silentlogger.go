package testutils

// SilentWriter is used to not print anything for a logger
type SilentWriter struct {
}

func (l SilentWriter) Write(p []byte) (n int, err error) {
	return
}
