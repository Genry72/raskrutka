package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"raskrutka/vktarget"

	"gopkg.in/yaml.v2"
)

var versionApiVK = "5.89"

//Подсветка ошибок и удачных сообщений
var colorRed = "\033[31m"

// var colorGreen = "\033[32m"
var reset = "\033[0m"

// var infoLog = log.New(os.Stdout, fmt.Sprint(string(colorGreen), "INFO\t"+reset), log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, fmt.Sprint(string(colorRed), "ERROR\t"+reset), log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	//Парсим настройки
	data, err := ioutil.ReadFile("./settings.yaml")
	if err != nil {
		err = fmt.Errorf("не удалось прочитать настройки %v", err)
		errorLog.Println(err)
		return
	}
	s := settingsStruct{}
	err = yaml.Unmarshal([]byte(data), &s)
	if err != nil {
		err = fmt.Errorf("не удалось прочитать настройки %v", err)
		errorLog.Println(err)
		return
	}

	vktarget.MainVktarget(s.LoginVK, s.PsswordVK, versionApiVK)
}

type settingsStruct struct {
	LoginVK   string `yaml:"loginVK"`
	PsswordVK string `yaml:"psswordVK"`
}
