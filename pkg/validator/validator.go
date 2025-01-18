package validator

import (
	"fmt"
	"math"
	"receipt-processor/internal/models"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ValidationError struct {
	ErrorMessage string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s", e.ErrorMessage)
}

type Validator struct {
	errors []ValidationError
}

func New() *Validator {
	return &Validator{errors: []ValidationError{}}
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) Errors() []ValidationError {
	return v.errors
}

func (v *Validator) addError(message string) {
	v.errors = append(v.errors, ValidationError{ErrorMessage: message})
}

// ValidateReceipt validates all receipt fields
func (v *Validator) ValidateReceipt(receipt *models.Receipt) bool {

	v.validateRetailer(receipt.Retailer)

	v.validatePurchaseDate(receipt.PurchaseDate)

	v.validatePurchaseTime(receipt.PurchaseTime, receipt.PurchaseDate)

	v.validateItems(receipt.Items)

	v.validateTotal(receipt.Total)

	v.validateTotalMatchesItems(receipt.Total, receipt.Items)

	return !v.HasErrors()
}

// Retailer validation
func (v *Validator) validateRetailer(retailer string) {
	if retailer == "" {
		v.addError("Retailer cannot be empty")
		return
	}

	// Check for valid characters: alphanumeric, spaces, hyphen and ampersand
	validRetailer := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !validRetailer.MatchString(retailer) {
		v.addError("Retailer contains invalid characters")
	}

	// At least one alphanumeric character
	hasAlphanumeric := false
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			hasAlphanumeric = true
			break
		}
	}
	if !hasAlphanumeric {
		v.addError("Retailer name must contain at least one alphanumeric character")
	}
}

// Purchase date validation
func (v *Validator) validatePurchaseDate(date string) {
	// Format check
	purchaseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		v.addError("Invalid purchase date.")
		return
	}

	// Check if date is in future
	if purchaseDate.After(time.Now()) {
		v.addError("Purchase date cannot be in the future")
		return
	}

	// Check for valid day in month (including leap year handling)
	year := purchaseDate.Year()
	month := purchaseDate.Month()
	day := purchaseDate.Day()

	// Get last day of month
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	if day > lastDay {
		v.addError(fmt.Sprintf("Invalid day for month (max %d)", lastDay))
	}
}

// Purchase time validation
func (v *Validator) validatePurchaseTime(timeStr string, date string) {
	// Check format
	purchaseTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		v.addError("Invalid time format (should be HH:MM in 24-hour format)")
		return
	}

	// Get current time
	now := time.Now()
	purchaseDate, err := time.Parse("2006-01-02", date)
	// If purchase date is today, ensure time is not in the future
	if !v.HasErrors() {
		purchaseDateTime := time.Date(
			purchaseDate.Year(), purchaseDate.Month(), purchaseDate.Day(),
			purchaseTime.Hour(), purchaseTime.Minute(), 0, 0, now.Location(),
		)

		if purchaseDateTime.After(now) {
			v.addError("Purchase time cannot be in the future for today's date")
		}
	}

	// Additional validations for hours and minutes
	hour := purchaseTime.Hour()
	minute := purchaseTime.Minute()

	if hour < 0 || hour > 23 {
		v.addError("Purchase time hours must be between 0 and 23")
	}

	if minute < 0 || minute > 59 {
		v.addError("Purchase time minutes must be between 0 and 59")
	}
}

// Items validations
func (v *Validator) validateItems(items []models.Item) {
	if len(items) == 0 {
		v.addError("Receipt must include at least one item")
		return
	}

	for i, item := range items {
		// Short description validation
		if strings.TrimSpace(item.ShortDescription) == "" {
			v.addError("Item short description cannot be empty")
		}

		validDesc := regexp.MustCompile(`^[\w\s\-]+$`)
		if !validDesc.MatchString(item.ShortDescription) {
			v.addError("Item short description contains invalid characters")
		}

		// Item Price Validation
		v.validatePrice(item.Price, fmt.Sprintf("items[%d].price", i))
	}
}

// Price validation
func (v *Validator) validatePrice(price, field string) {
	// Check format
	validPrice := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !validPrice.MatchString(price) {
		v.addError("Invalid price format for item (should be X.XX)")
		return
	}

	// Check if valid number
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		v.addError("Invalid price value for item")
		return
	}

	// Check if price  positive
	if priceFloat <= 0 {
		v.addError("Item price must be greater than 0")
	}
}

// Total validation
func (v *Validator) validateTotal(total string) {
	v.validatePrice(total, "total")
}

// Validation for total matches sum of items
func (v *Validator) validateTotalMatchesItems(total string, items []models.Item) {
	totalFloat, _ := strconv.ParseFloat(total, 64)

	var sum float64
	for _, item := range items {
		itemPrice, _ := strconv.ParseFloat(item.Price, 64)
		sum += itemPrice
	}

	// Round to 2 decimal places to avoiding float decimal issues
	sum = math.Round(sum*100) / 100
	totalFloat = math.Round(totalFloat*100) / 100

	if sum != totalFloat {
		v.addError(fmt.Sprintf("Total (%.2f) does not match sum of items (%.2f)", totalFloat, sum))
	}
}

// Receipt ID validation
func (v *Validator) ValidateReceiptID(id string) bool {
	validID := regexp.MustCompile(`^\S+$`)
	if !validID.MatchString(id) {
		v.addError("id must not contain whitespace or be empty")
		return false
	}

	return true
}
