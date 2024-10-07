package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var fileName string = "main.txt"

func controlFile() bool {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return false
		}
		file.Close()
	}
	return true
}

func Readfile(urlport map[string]string, json chan bool) error {
	if isErr := controlFile(); !isErr {
		return errors.New("Exit")
	}
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		fmt.Print(string(line))
		if line == nil {
			fmt.Print("selamlar")
			break
		}
		lineArray := strings.Split(string(line), "|")
		urlport[lineArray[0]] = lineArray[1]
	}
	file.Close()
	return nil
}

func Writefile(key string, value string) error {
	if isErr := controlFile(); !isErr {
		return errors.New("notWrite")
	}
	file, err := os.OpenFile(fileName, os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	file.WriteString(key + "|" + value + "\n")
	file.Close()
	return nil
}

func ControlMap(origin string, port string, urlPort map[string]string) bool {
	for key, value := range urlPort {
		if key == origin && port == value {
			return true
		}
	}
	return false
}
func Deletefile() {
	if controlFile() {
		err := os.Remove(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
