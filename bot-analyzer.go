package main

import (
	"flag"
	"strings"

	"github.com/djangarov/go-git-bot-analyzer/api"
	"github.com/djangarov/go-git-bot-analyzer/utils"
)

func main() {
	host := flag.String("gitlab-host", "https://gitlab.com", "provide the gitlab host of your repo")
	port := flag.Int("port", 8080, "provide the port that the gitlab bot must run on")
	privateToken := flag.String("private-token", "", "provide the token of your gitlab account")

	flag.Parse()

	// if the provided host doesn't end with '/', append '/'
	if strings.LastIndex(*host, utils.URL_SLASH) != len(*host)-1 {
		*host += utils.URL_SLASH
	}

	// append gitlab api prefix api/v4/projects
	*host += utils.GITLAB_API_PREFIX

	app := api.Api{
		Port:         *port,
		PrivateToken: *privateToken,
		GitlabHost:   *host,
	}

	app.Run()
}
