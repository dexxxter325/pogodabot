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

type WeatherData struct {
	City        string
	Humidity    float64
	Pressure    float64
	WindSpeed   float64
	Sunrise     time.Time
	Sunset      time.Time
	CurWeather  float64
	wd          string
	LengthOfDay time.Duration
}

func GetWeather(city string, openWeatherMapAPI string) WeatherData {
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + openWeatherMapAPI + "&units=metric"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Проверьте название города")

	}
	defer res.Body.Close()

	data := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", data)

	weather := data["main"].(map[string]interface{})
	weatherDescription := data["weather"].([]interface{})[0].(map[string]interface{})["main"].(string)
	var _ string
	if weatherDescription, ok := codeToSmile[weatherDescription]; ok {
		_ = weatherDescription
	} else {
		_ = "Посмотри в окно, не пойму что там за погода!"
	}
	sunriseTimestamp := time.Unix(int64(data["sys"].(map[string]interface{})["sunrise"].(float64)), 0)
	sunsetTimestamp := time.Unix(int64(data["sys"].(map[string]interface{})["sunset"].(float64)), 0)
	weatherData := WeatherData{
		City:        data["name"].(string),
		Humidity:    weather["humidity"].(float64),
		Pressure:    weather["pressure"].(float64),
		WindSpeed:   data["wind"].(map[string]interface{})["speed"].(float64),
		Sunrise:     time.Unix(int64(data["sys"].(map[string]interface{})["sunrise"].(float64)), 0),
		Sunset:      time.Unix(int64(data["sys"].(map[string]interface{})["sunset"].(float64)), 0),
		CurWeather:  weather["temp"].(float64),
		wd:          "Посмотри в окно, не пойму что там за погода!",
		LengthOfDay: sunsetTimestamp.Sub(sunriseTimestamp),
	}

	// Возвращаем данные о погоде
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
			city :=//////////////////////////////////
			weatherData := GetWeather(city, "e71cc7509ba7040322d574ebdad1b5c3")
			go GetWeather(city, "e71cc7509ba7040322d574ebdad1b5c3")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("***%s***\nПогода в городе: %s\nТемпература: %.2fC° %s\nВлажность: %.0f%%\nДавление: %.0f мм.рт.ст\nВетер: %.2f м/с\nВосход солнца: %s\nЗакат солнца: %s\nПродолжительность дня: %s\nХорошего дня!",
				time.Now().Format("2006-01-02 15:04"),
				city,
				weatherData.CurWeather,
				weatherData.wd,
				weatherData.Humidity,
				weatherData.Pressure,
				weatherData.WindSpeed,
				weatherData.Sunrise.Format("2006-01-02 15:04:05"),
				weatherData.Sunset.Format("2006-01-02 15:04:05"),
				weatherData.LengthOfDay.String(),
			))
			bot.Send(msg)

		}
	}
}

