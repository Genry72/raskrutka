package vktarget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"raskrutka/all"
)

//getDjob получаем список заданий. Возвращает список мап, ключь id задания, значение - список из 0. Имя 1. Ссылка
func GetDjob(loginVK, passVK string) (jobList map[string][]string, err error) {
	site := "vktarget"
	sessID, err := all.GetphpSessID(site)
	if err != nil {
		err = fmt.Errorf("ошибка получения sessID: %v", err)
		return
	}
	//Делаем куку рабочей
	err = all.GetUloginToken(loginVK, passVK, sessID, site)
	if err != nil {
		err = fmt.Errorf("ошибка получения токена %v", err)
		return
	}
	url := "https://vktarget.ru/api/all.php?action=get_list&v=1.2&offset=0"
	method := "POST"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("authority", "vktarget.ru")
	req.Header.Add("content-length", "0")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
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
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	t := JobStruct{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		err = fmt.Errorf("ошибка парсинга боди: %v %v", err, string(body))
		return
	}
	jobList = make(map[string][]string)
	switch t.Tasks.(type) { //Проверяем тип интерфейса. Может возвращать "tasks": [] - не задач и мапу, если задачи есть
	case map[string]interface{}:
		for _, m := range t.Tasks.(map[string]interface{}) { //Идим по задачам
			var jobID string
			var typeName1 string
			var typeName2 string
			var uri string
			for k, value := range m.(map[string]interface{}) { //Идем по значениям в каждой задаче
				switch k {
				case "id":
					jobID = value.(string)
				case "type_name":
					typeName1 = value.(string)
				case "type_name_link":
					typeName2 = value.(string)
				case "url":
					uri = value.(string)
				}
			}
			jobList[jobID] = []string{fmt.Sprintf("%v %v", typeName1, typeName2), uri}
		}
		return
	case []interface{}:
		err = fmt.Errorf("пустой список заданий")
		return
	default:
		err = fmt.Errorf("ожидался другой тип данных в интерфейсе %T", t.Tasks)
		return
	}
}

//JobStruct структура ответа запроса на список заданий
type JobStruct struct {
	UserBalance      float64       `json:"user_balance"`
	Social           Social        `json:"social"`
	ListNotification int           `json:"list_notification"`
	Permissions      Permissions   `json:"permissions"`
	CountErrors      int           `json:"count_errors"`
	Seconds          int           `json:"seconds"`
	CountDone        int           `json:"count_done"`
	CountInvalid     int           `json:"count_invalid"`
	Tip              Tip           `json:"tip"`
	CountWait        int           `json:"count_wait"`
	Available        int           `json:"available"`
	VktCaptcha       int           `json:"vkt_captcha"`
	UID              string        `json:"uid"`
	Timers           []interface{} `json:"timers"`
	// Tasks            map[string]map[string]interface{} `json:"tasks"`
	Tasks interface{} `json:"tasks"`
}
type Social struct {
	Vk         bool   `json:"vk"`
	VkID       string `json:"vk_id"`
	Fb         bool   `json:"fb"`
	Tw         bool   `json:"tw"`
	In         bool   `json:"in"`
	Yt         bool   `json:"yt"`
	Ok         bool   `json:"ok"`
	Quora      bool   `json:"quora"`
	Tiktok     bool   `json:"tiktok"`
	Tumblr     bool   `json:"tumblr"`
	Vimeo      bool   `json:"vimeo"`
	Mixcloud   bool   `json:"mixcloud"`
	Soundcloud bool   `json:"soundcloud"`
	Likee      bool   `json:"likee"`
	Reddit     bool   `json:"reddit"`
	Zen        bool   `json:"zen"`
	Telegram   bool   `json:"telegram"`
}
type Notifications struct {
	Type     int    `json:"type"`
	Value    int    `json:"value"`
	Permname string `json:"permname"`
}
type Types struct {
	Type     int    `json:"type"`
	Value    int    `json:"value"`
	Permname string `json:"permname"`
}
type Vk struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Yt struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Ok struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Android struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Quora struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Tiktok struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Tumblr struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Vimeo struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Soundcloud struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Reddit struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Zen struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type Telegram struct {
	Type     int     `json:"type"`
	Wtype    int     `json:"wtype"`
	Value    int     `json:"value"`
	Permname string  `json:"permname"`
	Types    []Types `json:"types"`
}
type SocNetworksPermissions struct {
	Vk         Vk         `json:"vk"`
	Yt         Yt         `json:"yt"`
	Ok         Ok         `json:"ok"`
	Android    Android    `json:"android"`
	Quora      Quora      `json:"quora"`
	Tiktok     Tiktok     `json:"tiktok"`
	Tumblr     Tumblr     `json:"tumblr"`
	Vimeo      Vimeo      `json:"vimeo"`
	Soundcloud Soundcloud `json:"soundcloud"`
	Reddit     Reddit     `json:"reddit"`
	Zen        Zen        `json:"zen"`
	Telegram   Telegram   `json:"telegram"`
}
type MinPrice struct {
	Type     int    `json:"type"`
	Value    string `json:"value"`
	Permname string `json:"permname"`
}
type MaxPrice struct {
	Type     int    `json:"type"`
	Value    string `json:"value"`
	Permname string `json:"permname"`
}
type TaskPricesPermissions struct {
	MinPrice MinPrice `json:"min_price"`
	MaxPrice MaxPrice `json:"max_price"`
}
type Permissions struct {
	Notifications          Notifications          `json:"notifications"`
	SocNetworksPermissions SocNetworksPermissions `json:"soc_networks_permissions"`
	TaskPricesPermissions  TaskPricesPermissions  `json:"task_prices_permissions"`
}
type Tip struct {
	Type int    `json:"type"`
	Tip  string `json:"tip"`
}
