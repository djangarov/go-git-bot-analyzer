package analyzer

import (
	"log"
	"strconv"
	"time"

	"github.com/djangarov/go-git-bot-analyzer/command-executor"
)

const DEFAULT_RULES = "./static-code-tool/java-rules/java-basic-rules.xml"
const PMD_EXEC_LOCATION = "sh ./static-code-tool/pmd/bin/run.sh"
const REPORTS = "./reports/"
const DEFAULT_REPORT_TYPE = "json"

func Analyze(pathToFile string) (Report, error) {
	currentReport := REPORTS + strconv.FormatInt(time.Now().Unix(), 10) + "-report.json"
	execCommand := pmdCommandBuilder(pathToFile, DEFAULT_RULES, currentReport)
	err, _, _ := command.Execute(execCommand)

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	return Parse(currentReport)
}

func pmdCommandBuilder(pathToFile string, pathToRules string, pathToReport string) command.Command {
	pathArg := command.Argument{
		Argument: "-d",
		Param:    pathToFile,
	}
	formatArg := command.Argument{
		Argument: "-f",
		Param:    DEFAULT_REPORT_TYPE,
	}

	rulesArg := command.Argument{
		Argument: "-R",
		Param:    pathToRules,
	}
	reportArg := command.Argument{
		Argument: "-r",
		Param:    pathToReport,
	}
	var arguments []command.Argument
	arguments = append(arguments, pathArg)
	arguments = append(arguments, formatArg)
	arguments = append(arguments, rulesArg)
	arguments = append(arguments, reportArg)
	command := command.Command{
		App:      PMD_EXEC_LOCATION + command.WHITESPACE + "pmd",
		Argument: arguments,
	}

	return command
}
