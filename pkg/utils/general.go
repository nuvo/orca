package utils

import (
	"log"
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
	_, err := exec.LookPath(binary)
	if err != nil {
		log.Fatal(err)
	}

	output, err := exec.Command(binary, args[1:]...).CombinedOutput()
	if err != nil {
		log.Fatal(string(output))
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

// SplitInTwo splits a string to two parts by a delimeter
func SplitInTwo(s, sep string) (string, string) {
	if !strings.Contains(s, sep) {
		log.Fatal(s, "does not contain", sep)
	}
	split := strings.Split(s, sep)
	return split[0], split[1]
}

// GetIntEnvVar returns 0 if the variable is empty or not int, else the value
func GetIntEnvVar(name string, defVal int) int {
	val := os.Getenv(name)
	if val == "" {
		return defVal
	}
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return defVal
	}
	return iVal
}

// GetStringEnvVar returns the default value if the variable is empty, else the value
func GetStringEnvVar(name, defVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defVal
	}
	return val
}

// GetBoolEnvVar returns the default value if the variable is empty or not true or false, else the value
func GetBoolEnvVar(name string, defVal bool) bool {
	val := os.Getenv(name)
	if strings.ToLower(val) == "true" {
		return true
	}
	if strings.ToLower(val) == "false" {
		return false
	}
	return defVal
}
