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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type AnkiConnect struct {
	client   HttpClient
	hostname string
	port     int64
}

func NewAnkiConnect(client HttpClient, host string, port int64) AnkiService {
	a := &AnkiConnect{client, host, port}
	return a
}

func (ac AnkiConnect) responseHasError(resp *http.Response, err error) bool {
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

func (ac AnkiConnect) GetEndpoint() string {
	return ac.hostname + ":" + strconv.FormatInt(ac.port, 10)
}

func (ac AnkiConnect) IsConnected() bool {
	r, err := ac.client.Get(ac.GetEndpoint())
	return (r != nil) && (r.StatusCode == 200) && (err == nil)
}

func (ac AnkiConnect) MaxSentencesPerCard() int {
	return 1
}

func (ac AnkiConnect) HasMiaCardModel() bool {
	request := map[string]interface{}{
		"action":  "modelNames",
		"version": 6,
	}
	resp, _ := ac.getAnkiConnectJSONResponse(request)
	modelNames, _ := resp["result"].([]string)
	return contains(modelNames, "MIA Japanese")
}

func (ac AnkiConnect) AddCard(card AnkiCard) bool {
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

	_, err := ac.getAnkiConnectJSONResponse(request)

	// It seems adding cards too quickly results in strange issues.
	// Make sure that doesn't happen here.
	time.Sleep(500)

	return err == nil
}

func (ac AnkiConnect) getAnkiConnectJSONResponse(request map[string]interface{}) (map[string]interface{}, error) {
	requestString, _ := json.Marshal(request)
	resp, err := ac.client.Post(ac.GetEndpoint(), "application/json", bytes.NewBuffer(requestString))
	if ac.responseHasError(resp, err) {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var ankiConnectResponse map[string]interface{}
	json.Unmarshal(bodyBytes, &ankiConnectResponse)
	return ankiConnectResponse, nil
}
