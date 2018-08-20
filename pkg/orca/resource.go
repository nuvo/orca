package orca

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"orca/pkg/utils"

	"github.com/spf13/cobra"
)

type resourceCmd struct {
	url            string
	headers        []string
	key            string
	value          string
	offset         int
	errorIndicator string
	printKey       string

	out io.Writer
}

// NewGetResourceCmd represents the get resource command
func NewGetResourceCmd(out io.Writer) *cobra.Command {
	r := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Get a resource from REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			var data []map[string]interface{}
			bytes := utils.PerformRequest("GET", r.url, r.headers, 200)
			if err := json.Unmarshal(bytes, &data); err != nil {
				log.Fatal(err)
			}
			if r.key == "" {
				if r.printKey != "" {
					fmt.Println(r.errorIndicator)
					return
				}
				fmt.Println(string(bytes))
				return
			}
			desiredIndex := -1
			for i := 0; i < len(data); i++ {
				if data[i][r.key] == r.value {
					desiredIndex = i
				}
			}
			if desiredIndex == -1 {
				fmt.Println(r.errorIndicator)
				return
			}
			if desiredIndex+r.offset > len(data) {
				fmt.Println(r.errorIndicator)
				return
			}

			result, _ := json.Marshal(data[desiredIndex+r.offset])
			if r.printKey != "" {
				result, _ = json.Marshal(data[desiredIndex+r.offset][r.printKey])
			}
			fmt.Println(strings.Trim(string(result), "\""))
		},
	}

	f := cmd.Flags()

	f.StringVar(&r.url, "url", "", "url to send the request to")
	f.StringSliceVar(&r.headers, "headers", []string{}, "headers of the request (supports multiple)")
	f.StringVar(&r.key, "key", "", "find the desired object according to this key")
	f.StringVar(&r.value, "value", "", "find the desired object according to to key`s value")
	f.IntVar(&r.offset, "offset", 0, "offset of the desired object from the reference key")
	f.StringVarP(&r.errorIndicator, "error-indicator", "e", "E", "string indicating an error in the request")
	f.StringVarP(&r.printKey, "print-key", "p", "", "key to print. If not specified - prints the response")

	return cmd
}

// NewDeleteResourceCmd represents the delete resource command
func NewDeleteResourceCmd(out io.Writer) *cobra.Command {
	r := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Delete a resource from REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PerformRequest("DELETE", r.url, r.headers, 204)
		},
	}

	f := cmd.Flags()

	f.StringVar(&r.url, "url", "", "url to send the request to")
	f.StringSliceVar(&r.headers, "headers", []string{}, "headers of the request (supports multiple)")

	return cmd
}
