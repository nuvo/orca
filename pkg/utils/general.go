package utils

import (
	"math/rand"
	"os"
	exec "os/exec"
	"strconv"
	"strings"
	"time"
)

// Exec takes a command as a string and executes it
func Exec(cmd string) string {
	args := strings.Split(cmd, " ")
	binary := args[0]
	_, lookErr := exec.LookPath(binary)
	if lookErr != nil {
		panic(lookErr)
	}

	output, err := exec.Command(binary, args[1:]...).CombinedOutput()
	if err != nil {
		panic(string(output))
	}
	return string(output)
}

// MkRandomDir creates a new directory with a random name made of numbers
func MkRandomDir() string {
	r := strconv.Itoa((rand.New(rand.NewSource(time.Now().UnixNano()))).Int())
	os.Mkdir(r, 0755)

	return r
}

// AddIfNotContained adds a string to a slice if it is not contained in it and not empty
func AddIfNotContained(s []string, e string) (sout []string) {
	if (!Contains(s, e)) && (e != "") {
		s = append(s, e)
	}

	return s
}

// Contains checks if a slice contains a given value
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Truncate(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
