package utils

import (
	"testing"
)

func TestAuthentication(t *testing.T) {
	_, err := TokenFromAuthentication("username1", "senha")
	if err != nil {
		t.Errorf("error")
	}
}
