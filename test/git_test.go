package test

import (
	"strconv"
	"testing"

	"github.com/maorfr/orca/pkg/utils"
)

func TestGetBuildTypeByPathFilters_NoChangedPaths(t *testing.T) {
	var changedPaths = []string{}
	var pathFilters = []string{"^src.*$=code", "^kubernetes.*$=chart"}
	bt := utils.GetBuildTypeByPathFilters("default", changedPaths, pathFilters, false)
	if bt != "default" {
		t.Errorf("Expected: default, Actual: " + bt)
	}
}
func TestGetBuildTypeByPathFilters_MultipleNotAllowed_AllMatch(t *testing.T) {
	changedPaths := []string{"src/file1.go", "kubernetes/Chart.yaml"}
	pathFilters := []string{"^src.*$=code", "^kubernetes.*$=chart"}
	bt := utils.GetBuildTypeByPathFilters("default", changedPaths, pathFilters, false)
	if bt != "default" {
		t.Errorf("Expected: default, Actual: " + bt)
	}
}
func TestGetBuildTypeByPathFilters_MultipleAllowed_AllMatch(t *testing.T) {
	changedPaths := []string{"src/file1.go", "kubernetes/Chart.yaml"}
	pathFilters := []string{"^src.*$=code", "^kubernetes.*$=chart"}
	bt := utils.GetBuildTypeByPathFilters("default", changedPaths, pathFilters, true)
	if bt != "code;chart" && bt != "chart;code" {
		t.Errorf("Expected: code;chart, Actual: " + bt)
	}
}
func TestGetBuildTypeByPathFilters_MultipleNotAllowed_NotAllMatch(t *testing.T) {
	changedPaths := []string{"src/file1.go", "kubernetes/Chart.yaml", "other/file"}
	pathFilters := []string{"^src.*$=code", "^kubernetes.*$=chart"}
	bt := utils.GetBuildTypeByPathFilters("default", changedPaths, pathFilters, false)
	if bt != "default" {
		t.Errorf("Expected: default, Actual: " + bt)
	}
}
func TestGetBuildTypeByPathFilters_MultipleAllowed_NotAllMatch(t *testing.T) {
	changedPaths := []string{"src/file1.go", "kubernetes/Chart.yaml", "other/file"}
	pathFilters := []string{"^src.*$=code", "^kubernetes.*$=chart"}
	bt := utils.GetBuildTypeByPathFilters("default", changedPaths, pathFilters, true)
	if bt != "default" {
		t.Errorf("Expected: default, Actual: " + bt)
	}
}
func TestIsMainlineOrReleaseRef_Mainline(t *testing.T) {
	isMainlineOrReleaseRef := utils.IsMainlineOrReleaseRef("master", "master", "^.*/rel-.*$")
	if isMainlineOrReleaseRef != true {
		t.Errorf("Expected: true, Actual: " + strconv.FormatBool(isMainlineOrReleaseRef))
	}
}
func TestIsMainlineOrReleaseRef_ReleaseRef(t *testing.T) {
	isMainlineOrReleaseRef := utils.IsMainlineOrReleaseRef("fda/rel-1", "master", "^.*/rel-.*$")
	if isMainlineOrReleaseRef != true {
		t.Errorf("Expected: true, Actual: " + strconv.FormatBool(isMainlineOrReleaseRef))
	}
}
func TestIsMainlineOrReleaseRef_Neither(t *testing.T) {
	isMainlineOrReleaseRef := utils.IsMainlineOrReleaseRef("develop", "master", "^.*/rel-.*$")
	if isMainlineOrReleaseRef != false {
		t.Errorf("Expected: false, Actual: " + strconv.FormatBool(isMainlineOrReleaseRef))
	}
}
func TestIsCommitError_NotError(t *testing.T) {
	isError := utils.IsCommitError("2f7444d674d79ea111483078e803cf3119c88e59", "E")
	if isError != false {
		t.Errorf("Expected: false, Actual: " + strconv.FormatBool(isError))
	}
}
func TestIsCommitError_Error(t *testing.T) {
	isError := utils.IsCommitError("E", "E")
	if isError != true {
		t.Errorf("Expected: false, Actual: " + strconv.FormatBool(isError))
	}
}
func TestCountLinesPerPathFilter_AllMatch(t *testing.T) {
	var changedPaths = []string{"src/file1.go", "kubernetes/Chart.yaml"}
	var pathFilters = []string{"^src.*$=code", "^kubernetes.*$=chart"}
	changedPathsPerFilter, changedPathsPerFilterCount := utils.CountLinesPerPathFilter(pathFilters, changedPaths)
	if changedPathsPerFilterCount != 2 {
		t.Errorf("Expected: 2, Actual: " + (string)(changedPathsPerFilterCount))
	}
	if changedPathsPerFilter["code"] != 1 {
		t.Errorf("Expected: 1, Actual: " + (string)(changedPathsPerFilter["code"]))
	}
	if changedPathsPerFilter["chart"] != 1 {
		t.Errorf("Expected: 2, Actual: " + (string)(changedPathsPerFilter["chart"]))
	}
}
func TestCountLinesPerPathFilter_NotAllMatch(t *testing.T) {
	var changedPaths = []string{"src/file1.go", "kubernetes/Chart.yaml", "other/file"}
	var pathFilters = []string{"^src.*$=code", "^kubernetes.*$=chart"}
	changedPathsPerFilter, changedPathsPerFilterCount := utils.CountLinesPerPathFilter(pathFilters, changedPaths)
	if changedPathsPerFilterCount != 2 {
		t.Errorf("Expected: 2, Actual: " + (string)(changedPathsPerFilterCount))
	}
	if changedPathsPerFilter["code"] != 1 {
		t.Errorf("Expected: 1, Actual: " + (string)(changedPathsPerFilter["code"]))
	}
	if changedPathsPerFilter["chart"] != 1 {
		t.Errorf("Expected: 2, Actual: " + (string)(changedPathsPerFilter["chart"]))
	}
}
