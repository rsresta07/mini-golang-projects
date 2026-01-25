package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CreditCardInfo contains basic validation metadata
type CreditCardInfo struct {
	Valid        bool
	CardType     string
	ErrorMessage string
	Length       int
	PassedLuhn   bool
}

// ValidateCreditCard performs full validation including Luhn check
func ValidateCreditCard(number string) CreditCardInfo {
	// Clean input
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	number = strings.TrimSpace(number)

	info := CreditCardInfo{
		Length: len(number),
	}

	// Length check (general PAN rule)
	if info.Length < 13 || info.Length > 19 {
		info.ErrorMessage = "Card number must be between 13–19 digits"
		return info
	}

	// Digit-only check (no int overflow risk)
	if !isNumeric(number) {
		info.ErrorMessage = "Card number must contain only digits"
		return info
	}

	// Detect the card type (most common ranges)
	info.CardType = detectCardType(number)

	// Network-specific length rules
	switch info.CardType {
	case "American Express":
		if info.Length != 15 {
			info.ErrorMessage = "American Express cards must be 15 digits"
			return info
		}

	case "Visa":
		if info.Length != 13 && info.Length != 16 && info.Length != 19 {
			info.ErrorMessage = "Visa cards must be 13, 16, or 19 digits"
			return info
		}

	case "Mastercard":
		if info.Length != 16 {
			info.ErrorMessage = "Mastercard cards must be 16 digits"
			return info
		}

	case "Discover":
		if info.Length != 16 && info.Length != 19 {
			info.ErrorMessage = "Discover cards must be 16 or 19 digits"
			return info
		}

	case "JCB":
		if info.Length < 16 || info.Length > 19 {
			info.ErrorMessage = "JCB cards must be 16–19 digits"
			return info
		}
	}

	// Luhn check
	info.PassedLuhn = luhnCheck(number)
	info.Valid = info.PassedLuhn

	if !info.PassedLuhn {
		info.ErrorMessage = "Failed Luhn check"
	}

	return info
}

// isNumeric checks if string contains only digits
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// detectCardType identifies card network using BIN/IIN rules
func detectCardType(number string) string {
	// Visa
	if number[0] == '4' {
		return "Visa"
	}

	if len(number) >= 2 {
		prefix2, _ := strconv.Atoi(number[:2])

		// American Express
		if prefix2 == 34 || prefix2 == 37 {
			return "American Express"
		}

		// Mastercard old range
		if prefix2 >= 51 && prefix2 <= 55 {
			return "Mastercard"
		}

		// Discover
		if prefix2 == 65 {
			return "Discover"
		}

		// JCB
		if prefix2 >= 35 && prefix2 <= 39 {
			return "JCB"
		}
	}

	// Mastercard new range (222100–272099)
	if len(number) >= 6 {
		prefix6, _ := strconv.Atoi(number[:6])
		if prefix6 >= 222100 && prefix6 <= 272099 {
			return "Mastercard"
		}
	}

	// Discover additional ranges
	if len(number) >= 4 && number[:4] == "6011" {
		return "Discover"
	}
	if len(number) >= 3 {
		prefix3, _ := strconv.Atoi(number[:3])
		if prefix3 >= 644 && prefix3 <= 649 {
			return "Discover"
		}
	}

	return "Unknown"
}

// luhnCheck implements the Luhn (mod 10) algorithm
func luhnCheck(number string) bool {
	sum := 0
	double := false

	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0 // Valid if the total sum is divisible by 10
}

func main() {
	fmt.Println("Credit Card Validator (Luhn + Network Rules)")
	fmt.Println("Enter card numbers (or 'q' to quit)\n")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("=> ")
		scanner.Scan()
		input := scanner.Text()

		// Exit condition
		if strings.ToLower(strings.TrimSpace(input)) == "q" {
			break
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		result := ValidateCreditCard(input)

		// Display results
		fmt.Printf("\nNumber:  %s\n", formatCardNumber(input))
		fmt.Printf("Type:    %s\n", result.CardType)
		fmt.Printf("Length:  %d\n", result.Length)
		fmt.Printf("Luhn:    %v\n", result.PassedLuhn)

		if result.Valid {
			fmt.Println("Status:  VALID")
		} else {
			fmt.Println("Status:  INVALID")
			if result.ErrorMessage != "" {
				fmt.Printf("Reason:  %s\n", result.ErrorMessage)
			}
		}
		fmt.Println()
	}
}

// formatCardNumber formats card number with spaces every 4 digits
func formatCardNumber(number string) string {
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")

	var b strings.Builder
	for i, r := range number {
		if i > 0 && i%4 == 0 {
			b.WriteRune(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}
