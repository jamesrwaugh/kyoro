package kyoro

import "net/http"
import "strconv"

type AnkiConnect struct {
	hostname string
	port     int64
}

func (this AnkiConnect) isConnected() bool {
	r, err := http.Get(this.hostname + ":" + strconv.FormatInt(this.port, 10))
	return (r.StatusCode == 200) && (err == nil)
}

func (this AnkiConnect) addCard(card AnkiCard) bool {
	/*note, _ := json.Marshal(card)
	requst := map[string]interface{}{
		"action":  "addNote",
		"version": 6,
		"params": map[string]interface{}{
			"note": note,
		},
	}
	_, err := gozenity.List(
		"Choose an option:",
		"One word",
		"Two",
		"Three things",
	)
	print(err)*/
	return true
}
