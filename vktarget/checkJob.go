package vktarget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CheckJob(jobID, sessID string) (err error) {
	infoLog.Printf("Проверяем корректность выполнения задания %v", jobID)
	url := "https://vktarget.ru/api/all.php?action=check_task&tid=" + jobID + "&host_state=vktarget.ru"
	method := "POST"
	payload := strings.NewReader(`action=check_task&tid=` + jobID + `&host_state=vktarget.ru`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return
	}
	req.Header.Add("authority", "vktarget.ru")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("content-type", "text/plain;charset=UTF-8")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://vktarget.ru")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://vktarget.ru/list/")
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9")
	req.Header.Add("Cookie", "PHPSESSID="+sessID)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("ошибка ")
	}
	t := VktargetAnswerStruct{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		err = fmt.Errorf("ошибка парсинга боди: %v %v", err, string(body))
		return
	}
	if t.Code != 200 {
		err = fmt.Errorf("ошибка проверки задания %v %v", jobID, string(body))
		return
	}
	infoLog.Printf("Задание %v выполнено корректно", jobID)
	return
}

type VktargetAnswerStruct struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}
