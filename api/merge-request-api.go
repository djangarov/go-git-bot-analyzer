package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/djangarov/go-git-bot-analyzer/analyzer"
	"github.com/djangarov/go-git-bot-analyzer/differ"
	"github.com/djangarov/go-git-bot-analyzer/utils"
)

const JAVA_EXTENSION = ".java"

type MergeRequestEvent struct {
	Project          Project          `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
}

type Project struct {
	Id int `json:"id"`
}

type ObjectAttributes struct {
	State string `json:"state"`
	Url   string `json:"url"`
}

/*
*
This method listens for merge request event
*/
func (api *Api) HandleMergeRequestEvent(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		fmt.Println(http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	api.mergeRequestEvent(w, req)

	body := "ScanMergeRequestHandler"
	fmt.Println(body)
	w.Write([]byte(body))
}

func (api *Api) mergeRequestEvent(w http.ResponseWriter, req *http.Request) error {
	var mergeRequest MergeRequestEvent

	err := json.NewDecoder(req.Body).Decode(&mergeRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// if the merge request is close, we can ignore the rest of the code
	if mergeRequest.ObjectAttributes.State != "opened" {
		return nil
	}

	// retrieve the mr id
	split := strings.Split(mergeRequest.ObjectAttributes.Url, utils.URL_SLASH)
	mergeRequestId := split[len(split)-1]

	// build the url for detailed mr info.
	// we use the following gitlab endpoint
	// https://:host/api/v4/projects/:project-id/merge_requests/:merge_id/changes
	mrDetailedInfoURL := mergeRequestURIBuilder(api.GitlabHost, strconv.Itoa(mergeRequest.Project.Id), mergeRequestId) + utils.MERGE_REQUEST_CHANGES

	info, err := api.getMRDetailedInfo(mrDetailedInfoURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// for every file change, we need to validate the file and post a comment if there is a violation
	processFileChanges(info, strconv.Itoa(mergeRequest.Project.Id), mergeRequestId, api)
	return nil

}

func processFileChanges(info MRInfo, projectId string, mergeRequestId string, api *Api) {
	for _, change := range info.Changes {
		// temporary ignore the files that contain 'Test' keyword
		// later the test files will be recognized based on the folder structure
		if strings.Contains(change.NewPath, "Test") {
			continue
		}
		// split the string in order to get the file name
		splittedPath := strings.Split(change.NewPath, utils.URL_SLASH)
		fileName := splittedPath[len(splittedPath)-1]

		// if the current file is not JAVA file, skip this entry
		if !strings.Contains(fileName, JAVA_EXTENSION) {
			continue
		}

		// copy the file into the local system and get its path
		fileSystemPath, fileSystemDirectory, err := api.GetRepositoryFile(change.NewPath, fileName, projectId, info.Sha)

		if err != nil {
			log.Fatal(err.Error())
		}

		report, err := analyzer.Analyze(fileSystemPath)

		if err != nil {
			log.Fatal(err.Error())
		}

		// if there are no files, skip this entry
		if len(report.Files) == 0 {
			continue
		}

		// create temp directories to store the files in order to analyze them
		// as of now, the temp directories are not deleted
		tempFilePath := "temp-" + strconv.FormatInt(time.Now().Unix(), 10)
		directory := filepath.Join(".", "temp", tempFilePath)
		path := filepath.Join(directory, fileName)

		out, err := Create(path)
		_, err = out.WriteString(change.Diff)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = out.Sync()
		if err != nil {
			log.Fatal(err.Error())
		}

		// right now, for the testing purposes we scan one file at a time
		// scanning multiple files would be much much faster
		violations := report.Files[0].Violations

		// if there are no violation skip this entry
		if len(violations) == 0 {
			purge(fileSystemDirectory, directory)
			continue
		}

		if change.NewFile {
			for _, violation := range violations {
				api.CommentOn(mergeRequestId, projectId, violation, change.NewPath, info.DiffRefs)
			}
			purge(fileSystemDirectory, directory)
			continue
		}

		violationMap := make(map[int]analyzer.Violation)
		lines := differ.ExtractModifiedLines(change.Diff)
		var violationLinesBetween []int

		for _, violation := range violations {
			fmt.Printf("Violation!!! Rule: %s  File: %s \n", violation.Rule, fileName)
			if differ.IsBetween(violation.BeginLine, lines) {
				violationLinesBetween = append(violationLinesBetween, violation.BeginLine)
				violationMap[violation.BeginLine] = violation
			}
		}

		addedLines := differ.GetAddedLines(violationLinesBetween, path)

		// iterate through the added lines and comment on the violation
		for line := range addedLines {
			api.CommentOn(mergeRequestId, projectId, violationMap[line], change.NewPath, info.DiffRefs)
		}

		purge(fileSystemDirectory, directory)
	}
}

func purge(fileSysDirectory, directory string) {
	err := os.RemoveAll(fileSysDirectory)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.RemoveAll(directory)
	if err != nil {
		log.Fatal(err.Error())
	}
}
