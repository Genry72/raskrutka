package vktarget

import (
	"fmt"
	"log"
	"os"
	"raskrutka/all"
	"raskrutka/vk"
	"strings"
	"time"
)

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var reset = "\033[0m"
var infoLog = log.New(os.Stdout, fmt.Sprint(string(colorGreen), "INFO\t"+reset), log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, fmt.Sprint(string(colorRed), "ERROR\t"+reset), log.Ldate|log.Ltime|log.Lshortfile)

func MainVktarget(loginVK, passVK, versionApiVK string) {
	site := "vktarget"
	sessID, err := all.GetphpSessID(site)
	if err != nil {
		errorLog.Printf("ошибка получения sessID: %v", err)
		return
	}
	//Делаем куку рабочей
	err = all.GetUloginToken(loginVK, passVK, sessID, site)
	if err != nil {
		errorLog.Printf("ошибка получения UloginToken %v", err)
		return
	}
	//Получаем токен пользователя
	tokenVK, _, err := vk.GetToken(loginVK, passVK)
	if err != nil {
		err = fmt.Errorf("ошибка получения токена для пользователя %v %v ", loginVK, err)
		errorLog.Println(err)
		return
	}
	//Получаем задания
	jobList, err := GetDjob(loginVK, passVK, sessID)
	if err != nil {
		errorLog.Printf("Ошибка получения информации по заданиям vktarget: %v", err)
		return
	}
	time.Sleep(5 * time.Second)
	for jobID, value := range jobList {
		switch value[0] { //Тип заданимя
		case "Вступите в сообщество":
			groupName := strings.Split(value[1], "/") //бьем слешами адрес группы для получения имени
			//Получаем id группы
			groupID, _, err := vk.GroupsGetById(groupName[1], tokenVK, versionApiVK)
			if err != nil {
				err = fmt.Errorf("ошибка получения id группы %v %v", groupName[1], err)
				errorLog.Println(err)
				continue
			}
			//Вступаем в группу
			err = vk.GroupsJoin(versionApiVK, tokenVK, groupID)
			if err != nil {
				err = fmt.Errorf("ошибка вступления в группу %v пользователем %v %v", groupID, loginVK, err)
				errorLog.Println(err)
			}
			time.Sleep(3 * time.Second)
			//Проверяем задание
			err = CheckJob(jobID, sessID)
			if err != nil {
				err = fmt.Errorf("ошибка проверки задания %v", jobID)
				errorLog.Println(err)
			}
		case "Нажмите поделиться записью":
			infoLog.Println("Репост")
		case "Поставьте лайк на странице":
			infoLog.Println("Лайк")
		case "Расскажите о группе":
			infoLog.Println("Рассказать о группе")
		default:
			err = fmt.Errorf("не известный тип задания %v", value[0])
			errorLog.Println(err)
		}
	}
}
