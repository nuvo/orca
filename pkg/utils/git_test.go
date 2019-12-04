package utils

import (
	"testing"
)

func TestIsCommitError(t *testing.T) {
	type args struct {
		commit               string
		commitErrorIndicator string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "commit is not an error",
			args: args{commit: "2f7444d674d79ea111483078e803cf3119c88e59", commitErrorIndicator: "E"},
			want: false,
		},
		{
			name: "commit is an error",
			args: args{commit: "E", commitErrorIndicator: "E"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCommitError(tt.args.commit, tt.args.commitErrorIndicator); got != tt.want {
				t.Errorf("IsCommitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMainlineOrReleaseRef(t *testing.T) {
	type args struct {
		currentRef string
		mainRef    string
		releaseRef string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "this is the mainline",
			args: args{"master", "master", "^./rel-.*$"},
			want: true,
		},
		{
			name: "this is a release branch",
			args: args{"fda/rel-1", "master", "^.*/rel-.*$"},
			want: true,
		},
		{
			name: "this is neither a release nor mainline",
			args: args{"develop", "master", "^.*/rel-.*$"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMainlineOrReleaseRef(tt.args.currentRef, tt.args.mainRef, tt.args.releaseRef); got != tt.want {
				t.Errorf("IsMainlineOrReleaseRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBuildTypeByPathFilters(t *testing.T) {
	type args struct {
		defaultType        string
		changedPaths       []string
		pathFilter         []string
		allowMultipleTypes bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "multiple not allowed, no changed paths",
			args: args{"default", []string{}, []string{"^src.*$=code", "^kubernetes.*$=chart"}, false},
			want: "default",
		},
		{
			name: "multiple allows, all paths match",
			args: args{"default", []string{"src/file1.go", "kubernetes/Chart.yaml"}, []string{"^src.*$=code", "^kubernetes.*$=chart"}, true},
			want: "code;chart",
		},
		{
			name: "multiple not allowed, not all paths match",
			args: args{"default", []string{"src/file1.go", "kubernetes/Chart.yaml", "other/file"}, []string{"^src.*$=code", "^kubernetes.*$=chart"}, false},
			want: "default",
		},
		{
			name: "multiple allowed, not all paths match",
			args: args{"default", []string{"src/file1.go", "kubernetes/Chart.yaml", "other/file"}, []string{"^src.*$=code", "^kubernetes.*$=chart"}, true},
			want: "default",
		},
		{
			name: "multiple not allowed, all paths match",
			args: args{"default", []string{"src/file1.go", "kubernetes/Chart.yaml"}, []string{"^src.*$=code", "^kubernetes.*$=chart"}, false},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBuildTypeByPathFilters(tt.args.defaultType, tt.args.changedPaths, tt.args.pathFilter, tt.args.allowMultipleTypes); got != tt.want {
				t.Errorf("GetBuildTypeByPathFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
