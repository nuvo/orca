package orca

import (
	"fmt"
	"io"
	"log"
	"sync"

	"orca/pkg/utils"

	"github.com/spf13/cobra"
)

type envCmd struct {
	chartsFile   string
	name         string
	packedValues []string
	set          []string
	kubeContext  string
	tls          bool
	helmTLSStore string
	museum       string

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

			desiredReleases := utils.ChartsYamlToStruct(e.chartsFile, e.name)
			installedReleases := utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, true)
			releasesToInstall := utils.GetReleasesDelta(desiredReleases, installedReleases)

			for len(releasesToInstall) > 0 {

				mutex.Lock()
				for _, c := range releasesToInstall {

					wg.Add(1)
					go func(c utils.ReleaseSpec) {
						defer wg.Done()

						// If there are (still) any dependencies - leave this chart for a later iteration
						if len(c.Dependencies) != 0 {
							return
						}

						// Find index of chart in slice
						// may have changed by now since we are using go routines
						// If chart was not found - another routine is taking care of it
						mutex.Lock()
						index := utils.GetChartIndex(releasesToInstall, c.ChartName)
						if index == -1 {
							mutex.Unlock()
							return
						}
						releasesToInstall = utils.RemoveChartFromCharts(releasesToInstall, index)
						mutex.Unlock()

						// deploy chart
						log.Println("deploying chart", c.ChartName, "version", c.ChartVersion)
						utils.DeployChartFromMuseum(c.ReleaseName, c.ChartName, c.ChartVersion, e.kubeContext, e.name, e.museum, e.helmTLSStore, e.tls, e.packedValues, e.set, false)
						log.Println("deployed chart", c.ChartName, "version", c.ChartVersion)

						// Deployment is done, remove chart from dependencies
						mutex.Lock()
						releasesToInstall = utils.RemoveChartFromDependencies(releasesToInstall, c.ChartName)
						mutex.Unlock()

					}(c)
				}
				mutex.Unlock()
			}
			wg.Wait()

			installedReleases = utils.GetInstalledReleases(e.kubeContext, e.name, e.helmTLSStore, e.tls, true)
			releasesToDelete := utils.GetReleasesDelta(installedReleases, desiredReleases)

			for _, c := range releasesToDelete {
				wg.Add(1)
				go func(c utils.ReleaseSpec) {
					defer wg.Done()
					log.Println("deleting", c.ReleaseName)
					utils.DeleteRelease(c.ReleaseName, e.kubeContext, e.tls, e.helmTLSStore, false)
					log.Println("deleted", c.ReleaseName)
				}(c)
			}
			wg.Wait()
		},
	}

	f := cmd.Flags()

	f.StringVarP(&e.chartsFile, "charts-file", "c", "", "path to file with list of Helm charts to install")
	f.StringVar(&e.name, "name", "", "name of environment (namespace) to deploy to")
	f.StringVar(&e.museum, "museum", "", "chart museum instance (name=url)")
	f.StringVar(&e.kubeContext, "kube-context", "", "kubernetes context to deploy to")
	f.StringSliceVarP(&e.packedValues, "values", "f", []string{}, "values file to use (packaged within the chart)")
	f.StringSliceVarP(&e.set, "set", "s", []string{}, "set additional parameters")
	f.BoolVar(&e.tls, "tls", false, "should use communication over TLS")
	f.StringVar(&e.helmTLSStore, "helm-tls-store", "", "directory with TLS certs and keys")

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
