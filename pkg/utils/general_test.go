package utils

import (
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice contains string",
			args: args{[]string{"cat", "lion", "dog"}, "lion"},
			want: true,
		},
		{
			name: "slice doesn't contain string",
			args: args{[]string{"cat", "lion", "dog"}, "moose"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToString(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple map",
			args: args{map[string]string{"animal": "lion"}},
			want: "animal=lion",
		},
		{
			name: "complex map",
			args: args{map[string]string{
				"animal": "lion",
				"tool":   "hammer",
				"car":    "honda",
			},
			},
			want: "animal=lion, car=honda, tool=hammer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToString(tt.args.m); got != tt.want {
				t.Errorf("MapToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
