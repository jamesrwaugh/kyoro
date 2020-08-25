package anki

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/jamesrwaugh/kyoro/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeAnkiConnectTestObjects() (service AnkiService, client *MockHTTPClient) {
	client = &MockHTTPClient{}
	logger := log.New(testutils.SilentWriter{}, "", log.LstdFlags)
	service = NewAnkiConnect(client, "の.の", 50, logger)
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

func Test_AddCardWithError_RetriesToAddCard(t *testing.T) {

}
