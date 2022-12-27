package jsontype

import (
	"encoding/json"
)

type String string

func (value *String) UnmarshalJSON(data []byte) (err error) {
	var val string
	if data[0] != QuotesByte {
		data = []byte(string(byte(QuotesByte)) + string(data) + string(byte(QuotesByte)))
	}
	if err = json.Unmarshal(data, &val); err == nil {
		*value = String(val)
	}
	return
}
