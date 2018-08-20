package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

// PerformRequest performs an HTTP request to a given url with an expected status code (to support testing) and returns the body
func PerformRequest(method, url string, headers []string, expectedStatusCode int) []byte {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	for _, header := range headers {
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
	if res.StatusCode != expectedStatusCode {
		log.Fatal(string(body))
	}

	return body
}
