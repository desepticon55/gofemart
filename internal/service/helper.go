package service

import (
	"hash/fnv"
	"strconv"
)

func IsValidOrderNumber(orderNumber string) bool {
	sum := 0
	needDouble := false

	if orderNumber == "" {
		return false
	}

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}
		if needDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		needDouble = !needDouble
	}
	return sum%10 == 0
}

func HashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
