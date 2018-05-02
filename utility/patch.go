package utility

import (
	"encoding/json"
	"fmt"
	"github.com/bitspill/json-patch"
)

func UnSquashPatch(spb []byte) (jsonpatch.Patch, error) {
	fmt.Println("unsquash")
	var sp map[string][]*json.RawMessage
	var up jsonpatch.Patch

	err := json.Unmarshal(spb, &sp)
	if err != nil {
		return up, err
	}

	if val, ok := sp["remove"]; ok {
		o := json.RawMessage([]byte(`"remove"`))
		for _, value := range val {
			var op = make(map[string]*json.RawMessage)
			op["op"] = &o
			op["path"] = value
			up = append(up, op)
		}
	}

	var op map[string]*json.RawMessage

	if val, ok := sp["add"]; ok {
		o := json.RawMessage([]byte(`"add"`))
		for _, value := range val {
			err = json.Unmarshal(*value, &op)
			op["op"] = &o
			up = append(up, op)
		}
	}

	if val, ok := sp["replace"]; ok {
		o := json.RawMessage([]byte(`"replace"`))
		for _, value := range val {
			err = json.Unmarshal(*value, &op)
			op["op"] = &o
			up = append(up, op)
		}
	}

	if val, ok := sp["move"]; ok {
		o := json.RawMessage([]byte(`"move"`))
		for _, value := range val {
			err = json.Unmarshal(*value, &op)
			op["op"] = &o
			up = append(up, op)
		}
	}

	if val, ok := sp["copy"]; ok {
		o := json.RawMessage([]byte(`"copy"`))
		for _, value := range val {
			err = json.Unmarshal(*value, &op)
			op["op"] = &o
			up = append(up, op)
		}
	}

	if val, ok := sp["test"]; ok {
		o := json.RawMessage([]byte(`"test"`))
		for _, value := range val {
			err = json.Unmarshal(*value, &op)
			op["op"] = &o
			up = append(up, op)
		}
	}

	return up, err
}
