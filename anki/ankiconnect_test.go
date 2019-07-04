package anki

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeAnkiConnectTestObjects() (service AnkiService, client *MockHttpClient) {
	client = &MockHttpClient{}
	service = NewAnkiConnect(client, "の.の", 50)
	return
}

func makeResponse(code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(bytes.NewBufferString("えぇぇ")),
	}
}

func Test_IsConnected_WhenConnected_ReturnsTrue(t *testing.T) {
	anki, mhc := makeAnkiConnectTestObjects()
	mhc.On("Get", mock.Anything).Return(makeResponse(200), nil)
	assert.True(t, anki.IsConnected())
}

func Test_IsConnected_WhenNon200Code_ReturnsFalse(t *testing.T) {
	anki, mhc := makeAnkiConnectTestObjects()
	mhc.On("Get", mock.Anything).Return(makeResponse(500), nil)
	assert.False(t, anki.IsConnected())
}

func Test_IsConnected_WhenRequestReturnsError_ReturnsFalse(t *testing.T) {
	anki, mhc := makeAnkiConnectTestObjects()
	mhc.On("Get", mock.Anything).Return(makeResponse(200), errors.New(":("))
	assert.False(t, anki.IsConnected())
}

func Test_MaxSentencesPerCard_Is_1(t *testing.T) {
	anki, _ := makeAnkiConnectTestObjects()
	assert.Equal(t, 1, anki.MaxSentencesPerCard())
}

func Test_AddCard_WithValidCardAndNoError_MakesCorrectPostResponse(t *testing.T) {
	anki, mhc := makeAnkiConnectTestObjects()
	mhc.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(makeResponse(200), nil)

	card := AnkiCard{
		DeckName:  "A",
		ModelName: "B",
		Fields: map[string]string{
			"F1": "の・の",
			"F2": "Field2",
		},
		Tags: []string{
			"T1", "T2",
		},
	}
	anki.AddCard(card)

	mhc.AssertCalled(
		t,
		"Post",
		"の.の:50",
		"application/json",
		bytes.NewBufferString("{\"action\":\"addNote\",\"params\":{\"note\":{\"deckName\":\"A\",\"fields\":{\"F1\":\"の・の\",\"F2\":\"Field2\"},\"modelName\":\"B\",\"options\":{\"allowDuplicates\":false},\"tags\":[\"T1\",\"T2\"]}},\"version\":6}"),
	)
}
