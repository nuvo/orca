package orca

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/maorfr/orca/pkg/utils"

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
	method         string
	update         bool

	out io.Writer
}

// NewCreateResourceCmd represents the create resource command
func NewCreateResourceCmd(out io.Writer) *cobra.Command {
	r := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Create or update a resource via REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			method := r.method
			if r.update {
				method = "PATCH"
			}
			utils.PerformRequest(utils.PerformRequestOptions{
				Method:             method,
				URL:                r.url,
				Headers:            r.headers,
				ExpectedStatusCode: 201,
				Data:               nil,
			})
		},
	}

	f := cmd.Flags()

	f.StringVar(&r.url, "url", os.Getenv("ORCA_URL"), "url to send the request to. Overrides $ORCA_URL")
	f.StringVar(&r.method, "method", utils.GetStringEnvVar("ORCA_METHOD", "POST"), "method to use in the request. Overrides $ORCA_METHOD")
	f.BoolVar(&r.update, "update", utils.GetBoolEnvVar("ORCA_UPDATE", false), "should method be PUT instead of POST. Overrides $ORCA_UPDATE")
	f.StringSliceVar(&r.headers, "headers", []string{}, "headers of the request (supports multiple)")

	return cmd
}

// NewGetResourceCmd represents the get resource command
func NewGetResourceCmd(out io.Writer) *cobra.Command {
	r := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Get a resource via REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			var data []map[string]interface{}
			bytes := utils.PerformRequest(utils.PerformRequestOptions{
				Method:             "GET",
				URL:                r.url,
				Headers:            r.headers,
				ExpectedStatusCode: 200,
				Data:               nil,
			})
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
					break
				}
			}
			if desiredIndex == -1 {
				fmt.Println(r.errorIndicator)
				return
			}
			if desiredIndex+r.offset >= len(data) {
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

	f.StringVar(&r.url, "url", os.Getenv("ORCA_URL"), "url to send the request to. Overrides $ORCA_URL")
	f.StringSliceVar(&r.headers, "headers", []string{}, "headers of the request (supports multiple)")
	f.StringVar(&r.key, "key", os.Getenv("ORCA_KEY"), "find the desired object according to this key. Overrides $ORCA_KEY")
	f.StringVar(&r.value, "value", os.Getenv("ORCA_VALUE"), "find the desired object according to to key`s value. Overrides $ORCA_VALUE")
	f.IntVar(&r.offset, "offset", utils.GetIntEnvVar("ORCA_OFFSET", 0), "offset of the desired object from the reference key. Overrides $ORCA_OFFSET")
	f.StringVarP(&r.errorIndicator, "error-indicator", "e", utils.GetStringEnvVar("ORCA_ERROR_INDICATOR", "E"), "string indicating an error in the request. Overrides $ORCA_ERROR_INDICATOR")
	f.StringVarP(&r.printKey, "print-key", "p", os.Getenv("ORCA_PRINT_KEY"), "key to print. If not specified - prints the response. Overrides $ORCA_PRINT_KEY")

	return cmd
}

// NewDeleteResourceCmd represents the delete resource command
func NewDeleteResourceCmd(out io.Writer) *cobra.Command {
	r := &resourceCmd{out: out}

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Delete a resource via REST API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			utils.PerformRequest(utils.PerformRequestOptions{
				Method:             "DELETE",
				URL:                r.url,
				Headers:            r.headers,
				ExpectedStatusCode: 204,
				Data:               nil,
			})
		},
	}

	f := cmd.Flags()

	f.StringVar(&r.url, "url", os.Getenv("ORCA_URL"), "url to send the request to. Overrides $ORCA_URL")
	f.StringSliceVar(&r.headers, "headers", []string{}, "headers of the request (supports multiple)")

	return cmd
}
