package acquisition

// Translation has enough information to describe a sentence
type Translation struct {
	Japanese      string
	Dictionary 	  string
	Reading       string
	KanjiReadings []string
	English       string
}
