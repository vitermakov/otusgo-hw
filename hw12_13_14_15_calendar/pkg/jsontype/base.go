package jsontype

import (
	"encoding/json"
	"errors"
)

const QuotesByte = 34

func UnmarshalJSON(data []byte, val interface{}) error {
	var err error
	if data[0] == QuotesByte {
		err = json.Unmarshal(data[1:len(data)-1], val)
	} else {
		err = json.Unmarshal(data, val)
	}
	if err != nil {
		return errors.New("Value " + string(data) + " invalid")
	}
	return nil
}
