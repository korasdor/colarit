package services

import (
	"time"
	"math/rand"
	"github.com/korasdor/colarit/model"
	"os"
	"bufio"
	"github.com/korasdor/colarit/utils"
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

func GetSerialFromFile(fileName string, ) ([]string, error) {
	var serials []string

	file, err := os.Open("serials/" + fileName)
	defer file.Close()

	if err != nil {
		utils.PrintOutput(err.Error())
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			serials = append(serials, scanner.Text())
		}
	}

	return serials, err
}
