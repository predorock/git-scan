package utils

import (
	"io/ioutil"
	"os"
	"predorock/gitscan/logstreamer"
	"testing"
)

const tmpPrefix = "git-scan-utils"

var log = logstreamer.GetInstance().Log

func provideTempEnv(prefix string, runner func(string)) {
	dir, err := ioutil.TempDir("/tmp", prefix)

	if err != nil {
		log.Fatal(err)
	}

	runner(dir)

	defer os.RemoveAll(dir)
}

func provideFakeGitEnv(prefix string, runner func(string)) {
	provideTempEnv(prefix, func(tmp string) {
		var path = tmp + "/" + ".git"
		if err := os.Mkdir(path, os.ModeTemporary); err != nil {
			log.Fatal(err)
		}
		runner(tmp)
	})
}

func unorderedEqual(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value] = true
	}
	for _, value := range second {
		if !exists[value] {
			return false
		}
	}
	return true
}

func TestIsNotGitFolder(t *testing.T) {

	provideTempEnv(tmpPrefix, func(folder string) {
		log.Printf("folder: %s\n", folder)
		if IsGitFolder(folder) {
			t.Errorf("%sis not a git folder\n", folder)
		}
	})
}

func TestIsGitFolder(t *testing.T) {

	provideFakeGitEnv(tmpPrefix, func(folder string) {
		log.Printf("folder: %s\n", folder)
		if !IsGitFolder(folder) {
			t.Errorf("%sis not a git folder\n", folder)
		}
	})
}

func TestGetDirs(t *testing.T) {
	provideTempEnv(tmpPrefix, func(s string) {

		expected := []string{
			"hello",
			"foo",
			"bar",
			"something",
		}

		for _, ff := range expected {
			var path = s + "/" + ff

			if err := os.Mkdir(path, os.ModeTemporary); err != nil {
				log.Fatal(err)
			}
		}

		var result, err = GetDirs(s)

		if err != nil {
			t.Errorf("getDirs has raised an error: %s\n", err.Error())
		}

		if len(result) != len(expected) {
			t.Errorf("getDirs got a different number of directories")
		}

		if !unorderedEqual(result, expected) {
			t.Errorf("getDirs returns different directories names\n- Expected: %v\n- Result: %v\n", expected, result)
		}

	})
}
