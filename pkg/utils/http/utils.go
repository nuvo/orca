package httputils

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func PerformRequest(url string, headers []string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	for _, header := range headers {
		headersSplit := strings.Split(header, ":")
		header, value := headersSplit[0], headersSplit[1]
		req.Header.Add(header, value)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Fatal(string(body))
	}

	return body
}
