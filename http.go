package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HTTPPost(url string, headers map[string] string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s",err.Error())
		return st,""
	}

	return st, string(body)
}

func HTTPGet(remoteURL string, headers map[string] string, payload map[string]string) (httpStatus int, response string) {

	var fields []string

	if payload != nil {

		for key, value := range payload {

			val := fmt.Sprintf("%s=%v", key, url.QueryEscape(value))

			fields = append(fields, val)
		}
	}

	params := strings.Join(fields, "&")

	endpoint := fmt.Sprintf("%s?%s", remoteURL, params)

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		log.Printf(" Wants to GET data to URL %s ", endpoint)

	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s",err.Error())
		return st,""
	}

	return st, string(body)
}
