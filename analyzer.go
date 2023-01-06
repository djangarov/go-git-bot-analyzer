package main

import (
	"flag"
	"fmt"
)

func main() {
	hostFlag := flag.String("gitlab-host", "https://gitlab.com", "provide the gitlab host of your repo")
	portFlag := flag.Int("port", 8080, "provide the port that the gitlab bot must run on")
	privateTokenFlag := flag.String("private-token", "", "provide the token of your gitlab account")

	flag.Parse()

	fmt.Println(*hostFlag)
	fmt.Println(*portFlag)
	fmt.Print(*privateTokenFlag)
}
