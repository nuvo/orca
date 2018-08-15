package orca

import (
	"fmt"
	"io"
	"sync"
	"time"

	"orca/pkg/utils"

	"github.com/spf13/cobra"
)

type envCmd struct {
	chartsFile string

	nada string

	out io.Writer
}

// NewGetEnvCmd represents the get env command
func NewGetEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Get list of Helm releases in an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get env called")
		},
	}

	f := cmd.Flags()

	f.StringVar(&e.nada, "nada", "", "nada help")

	return cmd
}

// NewDeployEnvCmd represents the deploy env command
func NewDeployEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Deploy a list of Helm charts to an environment (Kubernetes namespace)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			var mutex = &sync.Mutex{}
			var wg sync.WaitGroup

			charts := utils.ChartsYamlToStruct(e.chartsFile)
			for len(charts) > 0 {

				mutex.Lock()
				for _, c := range charts {

					wg.Add(1)
					go func(c utils.ChartSpec) {
						defer wg.Done()

						// If there are (still) any dependencies - leave this chart for a later iteration
						if len(c.Dependencies) != 0 {
							return
						}

						mutex.Lock()
						// Find index of chart in slice
						// may have changed by now since we are using go routines
						// If chart was not found - another routine is taking care of it
						index := utils.GetChartIndex(charts, c.Name)
						if index == -1 {
							mutex.Unlock()
							return
						}

						charts = utils.RemoveChartFromCharts(charts, index)
						mutex.Unlock()

						// deploy chart
						fmt.Println(c.Name, "deployment: In progress")
						time.Sleep(5 * time.Second)
						fmt.Println(c.Name, "deployment: Done")

						// Deployment is done, remove chart from dependencies
						mutex.Lock()
						charts = utils.RemoveChartFromDependencies(charts, c.Name)
						mutex.Unlock()

					}(c)
				}
				mutex.Unlock()
			}
			wg.Wait()
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.chartsFile, "charts-file", "c", "", "path to file with list of Helm charts to install")

	return cmd
}

// NewDeleteEnvCmd represents the delete env command
func NewDeleteEnvCmd(out io.Writer) *cobra.Command {
	e := &envCmd{out: out}

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Delete an environment (Kubernetes namespace) along with all Helm releases in it",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete env called")

		},
	}

	f := cmd.Flags()

	f.StringVar(&e.nada, "nada", "", "nada help")

	return cmd
}
