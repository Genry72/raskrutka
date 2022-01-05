package all

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sclevine/agouti"
)

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var reset = "\033[0m"
var infoLog = log.New(os.Stdout, fmt.Sprint(string(colorGreen), "INFO\t"+reset), log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, fmt.Sprint(string(colorRed), "ERROR\t"+reset), log.Ldate|log.Ltime|log.Lshortfile)

//GetphpSessID получаем из кук phpsessid (аутентификация через соцсети)
func GetphpSessID(site string) (sessID string, err error) {
	var url string
	switch site {
	case "vktarget":
		url = "https://vktarget.ru/"
	}
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return sessID, err
	}
	res, err := client.Do(req)
	if err != nil {
		return sessID, err
	}
	defer res.Body.Close()
	if err != nil {
		return sessID, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return sessID, err
	}
	for _, kuka := range res.Cookies() {
		if kuka.Name == "PHPSESSID" {
			sessID = kuka.Value
		}
	}
	if sessID == "" {
		err = fmt.Errorf("пустаяя кука:%v %v", res.Status, string(body))
	}
	return sessID, err
}

//2. 3. getUloginToken Получает токен и апдейтит sessID (аутентификация через соц-сети)
func GetUloginToken(login, pass, sessID, site string) (err error) {
	infoLog.Printf("Выполняем получение токена Utoken для пользователя %v", login)

	var token string
	// driver := agouti.PhantomJS()
	// driver := agouti.ChromeDriver()
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
		// agouti.ChromeOptions("args", []string{"--disable-gpu", "--no-sandbox"}),
	)

	err = driver.Start()
	if err != nil {
		return err
	}

	page, err := driver.NewPage()
	if err != nil {
		return err
	}
	// page.Size(1920, 1080)
	page.SetImplicitWait(15000000) //Таймаут ожиданя вейта

	if err := page.Navigate("https://oauth.vk.com/authorize?v=5.62&client_id=3280318&scope=friends,schools,email&display=page&response_type=code&redirect_uri=https://ulogin.ru/auth.php?name=vkontakte"); err != nil {
		return err
	}
	_, err = page.FindByID(`install_allow`).Visible()
	if err != nil {
		err = fmt.Errorf("не дождались загрузки кнопки войти:%v", err)
		return err
	}
	//Вводим логин
	err = page.FindByXPath(`/html/body/div/div/div/div[2]/form/div/div/input[6]`).SendKeys(login)
	if err != nil {
		return err
	}
	//Вводим пароль
	err = page.FindByXPath(`/html/body/div/div/div/div[2]/form/div/div/input[7]`).SendKeys(pass)
	if err != nil {
		return err
	}
	//Нажимаем вход
	err = page.FindByID(`install_allow`).Submit()
	if err != nil {
		return err
	}
	//Ждем прогрузки страницы. Когда появится кнопка закрыть окно
	_, err = page.FindByXPath(`/html/body/div/button`).Visible()
	if err != nil {
		err = fmt.Errorf("не дождались загрузки кнопки войти:%v", err)
		return err
	}

	//Забираем боди, там токен
	z, err := page.HTML()
	if err != nil {
		return err
	}

	//Парсим боди
	scanner := bufio.NewScanner(strings.NewReader(z)) //Построчно читаем боди
	for scanner.Scan() {
		var out string
		out = scanner.Text()
		var clearstring string
		ansiCode := []string{
			`</script>`,
			`;`,
			"'",
			"=",
			` `,
			// `"`,
			// `>`,
		}
		//Идем по циклу с цветами и удаляем их
		for _, color := range ansiCode {
			// fmt.Println(color)
			clearstring = strings.Replace(out, color, "", -1)
			// fmt.Println("Нашли")
			out = clearstring

		}
		if strings.Contains(out, "token") {
			s := strings.Split(out, "token")
			token = s[1]
		}
	}
	err = driver.Stop()
	if err != nil {
		return err
	}
	//Продлеваем имеющуюся куку
	err = updatePhpSessID(token, sessID, site)
	if err != nil {
		return err
	}
	infoLog.Printf("Получение токена Utoken для пользователя %v выполнено", login)
	return err
}

//Продлеваем имеющуюся куку
func updatePhpSessID(uloginToken, sessID, site string) (err error) {
	var url string
	var host string
	var origin string
	var referer string
	switch site {
	case "vktarget":
		url = "https://vktarget.ru/app_handlers/update_vk.php"
		host = "vktarget.ru"
		origin = "https://vktarget.ru"
		referer = origin + "/"
	}
	method := "POST"
	payload := strings.NewReader("token=" + uloginToken)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Host", host)
	req.Header.Add("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0")
	req.Header.Add("Accept", " text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", " ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Accept-Encoding", " gzip, deflate")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", " 38")
	req.Header.Add("Origin", origin)
	req.Header.Add("Referer", referer)
	req.Header.Add("Cookie", " PHPSESSID="+sessID)
	req.Header.Add("Upgrade-Insecure-Requests", " 1")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), "Undefined index") {
		err = fmt.Errorf("ошибка обновления PHPSESSID")
		return err
	}
	return err
}
