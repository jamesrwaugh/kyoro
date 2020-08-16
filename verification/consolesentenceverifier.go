package verification

import (
	"bufio"
	"fmt"
	"os"
)

// NewConsoleSentenceVerifier creates a new ConsoleSentenceVerifier
// to verifiy sentences using the console.
func NewConsoleSentenceVerifier() *ConsoleSentenceVerifier {
	v := ConsoleSentenceVerifier{}
	return &v
}

// ConsoleSentenceVerifier retrieves a sentence from Jisho.org
type ConsoleSentenceVerifier struct {
}

// UserConfirmSentence Confirms a sentece using the console.
func (c ConsoleSentenceVerifier) UserConfirmSentence(info *SentenceVerificationInfo) bool {
	fmt.Printf("(Got %d) %d. [Y/n] '%s' '%s' ",
		len(info.AcceptedSentences),
		info.SentenceIndex,
		info.Sentence.Japanese,
		info.Sentence.English)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	return (len(text) == 0 || text == "y")
}
