package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jamesrwaugh/kyoro/acquisition"
	"github.com/jamesrwaugh/kyoro/anki"
	"github.com/jamesrwaugh/kyoro/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type KyoroTestEnvironment struct {
	Kyoro       Kyoro
	AnkiClient  *anki.MockHttpClient
	Anki        anki.AnkiService
	AcquiMockRC *acquisition.MockResourceClient
	Sentences   acquisition.SentenceRetriever
	Meanings    acquisition.MeaningRetriever
	Verifier    verification.SentenceVerifier
}

func (e *KyoroTestEnvironment) RunKyoro(options *Options) bool {
	return e.Kyoro.Kyoro(
		*options,
		e.Anki,
		e.Sentences,
		e.Meanings,
		e.Verifier,
	)
}

func makeKyoroTestEnvironment() (env *KyoroTestEnvironment) {
	//
	mrc := &acquisition.MockResourceClient{}
	jisho := acquisition.NewJishoSentenceretriever(mrc)
	jmdict := acquisition.NewJdictMeaningRetriever(mrc)
	mvf := &verification.MockSentenceVerifier{}

	//
	ankiConnectClient := &anki.MockHttpClient{}
	ankiConnect := anki.NewAnkiConnect(ankiConnectClient, "の.の", 50)

	return &KyoroTestEnvironment{
		Kyoro:       NewKyoro(),
		AnkiClient:  ankiConnectClient,
		Anki:        ankiConnect,
		AcquiMockRC: mrc,
		Sentences:   jisho,
		Meanings:    jmdict,
		Verifier:    mvf,
	}
}

func makeKyoroSentenceModeEnvironment() (env *KyoroTestEnvironment) {
	env = makeKyoroTestEnvironment()
	env.AnkiClient.On("Get", mock.Anything).Return(makeResponse(200), nil)
	env.AnkiClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(makeResponse(200), nil)
	env.AcquiMockRC.On("Get", mock.Anything).Return(MomTest)
	return
}

func makeResponse(code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(bytes.NewBufferString("残念")),
	}
}

func Test_Kyoro_AnkiNotConnected_ReturnsFalse(t *testing.T) {
	env := makeKyoroTestEnvironment()
	env.AnkiClient.On("Get", mock.Anything).Return(makeResponse(500), nil)
	result := env.RunKyoro(&Options{})
	assert.False(t, result)
}

var MomTest = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<HTML>
	<body>
		<ul class="sentences">
			<li class="entry sentence clearfix">
				<div class="sentence_content">
					<ul class="japanese_sentence japanese japanese_gothic clearfix" lang="ja">
						<li class="clearfix">
							<span class="unlinked">母ちゃん</span>
						</li>
					</ul>
					<div class="english_sentence clearfix">
						<span class="english">Mom</span>
					</div>
				</div>
			</li>
		</ul>
	</body>
</HTML>`

func Test_Kyoro_SentencesOnFrontMode_AddsExpectedCardCount(t *testing.T) {
	env := makeKyoroSentenceModeEnvironment()
	opts := Options{
		InputPhrase:          "ルナ",
		SentencesOnFrontMode: true,
		MaxSentences:         5,
	}
	result := env.RunKyoro(&opts)
	assert.Equal(t, true, result)
	env.AnkiClient.AssertNumberOfCalls(t, "Post", opts.MaxSentences)
}

func Test_Kyoro_NoMonolingualMode_EnglishSubmitted(t *testing.T) {
	env := makeKyoroSentenceModeEnvironment()
	opts := Options{
		InputPhrase:          "ルナ",
		SentencesOnFrontMode: true,
		MonoligualMode:       false,
		MaxSentences:         1,
	}
	result := env.RunKyoro(&opts)
	assert.Equal(t, true, result)
	env.AnkiClient.AssertCalled(
		t,
		"Post",
		"の.の:50",
		"application/json",
		bytes.NewBufferString("{\"action\":\"addNote\",\"params\":{\"note\":{\"deckName\":\"\",\"fields\":{\"english\":\"Mom\",\"japanese\":\"母ちゃん\",\"reading\":\"母ちゃん\"},\"modelName\":\"\",\"options\":{\"allowDuplicates\":false},\"tags\":[\"kyoro\"]}},\"version\":6}"),
	)
}

func Test_Kyoro_MonolingualMode_NoEnglishSubmitted(t *testing.T) {
	env := makeKyoroSentenceModeEnvironment()
	opts := Options{
		InputPhrase:          "ルナ",
		SentencesOnFrontMode: true,
		MonoligualMode:       true,
		MaxSentences:         1,
	}
	result := env.RunKyoro(&opts)
	assert.Equal(t, true, result)
	env.AnkiClient.AssertCalled(
		t,
		"Post",
		"の.の:50",
		"application/json",
		bytes.NewBufferString("{\"action\":\"addNote\",\"params\":{\"note\":{\"deckName\":\"\",\"fields\":{\"japanese\":\"母ちゃん\",\"reading\":\"母ちゃん\"},\"modelName\":\"\",\"options\":{\"allowDuplicates\":false},\"tags\":[\"kyoro\"]}},\"version\":6}"),
	)
}

var 文がないテスト string = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<HTML>
    <HEAD>
        <META http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <TITLE>WWWJDIC: Word Display</TITLE>
    </HEAD>
    <BODY>
        <br>&nbsp;
        <br>
        <pre></pre>
    </BODY>
</HTML>`

func Test_Kyoro_NotSentencesOnFrontMode_AddsExpectedCardCount(t *testing.T) {
	//env := makeKyoroTestEnvironment()
}

/*func Test_Kyoro_SentencesOnFrontMode_CardHasExpectedFields(t *testing.T) {
	env := makeKyoroTestEnvironment()
}

func Test_Kyoro_NotSentencesOnFrontMode_CardHasExpectedFields(t *testing.T) {
	env := makeKyoroTestEnvironment()
}
*/
