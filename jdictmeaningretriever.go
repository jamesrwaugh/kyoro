package kyoro

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
)

func NewJdictMeaningRetriever(c ResourceClient) *JdictMeaningRetriever {
	r := JdictMeaningRetriever{c}
	return &r
}

type JdictMeaningRetriever struct {
	client ResourceClient
}

func (this JdictMeaningRetriever) getJDicResults(word string) []string {
	// URL query format from http://nihongo.monash.edu/wwwjdicinf.html#backdoor_tag:
	// 	1: Dictionary type: Dictionary to use, EDICT
	// 	Z: Return format:   Backdoor entry, raw data only
	// 	U: Search Type:     For dictionary lookups, text is UTF-8 (Must be URL escapred)
	// 	J: the key type:    "J" Required for Japanese keys
	// The search query is appended to the end as URL-escapted UTF-8 in this case.
	baseURL := "http://nihongo.monash.edu/cgi-bin/wwwjdic?1ZUJ"
	url := baseURL + url.QueryEscape(word)
	html, _ := this.client.Get(url)
	doc := soup.HTMLParse(html)
	results := doc.Find("pre")
	if results.Error != nil {
		log.Println("Could not find WWWJDIC entries for \"", word, "\"")
		return []string{}
	}
	resultLines := strings.Split(results.Text(), "\n")
	return resultLines
}

func (this JdictMeaningRetriever) parseDictionaryEntries(word string, entries []string) Translation {
	r, _ := regexp.Compile("(.*) \\[(.*)\\] ?\\/\\((.*?)\\) ?(.*)")
	for _, entry := range entries {
		matches := r.FindStringSubmatch(entry)
		if len(matches) < 5 || matches[1] != word {
			continue
		}
		log.Println("WWWJDIC: Found match for \"", word, "\" as ", entry)
		return Translation{
			Japanese: matches[1],
			Reading:  matches[2],
			English:  matches[4],
		}
	}
	log.Println("Could not find a WWWJDIC entry for ", word)
	return Translation{}
}

func (this JdictMeaningRetriever) GetMeaningforKanji(word string) Translation {
	dictionaryEntries := this.getJDicResults(word)
	return this.parseDictionaryEntries(word, dictionaryEntries)
}
