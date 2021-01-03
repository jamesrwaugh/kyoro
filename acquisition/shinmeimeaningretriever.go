package acquisition

import (
	"log"
	"net/url"
	"strings"
	"encoding/json"
)

// NewShinMeiMeaningRetriever creates a new ShinMeiMeaningRetriever
// To retrieve a word's meaning from JDict
func NewShinMeiMeaningRetriever(c ResourceClient, logger *log.Logger) *ShinMeiMeaningRetriever {
	r := ShinMeiMeaningRetriever{c, logger}
	return &r
}

type yomichanDictEntry struct {
	Expression     string     `json:"expression"`
	Reading        string     `json:"reading"`
	DefinitionTags string     `json:"definitionTags"`
	Rules          string     `json:"rules"`
	Score          int        `json:"score"`
	Glossary       [][]string `json:"glossary"`
	Sequence       int        `json:"sequence"`
	TermTags       string     `json:"termTags"`
}

// ShinMeiMeaningRetriever retrieves a word's meaning from JDict
type ShinMeiMeaningRetriever struct {
	client ResourceClient
	logger *log.Logger
}

func (dict ShinMeiMeaningRetriever) getResults(word string) []string {
	baseURL := "http://localhost:8000/book/"
	url := baseURL + url.QueryEscape(word)
	resp, _ := dict.client.Get(url)
	var jibikiResponse yomichanDictEntry
	json.Unmarshal([]byte(resp), &jibikiResponse)
	firstEntry := jibikiResponse.Glossary[0]
	joined := strings.Join(firstEntry, "<br/>")
	return []string{joined}
}

// GetMeaningforKanji retrieves a word's meaning from JDict
func (dict ShinMeiMeaningRetriever) GetMeaningforKanji(word string) Translation {
	dictionaryEntries := dict.getResults(word)
	return Translation {
		Dictionary: dictionaryEntries[0],
	}
}
