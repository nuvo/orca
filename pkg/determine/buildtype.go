package determine

import (
	"fmt"
	"io"

	gitutils "orca/pkg/utils/git"

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

// Buildtype represents the determine buildtype command
func Buildtype(out io.Writer) *cobra.Command {
	s := &determineCmd{out: out}

	cmd := &cobra.Command{
		Use:   "buildtype",
		Short: "Determine build type based on path filters",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			if !gitutils.IsMainlineOrReleaseRef(s.currentRef, s.mainRef, s.releaseRef) {
				fmt.Println(s.defaultType)
				return
			}

			if gitutils.IsCommitError(s.previousCommit, s.previousCommitErrorIndicator) {
				fmt.Println(s.defaultType)
				return
			}

			// If no path filters are defined - default type
			if len(s.pathFilter) == 0 {
				fmt.Println(s.defaultType)
				return
			}

			// Get changed paths
			changedPaths := gitutils.GetChangedPaths(s.previousCommit)

			// Some paths changed, check against path filters
			buildTypeByPathFilters := gitutils.GetBuildTypeByPathFilters(s.defaultType, changedPaths, s.pathFilter, s.allowMultipleTypes)
			fmt.Println(buildTypeByPathFilters)
		},
	}

	f := cmd.Flags()

	f.StringVar(&s.defaultType, "default-type", "default", "default build type (default: default)")
	f.StringSliceVar(&s.pathFilter, "path-filter", []string{}, "path filter (supports multiple), can use regex: path=buildtype")
	f.BoolVar(&s.allowMultipleTypes, "allow-multiple-types", false, "allow multiple build types (default: false)")
	f.StringVar(&s.mainRef, "main-ref", "master", "name of the reference which is the main line (default: master)")
	f.StringVar(&s.releaseRef, "rel-ref", "^.*/rel-.*$", "release reference name (or regex) (default: ^.*/rel-.*$)")
	f.StringVar(&s.currentRef, "curr-ref", "", "current reference name")
	f.StringVar(&s.previousCommit, "prev-commit", "", "previous commit for paths comparison")
	f.StringVar(&s.previousCommitErrorIndicator, "prev-commit-error", "E", "identify an error with the previous commit by this string (default: E)")

	return cmd
}
