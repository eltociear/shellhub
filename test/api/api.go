package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shellhub-io/shellhub/pkg/models"
	log "github.com/sirupsen/logrus"
)

type API struct {
	Address string
}

func (a *API) IsOnline() bool {
	log.Println("Pinging the API...")

	r, err := http.Get(fmt.Sprintf("http://%s/info", a.Address)) // TODO This is only true to ShellHub community version. Fix it!
	if err != nil {
		return false
	}

	log.Traceln("API has responded!")

	log.Traceln("Checking response status code...")

	log.Traceln(fmt.Sprintf("status code from the IsOnline request to %s is %d", a.Address, r.StatusCode))
	if r.StatusCode != 200 {
		return false
	}

	log.Traceln("The response status code is valid!")

	log.Println("API is online!")

	return true
}
func (a *API) Test(tests []Test) error {
	log.Println("Starting tests...")

	for i, t := range tests {
		log.Println(fmt.Sprintf("[Starting test %s...]", t.Description))
		log.Println(fmt.Sprintf("\tmethod: %s", t.Method))
		log.Println(fmt.Sprintf("\turl: %s", t.Url))

		client := http.Client{}

		log.Traceln("Creating request...")

		request, err := http.NewRequest(t.Method, t.Url, strings.NewReader(t.Cases[i].Body))
		if err != nil {
			err := fmt.Errorf("test: %w", err)
			return err
		}

		log.Traceln("Setting content type...")

		request.Header.Set("Content-Type", t.Kind)

		for j, c := range t.Cases {
			log.Println(fmt.Sprintf("Testing case %d", j))
			log.Println(fmt.Sprintf("\tpath: %s", c.Path))
			log.Println(fmt.Sprintf("\tquery: %s", c.Query))
			log.Println(fmt.Sprintf("\theader: %s", c.Header))
			log.Println(fmt.Sprintf("\tbody: %s", c.Body))

			log.Traceln("Setting headers...")

			for _, h := range c.Header {
				request.Header.Set(h.Key, h.Value)
			}

			log.Traceln("Creating JSON auth structure...")

			auth, err := json.Marshal(Auth{Username: c.Auth.Username, Password: c.Auth.Password})
			if err != nil {
				err := fmt.Errorf("test: %w", err)
				return err
			}

			log.Traceln("Sending auth request to login route...")

			response, err := http.Post(fmt.Sprintf("http://%s/api/login", a.Address), t.Kind, bytes.NewReader(auth))
			if err != nil {
				err := fmt.Errorf("test: %w", err)
				return err
			}

			log.Traceln("Reading from response body to variable...")

			read, err := io.ReadAll(response.Body)
			if err != nil {
				err := fmt.Errorf("test: %w", err)
				return err
			}

			log.Traceln("Converting response body to Golang structure")

			var data models.UserAuthResponse
			err = json.Unmarshal(read, &data)
			if err != nil {
				log.Errorln(err)
				err := fmt.Errorf("test: %w", err)
				return err
			}

			log.Traceln("Success to get authorization data!")

			log.Traceln("Setting authorization header...")

			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.Token))

			log.Traceln("Authorization set!")

			log.Traceln("Doing HTTP request...")

			response, err = client.Do(request)
			if err != nil {
				err := fmt.Errorf("test: %w", err)
				return err
			}

			log.Traceln("Success HTTP request!")

			log.Traceln("Checking status code...")

			if response.StatusCode != c.Expected {
				b, _ := io.ReadAll(response.Body)
				log.Debugln(string(b))
				log.Debugln(fmt.Sprintf("Status code from the test request is %d", response.StatusCode))
				err := fmt.Errorf("test: %w", fmt.Errorf("test %s, failed in case %d", t.Description, j))
				return err
			}

			log.Traceln("Status code expected!")

			log.Println(fmt.Sprintf("Test %s success!", t.Description))
		}
	}

	return nil
}
