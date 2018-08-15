package env

import (
	"fmt"
	"io"
	"sync"
	"time"

	chartutils "orca/pkg/utils/chart"

	"github.com/spf13/cobra"
)

type envCmd struct {
	chartsFile string

	nada string

	out io.Writer
}

// NewGetCmd represents the get env command
func NewGetCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Get list of Helm releases in an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}

// NewDeployCmd represents the deploy env command
func NewDeployCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Deploy a list of Helm charts to an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			charts := chartutils.ChartsYamlToStruct(s.chartsFile)

			var mutex = &sync.Mutex{}

			var wg sync.WaitGroup
			for len(charts) > 0 {

				mutex.Lock()
				for _, c := range charts {

					wg.Add(1)
					go func(c chartutils.ChartSpec) {
						defer wg.Done()
						if len(c.Dependencies) != 0 {
							return
						}

						mutex.Lock()
						// Find index of chart in slice (may have changed by now since we are using go routines)
						index := -1
						for i, elem := range charts {
							if elem.Name == c.Name {
								index = i
							}
						}
						// If chart was not found - another routine is taking care of it
						if index == -1 {
							mutex.Unlock()
							return
						}

						// Remove chart from charts list
						charts[index] = charts[len(charts)-1]
						charts = charts[:len(charts)-1]

						mutex.Unlock()

						// deploy chart
						fmt.Println(c.Name, "deployment: In progress")
						time.Sleep(5 * time.Second)
						fmt.Println(c.Name, "deployment: Done")

						// Deployment is done, remove chart from dependencies
						mutex.Lock()
						charts = chartutils.RemoveChartFromDependencies(charts, c.Name)
						mutex.Unlock()

					}(c)
				}
				mutex.Unlock()
			}
			wg.Wait()
		},
	}

	f := cmd.Flags()

	f.StringVarP(&s.chartsFile, "charts-file", "c", "", "path to file with list of Helm charts to install")

	return cmd
}

// NewDeleteCmd represents the delete env command
func NewDeleteCmd(out io.Writer) *cobra.Command {
	s := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Delete an environment (Kubernetes namespace) along with all Helm releases in it",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete env called")

		},
	}

	f := cmd.Flags()

	f.StringVar(&s.nada, "nada", "", "nada help")

	return cmd
}
