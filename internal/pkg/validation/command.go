package validation

import "encoding/json"

func IsStringJSON(data string) bool {
	var marshaled map[string]interface{}

	return json.Unmarshal([]byte(data), &marshaled) == nil
}
