package differ

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/djangarov/go-git-bot-analyzer/utils"
)

func ExtractModifiedLines(str string) map[int]int {
	var changedSpecter []string

	index := strings.Index(str, utils.GIT_HEAD_PREFIX)
	firstSegment := str[index+1:]
	index = strings.Index(firstSegment, utils.GIT_HEAD_PREFIX)
	secondSegment := firstSegment[:index+1]
	changedSpecter = append(changedSpecter, secondSegment)

	for index >= 0 {
		str = firstSegment[index+2:]
		index = strings.Index(str, utils.GIT_HEAD_PREFIX)
		firstSegment = str[index+1:]
		index = strings.Index(firstSegment, utils.GIT_HEAD_PREFIX)
		secondSegment = firstSegment[:index+1]
		changedSpecter = append(changedSpecter, secondSegment)
	}

	return createDiapasonMap(changedSpecter)
}

func createDiapasonMap(changedSpecter []string) map[int]int {
	diapasonMap := make(map[int]int)

	for _, specter := range changedSpecter {
		if len(specter) < 2 {
			continue
		}
		formattedSpecter := specter[1 : len(specter)-2]
		formattedSpecter = strings.TrimSpace(formattedSpecter)
		fields := strings.Fields(formattedSpecter)
		if len(fields) != 2 {
			return diapasonMap
		}

		addedSpecter := (fields[1])[1:]

		diapasons := strings.Split(addedSpecter, ",")
		if len(diapasons) != 2 {
			return diapasonMap
		}

		startLine, err := strconv.Atoi(diapasons[0])

		if err != nil {
			log.Fatal(err.Error())
		}

		lenght, err := strconv.Atoi(diapasons[1])

		if err != nil {
			log.Fatal(err.Error())
			return diapasonMap
		}

		diapasonMap[startLine] = startLine + lenght
	}
	return diapasonMap
}

func IsBetween(modifiedLine int, diapason map[int]int) bool {
	isLineBetween := func(start, finish, target int) bool {
		return target >= start && target <= finish
	}

	result := false

	for startLine, endLine := range diapason {
		if isLineBetween(startLine, endLine, modifiedLine) {
			result = true
		}
	}
	return result
}

func GetAddedLines(lines []int, path string) map[int]bool {
	file, err := os.Open(path)
	availabilityMap := make(map[int]bool)

	if err != nil {
		log.Fatal(err)
		return availabilityMap
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// set line 1 as a starting line
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()

		// if its added line, it must start with '+'
		isAddedLine := strings.Index(line, utils.GIT_LINE_ADD) == 0

		if isAddedLine && isLineInArr(lineNumber, lines) {
			availabilityMap[lineNumber] = true
		}
		// if the line starts with @@ change the line number
		if strings.Index(line, utils.GIT_HEAD_PREFIX) == 0 {
			lineNumber, _ = getNextStartingLine(line)
			// decrement the line number by one, to get the actual line number
			lineNumber = lineNumber - 1
		} else if strings.Index(line, utils.GIT_LINE_REMOVE) == 0 {
			// if the line starts with '-', it means that the line is deleted,
			// so we dont take in into consideration
			continue
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return availabilityMap
	}

	return availabilityMap
}

func getNextStartingLine(line string) (int, error) {
	// slice the string with two positions to ignore the GIT_HEAD_PREFIX
	line = line[2:]
	index := strings.Index(line, utils.GIT_HEAD_PREFIX)
	// slice the string with one position to ignore the GIT_HEAD_PREFIX
	line = line[:index-1]
	line = strings.TrimSpace(line)
	// trim by whitespace
	fields := strings.Fields(line)
	addedLine := fields[1]
	addedLine = addedLine[1:]
	nextStartingLine := strings.Split(addedLine, utils.COMMA_SEPARATOR)[0]
	return strconv.Atoi(nextStartingLine)
}

// since the lines array is sorted, we can do binary search to optimize the search
func isLineInArr(line int, lines []int) bool {
	for _, l := range lines {
		if line == l {
			return true
		}
	}
	return false
}
