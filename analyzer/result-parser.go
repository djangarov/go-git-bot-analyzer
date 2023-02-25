package analyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Report struct {
	Files []File
}

type Violation struct {
	BeginLine       int    `json:"beginline"`
	BeginColumn     int    `json:"begincolumn"`
	EndLine         int    `json:"endline"`
	Description     string `json:"description"`
	Rule            string `json:"rule"`
	Ruleset         string `json:"ruleset"`
	Priority        int    `json:"priority"`
	ExternalInfoUrl string `json:"externalInfoUrl"`
}

type File struct {
	Filename   string      `json:"filename"`
	Violations []Violation `json:"violations"`
}

func Parse(pathToReport string) (Report, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(pathToReport)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return Report{}, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var report Report

	json.Unmarshal(byteValue, &report)
	defer os.Remove(pathToReport)
	return report, err
}
