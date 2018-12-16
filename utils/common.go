package utils

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"io/ioutil"
	"net/http"
	"bytes"
	"errors"
)

var (
	GEO_LITE_FILE_PATH string = "static/data/GeoLite2-Country.mmdb"
	TEMPLATES_URL      string = "http://colarit.com/pearkid/templates/"
)

func FormatSerials(serials []string, serialFormat string) (string, error) {
	var resultString string
	var err error

	if serialFormat == "csv" {
		var buffer bytes.Buffer
		for i := 0; i < len(serials); i++ {
			buffer.WriteString(serials[i] + "\n")
		}

		resultString = buffer.String()
	} else {
		err = errors.New("Unsupported format")
	}

	return resultString, err
}

func GetDealerName(dealerId int) string {
	dealerNamesMap := map[int]string{0: "test", 1: "uz", 2: "tj", 3: "ru"}

	return dealerNamesMap[dealerId]
}

func GetCountry(clientIp string) string {
	country := "Us"

	db, err := geoip2.Open(GEO_LITE_FILE_PATH)
	if err != nil {
		PrintOutput(err.Error())
	}
	defer db.Close()

	ip := net.ParseIP(clientIp)
	if ip != nil {
		if db != nil {
			record, err := db.City(ip)
			if err != nil {
				PrintOutput(err.Error())
			} else {
				country = record.Country.IsoCode
			}
		}
	}

	return country
}

func UpdateAllBooksTemplates() bool {
	var result = true

	result = UpdateBooksTemplate()
	if result == false {
		return result
	}

	langs := [4]string{"ru", "en", "uz", "kz"}
	for i := 0; i < len(langs); i++ {
		result = UpdateBooksLangTemplate(langs[i])

		if result == false {
			break
		}
	}

	return result
}

func UpdateBooksTemplate() bool {
	var result = true
	response, err := http.Get(TEMPLATES_URL + "books.json")
	if err != nil {
		result = false
		PrintOutput(err.Error())
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			result = false
			PrintOutput(err.Error())
		}

		err = ioutil.WriteFile("templates/books.json", contents, 0644)
		if err != nil {
			result = false
			PrintOutput(err.Error())
		}
	}

	return result
}

func UpdateBooksLangTemplate(lang string) bool {
	var result = true
	response, err := http.Get(TEMPLATES_URL + lang + "/books.json")
	if err != nil {
		result = false
		PrintOutput(err.Error())
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			result = false
			PrintOutput(err.Error())
		}

		err = ioutil.WriteFile("templates/"+lang+"/books.json", contents, 0644)
		if err != nil {
			result = false
			PrintOutput(err.Error())
		}
	}

	return result
}

func GetBooksJson(lang string) []byte {
	//var bookJsonStr []byte

	var filepath string
	if lang == "" {
		filepath = "templates/books.json"
	} else {
		filepath = "templates/" + lang + "/books.json"
	}

	booksBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		PrintOutput(err.Error())
	}
	//var dat map[string]interface{}
	//if err := json.Unmarshal(booksBytes, &dat); err != nil {
	//	PrintOutput(err.Error())
	//} else {
	//
	//	if lang == "UZ" {
	//		dat["assets_server"] = "http://colarit.com/colar"
	//		dat["supported_langs"] = []string{"en", "ru", "uz"}
	//		// dat["assets_server"] = "http://colar.uz"
	//	} else {
	//		dat["assets_server"] = "http://colarit.com/colar"
	//		dat["supported_langs"] = []string{"en", "ru"}
	//	}
	//
	//	bookJsonStr, err = json.Marshal(dat)
	//	if err != nil {
	//		PrintOutput(err.Error())
	//	}
	//}

	return booksBytes
}

func PrintOutput(str string) {
	fmt.Println(str)
}
