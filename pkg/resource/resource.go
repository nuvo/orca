package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	httputils "orca/pkg/utils/http"

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

// NewGetCmd represents the get resource command
func NewGetCmd(out io.Writer) *cobra.Command {
	s := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Get a resource from REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			var data []map[string]interface{}
			bytes := httputils.PerformRequest("GET", s.url, s.headers, 200)
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
	f.StringSliceVar(&s.headers, "headers", []string{}, "headers of the request (supports multiple)")
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
