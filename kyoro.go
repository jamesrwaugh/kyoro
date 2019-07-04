package kyoro

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/jamesrwaugh/kyoro/acquisition"
	"github.com/jamesrwaugh/kyoro/anki"
)

type KyoroProduction struct {
}

func (this KyoroProduction) makeSentenceAnkiCard(
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

func (this KyoroProduction) generateKeywordBackHTML(
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

func (this KyoroProduction) makeKeywordAnkiCard(
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
	backMatter, _ := this.generateKeywordBackHTML(sentences, options.MonoligualMode)
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

func (this KyoroProduction) Kyoro(
	options Options,
	anki anki.AnkiService,
	sentenceSource acquisition.SentenceRetriever,
	meaningSource acquisition.MeaningRetriever,
) bool {
	if !anki.IsConnected() {
		log.Fatal("Could not connect to Anki. Failing.")
		return false
	}
	sentences := sentenceSource.GetSentencesforKanji(options.InputPhrase, options.MaxSentences)
	if options.SentencesOnFrontMode {
		for _, sentence := range sentences {
			card := this.makeSentenceAnkiCard(sentence, options)
			anki.AddCard(card)
		}
	} else {
		meaning := meaningSource.GetMeaningforKanji(options.InputPhrase)
		card := this.makeKeywordAnkiCard(options, meaning, sentences)
		anki.AddCard(card)
	}

	return true
}
