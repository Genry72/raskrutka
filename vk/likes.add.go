package vk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func LikesADD(tokenVK, versionApiVK, typeObject, ownerID, itemID string) (err error) {
	infoLog.Printf("Лайкаем %v %v %v", typeObject, ownerID, itemID)
	url := "https://api.vk.com/method/likes.add?v=" + versionApiVK + "&access_token=" + tokenVK + "&type=" + typeObject + "&owner_id=" + ownerID + "&item_id=" + itemID
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
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
		err = fmt.Errorf("не корректный код ответа %v %v %v", ownerID, itemID, string(body))
		return
	}
	if strings.Contains(string(body), "failed") || strings.Contains(string(body), "error") {
		err = fmt.Errorf("не корректный ответ %v %v %v", ownerID, itemID, string(body))
		return
	}
	infoLog.Printf("Лайкнули %v %v", ownerID, itemID)
	return
}
