package tools

import (
	"encoding/json"
	"fmt"
)

type StringOrStringSlice struct {
	Value []string
}

func (s *StringOrStringSlice) UnmarshalJSON(data []byte) error {
	var singleString string
	if err := json.Unmarshal(data, &singleString); err == nil {
		s.Value = []string{singleString}
		return nil
	}

	var stringSlice []string
	if err := json.Unmarshal(data, &stringSlice); err == nil {
		s.Value = stringSlice
		return nil
	}

	return fmt.Errorf("invalid data for StringOrStringSlice: %s", data)
}
