package library

import (
	"bitbucket.org/dnda-tech/api-and-billing/app/logger"
	"fmt"
	"github.com/go-cmd/cmd"
	"path"
	"strconv"
	"strings"
)

func NumberOfLines(file string) int {

	envCmd := cmd.NewCmd("wc", "-l", file)

	// Run and wait for Cmd to return Status
	status := <-envCmd.Start()

	var response1 string

	// Print each line of STDOUT from Cmd
	for _, line := range status.Stdout {
		response1 = line
	}

	words := strings.Fields(response1)

	fmt.Println(words[0])

	// Create Cmd, buffered output
	//envCmd := cmd.NewCmd("awk","'END{print NR}' ",file)
	envCmd = cmd.NewCmd("wc", "-l", file)

	// Run and wait for Cmd to return Status
	status = <-envCmd.Start()

	var response string

	// gets each line of STDOUT from Cmd
	for i, line := range status.Stdout {
		logger.Info("got index %d is %s", i, line)
		response = line
	}

	words = strings.Fields(response)

	count, _ := strconv.Atoi(words[0])

	logger.Info("NumberOfLines checking number of files for file %s, got %d lines from %s", file, count, words[0])

	return count
}

func GetFileExtension(fn string) string {

	return strings.TrimSpace(strings.ReplaceAll(path.Ext(fn),".",""))

}

func RandomFileName(length int) (string, error) {

	letters := fmt.Sprintf("%s%s%s", LowerLetters, UpperLetters, Digits)

	code := ""

	for i := 0; i < length; i++ {

		sym, err := RandomElement(letters)

		if err != nil {
			return "", err
		}

		code, err = RandomInsert(code, sym)
		if err != nil {

			return "", err
		}
	}

	return code, nil
}