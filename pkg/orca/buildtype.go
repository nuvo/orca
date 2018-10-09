package orca

import (
	"fmt"
	"io"
	"os"

	"github.com/maorfr/orca/pkg/utils"

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

	f.StringVar(&d.defaultType, "default-type", utils.GetStringEnvVar("ORCA_DEFAULT_TYPE", "default"), "default build type. Overrides $ORCA_DEFAULT_TYPE")
	f.StringSliceVar(&d.pathFilter, "path-filter", []string{}, "path filter (supports multiple) in the path=buildtype form (supports regex)")
	f.BoolVar(&d.allowMultipleTypes, "allow-multiple-types", utils.GetBoolEnvVar("ORCA_ALLOW_MULTIPLE_TYPES", false), "allow multiple build types. Overrides $ORCA_ALLOW_MULTIPLE_TYPES")
	f.StringVar(&d.mainRef, "main-ref", os.Getenv("ORCA_MAIN_REF"), "name of the reference which is the main line. Overrides $ORCA_MAIN_REF")
	f.StringVar(&d.releaseRef, "rel-ref", os.Getenv("ORCA_REL_REF"), "release reference name (or regex). Overrides $ORCA_REL_REF")
	f.StringVar(&d.currentRef, "curr-ref", os.Getenv("ORCA_CURR_REF"), "current reference name. Overrides $ORCA_CURR_REF")
	f.StringVar(&d.previousCommit, "prev-commit", os.Getenv("ORCA_PREV_COMMIT"), "previous commit for paths comparison. Overrides $ORCA_PREV_COMMIT")
	f.StringVar(&d.previousCommitErrorIndicator, "prev-commit-error", utils.GetStringEnvVar("ORCA_PREV_COMMIT_ERROR", "E"), "identify an error with the previous commit by this string. Overrides $ORCA_PREV_COMMIT_ERROR")

	return cmd
}
