package orca

import (
	"fmt"
	"io"

	"orca/pkg/utils"

	"github.com/spf13/cobra"
)

type determineCmd struct {
	defaultType                  string
	pathFilter                   []string
	allowMultipleTypes           bool
	mainRef                      string
	releaseRef                   string
	currentRef                   string
	previousCommit               string
	previousCommitErrorIndicator string

	out io.Writer
}

// NewDetermineBuildtype represents the determine buildtype command
func NewDetermineBuildtype(out io.Writer) *cobra.Command {
	d := &determineCmd{out: out}

	cmd := &cobra.Command{
		Use:   "buildtype",
		Short: "Determine build type based on path filters",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			if !utils.IsMainlineOrReleaseRef(d.currentRef, d.mainRef, d.releaseRef) {
				fmt.Println(d.defaultType)
				return
			}

			if utils.IsCommitError(d.previousCommit, d.previousCommitErrorIndicator) {
				fmt.Println(d.defaultType)
				return
			}

			// If no path filters are defined - default type
			if len(d.pathFilter) == 0 {
				fmt.Println(d.defaultType)
				return
			}

			// Get changed paths
			changedPaths := utils.GetChangedPaths(d.previousCommit)

			// Some paths changed, check against path filters
			buildTypeByPathFilters := utils.GetBuildTypeByPathFilters(d.defaultType, changedPaths, d.pathFilter, d.allowMultipleTypes)
			fmt.Println(buildTypeByPathFilters)
		},
	}

	f := cmd.Flags()

	f.StringVar(&d.defaultType, "default-type", "default", "default build type")
	f.StringSliceVar(&d.pathFilter, "path-filter", []string{}, "path filter (supports multiple), can use regex: path=buildtype")
	f.BoolVar(&d.allowMultipleTypes, "allow-multiple-types", false, "allow multiple build types")
	f.StringVar(&d.mainRef, "main-ref", "master", "name of the reference which is the main line")
	f.StringVar(&d.releaseRef, "rel-ref", "^.*/rel-.*$", "release reference name (or regex)")
	f.StringVar(&d.currentRef, "curr-ref", "", "current reference name")
	f.StringVar(&d.previousCommit, "prev-commit", "", "previous commit for paths comparison")
	f.StringVar(&d.previousCommitErrorIndicator, "prev-commit-error", "E", "identify an error with the previous commit by this string")

	return cmd
}
