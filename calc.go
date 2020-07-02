package LuhnCalculator

func GenerateLastDigit(num int) int {
	num *= 10
	for i := 0; i < 10; i++ {
		if IsValidLuhn(num + i) {
			return num + i
		}
	}
	return 0
}

func IsValidLuhn(num int) bool {
	sum := num % 10
	num /= 10
	isSecond := true
	for num != 0 {
		digit := num % 10
		if isSecond {
			digit *= 2
			isSecond = false
		} else {
			isSecond = true
		}
		if digit > 9 {
			digit = digit%10 + digit/10
		}
		sum += digit
		num /= 10
	}
	if sum%10 == 0 {
		return true
	}
	return false
}
