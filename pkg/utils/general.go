package utils

import (
	"fmt"
	"log"
	exec "os/exec"
	"strings"
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
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	var output string
	for _, k := range keys {
		output += fmt.Sprintf("%s=%s, ", k, m[k])
	}
	return strings.TrimRight(output, ", ")
}
