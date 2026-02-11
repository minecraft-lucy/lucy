package tools

import (
	"encoding/json"
)

type OneOrMore[T any] []T

func (s *OneOrMore[T]) UnmarshalJSON(data []byte) error {
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []T{single}
		return nil
	}

	var multiple []T
	if err := json.Unmarshal(data, &multiple); err != nil {
		return err
	}
	*s = multiple
	return nil
}
