package check

import (
	"errors"
	"fmt"
	"time"
	"unicode"
)

//driver phone check

func PhoneNumber(phone string) bool {
	for _, p := range phone {
		if p == '+' {
			continue
		} else if !unicode.IsNumber(p) {
			return false
		}
	}
	return true
}

//car year check

func Year(year int) error {
	fmt.Println(year)
	if year <= 0 || year > time.Now().Year() {
		return errors.New("year is not correct for car!")
	}
	return nil
}