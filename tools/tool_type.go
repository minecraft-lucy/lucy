/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
