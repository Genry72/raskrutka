package vk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//WallRepost делаем репост указанной записи
func WallRepost(tokenVK, versionApiVK, object string) (err error) {
	infoLog.Printf("Делаем репост %v", object)
	url := "https://api.vk.com/method/wall.repost?v=" + versionApiVK + "&access_token=" + tokenVK + "&object=" + object
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
		err = fmt.Errorf("Ошибка репоста " + res.Status + " " + string(body))
		return err
	}
	if strings.Contains(string(body), "failed") || strings.Contains(string(body), "error") {
		err = fmt.Errorf("ошибка репоста: %v", string(body))
		return err
	}
	infoLog.Printf("Репост сделан %v", object)
	return
}
