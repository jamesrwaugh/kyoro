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

// AnkiConnect is interface to Anki Connect
// https://foosoft.net/projects/anki-connect/
type AnkiConnect struct {
	client   HTTPClient
	hostname string
	port     int64
	logger   *log.Logger
}

// NewAnkiConnect does what it says
func NewAnkiConnect(client HTTPClient, host string, port int64, logger *log.Logger) AnkiService {
	a := &AnkiConnect{client, host, port, logger}
	return a
}

func (ac AnkiConnect) getResponseError(resp *http.Response) (hasError bool, errorMessage string) {
	// Check AnkiConnect's response for the "error" property.
	// On error, it is a string, and null otherwise.
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var ankiConnectResponse map[string]interface{}
	json.Unmarshal(bodyBytes, &ankiConnectResponse)
	if errorString, ok := ankiConnectResponse["error"].(string); ok {
		return true, errorString
	}
	return false, ""
}

func (ac AnkiConnect) responseHasError(resp *http.Response, err error) (hasError bool, errorMessage string) {
	// Fail early if the requst itself failed.
	if err != nil {
		return true, err.Error()
	}
	hasError, errorMessage = ac.getResponseError(resp)
	return
}

// GetEndpoint returns the endpoint of where AnkiConnect is connecting a a display-full string
func (ac AnkiConnect) GetEndpoint() string {
	return ac.hostname + ":" + strconv.FormatInt(ac.port, 10)
}

// IsConnected tests to see if it is, in fact, connected
func (ac AnkiConnect) IsConnected() bool {
	r, err := ac.client.Get(ac.GetEndpoint())
	return (r != nil) && (r.StatusCode == 200) && (err == nil)
}

// HasMiaCardModel returns true if we have the MIA Japanese card from the
// Mass Immersion Approach. This is to provide rich support for this type as a standard format.
func (ac AnkiConnect) HasMiaCardModel() bool {
	request := map[string]interface{}{
		"action":  "modelNames",
		"version": 6,
	}
	resp, _ := ac.getAnkiConnectJSONResponse(request, 3)
	modelNames, _ := resp["result"].([]string)
	return contains(modelNames, "Migaku Japanese")
}

// AddCard will add the given card to Anki
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

	_, err := ac.getAnkiConnectJSONResponse(request, 3)

	// It seems adding cards too quickly results in strange issues.
	// Make sure that doesn't happen here.
	time.Sleep(2000)

	return err == nil
}

func (ac AnkiConnect) getAnkiConnectJSONResponse(request map[string]interface{}, retryCount int) (map[string]interface{}, error) {
	var err error
	requestString, _ := json.Marshal(request)
	for i := 0; i < retryCount; i++ {
		resp, err := ac.client.Post(ac.GetEndpoint(), "application/json", bytes.NewBuffer(requestString))
		defer resp.Body.Close()
		if hasError, errorMessage := ac.responseHasError(resp, err); hasError {
			if errorMessage == "cannot create note because it is a duplicate" {
				ac.logger.Println(errorMessage)
				break
			}
			if i < retryCount {
				ac.logger.Println("Received Error, Retying: " + errorMessage)
				time.Sleep(500)
				continue
			} else {
				break
			}
		}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var ankiConnectResponse map[string]interface{}
		json.Unmarshal(bodyBytes, &ankiConnectResponse)
		return ankiConnectResponse, nil
	}
	return nil, err
}
