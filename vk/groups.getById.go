package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GroupsGetById возвращает информацию о группе. На вход принимает id или псевдоним группы
func GroupsGetById(groupName, token, versionApiVK string) (id, name string, err error) {
	url := "https://api.vk.com/method/groups.getById?v=" + versionApiVK + "&access_token=" + token + "&group_id=" + groupName
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
		err = fmt.Errorf(string(body))
		return
	}
	s := GroupsGetByIdStruct{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		err = fmt.Errorf("ошибка парсинга:%v Боди: %v", err, string(body))
		return
	}
	if len(s.Response) == 0 {
		err = fmt.Errorf("хрень в боди %v", string(body))
		return
	}
	id = fmt.Sprint(s.Response[0].ID)
	name = s.Response[0].Name
	return
}

type GroupsGetByIdStruct struct {
	Response []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		ScreenName   string `json:"screen_name"`
		IsClosed     int    `json:"is_closed"`
		Type         string `json:"type"`
		IsAdmin      int    `json:"is_admin"`
		IsMember     int    `json:"is_member"`
		IsAdvertiser int    `json:"is_advertiser"`
	} `json:"response"`
}
