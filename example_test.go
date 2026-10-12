package LuhnCalculator_test

import (
	"fmt"

	luhn "github.com/Herts/LuhnCalculator"
)

func ExampleAppendCheckDigit() {
	number, err := luhn.AppendCheckDigit("0007992739871")
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
	fmt.Println(luhn.Validate(number))
	// Output:
	// 00079927398713
	// true
}

func ExampleCheckDigit() {
	digit, err := luhn.CheckDigit("7992739871")
	if err != nil {
		panic(err)
	}
	fmt.Println(digit)
	// Output: 3
}
