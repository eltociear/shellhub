package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/shellhub-io/shellhub/pkg/models"
	log "github.com/sirupsen/logrus"
)

type Loader struct{}

func load(p string) []byte {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		log.Errorln(fmt.Errorf("loader: %w", err))
		return nil
	}

	return f
}

func (l *Loader) LoadUsers(path string) ([]models.User, error) {
	b := load(path)

	us := new([]models.User)
	err := json.Unmarshal(b, &us)
	if err != nil {
		log.Errorln(fmt.Errorf("loader: %w", err))
		return nil, err
	}

	return *us, nil
}

func (l *Loader) LoadNamespaces(path string) ([]models.Namespace, error) {
	b := load(path)

	us := new([]models.Namespace)
	err := json.Unmarshal(b, &us)
	if err != nil {
		log.Errorln(fmt.Errorf("loader: %w", err))
		return nil, err
	}

	return *us, nil
}
