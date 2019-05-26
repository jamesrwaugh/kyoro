package kyoro

import (
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
)

type JishoSentenseRetreiver struct {
}

func (this JishoSentenseRetreiver) buildJapaneseAndReadingStrings(japaneseSentence soup.Root) (string, string) {
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

func (this JishoSentenseRetreiver) GetSentencesforKanji(kanji string, maxSentences int) []Sentence {
	var sentences []Sentence
	url := fmt.Sprintf("https://jisho.org/search/%s %%23sentences", kanji)
	print(url)
	resp, err := soup.Get(url)
	if err != nil {
		print("Didn't Work")
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	foundSentences := doc.FindAll("div", "class", "sentence_content")
	for _, sentence := range foundSentences {
		japaneseSentence := sentence.Find("ul", "class", "japanese_sentence")
		japanseString, readingString := this.buildJapaneseAndReadingStrings(japaneseSentence)
		englishSentence := sentence.Find("div", "class", "english_sentence").Find("span", "class", "english")
		sentences = append(sentences, Sentence{
			Japanese: japanseString,
			English:  englishSentence.Text(),
			Reading:  readingString,
		})
	}
	return sentences
}
