package acquisition

// AudioRetriever will somehow get Audio for Japanese text
type AudioRetriever interface {
	GetAudioForJapanese(text string) []byte
}
