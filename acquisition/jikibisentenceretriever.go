package acquisition

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// NewJibikiSentenceretriever creates a new JibikiSentenceretriever
// To retrieve a sentence from jibiki.app
func NewJibikiSentenceretriever(c ResourceClient, m MeaningRetriever, logger *log.Logger) *JibikiSentenceretriever {
	j := JibikiSentenceretriever{c, m, logger}
	return &j
}

// JibikiSentenceretriever retrieves a sentence from jibiki.app
type JibikiSentenceretriever struct {
	client ResourceClient
	meaning MeaningRetriever
	logger *log.Logger
}

// JibikiRestSentenceResponse is a return type from the Jibiki API.
type JibikiRestSentenceResponse []struct {
	ID           int    `json:"id"`
	Language     string `json:"language"`
	Sentence     string `json:"sentence"`
	Translations []struct {
		ID       int    `json:"id"`
		Language string `json:"language"`
		Sentence string `json:"sentence"`
	} `json:"translations"`
}

// GetSentencesforKanji gets a number of sentences for a given kanji.
func (jibiki JibikiSentenceretriever) GetSentencesforKanji(kanji string, maxSentences int) []Translation {
	var sentences []Translation
	escapedQuery := url.QueryEscape(kanji)
	url := fmt.Sprintf("https://api.jibiki.app/sentences?query=%s", escapedQuery)
	jibiki.logger.Println("[Jibiki] Looking for sentences on " + url)
	resp, _ := jibiki.client.Get(url)

	var jibikiResponse JibikiRestSentenceResponse
	json.Unmarshal([]byte(resp), &jibikiResponse)

	if len(jibikiResponse) > maxSentences {
		jibikiResponse = jibikiResponse[0:maxSentences]
	}

	dictionaryMeaning := jibiki.meaning.GetMeaningforKanji(kanji)

	for _, sentence := range jibikiResponse {
		firstEntry := sentence.Translations[0]
		sentences = append(sentences, Translation{
			Japanese: firstEntry.Sentence,
			English:  sentence.Sentence,
			Dictionary: dictionaryMeaning.Dictionary,
		})
	}

	return sentences
}
