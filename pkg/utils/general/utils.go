package genutils

import (
	"fmt"
	"math/rand"
	"os"
	exec "os/exec"
	"strconv"
	"strings"
	"time"
)

// Exec takes a command as a string and executes it
func Exec(cmd string) {
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
	fmt.Print(string(output))
}

// MkRandomDir creates a new directory with a random name made of numbers
func MkRandomDir() string {
	r := strconv.Itoa((rand.New(rand.NewSource(time.Now().UnixNano()))).Int())
	os.Mkdir(r, 0755)

	return r
}

// AddIfNotContained adds a string to a slice if it is not contained in it and not empty
func AddIfNotContained(s []string, e string) (sout []string) {
	if (!contains(s, e)) && (e != "") {
		s = append(s, e)
	}

	return s
}

// contains checks if a slice contains a given value
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
