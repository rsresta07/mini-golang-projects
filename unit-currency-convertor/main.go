package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ExchangeRates struct {
	Base    string             `json:"base"`
	Date    string             `json:"date"`
	Rates   map[string]float64 `json:"rates"`
	Success bool               `json:"success"`
}

var rates ExchangeRates
var lastFetch time.Time

func main() {
	fmt.Println("Mini Unit and Currency Converter")

	fetchCurrencyRates()

	for {
		showMenu()
		choice := readChoice()
		switch choice {
		case 1:
			unitConverter()
		case 2:
			currencyConverter()
		case 3:
			fetchCurrencyRates()
			fmt.Println("Exchange rates updated!")
		case 4:
			fmt.Println("\nBye!")
			os.Exit(4)
		default:
			fmt.Println("Invalid choice!")
		}
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("1. Unit Converter")
	fmt.Println("2. Currency Converter")
	fmt.Println("3. Update Exchange Rates")
	fmt.Println("4. Exit")
	fmt.Print("Enter your choice: ")
}

func readChoice() int {
	var input string
	fmt.Scanln(&input)
	choice, _ := strconv.Atoi(input)
	return choice
}

// Unit Converter
func unitConverter() {
	conversions := map[string]map[string]float64{
		"length": {
			"m-cm": 100, "cm-m": 0.01,
			"m-km": 0.001, "km-m": 1000,
			"m-inch": 39.3701, "inch-m": 0.0254,
			"m-ft": 3.28084, "ft-m": 0.3048,
		},
		"weight": {
			"kg-g": 1000, "g-kg": 0.001,
			"lb-kg": 0.453592, "kg-lb": 2.20462,
			"oz-g": 28.3495, "g-oz": 0.035274,
			"kg-oz": 35.274, "oz-kg": 0.002835,
		},
		"temperature": {},
		"area": {
			"m^2-ft^2": 10.7639, "ft^2-m^2": 0.092903,
			"m²→acre": 0.000247105, "acre→m²": 4046.86,
		},
	}

	fmt.Println("\nUnit Converter: length, weight, temperature, area")
	fmt.Print("Choose a unit: ")
	category:= strings.
}
