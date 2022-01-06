package vk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//GroupsJoin Вступление в группу
func GroupsJoin(versinonAPIVk, token, groupID string) (err error) {
	infoLog.Printf("Подписываемся к группе %v\n", groupID)
	url := "https://api.vk.com/method/groups.join?v=" + versinonAPIVk + "&access_token=" + token + "&group_id=" + groupID
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("ошибка вступления в группу %v %v", groupID, string(body))
		return
	}
	if strings.Contains(string(body), "failed") || strings.Contains(string(body), "error") {
		err = fmt.Errorf("ошибка вступления в группу %v: %v", groupID, string(body))
		return
	}
	infoLog.Printf("Подписались к группе %v\n", groupID)
	return err
}
