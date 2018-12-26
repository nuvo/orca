package orca

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/nuvo/orca/pkg/utils"

	"github.com/spf13/cobra"
)

type artifactCmd struct {
	url   string
	token string
	file  string

	out io.Writer
}

// NewGetArtifactCmd represents the get artifact command
func NewGetArtifactCmd(out io.Writer) *cobra.Command {
	a := &artifactCmd{out: out}

	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Get an artifact from Artifactory",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if a.url == "" {
				return errors.New("url of file to get has to be defined")
			}
			if a.token == "" {
				return errors.New("artifactory token to use has to be defined")
			}
			if a.file == "" {
				return errors.New("path of file to write has to be defined")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			artifact := utils.PerformRequest(utils.PerformRequestOptions{
				Method:             "GET",
				URL:                a.url,
				Headers:            []string{"X-JFrog-Art-Api:" + a.token},
				ExpectedStatusCode: 200,
			})
			err := ioutil.WriteFile(a.file, artifact, 0644)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	f := cmd.Flags()

	f.StringVar(&a.url, "url", os.Getenv("ORCA_URL"), "url of file to get. Overrides $ORCA_URL")
	f.StringVar(&a.token, "token", os.Getenv("ORCA_TOKEN"), "artifactory token to use. Overrides $ORCA_TOKEN")
	f.StringVar(&a.file, "file", os.Getenv("ORCA_FILE"), "path of file to write. Overrides $ORCA_FILE")

	return cmd
}

// NewDeployArtifactCmd represents the deploy artifact command
func NewDeployArtifactCmd(out io.Writer) *cobra.Command {
	a := &artifactCmd{out: out}

	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Deploy an artifact to Artifactory",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if a.url == "" {
				return errors.New("url of file to deploy has to be defined")
			}
			if a.token == "" {
				return errors.New("artifactory token to use has to be defined")
			}
			if a.file == "" {
				return errors.New("path of file to deploy has to be defined")
			}
			if _, err := os.Stat(a.file); err != nil {
				return errors.New("artifact to deploy does not exist")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.Open(a.file)
			if err != nil {
				log.Fatal(err)
			}
			utils.PerformRequest(utils.PerformRequestOptions{
				Method:             "PUT",
				URL:                a.url,
				Headers:            []string{"X-JFrog-Art-Api:" + a.token},
				ExpectedStatusCode: 201,
				Data:               data,
			})
		},
	}

	f := cmd.Flags()

	f.StringVar(&a.url, "url", os.Getenv("ORCA_URL"), "url of file to deploy. Overrides $ORCA_URL")
	f.StringVar(&a.token, "token", os.Getenv("ORCA_TOKEN"), "artifactory token to use. Overrides $ORCA_TOKEN")
	f.StringVar(&a.file, "file", os.Getenv("ORCA_FILE"), "path of file to deploy. Overrides $ORCA_FILE")

	return cmd
}
