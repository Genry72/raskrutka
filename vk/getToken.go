package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)
//GetToken Получаем токен для пользователя
func GetToken(username, password string) (token, userID string, err error) {
	url := "https://oauth.vk.com/token?grant_type=password&scope=notify,friends&client_id=3140623&client_secret=VeWdmVclDCtn6ihuP1nt&username=" + username + "&password=" + url.QueryEscape(password)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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
		return token, userID, err
	}
	if strings.Contains(string(body), "error") { //Если в боди вернулась ошибка то прекращаем
		err = fmt.Errorf(string(body))
		return
	}
	gn := tokenStruct{}
	err = json.Unmarshal(body, &gn)
	if err != nil {
		err = fmt.Errorf("ошибка парсинга боди на запрос getToken: %v/ Боди: %v", err, string(body))
		return
	}
	token = gn.AccessToken
	userID = strconv.FormatInt(int64(gn.UserID), 10)
	return
}

type tokenStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Error       string `json:"error"`
	CaptchaSid  string `json:"captcha_sid"`
	CaptchaImg  string `json:"captcha_img"`
}
