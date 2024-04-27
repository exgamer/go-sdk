package validation

import "regexp"

// CheckValidPhone - проверка номера телефона
func CheckValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\+7 \(\d{3}\)-\d{3}-\d{2}-\d{2}$`)

	return re.MatchString(phone)
}
