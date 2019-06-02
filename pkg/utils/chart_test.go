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
