package acquisition

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// NewJibikiSentenceretriever creates a new JibikiSentenceretriever
// To retrieve a sentence from jibiki.app
func NewJibikiSentenceretriever(c ResourceClient) *JibikiSentenceretriever {
	j := JibikiSentenceretriever{c}
	return &j
}

// JibikiSentenceretriever retrieves a sentence from jibiki.app
type JibikiSentenceretriever struct {
	client ResourceClient
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
	log.Println("[Jibiki] Looking for sentences on " + url)
	resp, _ := jibiki.client.Get(url)

	var jibikiResponse JibikiRestSentenceResponse
	json.Unmarshal([]byte(resp), &jibikiResponse)

	if len(jibikiResponse) > maxSentences {
		jibikiResponse = jibikiResponse[0:maxSentences]
	}

	for _, sentence := range jibikiResponse {
		translation := sentence.Translations[0].Sentence
		sentences = append(sentences, Translation{
			Japanese: sentence.Sentence,
			English:  translation,
		})
	}

	return sentences
}
