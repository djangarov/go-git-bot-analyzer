package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/djangarov/go-git-bot-analyzer/utils"
)

type MRInfo struct {
	DiffRefs DiffRefs  `json:"diff_refs"`
	Changes  []Changes `json:"changes"`
	Sha      string    `json:"sha"`
}

type DiffRefs struct {
	BaseSha  string `json:"base_sha"`
	HeadSha  string `json:"head_sha"`
	StartSha string `json:"start_sha"`
}

type Changes struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
	Diff        string `json:"diff"`
}

func (api *Api) getMRDetailedInfo(url string) (MRInfo, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return MRInfo{}, err
	}

	request.Header.Set(utils.PRIVATE_TOKEN_KEY, api.PrivateToken)

	client := &http.Client{}

	// https://:host/api/v4/projects/:project-id/merge_requests/:merge_id/changes
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
		return MRInfo{}, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close() // Remember to close the response body

	var mrInfo MRInfo
	err = json.Unmarshal(responseData, &mrInfo)

	if err != nil {
		log.Fatal(err.Error())
		return MRInfo{}, err
	}

	return mrInfo, nil
}

func mergeRequestURIBuilder(host string, projectId string, mergeId string) string {
	return host + projectId + utils.URL_SLASH + utils.MERGE_REQUEST + utils.URL_SLASH +
		mergeId + utils.URL_SLASH
}
