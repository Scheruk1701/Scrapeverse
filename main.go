package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Stock struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Change string  `json:"change"`
}

var (
	stocks          []Stock
	hasScrapedStock = false // Flag to track scraping status
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Static("/static", "./static")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.GET("/analyze", func(c *gin.Context) {
		c.HTML(http.StatusOK, "analyze.html", nil)
	})

	// API Routes
	r.GET("/api/stocks/scrape", scrapeStockData)
	r.GET("/api/stocks/reset-scraping", resetScrapeStatus)
	r.GET("/api/stocks", getStockData)
	r.GET("/api/stocks/search", searchStock)
	r.GET("/api/stocks/sort", sortStocks)
	r.GET("/api/stocks/top-gainer-loser", getTopGainerLoser)
	r.GET("/api/stocks/download", downloadCSV)

	r.Run(":8080")
}

func scrapeStockData(c *gin.Context) {
	// Check if stocks have already been scraped
	fmt.Println("Scrape request received. HasScrapedStock:", hasScrapedStock)
	if hasScrapedStock {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Stock data has already been scraped. Please refresh the data before re-scraping.",
		})
		return
	}

	// Clear existing stocks before scraping
	stocks = []Stock{}

	// Call the loadStockData function
	err := loadStockData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to scrape stock data: " + err.Error(),
		})
		return
	}

	// Set the flag to true after successful scraping
	hasScrapedStock = true

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Stock data scraped successfully!",
	})
}

func resetScrapeStatus(c *gin.Context) {
	// Reset the scraping status and clear current stock data
	hasScrapedStock = false
	stocks = []Stock{}
	c.JSON(http.StatusOK, gin.H{
		"message": "Scraping status reset successfully. You can now scrape again.",
	})
}

func loadStockData() error {
	stockSymbols := []string{
		"MSFT", "AAPL", "GOOGL", "AMZN", "TSLA", "NVDA", "ORCL", "AMD",
		"SNOW", "CRWD", "MSTR", "INOD", "APLD", "ADBE", "AVGO",
		"TXN", "QCOM", "V", "MA", "PYPL",
		"CAT", "GE", "HON", "MMM",
		"BABA", "NFLX", "AMT", "MDLZ", "MRK", "NKE",
	}

	file, err := os.Create("stocks_data.csv")
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Symbol", "Price", "Change"}
	writer.Write(headers)

	c := colly.NewCollector()

	uniqueStocks := make(map[string]Stock)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("fin-streamer[data-field='regularMarketPrice']", func(e *colly.HTMLElement) {
		symbol := strings.TrimPrefix(e.Request.URL.Path, "/quote/")
		symbol = strings.ReplaceAll(symbol, "/", "")
		price, _ := strconv.ParseFloat(e.Attr("data-value"), 64)

		if _, exists := uniqueStocks[symbol]; !exists && price > 0 {
			uniqueStocks[symbol] = Stock{Symbol: symbol, Price: price}
		}
	})

	c.OnHTML("fin-streamer[data-field='regularMarketChangePercent']", func(e *colly.HTMLElement) {
		symbol := strings.TrimPrefix(e.Request.URL.Path, "/quote/")
		symbol = strings.ReplaceAll(symbol, "/", "")
		change := e.Text

		if stock, exists := uniqueStocks[symbol]; exists && stock.Change == "" {
			stock.Change = change
			uniqueStocks[symbol] = stock
		}
	})

	for _, symbol := range stockSymbols {
		c.Visit("https://finance.yahoo.com/quote/" + symbol)
	}

	for _, stock := range uniqueStocks {
		stocks = append(stocks, stock)
		writer.Write([]string{stock.Symbol, fmt.Sprintf("%.2f", stock.Price), stock.Change})
	}

	fmt.Println("Stock data successfully written to CSV!")
	return nil
}

func getStockData(c *gin.Context) {
	c.JSON(http.StatusOK, stocks)
}

func searchStock(c *gin.Context) {
	query := c.Query("symbol")
	query = strings.ToUpper(query)

	for _, stock := range stocks {
		if stock.Symbol == query {
			c.JSON(http.StatusOK, stock)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Stock not found"})
}

func sortStocks(c *gin.Context) {
	sortOrder := c.DefaultQuery("order", "asc")

	if sortOrder == "asc" {
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].Price < stocks[j].Price
		})
	} else {
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].Price > stocks[j].Price
		})
	}

	c.JSON(http.StatusOK, stocks)
}

func getTopGainerLoser(c *gin.Context) {
	if len(stocks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No stock data available"})
		return
	}

	var topGainer, topLoser Stock
	topGainerChange := -100.0
	topLoserChange := 100.0

	for _, stock := range stocks {
		changeStr := strings.Trim(stock.Change, " ()%")
		changePercentage, err := strconv.ParseFloat(changeStr, 64)
		if err != nil {
			continue
		}

		if changePercentage > topGainerChange {
			topGainerChange = changePercentage
			topGainer = stock
		}

		if changePercentage < topLoserChange {
			topLoserChange = changePercentage
			topLoser = stock
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"top_gainer": topGainer,
		"top_loser":  topLoser,
	})
}

func downloadCSV(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=stocks_data.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"Symbol", "Price", "Change"})

	for _, stock := range stocks {
		writer.Write([]string{stock.Symbol, fmt.Sprintf("%.2f", stock.Price), stock.Change})
	}
}
