package services

import (
	"time"
	"math/rand"
	"github.com/korasdor/colarit/model"
)

var (
	TEMPLATE_CHARS  string = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
	TEMPLATE_DIGITS string = "987654321123456789"
)

func CreateSerials(serialCount int) []string {
	rand.Seed(time.Now().UTC().UnixNano())

	var serials []string
	for i := 0; i < serialCount; i++ {
		serial := GenerateSerial(model.SerialKeyLength)

		serials = append(serials, serial)
	}

	return serials
}

func GenerateSerial(size int) string {
	result := string(TEMPLATE_CHARS[rand.Intn(len(TEMPLATE_CHARS))])
	for i := 0; i < size-1; i++ {
		pos := rand.Intn(len(TEMPLATE_DIGITS))
		char := TEMPLATE_DIGITS[pos]

		result += string(char)
	}

	return result
}
