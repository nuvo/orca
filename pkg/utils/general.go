package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	exec "os/exec"
	"strconv"
	"strings"
	"time"
)

// PrintExec takes a command and executes it, with or without printing
func PrintExec(cmd []string, print bool) error {
	if print {
		fmt.Println(cmd)
	}
	output, err := Exec(cmd)
	if err != nil {
		return err
	}
	if print {
		fmt.Print(output)
	}
	return nil
}

// Exec takes a command as a string and executes it
func Exec(cmd []string) (string, error) {
	binary := cmd[0]
	_, err := exec.LookPath(binary)
	if err != nil {
		return "", err
	}

	output, err := exec.Command(binary, cmd[1:]...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(output))
	}
	return string(output), nil
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

// SplitInTwo splits a string to two parts by a delimeter
func SplitInTwo(s, sep string) (string, string) {
	if !strings.Contains(s, sep) {
		log.Fatal(s, "does not contain", sep)
	}
	split := strings.Split(s, sep)
	return split[0], split[1]
}

// MapToString returns a string representation of a map
func MapToString(m map[string]string) string {
	var output string
	for k, v := range m {
		output += fmt.Sprintf("%s=%s, ", k, v)
	}
	return strings.TrimRight(output, ", ")
}
