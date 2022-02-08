package api

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Init(address string, tests []Test) error {
	a := API{
		Address: address,
	}

	if !a.IsOnline() {
		err := fmt.Errorf("api: %w", errors.New("the API is not online"))
		log.Errorln(err)
		return err
	}

	err := a.Test(tests)
	if err != nil {
		err := fmt.Errorf("api: %w", err)
		log.Errorln(err)
		return err
	}

	return nil
}
