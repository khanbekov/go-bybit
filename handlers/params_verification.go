package handlers

import (
	"fmt"
)

func ValidateParams(Params map[string]interface{}) error {
	for key, value := range Params {
		if key == "" {
			return fmt.Errorf("empty key found in parameters")
		}
		if value == nil {
			return fmt.Errorf("parameter for key '%s' is nil", key)
		}
	}
	return nil
}
