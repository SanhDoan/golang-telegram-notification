package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	LISTEN string
	TELEGRAM_BOT_TOKEN string
	TELEGRAM_CHAT_ID string
}

var (
	config = &Config{}
)

func loadConfig() *Config {
	conf := Config{}
	content, err := ioutil.ReadFile("./config.json")
	checkError(err)
	err = json.Unmarshal(content, &conf)
	checkError(err)
	return &conf
}

func doSendTelegram(user string, remoteip string) {
	text := fmt.Sprintf("%s - %s", user, remoteip)
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage?chat_id=%v&text=%v", config.TELEGRAM_BOT_TOKEN, config.TELEGRAM_CHAT_ID, text)

	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	res, err := http.DefaultClient.Do(req)
	checkError(err)

	defer res.Body.Close()
}

func notifyLogin(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	user := keys.Get("user")
	remoteIp := keys.Get("remoteip")
	doSendTelegram(user, remoteIp)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config = loadConfig()
	http.HandleFunc("/notify/login", notifyLogin)
	err := http.ListenAndServe(config.LISTEN, nil)
	checkError(err)
}
