package pkg

import (
	"regexp"
	"strings"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	// Удаление всех символов, кроме цифр
	re := regexp.MustCompile("[^0-9]")
	normalizedPhone = re.ReplaceAllString(phone, "")

	// Проверка наличия кода страны
	if len(normalizedPhone) >= 11 && strings.HasPrefix(normalizedPhone, "8") {
		normalizedPhone = "7" + normalizedPhone[1:]
	}

	return normalizedPhone, nil
}
