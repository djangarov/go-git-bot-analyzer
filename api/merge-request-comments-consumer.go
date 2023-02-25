package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/djangarov/go-git-bot-analyzer/analyzer"
	"github.com/djangarov/go-git-bot-analyzer/utils"
)

func (api *Api) CommentOn(mergeId string, projectId string, violation analyzer.Violation, filePath string, refs DiffRefs) {

	values := url.Values{
		"position[base_sha]":      {refs.BaseSha},
		"position[start_sha]":     {refs.StartSha},
		"position[head_sha]":      {refs.HeadSha},
		"position[new_path]":      {filePath},
		"position[position_type]": {"text"},
		"position[new_line]":      {strconv.Itoa(violation.BeginLine)},
		"body":                    {commentBodyBuilder(violation)},
	}
	commentUrl := commentURLBuilder(api.GitlabHost, mergeId, projectId)
	strReader := strings.NewReader(values.Encode())
	request, err := http.NewRequest("POST", commentUrl, strReader)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set(utils.PRIVATE_TOKEN_KEY, api.PrivateToken)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		log.Fatal("Error on comment merge request:: " + string(bodyBytes))
	}
}

func commentURLBuilder(host string, mergeId string, projectId string) string {
	return host + projectId + utils.URL_SLASH + utils.MERGE_REQUEST + utils.URL_SLASH +
		mergeId + utils.URL_SLASH + utils.DISCUSSIONS
}

func commentBodyBuilder(violation analyzer.Violation) string {
	return violation.Rule + ": " + violation.Description
}
