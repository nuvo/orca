package determine

import (
	"fmt"
	"io"
	"regexp"
	"strings"

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

// Buildtype represents the determine buildtype command
func Buildtype(out io.Writer) *cobra.Command {
	s := &determineCmd{out: out}

	cmd := &cobra.Command{
		Use:   "buildtype",
		Short: "Determine build type based on path filters",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			relPattern, _ := regexp.Compile(s.releaseRef)

			// If this is not the mainline or a release branch - default type
			if (s.currentRef != s.mainRef) && !relPattern.MatchString(s.currentRef) {
				fmt.Println(s.defaultType)
				return
			}

			// If previous commit returned error - default type
			if s.previousCommit == s.previousCommitErrorIndicator {
				fmt.Println(s.defaultType)
				return
			}

			// Get changed paths
			changedPaths := utils.GetChangedPaths(s.previousCommit)

			// If no paths were changed - default type
			if len(changedPaths) == 0 {
				fmt.Println(s.defaultType)
				return
			}

			// If no path filters are defined - default type
			if len(s.pathFilter) == 0 {
				fmt.Println(s.defaultType)
				return
			}

			// Count lines per path filter
			changedPathsPerFilter, changedPathsPerFilterCount := utils.CountLinesPerPathFilter(s.pathFilter, changedPaths)

			// If not all paths matched filters - default type
			if changedPathsPerFilterCount != len(changedPaths) {
				fmt.Println(s.defaultType)
				return
			}

			// All paths matched a filter
			multipleTypes := ""
			for bt, btPathCount := range changedPathsPerFilter {
				if (!strings.Contains(multipleTypes, bt)) && (btPathCount != 0) {
					multipleTypes = multipleTypes + bt + ";"
				}
			}
			multipleTypes = strings.TrimRight(multipleTypes, ";")

			// If multiple is not allowed and there are multiple - default type
			if (s.allowMultipleTypes == false) && (strings.Contains(multipleTypes, ";")) {
				fmt.Println(s.defaultType)
				return
			}

			fmt.Println(multipleTypes)
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
