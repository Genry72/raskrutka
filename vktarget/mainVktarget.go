package vktarget

import (
	"fmt"
	"log"
	"math/rand"
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
	//Получаем токен пользователя
	tokenVK, _, err := vk.GetToken(loginVK, passVK)
	if err != nil {
		err = fmt.Errorf("ошибка получения токена для пользователя %v %v ", loginVK, err)
		errorLog.Println(err)
		return
	}
	//Делаем куку рабочей
	err = all.GetUloginToken(loginVK, passVK, sessID, site)
	if err != nil {
		errorLog.Printf("ошибка получения UloginToken %v", err)
		return
	}
	//Раз в 30 минут продлеваем куку
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			//Делаем куку рабочей
			err = all.GetUloginToken(loginVK, passVK, sessID, site)
			if err != nil {
				errorLog.Printf("ошибка получения UloginToken %v", err)
				return
			}
		}
	}()
	//Раз в 30 секунд проверяем задания
	for {
		rand.Seed(time.Now().Unix())                          //Seed для рандома
		son := time.Duration(30+rand.Intn(300)) * time.Second // Сон, рандом от 1 мин до 5 минут
		//Получаем задания
		jobList, err := GetDjob(loginVK, passVK, sessID)
		if err != nil {
			errorLog.Printf("Ошибка получения информации по заданиям vktarget: %v", err)
			infoLog.Printf("Спим %v перед следующей проверкой задания из-за ошибки", son)
			time.Sleep(son)
			continue
		}
		time.Sleep(5 * time.Second)
		//Выполняем задания
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
					err = fmt.Errorf("ошибка проверки задания %v %v", jobID, err)
					errorLog.Println(err)
				}
			case "Нажмите поделиться записью":
				object := strings.Split(value[1], "/")
				err = vk.WallRepost(tokenVK, versionApiVK, object[1])
				if err != nil {
					err = fmt.Errorf("ошибка репоста %v", err)
					errorLog.Println(err)
				}
				time.Sleep(3 * time.Second)
				//Проверяем задание
				err = CheckJob(jobID, sessID)
				if err != nil {
					err = fmt.Errorf("ошибка проверки задания %v %v", jobID, err)
					errorLog.Println(err)
				}
			case "Поставьте лайк на странице":
				infoLog.Printf("Лайк %v", value[1])
			case "Расскажите о группе":
				infoLog.Printf("Рассказать о группе %v", value[1])
			default:
				err = fmt.Errorf("не известный тип задания %v", value[0])
				errorLog.Println(err)
			}
		}
		infoLog.Printf("Спим %v перед следующей проверкой задания", son)
		// infoLog.Printf("Спим 10c перед следующей проверкой задания", son)
		time.Sleep(son)
	}
}
