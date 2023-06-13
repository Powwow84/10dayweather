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
	Day  string
	High string
	Low  string
}

func transformCSVToPrometheus(weatherData []WeatherData) (string, error) {
	output := ""

	for _, data := range weatherData {
		output += fmt.Sprintf("high_temp{day=\"%s\"} %s\n", data.Day, data.High)
		output += fmt.Sprintf("low_temp{day=\"%s\"} %s\n", data.Day, data.Low)
	}

	return output, nil
}

func main() {
	weatherData := make([]WeatherData, 0)
	c := colly.NewCollector(colly.AllowedDomains("weather.com", "www.weather.com"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML(".DetailsSummary--fadeOnOpen--KnNyF", func(e *colly.HTMLElement) {
		day := e.ChildText("h3.DetailsSummary--daypartName--kbngc")
		high := e.ChildText("span.DetailsSummary--highTempValue--3PjlX")
		low := e.ChildText("span.DetailsSummary--lowTempValue--2tesQ")
		weatherData = append(weatherData, WeatherData{Day: day, High: high, Low: low})
	})

	err := c.Visit("https://weather.com/weather/tenday/l/fe3c78e80c47c404a4e64ec7c86ceccdb814894cedefdb528f9d8d95c3e4eb74")
	if err != nil {
		log.Fatal(err)
	}

	// Convert weather data to Prometheus-compatible format
	prometheusOutput, err := transformCSVToPrometheus(weatherData)
	if err != nil {
		log.Fatal(err)
	}

	// Save weather data to a CSV file
	t := time.Now()
	csvFileName := fmt.Sprintf("weather_%s.csv", t.Format("20060102150405"))
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("Cannot create CSV file:", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	csvWriter.Write([]string{"Day", "HighTemp", "LowTemp"})

	// Write weather data to CSV
	for _, w := range weatherData {
		csvWriter.Write([]string{w.Day, w.High, w.Low})
	}

	fmt.Println("Weather data saved to", csvFileName)

	// Save Prometheus-compatible output to a text file
	prometheusFileName := fmt.Sprintf("weather_%s.txt", t.Format("20060102150405"))
	prometheusFile, err := os.Create(prometheusFileName)
	if err != nil {
		log.Fatal("Cannot create Prometheus file:", err)
	}
	defer prometheusFile.Close()

	_, err = prometheusFile.WriteString(prometheusOutput)
	if err != nil {
		log.Fatal("Cannot write to Prometheus file:", err)
	}

	fmt.Println("Prometheus-compatible output saved to", prometheusFileName)
}
