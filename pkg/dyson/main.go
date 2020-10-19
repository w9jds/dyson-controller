package dyson

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const baseAPI = "https://appapi.cp.dyson.com"

type authorization struct {
	Account  string
	Password string
}

// Client is an API Client for the Dyson Api
type Client struct {
	account *authorization
	Client  *http.Client
}

func (dysonAPI *Client) setHeaders(request *http.Request) {
	request.Header.Add("User-Agent", "Mozilla/5.0")
	request.Header.Set("Content-Type", "application/json")

	if dysonAPI.account != nil {
		encoded := base64.StdEncoding.EncodeToString(
			[]byte(dysonAPI.account.Account + ":" + dysonAPI.account.Password),
		)

		request.Header.Add("Authorization", "Basic "+encoded)
	}
}

// Login to the dyson api with your credentials
func (dysonAPI *Client) Login(email string, password string, country string) {
	payload := &map[string]string{
		"Email":    email,
		"Password": password,
	}

	buffer, err := json.Marshal(payload)
	if err != nil {
		log.Panic("Failed to build auth payload")
	}

	url := fmt.Sprintf("%s/v1/userregistration/authenticate?country=%s", baseAPI, country)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(buffer))
	if err != nil {
		log.Panic(err)
	}

	resp, err := dysonAPI.do(request)
	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(resp, &dysonAPI.account); err != nil {
		log.Panic(err)
	}
}

// GetDevices returns a list of devices attached to your dyson account
func (dysonAPI *Client) GetDevices() ([]map[string]interface{}, error) {
	request, err := http.NewRequest("GET", baseAPI+"/v2/provisioningservice/manifest", nil)
	if err != nil {
		return nil, err
	}

	resp, err := dysonAPI.do(request)
	if err != nil {
		return nil, err
	}

	var devices []map[string]interface{}
	if err := json.Unmarshal(resp, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

func (dysonAPI *Client) do(request *http.Request) ([]byte, error) {
	dysonAPI.setHeaders(request)

	response, err := dysonAPI.Client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode < 200 || response.StatusCode > 299 {
		if response.StatusCode == 403 || response.StatusCode == 401 {
			return nil, fmt.Errorf("Status %d: Unauthorized", response.StatusCode)
		}

		message, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.New("Request to dyson api failed")
		}

		return nil, errors.New(string(message))
	} else {
		return ioutil.ReadAll(response.Body)
	}
}
