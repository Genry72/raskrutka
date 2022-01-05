package vktarget

import (
	"fmt"
	"log"
	"os"
)

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var reset = "\033[0m"
var infoLog = log.New(os.Stdout, fmt.Sprint(string(colorGreen), "INFO\t"+reset), log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, fmt.Sprint(string(colorRed), "ERROR\t"+reset), log.Ldate|log.Ltime|log.Lshortfile)

func MainVktarget(loginVK, passVK string) {
	jobsList, err := GetDjob(loginVK, passVK)
	if err != nil {
		errorLog.Println(err)
		return
	}
	infoLog.Println(jobsList)
}
