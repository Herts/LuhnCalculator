// Package LuhnCalculator validates decimal strings and generates Luhn check digits.
package LuhnCalculator

import (
	"errors"
	"strconv"
)

// ErrInvalidDigits indicates an empty input or a character other than ASCII 0–9.
var ErrInvalidDigits = errors.New("luhn: input must contain one or more ASCII digits")

// Validate reports whether number consists of ASCII digits with a valid Luhn
// checksum. Leading zeros are preserved; spaces and separators are rejected.
// This validates the checksum only, not an identifier's length or issuer.
func Validate(number string) bool {
	sum, err := checksum(number, false)
	return err == nil && sum == 0
}

// CheckDigit returns the digit (0–9) to append to payload to make its checksum
// valid. Payload must contain one or more ASCII digits.
func CheckDigit(payload string) (int, error) {
	sum, err := checksum(payload, true)
	if err != nil {
		return 0, err
	}
	return (10 - sum) % 10, nil
}

// AppendCheckDigit returns payload followed by its Luhn check digit.
func AppendCheckDigit(payload string) (string, error) {
	digit, err := CheckDigit(payload)
	if err != nil {
		return "", err
	}
	return payload + string(byte('0'+digit)), nil
}

func checksum(digits string, double bool) (int, error) {
	if digits == "" {
		return 0, ErrInvalidDigits
	}
	sum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < '0' || digits[i] > '9' {
			return 0, ErrInvalidDigits
		}
		digit := int(digits[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum = (sum + digit) % 10
		double = !double
	}
	return sum, nil
}

// GenerateLastDigit returns the complete integer with a check digit appended,
// preserving the historical behavior despite the function's name. It returns
// zero for negative inputs or when the result would overflow int.
//
// Deprecated: Use AppendCheckDigit for the complete string, or CheckDigit for
// just the digit. The string APIs preserve leading zeros and report errors.
func GenerateLastDigit(num int) int {
	if num < 0 {
		return 0
	}
	complete, _ := AppendCheckDigit(strconv.Itoa(num))
	result, err := strconv.Atoi(complete)
	if err != nil {
		return 0
	}
	return result
}

// IsValidLuhn reports whether a nonnegative integer has a valid Luhn checksum.
//
// Deprecated: Use Validate to preserve leading zeros and accept long numbers.
func IsValidLuhn(num int) bool {
	return num >= 0 && Validate(strconv.Itoa(num))
}
