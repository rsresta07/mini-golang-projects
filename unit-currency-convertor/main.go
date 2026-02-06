package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

/* to run it first,
go build -o convert
./convert -t curr -from USD -to JPY -v 50
./convert -t length -from m -to ft -v 1.85
*/

const API_URL = "https://open.er-api.com/v6/latest/" // Free, no-key-required API for this demo

type ExchangeResponse struct {
	Result   string             `json:"result"`
	Rates    map[string]float64 `json:"rates"`
	BaseCode string             `json:"base_code"`
}

var lengthUnits = map[string]float64{"m": 1, "cm": 100, "km": 0.001, "in": 39.37, "ft": 3.28, "mi": 0.000621}
var weightUnits = map[string]float64{"kg": 1, "g": 1000, "lb": 2.204, "oz": 35.27}

// main is the entry point for the application.
// It parses command line flags, calls the appropriate conversion
// function based on the mode flag, and prints the result to
// stdout.
func main() {
	cmd := flag.NewFlagSet("convert", flag.ExitOnError)
	mode := cmd.String("t", "length", "Conversion type: length, weight, temp, curr")
	from := cmd.String("from", "", "Source unit (e.g., USD, kg, m, c)")
	to := cmd.String("to", "", "Target unit (e.g., EUR, lb, cm, f)")
	val := cmd.Float64("v", 0.0, "The numeric value to convert")

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd.Parse(os.Args[1:])

	switch strings.ToLower(*mode) {
	case "curr", "currency":
		convertCurrency(*from, *to, *val)
	case "length":
		convertStandard(*val, *from, *to, lengthUnits)
	case "weight":
		convertStandard(*val, *from, *to, weightUnits)
	case "temp":
		convertTemp(*val, *from, *to)
	default:
		fmt.Printf("❌ Unknown mode: %s\n", *mode)
	}
}

// convertCurrency converts a given value from one currency to another.
// It fetches live rates from the Open Exchange Rates API and then
// performs the conversion. If the target currency is not supported,
// it prints an error message and returns.
func convertCurrency(from, to string, val float64) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	fmt.Printf("🌐 Fetching live rates for %s...\n", from)
	resp, err := http.Get(API_URL + from)
	if err != nil {
		fmt.Println("❌ Network error. Check your connection.")
		return
	}
	defer resp.Body.Close()

	var data ExchangeResponse
	json.NewDecoder(resp.Body).Decode(&data)

	if data.Result != "success" {
		fmt.Printf("❌ Could not find currency: %s\n", from)
		return
	}

	if rate, ok := data.Rates[to]; ok {
		fmt.Printf("✅ %0.2f %s = %0.2f %s (Rate: %0.4f)\n", val, from, val*rate, to, rate)
	} else {
		fmt.Printf("❌ Target currency %s not supported.\n", to)
	}
}

// convertStandard converts a value from one standard unit to another.
//
// Supported units for length are:
//
//   - 'km' for kilometers
//   - 'mi' for miles
//   - 'm' for meters
//   - 'cm' for centimeters
//   - 'mm' for millimeters
//   - 'in' for inches
//   - 'ft' for feet
//   - 'yd' for yards
//
// Supported units for weight are:
//
//   - 'kg' for kilograms
//   - 'lb' for pounds
//   - 'g' for grams
//   - 'oz' for ounces
//
// The function takes a float64 value, two strings for the from and to units,
// and a map of strings to float64 values for the conversion rates.
// It returns the converted value as a float64.
// If the units are not recognized, it prints an error message and returns nothing.
func convertStandard(val float64, from, to string, units map[string]float64) {
	fRate, ok1 := units[from]
	tRate, ok2 := units[to]

	if !ok1 || !ok2 {
		fmt.Printf("❌ Invalid units. Use: %v\n", getKeys(units))
		return
	}

	result := (val / fRate) * tRate
	fmt.Printf("✅ %0.2f %s = %0.4f %s\n", val, from, result, to)
}

// convertTemp converts a temperature from one unit to another.
//
// Supported units:
//
//   - 'c' for Celsius
//   - 'f' for Fahrenheit
//
// The function takes a float64 value, a string for the from unit and a string for the to unit.
// It returns the converted value as a float64.
// If the units are not recognized, it prints an error message and returns nothing.
func convertTemp(val float64, from, to string) {
	from, to = strings.ToLower(from), strings.ToLower(to)
	var res float64
	if from == "c" && to == "f" {
		res = (val * 9 / 5) + 32
	} else if from == "f" && to == "c" {
		res = (val - 32) * 5 / 9
	} else {
		fmt.Println("❌ Use 'c' for Celsius or 'f' for Fahrenheit")
		return
	}
	fmt.Printf("✅ %0.2f°%s = %0.2f°%s\n", val, strings.ToUpper(from), res, strings.ToUpper(to))
}

// getKeys returns a slice of strings containing the keys of the given map.
// It is used to print the supported units for a given conversion type.
// The returned slice is not sorted.
func getKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// printUsage prints the usage of the program to the console.
func printUsage() {
	fmt.Println("Usage: converter -t [type] -from [unit] -to [unit] -v [value]")
	fmt.Println("Types: length, weight, temp, curr")
	fmt.Println("\nExamples:")
	fmt.Println("  ./converter -t curr -from USD -to EUR -v 100")
	fmt.Println("  ./converter -t length -from km -to mi -v 5")
}
