package utils

import "testing"

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
