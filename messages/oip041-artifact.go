package messages

import "encoding/json"

func (o Oip041) GetJSON() (string, error) {
	// ToDo: remove redundant Storage items, potentially cache?
	var s string

	s, err := json.Marshal(o)

	return s, err
}
