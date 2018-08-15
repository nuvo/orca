package utils

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// GetBuildTypeByPathFilters determines the build type according to path filters
func GetBuildTypeByPathFilters(defaultType string, changedPaths, pathFilter []string, allowMultipleTypes bool) string {

	// If no paths were changed - default type
	if len(changedPaths) == 0 {
		return defaultType
	}

	// Count lines per path filter
	changedPathsPerFilter, changedPathsPerFilterCount := countLinesPerPathFilter(pathFilter, changedPaths)

	// If not all paths matched filters - default type
	if changedPathsPerFilterCount != len(changedPaths) {
		return defaultType
	}

	multipleTypes := ""
	for bt, btPathCount := range changedPathsPerFilter {
		if (!strings.Contains(multipleTypes, bt)) && (btPathCount != 0) {
			multipleTypes = multipleTypes + bt + ";"
		}
	}
	multipleTypes = strings.TrimRight(multipleTypes, ";")

	// If multiple is not allowed and there are multiple - default type
	if (allowMultipleTypes == false) && (strings.Contains(multipleTypes, ";")) {
		return defaultType
	}

	return multipleTypes
}

// GetChangedPaths compares the current commit (HEAD) with the given commit and returns a list of the paths that were changed between them
func GetChangedPaths(previousCommit string) []string {
	r, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	head, err := r.Head()
	if err != nil {
		panic(err)
	}

	currentCommitTree := getTreeFromHash(head.Hash(), r)
	previousCommitTree := getTreeFromStr(previousCommit, r)
	changes, err := currentCommitTree.Diff(previousCommitTree)
	if err != nil {
		panic(err)
	}

	var changedFiles []string

	for _, change := range changes {
		changedFiles = AddIfNotContained(changedFiles, change.From.Name)
		changedFiles = AddIfNotContained(changedFiles, change.To.Name)
	}

	return changedFiles
}

// IsMainlineOrReleaseRef returns true if this is the mainline or a release branch
func IsMainlineOrReleaseRef(currentRef, mainRef, releaseRef string) bool {
	relPattern, _ := regexp.Compile(releaseRef)
	return (currentRef == mainRef) || relPattern.MatchString(currentRef)
}

// IsCommitError returns true if the commit string equals the error indicator
func IsCommitError(commit, commitErrorIndicator string) bool {
	return commit == commitErrorIndicator
}

// countLinesPerPathFilter get a list of path filters (regex=type) and counts matches from the paths that were changed
func countLinesPerPathFilter(pathFilter []string, changedPaths []string) (changedPathsPerFilter map[string]int, changedPathsPerFilterCount int) {

	changedPathsPerFilter = map[string]int{}
	changedPathsPerFilterCount = 0

	for _, pf := range pathFilter {
		pfSplit := strings.Split(pf, "=")
		pfPath, _ := regexp.Compile(pfSplit[0])
		pfBuildtype := pfSplit[1]

		changedPathsPerFilter[pfBuildtype] = 0

		for _, path := range changedPaths {
			if pfPath.MatchString(path) {
				changedPathsPerFilter[pfBuildtype]++
				changedPathsPerFilterCount++
			}
		}
	}

	return changedPathsPerFilter, changedPathsPerFilterCount
}

func getTreeFromStr(hash string, r *git.Repository) *object.Tree {
	commitHash := plumbing.NewHash(hash)

	return getTreeFromHash(commitHash, r)
}

func getTreeFromHash(hash plumbing.Hash, r *git.Repository) *object.Tree {
	commitObject, err := r.CommitObject(hash)
	if err != nil {
		panic(err)
	}
	commitTree, err := commitObject.Tree()
	if err != nil {
		panic(err)
	}

	return commitTree
}
