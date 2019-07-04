package acquisition

// Sentence has enough information to describe a sentence
type Translation struct {
	Japanese      string
	Reading       string
	KanjiReadings []string
	English       string
}
