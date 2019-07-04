package kyoro

import (
	"fmt"
	"log"
	"strings"

	"github.com/anaskhan96/soup"
)

func NewJishoSentenceRetreiver(c ResourceClient) *JishoSentenceRetreiver {
	j := JishoSentenceRetreiver{c}
	return &j
}

type JishoSentenceRetreiver struct {
	client ResourceClient
}

func (this JishoSentenceRetreiver) buildJapaneseAndReadingStrings(japaneseSentence soup.Root) (japanese string, reading string, kaniReadings []string) {
	for _, element := range japaneseSentence.Children() {
		nodeValue := strings.TrimSpace(element.NodeValue)
		if nodeValue == "li" {
			elementText := element.Find("span", "class", "unlinked").Text()
			japanese += elementText
			reading += elementText
			furigana := element.Find("span", "class", "furigana")
			if furigana.Error == nil {
				quotedFurigana := "「" + furigana.Text() + "」"
				reading += quotedFurigana
				kaniReadings = append(kaniReadings, elementText+quotedFurigana)
			}
		} else if len(nodeValue) > 0 {
			// We expect to get punctuation here. jisho adds things like 、 and 。
			// outside of <li> elements.
			japanese += nodeValue
			reading += nodeValue
		}
	}
	return
}

func (this JishoSentenceRetreiver) addSentencesFromPage(foundSentences []soup.Root, sentences *[]Translation, maxSentences int) {
	for _, sentence := range foundSentences {
		if len(*sentences) >= maxSentences {
			break
		}
		japaneseSentence := sentence.Find("ul", "class", "japanese_sentence")
		japanseString, readingString, kaniReadings := this.buildJapaneseAndReadingStrings(japaneseSentence)
		englishSentence := sentence.Find("div", "class", "english_sentence").Find("span", "class", "english")
		*sentences = append(*sentences, Translation{
			Japanese:      japanseString,
			English:       englishSentence.Text(),
			Reading:       readingString,
			KanjiReadings: kaniReadings,
		})
	}
}

func (this JishoSentenceRetreiver) GetSentencesforKanji(kanji string, maxSentences int) []Translation {
	var sentences []Translation
	for pageNumber := 1; len(sentences) < maxSentences; pageNumber++ {
		url := fmt.Sprintf("https://jisho.org/search/%s %%23sentences?page=%d", kanji, pageNumber)
		log.Println("Looking for sentences on " + url)
		resp, _ := this.client.Get(url)
		doc := soup.HTMLParse(resp)
		foundSentences := doc.FindAll("div", "class", "sentence_content")
		if len(foundSentences) == 0 {
			break
		}
		this.addSentencesFromPage(foundSentences, &sentences, maxSentences)
	}
	return sentences
}
