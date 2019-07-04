package acquisition

type MeaningRetriever interface {
	GetMeaningforKanji(kanji string) Translation
}
