package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"time"
)

const (
	help              = "Здравствуйте сэр,я бот,который показывает прогноз погоды в любом известном для вас городе. Для этого напишите название города👀"
	openWeatherMapAPI = "e71cc7509ba7040322d574ebdad1b5c3"
)

func GetWeather(city string, openWeatherMapAPI string) string {
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + openWeatherMapAPI + "&units=metric"
	res, err := http.Get(url) //отправляет запрос к url,res-ответ серв.
	if err != nil {
		fmt.Println(err)
		fmt.Println("Проверьте название города")
		return ""

	}
	defer res.Body.Close() //res.body-ответ серв,закрываем его в конце,избегая утечку данных

	data := make(map[string]interface{})          //данные из json объекта
	err = json.NewDecoder(res.Body).Decode(&data) //декодируем данные из res.body и записываем их в data
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", data) //содержимое data в консоль

	city = data["name"].(string)
	weather := data["main"].(map[string]interface{})
	curWeather := weather["temp"].(float64)
	weatherDescription := data["weather"].([]interface{})[0].(map[string]interface{})["main"].(string)
	wd := ""                                                           //смайлик
	if weatherDescription, ok := codeToSmile[weatherDescription]; ok { //если ok tr,то присваиваем
		wd = weatherDescription
	} else {
		wd = "Посмотри в окно, не пойму что там за погода!"
	}

	humidity := weather["humidity"].(float64)
	pressure := weather["pressure"].(float64)
	wind := data["wind"].(map[string]interface{})["speed"].(float64)

	sunriseTimestamp := time.Unix(int64(data["sys"].(map[string]interface{})["sunrise"].(float64)), 0)
	sunsetTimestamp := time.Unix(int64(data["sys"].(map[string]interface{})["sunset"].(float64)), 0)
	lengthOfDay := sunsetTimestamp.Sub(sunriseTimestamp)

	weatherData := fmt.Sprintf("***%s***\nПогода в городе: %s🌌\nТемпература: %.2fC°🌡 %s\nВлажность: %.0f%%💦\nДавление: %.0f мм.рт.ст\nВетер: %.2f м/с💨️\nВосход солнца: %s☀️\nЗакат солнца: %s🌇\nПродолжительность дня: %s🌍\nХорошего дня!👋",
		time.Now().Format("2006-01-02 15:04"),
		city,
		curWeather, //темпа
		wd,         //тип погоды
		humidity,   //влажность
		pressure,   //давление
		wind,       //ветер
		sunriseTimestamp.Format("2006-01-02 15:04:05"), //восход с
		sunsetTimestamp.Format("2006-01-02 15:04:05"),  //закат с
		lengthOfDay.String(),
	)

	return weatherData
}

var codeToSmile = map[string]string{
	"Clear":        "Ясно ☀️",
	"Clouds":       "Облачно ☁️",
	"Drizzle":      "Морось 🌧",
	"Rain":         "Дождь 🌧",
	"Thunderstorm": "Гроза ⛈",
	"Snow":         "Снег ❄️",
	"Mist":         "Туман 🌫",
}

func main() {

	bot, err := tgbotapi.NewBotAPI("6537307160:AAEizhtTdKGu1ez5Jeb_uWjvfZg43GxTDaI")
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, help)
			bot.Send(msg)

		} else if update.Message.Text == "/help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Просто отправь мне название города и я покажу тебе погоду!")
			bot.Send(msg)
		} else {
			city := update.Message.Text
			go GetWeather(city, openWeatherMapAPI)
			weatherData := GetWeather(city, openWeatherMapAPI)

			if len(weatherData) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю, попробуйте отправить название города.")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherData)
				bot.Send(msg)
			}

		}
	}
}
