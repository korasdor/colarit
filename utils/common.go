package utils

import (
	"fmt"
	"os"
	"github.com/oschwald/geoip2-golang"
	"net"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

var (
	ROOT_PATH string = ""
	GEO_LITE_FILE_PATH  string = ROOT_PATH + "static/data/GeoLite2-Country.mmdb"
	TEMPLATE_LOCAL_PATH string = ROOT_PATH + "templates/books.json"
	TEMPLATE_REMOTE_URL string = "http://colarit.com/colar/templates/books.json"
)

func GetDBAddress() string {
	var dbAddress string

	// korasdor:19841986aA@tcp(127.0.0.1:3306)/im
	if os.Getenv("OPENSHIFT_MYSQL_DB_USERNAME") == "" {
		dbAddress = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "korasdor", "19841986aA", "127.0.0.1", "3306", "colar_db")
	} else {
		dbAddress = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("OPENSHIFT_MYSQL_DB_USERNAME"),
			os.Getenv("OPENSHIFT_MYSQL_DB_PASSWORD"),
			os.Getenv("OPENSHIFT_MYSQL_DB_HOST"),
			os.Getenv("OPENSHIFT_MYSQL_DB_PORT"),
			"ar")
	}

	return dbAddress
}

func GetTempDir() string {
	return os.Getenv("OPENSHIFT_TMP_DIR")
}


func GetCountry(clientIp string) string {
	country := "Us"

	fmt.Println(GEO_LITE_FILE_PATH)

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

func UpdateBooksTemplate() bool {
	var result = true
	response, err := http.Get(TEMPLATE_REMOTE_URL)
	if err != nil {
		result = false
		fmt.Println(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			result = false
			fmt.Println(err)
		}

		err = ioutil.WriteFile(TEMPLATE_LOCAL_PATH, contents, 0644)
		if err != nil {
			result = false
			fmt.Println(err)
		}
	}

	return result
}

func GetBooksJson(country string) []byte {
	var bookJsonStr []byte

	b, err := ioutil.ReadFile(TEMPLATE_LOCAL_PATH)
	if err != nil {
		PrintOutput(err.Error())
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(b, &dat); err != nil {
		PrintOutput(err.Error())
	} else {

		if country == "UZ" {
			dat["assets_server"] = "http://colarit.com/colar"
			dat["supported_langs"] = []string{"en", "ru", "uz"}
			// dat["assets_server"] = "http://colar.uz"
		} else {
			dat["assets_server"] = "http://colarit.com/colar"
			dat["supported_langs"] = []string{"en", "ru"}
		}

		bookJsonStr, err = json.Marshal(dat)
		if err != nil {
			PrintOutput(err.Error())
		}
	}

	return bookJsonStr

}

func PrintOutput(str string) {
	fmt.Println(str)
}
