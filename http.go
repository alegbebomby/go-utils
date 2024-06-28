package library

import (
	"bytes"
	"context"
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

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] =  "application/json"
	logHeaders["Accept"] =  "application/json"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
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

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",url,strings.Join(heads,"\n"),string(jsonData),st,strings.Join(rheads,","),string(body))
	}

	return st, string(body)
}

func HTTPPostWithContext(ctx context.Context, url string, headers map[string] string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] =  "application/json"
	logHeaders["Accept"] =  "application/json"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
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

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",url,strings.Join(heads,"\n"),string(jsonData),st,strings.Join(rheads,","),string(body))
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

		log.Printf(" Request endpoint %s ", endpoint)

	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] =  "application/json"
	logHeaders["Accept"] =  "application/json"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
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

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",endpoint,strings.Join(heads,"\n"),"",st,strings.Join(rheads,","),string(body))
	}


	return st, string(body)
}

func HTTPGetWithContext(ctx context.Context,remoteURL string, headers map[string] string, payload map[string]string) (httpStatus int, response string) {

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

	req, err := http.NewRequestWithContext(ctx,"GET", endpoint, nil)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] =  "application/json"
	logHeaders["Accept"] =  "application/json"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
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

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",endpoint,strings.Join(heads,"\n"),"",st,strings.Join(rheads,","),string(body))
	}

	return st, string(body)
}

func HTTPFormPost(endpoint string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	method := "POST"

	var stringPayload []string

	if payload != nil {

		for key, value := range payload {

			stringPayload = append(stringPayload, fmt.Sprintf("%s=%v", key, value))

		}

	}

	requestPayload := strings.NewReader(strings.Join(stringPayload, "&"))

	req, err := http.NewRequest(method, endpoint, requestPayload)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	logHeaders := make(map[string]string)
	logHeaders["Content-Type"] =  "application/x-www-form-urlencoded"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",endpoint,strings.Join(heads,"\n"),strings.Join(stringPayload, "&"),st,strings.Join(rheads,","),string(body))
	}

	return st, string(body)
}

func HTTPFormPostWithContext(ctx context.Context,endpoint string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	method := "POST"

	var stringPayload []string

	if payload != nil {

		for key, value := range payload {

			stringPayload = append(stringPayload, fmt.Sprintf("%s=%v", key, value))

		}

	}

	requestPayload := strings.NewReader(strings.Join(stringPayload, "&"))

	req, err := http.NewRequestWithContext(ctx, method, endpoint, requestPayload)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	logHeaders := make(map[string]string)
	logHeaders["Content-Type"] =  "application/x-www-form-urlencoded"

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
			logHeaders[k] = v
		}
	}


	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		responseHeaders := make(map[string]string)

		for k,v := range resp.Header {

			responseHeaders[k] = strings.Join(v,",")

		}

		var heads,rheads []string
		for k,v := range logHeaders {

			heads = append(heads, fmt.Sprintf("\t%s : %s",k,v))
		}

		for k,v := range responseHeaders {

			rheads = append(rheads, fmt.Sprintf("\t%s : %s",k,v))
		}

		log.Printf("**** BEGIN HTTP REQUEST ****\n" +
			"Remote Url : %s\n" +
			"Request Headers:\n" +
			"%s\n" +
			"Request Payload\n" +
			"\t%s\n" +
			"Response Status: %d\n" +
			"Response Headers\n" +
			"%s\n" +
			"Response Body\n" +
			"**** BEGIN HTTP REQUEST ****" +
			"\t%s",endpoint,strings.Join(heads,"\n"),strings.Join(stringPayload, "&"),st,strings.Join(rheads,","),string(body))
	}

	return st, string(body)
}
