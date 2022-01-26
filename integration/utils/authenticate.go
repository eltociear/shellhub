package utils

import (
	"bytes"
	"fmt"
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const BaseURL = "http://127.0.0.1/api"

func TokenFromAuthentication(username, password string) (string, error) {
	log.Println("Getting authorization bearer token...")

	bytesPayload, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/login", BaseURL), "application/json", bytes.NewBuffer(bytesPayload))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token, err := GetFieldValueFromJSON("token", string(body))
	if err != nil {
		return "", err
	}

	return token, nil
}

func TenantSwitch(token, tenant string) (string, error) {
	log.Println("switching to namespace...")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/auth/token/%s", BaseURL, tenant), nil)

	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Content-type":  []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        log.Println("deu pau 1")
		return "", err
	}

	token, err = GetFieldValueFromJSON("token", string(body))
	if err != nil {
        log.Println("deu pau 2")
        log.Println(string(body))
		return "", err
	}

	return token, nil
}
