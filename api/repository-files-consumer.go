package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/djangarov/go-git-bot-analyzer/utils"
)

func (api *Api) GetRepositoryFile(filePath string, fileName string, projectId string, commitSha string) (string, string, error) {
	fileUrl := repositoryFilesURLBuilder(api.GitlabHost, projectId, filePath, commitSha)

	// https://:host/api/v4/projects/:project-id/repository/files/:file-path/raw?ref=:commit_sha
	request, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set(utils.PRIVATE_TOKEN_KEY, api.PrivateToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
		return "", "", err
	}
	tempFilePath := "temp-" + strconv.FormatInt(time.Now().Unix(), 10)
	directory := filepath.Join(".", "temp", tempFilePath)
	path := filepath.Join(directory, fileName)

	out, err := Create(path)
	defer out.Close() // Remember to close the resource

	_, err = io.Copy(out, response.Body)
	defer response.Body.Close() // Remember to close the response body

	if err != nil {
		log.Fatal("Couldn't extract file for validation", err.Error())
		return "", "", err
	}

	return path, directory, nil
}

func repositoryFilesURLBuilder(host string, projectId string, filepath string, commitSha string) string {
	encodedFilePath := url.PathEscape(filepath)
	return host + projectId + utils.URL_SLASH + utils.REPOSITORY +
		utils.URL_SLASH + utils.FILES + utils.URL_SLASH + encodedFilePath +
		utils.URL_SLASH + utils.RAW + "?" + utils.REF + "=" + commitSha
}

func Create(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}
