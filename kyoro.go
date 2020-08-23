package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/jamesrwaugh/kyoro/acquisition"
	"github.com/jamesrwaugh/kyoro/anki"
	"github.com/jamesrwaugh/kyoro/verification"
)

// KyoroProduction is the actual real thing.
type KyoroProduction struct {
	ankiService    anki.AnkiService
	sentenceSource acquisition.SentenceRetriever
	meaningSource  acquisition.MeaningRetriever
	verifier       verification.SentenceVerifier
	logger         *log.Logger
}

// NewKyoro creates a new Kyoro object.
func NewKyoro(
	ankiService anki.AnkiService,
	sentenceSource acquisition.SentenceRetriever,
	meaningSource acquisition.MeaningRetriever,
	verifier verification.SentenceVerifier,
	logger *log.Logger) (instance Kyoro) {
	return &KyoroProduction{ankiService, sentenceSource, meaningSource, verifier, logger}
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

// Creates a sentence card for the MIA Sentence Card format.
// https://massimmersionapproach.com/table-of-contents/anki/mia-japanese-addon/
func (kyoro KyoroProduction) makeMiaSentenceAnkiCard(sentence acquisition.Translation, options Options) anki.AnkiCard {
	boldInputPhrase := fmt.Sprintf("<b>%s</b>", options.InputPhrase)
	boldJapanese := strings.Replace(sentence.Japanese, options.InputPhrase, boldInputPhrase, -1)
	cardFields := map[string]string{
		"Expression": boldJapanese,
		"Meaning":    sentence.English,
	}
	return anki.AnkiCard{
		DeckName:  options.DeckName,
		ModelName: "MIA Japanese",
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

// Acquire sentences and have the user confirm each. Return the ones that were accetped.
// These accepted sentences should be added to Anki.
func (kyoro KyoroProduction) getUserConfirmedSentences(options Options) []acquisition.Translation {
	var acceptedSentences []acquisition.Translation
	for len(acceptedSentences) < options.MaxSentences {
		// TODO: Fix Infinite loop when there are sentences, but the user has rejected
		// all of them. In this case, the "len(sentences) == 0" will never be true, and
		// the above condition will also never be true.
		// To solve this and other issues, it would be best to have a "start from" index
		// in the sentence interface.
		sentences := kyoro.sentenceSource.GetSentencesforKanji(options.InputPhrase, 10*options.MaxSentences)
		if len(sentences) == 0 {
			break
		}
		for index, sentence := range sentences {
			confirmed := kyoro.verifier.UserConfirmSentence(&verification.SentenceVerificationInfo{
				AcceptedSentences: acceptedSentences,
				Sentence:          sentence,
				SentenceIndex:     index,
			})
			if confirmed {
				acceptedSentences = append(acceptedSentences, sentence)
			}
			if len(acceptedSentences) >= options.MaxSentences {
				break
			}
		}
	}
	return acceptedSentences
}

// Run runs the main procedure of Kyoro from the command line,
// and adds cards accordingly.
func (kyoro KyoroProduction) Run(options Options) bool {
	if !kyoro.ankiService.IsConnected() {
		kyoro.logger.Println("Could not connect to Anki. Failing.")
		return false
	}
	sentences := kyoro.getUserConfirmedSentences(options)
	if options.SentencesOnFrontMode {
		for _, sentence := range sentences {
			var card anki.AnkiCard
			// TODO: This should be handled elsewhere
			// - Default the model name to MIA?
			// - Move this to "makeSentenceAnkiCard"?
			if options.ModelName == "MIA Japanese" {
				card = kyoro.makeMiaSentenceAnkiCard(sentence, options)
			} else {
				card = kyoro.makeSentenceAnkiCard(sentence, options)
			}

			kyoro.ankiService.AddCard(card)
		}
	} else {
		meaning := kyoro.meaningSource.GetMeaningforKanji(options.InputPhrase)
		card := kyoro.makeKeywordAnkiCard(options, meaning, sentences)
		kyoro.ankiService.AddCard(card)
	}

	return true
}
