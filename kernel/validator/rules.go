package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// mobile 验证手机号码
func mobile(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^1\d{10}$`, fl.Field().String())
	return ok
}

// idCard 验证身份证号码
func idCard(fl validator.FieldLevel) bool {
	id := fl.Field().String()

	var a1Map = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}

	var idStr = strings.ToUpper(id)
	var reg, err = regexp.Compile(`^\d{17}[\dX]$`)
	if err != nil {
		return false
	}
	if !reg.Match([]byte(idStr)) {
		return false
	}
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, err := strconv.Atoi(string(c)); err == nil {
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = a1Map[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}
	return a1Str == signChar
}

func dir(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^(/?[\da-zA-Z]+){1,3}$`, fl.Field().String())
	return ok
}

func username(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z\d\-_]{4,20}$`, fl.Field().String())
	return ok
}

func password(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z\d\-_@$&%!]{6,32}$`, fl.Field().String())
	return ok
}

func snowflake(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^\d{16,64}$`, fl.Field().String())
	return ok
}
