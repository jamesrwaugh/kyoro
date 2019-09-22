package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/jamesrwaugh/kyoro/acquisition"
	"github.com/jamesrwaugh/kyoro/anki"
)

// KyoroProduction is the actual real thing.
type KyoroProduction struct {
}

// NewKyoro creates a new Kyoro object.
func NewKyoro() (instance Kyoro) {
	return &KyoroProduction{}
}

func (kyoro KyoroProduction) makeSentenceAnkiCard(
	sentence acquisition.Translation,
	options Options,
) anki.AnkiCard {
	boldInputPhrase := fmt.Sprintf("<b>%s</b>", options.InputPhrase)
	boldJapanese := strings.Replace(sentence.Japanese, options.InputPhrase, boldInputPhrase, -1)
	cardFields := map[string]string{
		"japanese": boldJapanese,
		"reading":  sentence.Reading,
	}
	if !options.MonoligualMode {
		cardFields["english"] = sentence.English
	}
	return anki.AnkiCard{
		DeckName:  options.DeckName,
		ModelName: options.ModelName,
		Fields:    cardFields,
		Tags: []string{
			"kyoro",
		},
	}
}

func (kyoro KyoroProduction) generateKeywordBackHTML(
	sentences []acquisition.Translation,
	monoligualMode bool,
) (result string, err error) {
	backHTMLTemplate := `
<ul>
  {{range .sentences}}
    <li style="text-align: left">{{.Japanese}}</li>
    <ul>
      {{range .KanjiReadings}}
        <li style="text-align: left"><sub>{{.}}</sub></li>
      {{end}}
      {{if .English }}
        <li style="text-align: left"><sub>{{.English}}</sub></li>
      {{end}}
    </ul>
  {{end}}
</ul>`
	data := map[string]interface{}{
		"sentences":   sentences,
		"monolingual": monoligualMode,
	}
	ankiTemplate, err := template.New("anki").Parse(backHTMLTemplate)
	var buffer bytes.Buffer
	if err := ankiTemplate.Execute(&buffer, data); err != nil {
		panic(err)
	}
	result = string(buffer.String())
	return
}

func (kyoro KyoroProduction) makeKeywordAnkiCard(
	options Options,
	inputPhrase acquisition.Translation,
	sentences []acquisition.Translation,
) anki.AnkiCard {
	cardFields := map[string]string{
		"japanese": inputPhrase.Japanese,
	}
	if !options.MonoligualMode {
		cardFields["english"] = inputPhrase.English
	}
	backMatter, _ := kyoro.generateKeywordBackHTML(sentences, options.MonoligualMode)
	cardFields["sentences"] = backMatter
	cardFields["reading"] = inputPhrase.Reading
	return anki.AnkiCard{
		DeckName:  options.DeckName,
		ModelName: options.ModelName,
		Fields:    cardFields,
		Tags: []string{
			"kyoro",
		},
	}
}

// Kyoro runs the main procedure of Kyoro from the command line,
// and adds cards accordingly.
func (kyoro KyoroProduction) Kyoro(
	options Options,
	anki anki.AnkiService,
	sentenceSource acquisition.SentenceRetriever,
	meaningSource acquisition.MeaningRetriever,
) bool {
	if !anki.IsConnected() {
		log.Println("Could not connect to Anki. Failing.")
		return false
	}
	sentences := sentenceSource.GetSentencesforKanji(options.InputPhrase, options.MaxSentences)
	if options.SentencesOnFrontMode {
		for _, sentence := range sentences {
			card := kyoro.makeSentenceAnkiCard(sentence, options)
			anki.AddCard(card)
		}
	} else {
		meaning := meaningSource.GetMeaningforKanji(options.InputPhrase)
		card := kyoro.makeKeywordAnkiCard(options, meaning, sentences)
		anki.AddCard(card)
	}

	return true
}
