package calculator

import (
	"math"
	"receipt-processor/internal/models"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func CalculatePoints(receipt *models.Receipt) int64 {
	var points int64

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)
	// ____Alternate slow performance method with Regex
	// alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	// points += int64(len(alphanumeric.FindAllString(receipt.Retailer, -1)))
	// _______________________________

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if math.Mod(total*100, 25) == 0 {
			points += 25
		}
	}

	// Rule 4: 5 points for every two items on the receipt
	points += int64(len(receipt.Items) / 2 * 5)

	// Rule 5: Description length multiple of 3
	for _, item := range receipt.Items {
		trimDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimDesc)%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int64(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: Not Applicable

	// Rule 7: 6 points if the day in the purchase date is odd
	if purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if purchaseDate.Day()%2 == 1 {
			points += 6
		}
	}

	// Rule 8: 10 points if time is between 2:00pm and 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		startTime, _ := time.Parse("15:04", "14:00") // 2:00 PM
		endTime, _ := time.Parse("15:04", "16:00")   // 4:00 PM

		if !purchaseTime.Before(startTime) && purchaseTime.Before(endTime) {
			points += 10
		}
	}

	return points
}

func countAlphanumeric(s string) int64 {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return int64(count)
}
