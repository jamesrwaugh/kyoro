package kyoro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type AnkiConnect struct {
	Hostname string
	Port     int64
}

func (this AnkiConnect) responseHasError(resp *http.Response, err error) bool {
	// Fail early if the requst itself failed.
	if err != nil {
		log.Fatal("Recieved network error: ", err.Error())
		return true
	}
	// Check AnkiConnect's response for the "error" property.
	// On error, it is a string, and null otherwise.
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var ankiConnectResponse map[string]interface{}
	json.Unmarshal(bodyBytes, &ankiConnectResponse)
	if errorString, ok := ankiConnectResponse["error"].(string); ok {
		log.Fatal("Recieved AnkiConnect error: ", errorString)
		return true
	}
	return false
}

func (this AnkiConnect) GetEndpoint() string {
	return this.Hostname + ":" + strconv.FormatInt(this.Port, 10)
}

func (this AnkiConnect) IsConnected() bool {
	r, err := http.Get(this.GetEndpoint())
	return (r.StatusCode == 200) && (err == nil)
}

func (this AnkiConnect) MaxSentencesPerCard() int {
	return 1
}

func (this AnkiConnect) AddCard(card AnkiCard) bool {
	request := map[string]interface{}{
		"action":  "addNote",
		"version": 6,
		"params": map[string]interface{}{
			"note": map[string]interface{}{
				"deckName":  card.DeckName,
				"modelName": card.ModelName,
				"options": map[string]interface{}{
					"allowDuplicates": false,
				},
				"fields": card.Fields,
				"tags":   card.Tags,
			},
		},
	}
	requestString, _ := json.Marshal(request)
	log.Println(string(requestString))

	resp, err := http.Post(this.GetEndpoint(), "application/json", bytes.NewBuffer(requestString))
	defer resp.Body.Close()

	return !this.responseHasError(resp, err)
}
