package kyoro_test

import (
	"fmt"
	"testing"

	"github.com/jamesrwaugh/kyoro"
	"github.com/jamesrwaugh/kyoro/resourceclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var EmptySentenceTest string = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<HTML>
<body>
<ul class="sentences">
    <li class="entry sentence clearfix"></li>
    <li class="entry sentence clearfix"></li>
</ul>
</body>
</HTML>`

var MomTest string = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<HTML>
<body>
<ul class="sentences">
    <li class="entry sentence clearfix">
        <div class="debug">74143</div>
        <div class="sentence_content">
            <ul class="japanese_sentence japanese japanese_gothic clearfix" lang="ja">
                <li class="clearfix">
                    <span class="furigana">かあ</span>
                    <span class="unlinked">母ちゃん</span>
                </li>
                <li class="clearfix">
                    <span class="unlinked">も</span>
                </li>
                <li class="clearfix">
                    <span class="furigana">おな</span>
                    <span class="unlinked">同じ</span>
                </li>
                <li class="clearfix">
                    <span class="furigana">こと</span>
                    <span class="unlinked">事</span>
                </li>
            </ul>
            <div class="english_sentence clearfix">
                <span class="english">Mum said the same thing. But, so what? It's got nothing to do with me.</span>
                <span class="inline_copyright">— 
                    <a href="http://tatoeba.org/eng/sentences/show/74143">Tatoeba</a>
                </span>
            </div>
        </div>
        <a class="light-details_link" href="//jisho.org/sentences/51866352d5dda7e9810000bc">Details ▸</a>
    </li>
    <li class="entry sentence clearfix">
        <div class="debug">75143</div>
        <div class="sentence_content">
            <ul class="japanese_sentence japanese japanese_gothic clearfix" lang="ja">
                <li class="clearfix">
                    <span class="furigana">わたし</span>
                    <span class="unlinked">私</span>
                </li>
                <li class="clearfix">
                    <span class="unlinked">も</span>
                </li>、
                <li class="clearfix">
                    <span class="unlinked">ほんとに</span>
                </li>。
            </ul>
            <div class="english_sentence clearfix">
                <span class="english">I also really had the feeling of having had a relaxed day with my family for the first time in a long while.</span>
                <span class="inline_copyright">— 
                    <a href="http://tatoeba.org/eng/sentences/show/75143">Tatoeba</a>
                </span>
            </div>
        </div>
        <a class="light-details_link" href="//jisho.org/sentences/5186637bd5dda7e9810004a2">Details ▸</a>
    </li>
</ul>
</body>
</HTML>`

func makeJishoURL(term string, pageNumber int) string {
	return fmt.Sprintf("https://jisho.org/search/%s %%23sentences?page=%d", term, pageNumber)
}

func newJishoTestObjects() (jisho *kyoro.JishoSentenceRetreiver, mrc *resourceclient.MockResourceClient) {
	mrc = &resourceclient.MockResourceClient{}
	jisho = kyoro.NewJishoSentenceRetreiver(mrc)
	return
}

func Test_GetSentencesforKanji_EmptyResponsesAreOK(t *testing.T) {
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", mock.AnythingOfType("string")).Return("")
	results := jisho.GetSentencesforKanji("何も", 5)
	assert.Equal(t, 0, len(results))
}

func Test_GetSentencesforKanji_ResponsesWithNoSentencesAreOK(t *testing.T) {
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", mock.AnythingOfType("string")).Return(EmptySentenceTest)
	results := jisho.GetSentencesforKanji("何も", 5)
	assert.Equal(t, 0, len(results))
}

func Test_GetSentencesforKanji_ParsesHTMLCorrectly(t *testing.T) {
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", mock.AnythingOfType("string")).Return(MomTest)
	results := jisho.GetSentencesforKanji("何も", 2)

	assert := assert.New(t)
	assert.Equal(2, len(results))
	assert.Equal("Mum said the same thing. But, so what? It's got nothing to do with me.", results[0].English)
	assert.Equal("母ちゃんも同じ事", results[0].Japanese)
	assert.Equal("母ちゃん「かあ」も同じ「おな」事「こと」", results[0].Reading)
	assert.Equal(3, len(results[0].KanjiReadings))
	assert.Equal("I also really had the feeling of having had a relaxed day with my family for the first time in a long while.", results[1].English)
	assert.Equal("私も、ほんとに。", results[1].Japanese)
	assert.Equal("私「わたし」も、ほんとに。", results[1].Reading)
	assert.Equal(1, len(results[1].KanjiReadings))
}

func Test_GetSentencesforKanji_MakesMultipageRequestsForSentences_EnoughAvailiable(t *testing.T) {
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", mock.AnythingOfType("string")).Return(MomTest, nil)
	jisho.GetSentencesforKanji("の", 5)
	mrc.AssertCalled(t, "Get", makeJishoURL("の", 1))
	mrc.AssertCalled(t, "Get", makeJishoURL("の", 2))
}

func Test_GetSentencesforKanji_MakesMultipageRequestsForSentences_NotEnoughSentencesExist(t *testing.T) {
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", makeJishoURL("何も", 1)).Return(MomTest)
	mrc.On("Get", makeJishoURL("何も", 2)).Return(EmptySentenceTest)
	jisho.GetSentencesforKanji("何も", 5)
	mrc.AssertCalled(t, "Get", makeJishoURL("何も", 1))
	mrc.AssertCalled(t, "Get", makeJishoURL("何も", 2))
}

func Test_Regression_Issue1(t *testing.T) {
	Issue1Test := `
    <ul class="sentences">
       <li class="entry sentence clearfix">
          <div class="debug">74182</div>
          <div class="sentence_content">
             <ul class="japanese_sentence japanese japanese_gothic clearfix" lang="ja">
                <li class="clearfix"><span class="furigana">けっきょく</span><span class="unlinked">結局</span></li>
                、
                <li class="clearfix"><span class="furigana">ほうあん</span><span class="unlinked">法案</span></li>
                <li class="clearfix"><span class="unlinked">は</span></li>
                <li class="clearfix"><span class="unlinked">提出断念</span></li>
                <li class="clearfix"><span class="unlinked">に</span></li>
                <li class="clearfix"><span class="furigana">おいこ</span><span class="unlinked">追い込まれた</span></li>
                <li class="clearfix"><span class="unlinked">のだった</span></li>
                。
             </ul>
             <div class="english_sentence clearfix">
                <span class="english">In the end the bill was forced into being withdrawn.</span>
                <span class="inline_copyright">— <a href="http://tatoeba.org/eng/sentences/show/74182">Tatoeba</a></span>
             </div>
          </div>
          <a class="light-details_link" href="//jisho.org/sentences/51866354d5dda7e9810000e3">Details ▸</a>
       </li>
    </ul>
    `
	jisho, mrc := newJishoTestObjects()
	mrc.On("Get", mock.AnythingOfType("string")).Return(Issue1Test, nil)
	results := jisho.GetSentencesforKanji("", 1)
	assert.Equal(t, "結局、法案は提出断念に追い込まれたのだった。", results[0].Japanese)
}
