package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/gosuri/uitable"
	yaml "gopkg.in/yaml.v2"
)

// ReleaseSpec holds data relevant to deploying a release
type ReleaseSpec struct {
	ReleaseName  string
	ChartName    string
	ChartVersion string
	Dependencies []string
}

// GetReleasesDelta returns the delta between two slices of ReleaseSpec
func GetReleasesDelta(fromReleases, toReleases []ReleaseSpec) []ReleaseSpec {
	var releasesDelta []ReleaseSpec
	var releasesExists []ReleaseSpec

	for _, fromRelease := range fromReleases {
		exists := false
		for _, toRelease := range toReleases {
			if fromRelease.Equals(toRelease) {
				exists = true
				releasesExists = append(releasesExists, toRelease)
				break
			}
		}
		if !exists {
			releasesDelta = append(releasesDelta, fromRelease)
		}
	}

	for _, releaseExists := range releasesExists {
		releasesDelta = RemoveChartFromDependencies(releasesDelta, releaseExists.ChartName)
	}

	return releasesDelta
}

// InitReleasesFromChartsFiles initializes a slice of ReleaseSpec from a set of yaml formatted charts files, overriding versions from left to right
func InitReleasesFromChartsFiles(files []string, env string) []ReleaseSpec {
	if len(files) == 0 {
		return nil
	}
	desiredReleases := InitReleasesFromChartsFile(files[0], env)
	for i := 1; i < len(files); i++ {
		releasesToMerge := InitReleasesFromChartsFile(files[i], env)
		desiredReleases = MergeReleases(desiredReleases, releasesToMerge, env)
	}

	return desiredReleases
}

// InitReleasesFromChartsFile initializes a slice of ReleaseSpec from a yaml formatted charts file
func InitReleasesFromChartsFile(file, env string) []ReleaseSpec {
	var releases []ReleaseSpec

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	var v map[string][]map[string]interface{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalln(err)
	}

	for _, chart := range v["charts"] {

		c := initReleaseSpec(env, chart["name"].(string), chart["version"].(string))

		if chart["depends_on"] != nil {
			for _, dep := range chart["depends_on"].([]interface{}) {
				depStr := dep.(string)
				c.Dependencies = append(c.Dependencies, depStr)
			}
		}
		releases = append(releases, c)
	}

	return releases
}

// InitReleases initializes a slice of ReleaseSpec from a string slice
func InitReleases(env string, releases []string) []ReleaseSpec {
	var outReleases []ReleaseSpec

	for _, release := range releases {
		chartName, chartVersion := SplitInTwo(release, "=")
		r := initReleaseSpec(env, chartName, chartVersion)
		outReleases = append(outReleases, r)
	}

	return outReleases
}

func initReleaseSpec(env, name, version string) ReleaseSpec {
	return ReleaseSpec{
		ReleaseName:  env + "-" + name,
		ChartName:    name,
		ChartVersion: version,
	}
}

// CheckCircularDependencies verifies that there are no circular dependencies between ReleaseSpecs
func CheckCircularDependencies(releases []ReleaseSpec) bool {

	startLen := len(releases)
	endLen := -1

	// while a release was processed
	for startLen != endLen && endLen != 0 {
		startLen = len(releases)
		var indexesToRemove []int
		// find releases to process
		for i := 0; i < len(releases); i++ {
			if len(releases[i].Dependencies) != 0 {
				continue
			}
			indexesToRemove = append(indexesToRemove, i)
		}
		// "process" the releases
		for i := len(indexesToRemove) - 1; i >= 0; i-- {
			releases = RemoveChartFromDependencies(releases, releases[indexesToRemove[i]].ChartName)
			releases = RemoveChartFromCharts(releases, indexesToRemove[i])
		}
		endLen = len(releases)
	}

	// if there are any releases left to process - there is a circular dependency
	if endLen != 0 {
		return true
	}
	return false
}

// OverrideReleases overrides versions of specified overrides
func OverrideReleases(releases []ReleaseSpec, overrides []string, env string) []ReleaseSpec {
	if len(overrides) == 0 {
		return releases
	}

	var outReleases []ReleaseSpec
	var overrideFound = make([]bool, len(overrides))

	for _, r := range releases {
		for i := 0; i < len(overrides); i++ {
			oChartName, oChartVersion := SplitInTwo(overrides[i], "=")

			if r.ChartName == oChartName && r.ChartVersion != oChartVersion {
				overrideFound[i] = true
				r.ChartName = oChartName
				r.ChartVersion = oChartVersion
			}
		}
		outReleases = append(outReleases, r)
	}

	for i := 0; i < len(overrides); i++ {
		if overrideFound[i] {
			continue
		}
		oChartName, oChartVersion := SplitInTwo(overrides[i], "=")
		r := ReleaseSpec{
			ReleaseName:  env + "-" + oChartName,
			ChartName:    oChartName,
			ChartVersion: oChartVersion,
		}
		outReleases = append(outReleases, r)
	}

	return outReleases
}

// MergeReleases will merge releases with releasesToMerge (currently will only override the version field)
func MergeReleases(releases, releasesToMerge []ReleaseSpec, env string) []ReleaseSpec {
	var overrides []string
	for _, rtm := range releasesToMerge {
		override := rtm.ChartName + "=" + rtm.ChartVersion
		overrides = append(overrides, override)
	}

	return OverrideReleases(releases, overrides, env)
}

// RemoveChartFromDependencies removes a release from other releases ReleaseSpec depends_on field
func RemoveChartFromDependencies(charts []ReleaseSpec, name string) []ReleaseSpec {

	var outCharts []ReleaseSpec

	for _, dependant := range charts {
		if Contains(dependant.Dependencies, name) {

			index := -1
			for i, elem := range dependant.Dependencies {
				if elem == name {
					index = i
				}
			}
			if index == -1 {
				log.Fatal("Could not find element in dependencies")
			}

			dependant.Dependencies[index] = dependant.Dependencies[len(dependant.Dependencies)-1]
			dependant.Dependencies[len(dependant.Dependencies)-1] = ""
			dependant.Dependencies = dependant.Dependencies[:len(dependant.Dependencies)-1]
		}
		outCharts = append(outCharts, dependant)
	}

	return outCharts
}

// GetChartIndex returns the index of a desired release by its name
func GetChartIndex(charts []ReleaseSpec, name string) int {
	index := -1
	for i, elem := range charts {
		if elem.ChartName == name {
			index = i
		}
	}
	return index
}

// RemoveChartFromCharts removes a ReleaseSpec from a slice of ReleaseSpec
func RemoveChartFromCharts(charts []ReleaseSpec, index int) []ReleaseSpec {
	charts[index] = charts[len(charts)-1]
	return charts[:len(charts)-1]
}

// UpdateChartVersion updates a chart version with desired append value
func UpdateChartVersion(path, append string) string {
	filePath := path + "Chart.yaml"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	var v map[string]interface{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalln(err)
	}

	version := v["version"].(string)
	if append == "" {
		return version
	}
	newVersion := fmt.Sprintf("%s-%s", version, append)
	v["version"] = newVersion

	data, err = yaml.Marshal(v)
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile(filePath, data, 0755)

	return newVersion
}

// ResetChartVersion resets a chart version to a desired value
func ResetChartVersion(path, version string) {
	filePath := path + "Chart.yaml"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	var v map[string]interface{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalln(err)
	}

	v["version"] = version

	data, err = yaml.Marshal(v)
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile(filePath, data, 0755)
}

// Print prints a ReleaseSpec
func (r ReleaseSpec) Print() {
	fmt.Println("release name: " + r.ReleaseName)
	fmt.Println("chart name: " + r.ChartName)
	fmt.Println("chart version: " + r.ChartVersion)
	for _, dep := range r.Dependencies {
		fmt.Println("depends_on: " + dep)
	}
}

// Equals compares two ReleaseSpecs
func (r ReleaseSpec) Equals(b ReleaseSpec) bool {
	equals := false
	if r.ReleaseName == b.ReleaseName &&
		r.ChartName == b.ChartName &&
		r.ChartVersion == b.ChartVersion {
		equals = true
	}

	return equals
}

// PrintReleasesYaml prints releases in yaml format
func PrintReleasesYaml(releases []ReleaseSpec) {
	if len(releases) == 0 {
		return
	}
	fmt.Println("charts:")
	for _, r := range releases {
		fmt.Println("- name:", r.ChartName)
		fmt.Println("  version:", r.ChartVersion)
	}
}

// PrintReleasesMarkdown prints releases in markdown format
func PrintReleasesMarkdown(releases []ReleaseSpec) {
	if len(releases) == 0 {
		return
	}
	fmt.Println("| Name | Version |")
	fmt.Println("|------|---------|")
	for _, r := range releases {
		fmt.Println(fmt.Sprintf("| %s | %s |", r.ChartName, r.ChartVersion))
	}
}

// PrintReleasesTable prints releases in table format
func PrintReleasesTable(releases []ReleaseSpec) {
	if len(releases) == 0 {
		return
	}
	tbl := uitable.New()
	tbl.MaxColWidth = 60
	tbl.AddRow("NAME", "VERSION")

	for _, r := range releases {
		tbl.AddRow(r.ChartName, r.ChartVersion)
	}
	fmt.Println(tbl.String())
}

// DiffOptions are options passed to PrintDiffTable
type DiffOptions struct {
	KubeContextLeft   string
	EnvNameLeft       string
	KubeContextRight  string
	EnvNameRight      string
	ReleasesSpecLeft  []ReleaseSpec
	ReleasesSpecRight []ReleaseSpec
	Output            string
}

type diff struct {
	chartName    string
	versionLeft  string
	versionRight string
}

// PrintDiff prints a table of differences between two environments
func PrintDiff(o DiffOptions) {
	if len(o.ReleasesSpecLeft) == 0 && len(o.ReleasesSpecRight) == 0 {
		return
	}
	diffs := getDiffs(o.ReleasesSpecLeft, o.ReleasesSpecRight)
	if len(diffs) == 0 {
		return
	}

	switch o.Output {
	case "yaml":
		printDiffYaml(diffs)
	case "table":
		printDiffTable(o, diffs)
	case "":
		printDiffYaml(diffs)
	}

}

func printDiffYaml(diffs []diff) {
	fmt.Println("charts:")
	for _, d := range diffs {
		fmt.Println("- name:", d.chartName)
		fmt.Println("  versionLeft:", d.versionLeft)
		fmt.Println("  versionRight:", d.versionRight)
	}
}

func printDiffTable(o DiffOptions, diffs []diff) {
	tbl := uitable.New()
	tbl.MaxColWidth = 60
	leftColHeader := initHeader(o.KubeContextLeft, o.EnvNameLeft)
	rightColHeader := initHeader(o.KubeContextRight, o.EnvNameRight)
	tbl.AddRow("chart", leftColHeader, rightColHeader)

	for _, d := range diffs {
		tbl.AddRow(d.chartName, d.versionLeft, d.versionRight)
	}
	fmt.Println(tbl.String())
}

func initHeader(kubeContext, envName string) string {
	if kubeContext != "" {
		kubeContext += "/"
	}
	return fmt.Sprintf("%s%s", kubeContext, envName)
}

func getDiffs(releasesLeft, releasesRight []ReleaseSpec) []diff {
	leftAndRight := mergeReleasesToCompare(releasesLeft, releasesRight)
	diffs := removeEquals(leftAndRight)

	return diffs
}

func mergeReleasesToCompare(releasesLeft, releasesRight []ReleaseSpec) []diff {
	// Initialize all left elements
	var left []diff
	for _, r := range releasesLeft {
		d := diff{
			chartName:   r.ChartName,
			versionLeft: r.ChartVersion,
		}
		left = append(left, d)
	}
	// Add right elements to existing elements from left
	var leftAndRight []diff
	for _, r := range releasesRight {
		found := false
		for i := 0; i < len(left); i++ {
			l := left[i]
			if l.chartName == r.ChartName {
				found = true
				l.versionRight = r.ChartVersion
				leftAndRight = append(leftAndRight, l)
				left = append(left[:i], left[i+1:]...)
				break
			}
		}
		// Add right elements which do not exist in left
		if !found {
			d := diff{
				chartName:    r.ChartName,
				versionRight: r.ChartVersion,
			}
			leftAndRight = append(leftAndRight, d)
		}
	}
	// Add left elements which do not exist in right
	for _, r := range left {
		leftAndRight = append(leftAndRight, r)
	}

	return leftAndRight
}

func removeEquals(leftAndRight []diff) []diff {
	var diffs []diff
	for _, lar := range leftAndRight {
		if lar.versionLeft == lar.versionRight {
			continue
		}
		diffs = append(diffs, lar)
	}

	sort.Slice(diffs[:], func(i, j int) bool {
		return strings.Compare(diffs[i].chartName, diffs[j].chartName) <= 0
	})

	return diffs
}
