package kyoro

import (
	"fmt"
	"log"

	"github.com/anaskhan96/soup"
)

type JishoSentenceRetreiver struct {
}

func (this JishoSentenceRetreiver) buildJapaneseAndReadingStrings(japaneseSentence soup.Root) (string, string) {
	var japanseString string
	var readingString string
	elements := japaneseSentence.FindAll("li")
	for _, element := range elements {
		elementText := element.Find("span", "class", "unlinked").Text()
		japanseString += elementText
		readingString += elementText
		furigana := element.Find("span", "class", "furigana")
		if furigana.Error == nil {
			readingString += "「" + furigana.Text() + "」"
		}
	}
	return japanseString, readingString
}

func (this JishoSentenceRetreiver) addSentencesFromPage(foundSentences []soup.Root, sentences *[]Sentence, maxSentences int) {
	for _, sentence := range foundSentences {
		if len(*sentences) >= maxSentences {
			break
		}
		japaneseSentence := sentence.Find("ul", "class", "japanese_sentence")
		japanseString, readingString := this.buildJapaneseAndReadingStrings(japaneseSentence)
		englishSentence := sentence.Find("div", "class", "english_sentence").Find("span", "class", "english")
		*sentences = append(*sentences, Sentence{
			Japanese: japanseString,
			English:  englishSentence.Text(),
			Reading:  readingString,
		})
	}
}

func (this JishoSentenceRetreiver) GetSentencesforKanji(kanji string, maxSentences int) []Sentence {
	var sentences []Sentence
	for pageNumber := 1; len(sentences) < maxSentences; pageNumber++ {
		url := fmt.Sprintf("https://jisho.org/search/%s %%23sentences?page=%d", kanji, pageNumber)
		log.Println("Looking for sentences on " + url)
		resp, _ := soup.Get(url)
		doc := soup.HTMLParse(resp)
		foundSentences := doc.FindAll("div", "class", "sentence_content")
		if len(foundSentences) == 0 {
			break
		}
		this.addSentencesFromPage(foundSentences, &sentences, maxSentences)
	}
	return sentences
}
