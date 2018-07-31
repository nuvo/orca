package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type resourceCmd struct {
	url            string
	headers        string
	key            string
	value          string
	offset         int
	errorIndicator string
	printKey       string

	out io.Writer
}

// NewGetCmd represents the get resource command
func NewGetCmd(out io.Writer) *cobra.Command {
	s := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Get a resource from REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			var data []map[string]interface{}
			bytes := performRequest(s.url, s.headers)
			if err := json.Unmarshal(bytes, &data); err != nil {
				panic(err)
			}
			if s.key == "" {
				if s.printKey != "" {
					fmt.Println(s.errorIndicator)
					return
				}
				fmt.Println(string(bytes))
				return
			}
			desiredIndex := -1
			for i := 0; i < len(data); i++ {
				if data[i][s.key] == s.value {
					desiredIndex = i
				}
			}
			if desiredIndex == -1 {
				fmt.Println(s.errorIndicator)
				return
			}
			if desiredIndex+s.offset > len(data) {
				fmt.Println(s.errorIndicator)
				return
			}

			result, _ := json.Marshal(data[desiredIndex+s.offset])
			if s.printKey != "" {
				result, _ = json.Marshal(data[desiredIndex+s.offset][s.printKey])
			}
			fmt.Println(strings.Trim(string(result), "\""))
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.url, "url", "", "url to send the request to")
	f.StringVar(&s.headers, "headers", "", "headers help")
	f.StringVar(&s.key, "key", "", "find the desired object according to this key")
	f.StringVar(&s.value, "value", "", "find the desired object according to to key`s value")
	f.IntVar(&s.offset, "offset", 0, "offset of the desired object from the reference key")
	f.StringVarP(&s.errorIndicator, "error-indicator", "e", "E", "url to send the request to")
	f.StringVarP(&s.printKey, "print-key", "p", "", "url to send the request to")

	return cmd
}

// NewDeleteCmd represents the delete resource command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	s := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Delete a resource from REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete resource called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.url, "url", "", "url help")

	return cmd
}

func performRequest(url string, headers string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	headersSplit := strings.Split(headers, ":")
	header, value := headersSplit[0], headersSplit[1]
	req.Header.Set(header, value)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	return body
}
