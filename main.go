package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/gocolly/colly"
)

type WeatherData struct {
	Day    string
	High   string
	Low    string
}

func main() {
	data := make([]WeatherData, 0)
	c := colly.NewCollector(colly.AllowedDomains("weather.com", "www.weather.com"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML(".DetailsSummary--fadeOnOpen--KnNyF", func(e *colly.HTMLElement) {
		day := e.ChildText("h3.DetailsSummary--daypartName--kbngc")
		high := e.ChildText("span.DetailsSummary--highTempValue--3PjlX")
		low := e.ChildText("span.DetailsSummary--lowTempValue--2tesQ")
		data = append(data, WeatherData{Day: day, High: high, Low: low})
	})

	err := c.Visit("https://weather.com/weather/tenday/l/fe3c78e80c47c404a4e64ec7c86ceccdb814894cedefdb528f9d8d95c3e4eb74")
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	fileName := fmt.Sprintf("weather_%s.csv", t.Format("20060102150405"))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Day", "HighTemp", "LowTemp"})

	// Write weather data to CSV
	for _, w := range data {
		writer.Write([]string{w.Day, w.High, w.Low})
	}
}
