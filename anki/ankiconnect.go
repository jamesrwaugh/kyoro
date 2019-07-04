package anki

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AnkiConnect struct {
	client   HttpClient
	Hostname string
	Port     int64
}

func NewAnkiConnect(client HttpClient, host string, port int64) AnkiService {
	a := &AnkiConnect{client, host, port}
	return a
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
		log.Println("Recieved AnkiConnect error: ", errorString)
		return true
	}
	return false
}

func (this AnkiConnect) GetEndpoint() string {
	return this.Hostname + ":" + strconv.FormatInt(this.Port, 10)
}

func (this AnkiConnect) IsConnected() bool {
	r, err := this.client.Get(this.GetEndpoint())
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

	resp, err := this.client.Post(this.GetEndpoint(), "application/json", bytes.NewBuffer(requestString))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// It seems adding cards too quickly results in strange issues.
	// Make sure that doesn't happen here.
	time.Sleep(500)

	return !this.responseHasError(resp, err)
}
