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
const MAX_CAPACITYBUFFER int = 20
const MAX_CAPACITYFILE int = 20

func controlFile() bool {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return false
		}
		err = file.Close()
		if err!=nil{
			log.Fatal(err)
		}
	}
	return true
}

func Readfile(port string ,url string) (string,error) {
	if isErr := controlFile(); !isErr {
		return "",errors.New("ControlFileError")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return "",err
	}
	// Scannerın bufferi yapıldı.
	scanner := bufio.NewScanner(file)
	// eğer bu fonksiyona verilen port ve url içinde varsa error dönüldü...
	for scanner.Scan(){
		temp := strings.Split(scanner.Text(),",")
		if url == temp[0]{
			fmt.Println("HAVE URL")
			return temp[1] ,nil
		}
		if port == temp[1]{
			return "" , errors.New("PORT IS NOT AVAILABLE")
		}
	}
	err = file.Close()
	if err !=nil{
		return "",err
	}
	return "",nil
}

func Writefile(url string, port string) error {
	if isErr := controlFile(); !isErr {
		return errors.New("notWrite")
	}
	file, err := os.OpenFile(fileName, os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	_,err = file.WriteString(url +","+ port + "\n")
	if err!=nil{
		return err
	}
	err = file.Close()
	if err!=nil{
		return err
	}
	return nil
}
func Deletefile(port string, url string)error{
	file , err := os.Open(fileName)
	if err!=nil{
		return err
	}
	scanner := bufio.NewScanner(file)
	// eğer bu fonksiyona verilen port ve url içinde varsa error dönüldü...
	var filetxt [MAX_CAPACITYFILE]string
	var i int = -1
	for scanner.Scan(){
		i++
		temp := strings.Split(scanner.Text(),",")
		//fmt.Println("Scanner temp 0 : "+temp[0]+ "\n Scanner temp 1"+ temp[1])
		if url == temp[0] && port == temp[1]{
			continue
		}
		filetxt[i] = scanner.Text()
	}
	err = file.Close()
	if err!=nil{
		return err
	}
	err = os.Remove(fileName)
	if err!=nil{
		return err
	}
	if isErr := controlFile();!isErr{
		return errors.New("ControlFileError")
	}
	file2,err := os.OpenFile(fileName,os.O_RDWR,0644)
	if err !=nil{
		return err
	}

	for j := 0; j < len(filetxt) ; j++{
		if filetxt[j] == "" && filetxt[j+1] == "" {
			continue
		}
		_,err = file2.WriteString(filetxt[j] + "\n")
		if err !=nil{
			return err
		}
	}
	err = file2.Close()
	if err!=nil{
		return err
	}

	return nil
}
