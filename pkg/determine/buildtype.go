package determine

import (
	"fmt"
	"io"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/maorfr/orca/pkg/utils"
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

			changedPaths := utils.GetChangedPaths(s.previousCommit)

			fmt.Println(changedPaths)
			// fmt.Println("ERROR")
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
