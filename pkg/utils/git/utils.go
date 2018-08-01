package gitutils

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	genutils "orca/pkg/utils/general"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CountLinesPerPathFilter get a list of path filters (regex=type) and counts matches from the paths that were changed
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
		changedFiles = genutils.AddIfNotContained(changedFiles, change.From.Name)
		changedFiles = genutils.AddIfNotContained(changedFiles, change.To.Name)
	}

	return changedFiles
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
