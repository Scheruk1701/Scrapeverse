package main

import (
	"fmt"
)

type Stock struct {
    Symbol string
    Price  float64
    Change string
}

var stocks []Stock

func main() {
	// fmt.Println("Compile Crew")

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


}