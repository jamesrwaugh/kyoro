package kyoro_test

import "testing"
import "github.com/jamesrwaugh/kyoro"
import "github.com/jamesrwaugh/kyoro/resourceclient"
import "github.com/stretchr/testify/assert"

var 先生テスト string = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<HTML>
    <HEAD>
        <META http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <TITLE>WWWJDIC: Word Display</TITLE>
    </HEAD>
    <BODY>
        <br>&nbsp;
        <br>
        <pre>
先生 [シーサン] /(n) (hon) boy (chi: xiānshēng)/
先生 [せんせい(P);せんじょう(ok)] /(n) (1) (hon) teacher/master/doctor/(suf) (2) (hon) with names of teachers, etc. as an honorific/(n) (3) (せんじょう only) (arch) (See 前生) previous existence/(P)/
先生に就く [せんせいにつく] /(exp,v5k) to study under (a teacher)/
先生の述 [せんせいのじゅつ] /(n) teachers statement (expounding)/
先生方 [せんせいがた] /(n) doctors/teachers/
</pre>
    </BODY>
</HTML>`

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

func TestGetMeaningforKanji(t *testing.T) {
	cases := []struct {
		in              string
		englishReading  string
		japaneseText    string
		japaneseReading string
	}{
		{"先生", "(hon) boy (chi: xiānshēng)/", "先生", "シーサン"}, // This isn't right. Fix it.
		{"なし", "", "", ""},
		{"ない", "", "", ""},
	}

	mrc := resourceclient.MockResourceClient{}
	mrc.On("Get", "http://nihongo.monash.edu/cgi-bin/wwwjdic?1ZUJ%E5%85%88%E7%94%9F").Return(先生テスト)
	mrc.On("Get", "http://nihongo.monash.edu/cgi-bin/wwwjdic?1ZUJ%E3%81%AA%E3%81%97").Return("")
	mrc.On("Get", "http://nihongo.monash.edu/cgi-bin/wwwjdic?1ZUJ%E3%81%AA%E3%81%84").Return(文がないテスト)
	jdict := kyoro.NewJdictMeaningRetriever(&mrc)

	for _, c := range cases {
		result := jdict.GetMeaningforKanji(c.in)
		assert.Equal(t, c.englishReading, result.English)
		assert.Equal(t, c.japaneseText, result.Japanese)
		assert.Equal(t, c.japaneseReading, result.Reading)
	}
}
