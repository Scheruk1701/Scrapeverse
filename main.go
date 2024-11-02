package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Stock struct {
	Symbol string
	Price  float64
	Change string
}

var stocks []Stock

func main() {
	// fmt.Println("Compile Crew")
	loadStockData()
}

func loadStockData() {

	//setup for csv file writing
	file, err := os.Create("stocks_data.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Symbol", "Price", "Change"}
	writer.Write(headers)

	c := colly.NewCollector()

	c.OnHTML("fin-streamer[data-field='regularMarketPrice']", func(e *colly.HTMLElement) {
		symbol := strings.TrimPrefix(e.Request.URL.Path, "/quote/")
		price, _ := strconv.ParseFloat(e.Text, 64)
		if price > 0 {
			stocks = append(stocks, Stock{Symbol: symbol, Price: price})
		}
	})

	c.OnHTML("fin-streamer[data-field='regularMarketChangePercent']", func(e *colly.HTMLElement) {
		symbol := strings.TrimPrefix(e.Request.URL.Path, "/quote/")
		change := e.Text
		for i, stock := range stocks {
			if stock.Symbol == symbol {
				stocks[i].Change = change
			}
		}
	})

	stockSymbols := []string{"MSFT", "AAPL", "GOOGL", "AMZN", "FB", "TSLA", "NFLX", "NVDA", "BABA", "V"}

	for _, symbol := range stockSymbols {
		c.Visit("https://finance.yahoo.com/quote/" + symbol)
	}

	for _, stock := range stocks {
		writer.Write([]string{stock.Symbol, fmt.Sprintf("%.2f", stock.Price), stock.Change})
	}
	fmt.Println("Stock data successfully written to CSV!")

}
