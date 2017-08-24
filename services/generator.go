package services

import (
	"time"
	"fmt"
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"github.com/korasdor/colarit/utils"
	"github.com/korasdor/colarit/model"
)

var (
	TEMPLATE_CHARS  string = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
	TEMPLATE_DIGITS string = "987654321123456789"
	SERIALS_PATH    string = utils.ROOT_PATH + "serials/"
)

func CreateSerialKeys(serialsName string, serialCount int) bool {
	result := false

	rand.Seed(time.Now().UTC().UnixNano())

	var serials []string
	for i := 0; i < serialCount; i++ {
		serial := GenerateSerial(model.SerialKeyLength)

		serials = append(serials, serial)
	}

	csvContent := strings.Join(serials, "\n")
	sJson, _ := json.Marshal(serials)

	if _, err := os.Stat(SERIALS_PATH); os.IsNotExist(err) {
		os.Mkdir("serials", 0644)
	}

	if _, err := os.Stat(SERIALS_PATH + serialsName); os.IsNotExist(err) {
		err := ioutil.WriteFile(SERIALS_PATH+serialsName, sJson, 0644)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(SERIALS_PATH+serialsName+".csv", []byte(csvContent), 0644)
		if err != nil {
			fmt.Println(err)
		}

		result = true
	}

	return result
}

func CreateTable(tableName string, serialsName string) bool {
	result := true

	res, err := ioutil.ReadFile(SERIALS_PATH + serialsName)
	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	var serials []string
	if err := json.Unmarshal(res, &serials); err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	if CreateSerialsTable(tableName) {
		result = FillSerialTable(tableName, serials, "NULL")
	} else {
		result = false
	}

	return result
}

func GetSerials(filename string) ([]string, error) {
	res, err := ioutil.ReadFile(SERIALS_PATH + filename)
	if err != nil {
		utils.PrintOutput(err.Error())
	}

	var serials []string
	if err := json.Unmarshal(res, &serials); err != nil {
		utils.PrintOutput(err.Error())
	}

	return serials, err
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
