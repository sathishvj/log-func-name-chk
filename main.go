package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var funcBegin = regexp.MustCompile(`func.* (\w+)\(.*\{$`)
var funcEnd = regexp.MustCompile(`^\}.*$`)
var logPrefixRegex = regexp.MustCompile(`\s*logger\..*\(.*$`)

func main() {
	dir := os.Args[1]

	files := getFilesInDir(dir)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") {
			continue
		}

		filePath := dir + string(os.PathSeparator) + f.Name()
		//fmt.Printf("File name is: %s\n", filePath)
		linesOfFile := readLinesOfFile(filePath)
		funcsOfFiles := getFuncsOfFile(linesOfFile)

		for _, oneFunc := range funcsOfFiles {
			checkLoggerFunctionsMatch(f.Name(), oneFunc)
		}
	}
}

func readLinesOfFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func getFuncsOfFile(lines []string) [][]string {

	var allFuncs [][]string

	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "//") {
			i++
			continue
		}

		if funcBegin.MatchString(lines[i]) {
			var oneFunc []string
			for {
				if i >= len(lines) {
					break
				}
				if strings.HasPrefix(strings.TrimSpace(lines[i]), "//") {
					i++
					continue
				}

				if funcEnd.MatchString(lines[i]) {
					oneFunc = append(oneFunc, lines[i])
					i++
					break
				}

				oneFunc = append(oneFunc, lines[i])
				i++
			}

			allFuncs = append(allFuncs, oneFunc)
		}
	}

	return allFuncs
}

func checkLoggerFunctionsMatch(filename string, funcLines []string) {
	//first line has to contain name of function
	if len(funcLines) < 2 {
		panic(fmt.Sprintf("Empty function given?\n%v\n", funcLines))
	}

	matchingNames := funcBegin.FindStringSubmatch(funcLines[0])
	if len(matchingNames) != 2 {
		panic(fmt.Sprintf("More than one matching func name? \n%v\n", matchingNames))
	}
	funcName := matchingNames[1]
	//fmt.Printf("\nFunc name is: %s\n", funcName)

	for i, oneLine := range funcLines {
		if !logPrefixRegex.MatchString(oneLine) {
			continue
		}

		// todo: make these two part of exceptions in the matcher
		if strings.TrimSpace(oneLine) == "logger.Errorf(msg)" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(oneLine), "//") {
			continue
		}

		// else it is a match
		if !strings.Contains(oneLine, funcName) {
			fmt.Printf("Mismatched log func prefix in %s:%s():~%d: %s\n", filename, funcName, i, oneLine)
		}
	}
}

func getFilesInDir(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Could not read dir contents: %+v", err))
	}

	return files
}
