package test

import (
	"orca/pkg/utils"
	"testing"
)

func TestGetReleasesDelta(t *testing.T) {
	rel1 := utils.ReleaseSpec{ChartName: "chart1", ChartVersion: "1.0.0", ReleaseName: "dev-chart1"}
	rel2 := utils.ReleaseSpec{ChartName: "chart2", ChartVersion: "2.0.0", ReleaseName: "dev-chart2"}

	fromReleases := []utils.ReleaseSpec{rel1, rel2}
	toReleases := []utils.ReleaseSpec{rel1}

	releasesDelta := utils.GetReleasesDelta(fromReleases, toReleases)

	if len(releasesDelta) != 1 {
		t.Errorf("Expected: 1, Actual: " + (string)(len(releasesDelta)))
	}

	if !releasesDelta[0].Equals(rel2) {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestChartsYamlToStruct(t *testing.T) {
	file := "data/charts.yaml"
	rel0 := utils.ReleaseSpec{ChartName: "cassandra", ChartVersion: "0.4.0", ReleaseName: "test-cassandra"}
	rel1 := utils.ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	rel2 := utils.ReleaseSpec{ChartName: "kaa", ChartVersion: "0.1.7", ReleaseName: "test-kaa"}

	releases := utils.ChartsYamlToStruct(file, "test")

	if len(releases) != 3 {
		t.Errorf("Expected: 3, Actual: " + (string)(len(releases)))
	}
	if !releases[0].Equals(rel0) {
		t.Errorf("Expected: true, Actual: false")
	}
	if !releases[1].Equals(rel1) {
		t.Errorf("Expected: true, Actual: false")
	}
	if !releases[2].Equals(rel2) {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestCheckCircularDependencies_False(t *testing.T) {
	file := "data/charts.yaml"
	circular := utils.CheckCircularDependencies(utils.ChartsYamlToStruct(file, "test"))
	if circular {
		t.Errorf("Expected: false, Actual: true")
	}
}
func TestCheckCircularDependencies_True(t *testing.T) {
	file := "data/circular.yaml"
	circular := utils.CheckCircularDependencies(utils.ChartsYamlToStruct(file, "test"))
	if !circular {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestOverrideReleases(t *testing.T) {
	rel0 := utils.ReleaseSpec{ChartName: "cassandra", ChartVersion: "0.4.0", ReleaseName: "test-cassandra"}
	rel1 := utils.ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	rel2 := utils.ReleaseSpec{ChartName: "kaa", ChartVersion: "0.1.7", ReleaseName: "test-kaa"}
	rel2override := utils.ReleaseSpec{ChartName: "kaa", ChartVersion: "7.1.0", ReleaseName: "test-kaa"}

	releases := []utils.ReleaseSpec{rel0, rel1, rel2}

	overrideReleases := utils.OverrideReleases(releases, []string{"kaa=7.1.0"})

	if !overrideReleases[2].Equals(rel2override) {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestRemoveChartFromDependencies(t *testing.T) {
	file := "data/charts.yaml"
	releases := utils.ChartsYamlToStruct(file, "test")
	releases = utils.RemoveChartFromDependencies(releases, "mariadb")

	if len(releases[2].Dependencies) != 1 {
		t.Errorf("Expected: 1, Actual: " + (string)(len(releases)))
	}
	if releases[2].Dependencies[0] != "cassandra" {
		t.Errorf("Expected: cassandra, Actual: " + releases[2].Dependencies[0])
	}
}
func TestGetChartIndex(t *testing.T) {
	file := "data/charts.yaml"
	releases := utils.ChartsYamlToStruct(file, "test")
	index := utils.GetChartIndex(releases, "cassandra")

	if index != 0 {
		t.Errorf("Expected: 0, Actual: " + (string)(index))
	}
}
func TestRemoveChartFromCharts(t *testing.T) {
	rel1 := utils.ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	rel0 := utils.ReleaseSpec{ChartName: "kaa", ChartVersion: "0.1.7", ReleaseName: "test-kaa"}
	file := "data/charts.yaml"
	releases := utils.ChartsYamlToStruct(file, "test")
	index := utils.GetChartIndex(releases, "cassandra")
	releases = utils.RemoveChartFromCharts(releases, index)

	if len(releases) != 2 {
		t.Errorf("Expected: 2, Actual: " + (string)(len(releases)))
	}
	if !releases[0].Equals(rel0) {
		t.Errorf("Expected: true, Actual: false")
	}
	if !releases[1].Equals(rel1) {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestUpdateChartVersion(t *testing.T) {
	newVersion := utils.UpdateChartVersion("data/", "1234")

	if newVersion != "0.1.1-1234" {
		t.Errorf("Expected: 0.1.1-1234, Actual: " + newVersion)
	}
}
func TestEquals_True(t *testing.T) {
	rel := utils.ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	equals := rel.Equals(rel)
	if !equals {
		t.Errorf("Expected: true, Actual: false")
	}
}
func TestEquals_False(t *testing.T) {
	rel := utils.ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	relDiff := utils.ReleaseSpec{ChartName: "mariadba", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	equals := rel.Equals(relDiff)
	if equals {
		t.Errorf("Expected: false, Actual: true")
	}
}
