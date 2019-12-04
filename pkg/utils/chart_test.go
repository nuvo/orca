package utils

import (
	"testing"
)

func TestCheckCircularDependencies(t *testing.T) {
	type args struct {
		releases []ReleaseSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no circular dependencies",
			args: args{InitReleasesFromChartsFile("./testdata/charts.yaml", "test")},
			want: false,
		},
		{
			name: "circular dependencies",
			args: args{InitReleasesFromChartsFile("./testdata/circular.yaml", "test")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckCircularDependencies(tt.args.releases); got != tt.want {
				t.Errorf("CheckCircularDependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChartIndex(t *testing.T) {
	type args struct {
		charts []ReleaseSpec
		name   string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "charts file has this chart",
			args: args{InitReleasesFromChartsFile("./testdata/charts.yaml", "test"), "kaa"},
			want: 2,
		},
		{
			name: "charts file doesn't have this chart",
			args: args{InitReleasesFromChartsFile("./testdata/charts.yaml", "test"), "rabbitmq"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChartIndex(tt.args.charts, tt.args.name); got != tt.want {
				t.Errorf("GetChartIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReleasesDelta(t *testing.T) {
	rel1 := ReleaseSpec{ChartName: "chart1", ChartVersion: "1.0.0", ReleaseName: "dev-chart1"}
	rel2 := ReleaseSpec{ChartName: "chart2", ChartVersion: "2.0.0", ReleaseName: "dev-chart2"}

	fromReleases := []ReleaseSpec{rel1, rel2}
	toReleases := []ReleaseSpec{rel1}

	releasesDelta := GetReleasesDelta(fromReleases, toReleases)

	if len(releasesDelta) != 1 {
		t.Errorf("Expected: 1, Actual: " + (string)(len(releasesDelta)))
	}

	if !releasesDelta[0].Equals(rel2) {
		t.Errorf("Expected: true, Actual: false")
	}
}

func TestChartsYamlToReleases(t *testing.T) {
	rel0 := ReleaseSpec{ChartName: "cassandra", ChartVersion: "0.4.0", ReleaseName: "test-cassandra"}
	rel1 := ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}
	rel2 := ReleaseSpec{ChartName: "kaa", ChartVersion: "0.1.7", ReleaseName: "test-kaa"}

	releases := InitReleasesFromChartsFile("testdata/charts.yaml", "test")

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

func TestReleaseSpec_Equals(t *testing.T) {
	type fields struct {
		ReleaseName  string
		ChartName    string
		ChartVersion string
		Dependencies []string
	}
	type args struct {
		b ReleaseSpec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "equals should be true",
			fields: fields{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"},
			args:   args{b: ReleaseSpec{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}},
			want:   true,
		},
		{
			name:   "equals should be false",
			fields: fields{ChartName: "mariadb", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"},
			args:   args{b: ReleaseSpec{ChartName: "cassandra", ChartVersion: "0.5.4", ReleaseName: "test-mariadb"}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ReleaseSpec{
				ReleaseName:  tt.fields.ReleaseName,
				ChartName:    tt.fields.ChartName,
				ChartVersion: tt.fields.ChartVersion,
				Dependencies: tt.fields.Dependencies,
			}
			if got := r.Equals(tt.args.b); got != tt.want {
				t.Errorf("ReleaseSpec.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
