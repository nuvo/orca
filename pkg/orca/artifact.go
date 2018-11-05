package orca

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/maorfr/orca/pkg/utils"

	"github.com/spf13/cobra"
)

type artifactCmd struct {
	url      string
	token    string
	artifact string

	out io.Writer
}

// NewDeployArtifactCmd represents the deploy artifact command
func NewDeployArtifactCmd(out io.Writer) *cobra.Command {
	a := &artifactCmd{out: out}

	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Deploy artifact to Artifactory",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if a.url == "" {
				return errors.New("url to deploy to has to be defined")
			}
			if a.token == "" {
				return errors.New("token to use for deployment has to be defined")
			}
			if a.artifact == "" {
				return errors.New("artifact to deploy has to be defined")
			}
			if _, err := os.Stat(a.artifact); err != nil {
				return errors.New("artifact to deploy does not exist")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.Open(a.artifact)
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

	f.StringVar(&a.url, "url", os.Getenv("ORCA_URL"), "url to deploy to. Overrides $ORCA_URL")
	f.StringVar(&a.token, "token", os.Getenv("ORCA_TOKEN"), "token to use for deployment. Overrides $ORCA_TOKEN")
	f.StringVar(&a.artifact, "artifact", os.Getenv("ORCA_FILE"), "path to artifact to deploy. Overrides $ORCA_ARTIFACT")

	return cmd
}
