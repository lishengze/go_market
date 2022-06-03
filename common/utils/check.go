package utils

import "regexp"

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyMobileFormat(mobile string) bool {
	pattern := `^(0|[1-9]\d*)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}
