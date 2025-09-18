package utils

import "encoding/json"

func DeepCopy(a, b any) error {
	bytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, b)
}
