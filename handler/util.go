package handler

import (
	"errors"
	"fmt"

	postory "github.com/andynaguyen/postory-server"
)

func ValidateInput(carrier string) error {
	if !postory.IsCarrierSupported(carrier) {
		msg := fmt.Sprintf("carrier is unsupported: %s", carrier)
		return errors.New(msg)
	}
	return nil
}
