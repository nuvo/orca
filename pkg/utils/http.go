package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type PerformRequestOptions struct {
	Method             string
	URL                string
	Headers            []string
	ExpectedStatusCode int
	Data               io.Reader
}

// PerformRequest performs an HTTP request to a given url with an expected status code (to support testing) and returns the body
func PerformRequest(o PerformRequestOptions) []byte {
	req, err := http.NewRequest(o.Method, o.URL, o.Data)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	for _, header := range o.Headers {
		header, value := SplitInTwo(header, ":")
		req.Header.Add(header, value)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != o.ExpectedStatusCode {
		log.Fatal(string(body))
	}

	return body
}
