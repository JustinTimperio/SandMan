package shared

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ContainsAny checks if a string exists within a list of strings.
func ContainsAny(str string, elements []string) bool {
	for element := range elements {
		e := elements[element]
		if strings.Contains(str, e) {
			return true
		}
	}

	return false
}

func DoesFileExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		fmt.Println(err)
	}

	return !os.IsNotExist(err)
}

func DoesFileContain(file *os.File, stringsToBeFound ...string) bool {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {

			if !os.IsTimeout(err) && err != io.EOF {
				fmt.Println(err)
			}

			return false
		}

		for _, stringToBeFound := range stringsToBeFound {
			if strings.Contains(line, stringToBeFound) {
				return true
			}
		}
	}
}
