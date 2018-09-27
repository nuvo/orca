package utils

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type ReleaseSpec struct {
	ReleaseName  string
	ChartName    string
	ChartVersion string
	Dependencies []string
}

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

func ChartsYamlToStruct(file, env string) []ReleaseSpec {
	var charts []ReleaseSpec

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

		var c ReleaseSpec

		c.ReleaseName = env + "-" + chart["name"].(string)
		c.ChartName = chart["name"].(string)
		c.ChartVersion = chart["version"].(string)

		if chart["depends_on"] != nil {
			for _, dep := range chart["depends_on"].([]interface{}) {
				depStr := dep.(string)
				c.Dependencies = append(c.Dependencies, depStr)
			}
		}
		charts = append(charts, c)
	}

	return charts
}

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

func OverrideReleases(releases []ReleaseSpec, overrides []string) []ReleaseSpec {
	var outReleases []ReleaseSpec

	for _, r := range releases {
		for _, override := range overrides {
			oChartName, oChartVersion := SplitInTwo(override, "=")

			if r.ChartName == oChartName && r.ChartVersion != oChartVersion {
				r.ChartName = oChartName
				r.ChartVersion = oChartVersion
			}
		}
		outReleases = append(outReleases, r)
	}

	return outReleases
}

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

func GetChartIndex(charts []ReleaseSpec, name string) int {
	index := -1
	for i, elem := range charts {
		if elem.ChartName == name {
			index = i
		}
	}
	return index
}

func RemoveChartFromCharts(charts []ReleaseSpec, index int) []ReleaseSpec {
	charts[index] = charts[len(charts)-1]
	return charts[:len(charts)-1]
}

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

	newVersion := fmt.Sprintf("%s-%s", v["version"], append)
	v["version"] = newVersion

	data, err = yaml.Marshal(v)
	ioutil.WriteFile(filePath, data, 0755)

	return newVersion
}

func (r ReleaseSpec) Print() {
	fmt.Println("release name: " + r.ReleaseName)
	fmt.Println("chart name: " + r.ChartName)
	fmt.Println("chart version: " + r.ChartVersion)
	for _, dep := range r.Dependencies {
		fmt.Println("depends_on: " + dep)
	}
}

func (a ReleaseSpec) Equals(b ReleaseSpec) bool {
	equals := false
	if a.ReleaseName == b.ReleaseName &&
		a.ChartName == b.ChartName &&
		a.ChartVersion == b.ChartVersion {
		equals = true
	}

	return equals
}

func PrintReleasesYaml(releases []ReleaseSpec) {
	if len(releases) != 0 {
		fmt.Println("charts:")
	}
	for _, r := range releases {
		fmt.Println("- name:", r.ChartName)
		fmt.Println("  vesrion:", r.ChartVersion)
	}
}

func PrintReleasesMarkdown(releases []ReleaseSpec) {
	if len(releases) != 0 {
		fmt.Println("| Name | Version |")
		fmt.Println("|------|---------|")
	}
	for _, r := range releases {
		fmt.Println(fmt.Sprintf("| %s | %s |", r.ChartName, r.ChartVersion))
	}
}
