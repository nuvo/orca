package gitutils

import (
	"regexp"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func CountLinesPerPathFilter(pathFilter []string, changedPaths []string) (changedPathsPerFilter map[string]int, changedPathsPerFilterCount int) {

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

func GetChangedPaths(previousCommit string) []string {
	r, _ := git.PlainOpen(".")
	head, _ := r.Head()

	currentCommitHash := head.Hash()
	currentCommitObject, _ := r.CommitObject(currentCommitHash)
	currentCommitTree, _ := currentCommitObject.Tree()

	previousCommitHash := plumbing.NewHash(previousCommit)
	previousCommitObject, _ := r.CommitObject(previousCommitHash)
	previousCommitTree, _ := previousCommitObject.Tree()

	changes, _ := currentCommitTree.Diff(previousCommitTree)

	var changedFiles []string

	for _, change := range changes {
		if (!contains(changedFiles, change.From.Name)) && (change.From.Name != "") {
			changedFiles = append(changedFiles, change.From.Name)
		}
		if (!contains(changedFiles, change.To.Name)) && (change.To.Name != "") {
			changedFiles = append(changedFiles, change.To.Name)
		}
	}

	return changedFiles
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
