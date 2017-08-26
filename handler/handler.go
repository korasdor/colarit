package handler

import (
	"net/http"
	"fmt"
	"gopkg.in/gomail.v2"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/korasdor/colarit/model"
	"github.com/korasdor/colarit/utils"
	"github.com/korasdor/colarit/services"
)

/**
 * вернуть статичный класс
 */
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

/**
 * индексная страница
 */
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "index page in construction...")
}

/**
 * состояние базы данных
 */
func DBStateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "db state is: %s", model.DbSuccess);
}

/**
 * получить файл настроек книг.
 */
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	clientIp := r.Header.Get("X-Forwarded-For")
	country := utils.GetCountry(clientIp)

	booksJsonStr := utils.GetBooksJson(country)

	fmt.Fprint(w, string(booksJsonStr))
}

func UpdateBooksHandler(w http.ResponseWriter, r *http.Request) {
	result := utils.UpdateBooksTemplate()

	if result == true {
		fmt.Fprint(w, "good")
	} else {
		fmt.Fprint(w, "bad")
	}
}

/**
 * создаем таблицу
 */
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	var output string

	vars := mux.Vars(r)
	tableName := vars["table_name"]
	accessToken := vars["access_token"]

	if accessToken == model.ACCESS_TOKEN {
		if services.CreateSerialsTable(tableName) {
			output = fmt.Sprintf("{ \"message\":\"Table %s is successfully created\"}", tableName)
		} else {
			output = fmt.Sprintf("{ \"error\": \"2\", \"message\":\"Error creating table with name %s\"}", tableName)
		}
	} else {
		output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Incorrect access token. Token %s incorrect\"}", accessToken)
	}

	fmt.Fprint(w, output)
}

/**
 * заполняем таблицу
 */
func FillTableHandler(w http.ResponseWriter, r *http.Request) {
	var output string

	vars := mux.Vars(r)
	accessToken := vars["access_token"]

	if accessToken == model.ACCESS_TOKEN {
		rangeId := vars["range_id"]
		tableName := vars["table_name"]
		dealerId := vars["dealer_id"]

		serialsCount, err := strconv.Atoi(vars["serials_count"])
		if err != nil {
			fmt.Fprintf(w, "%s", "table fill complete")
		} else {
			serials := services.CreateSerials(serialsCount)
			if services.FillSerialsTable(tableName, serials, dealerId, rangeId) {
				output = fmt.Sprintf("{ \"message\":\"Table %s is filled successfully\"}", tableName)
			} else {
				output = fmt.Sprintf("{ \"error\": \"2\", \"message\":\"Table %s filling error\"}", tableName)
			}
		}
	} else {
		output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Incorrect access token. Token %s incorrect\"}", accessToken)
	}

	fmt.Fprint(w, output)
}

func FillTableWithFileHandler(w http.ResponseWriter, r *http.Request) {
	var output string

	vars := mux.Vars(r)
	accessToken := vars["access_token"]

	if accessToken == model.ACCESS_TOKEN {
		fileName := vars["file_name"]
		rangeId := vars["range_id"]
		tableName := vars["table_name"]
		dealerId := vars["dealer_id"]

		serials, err := services.GetSerialFromFile(fileName)
		if err != nil {
			output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Error reading the file %s\"}", fileName)
		} else {
			if services.FillSerialsTable(tableName, serials, dealerId, rangeId) {
				output = fmt.Sprintf("{ \"message\":\"Table %s is filled successfully\"}", tableName)
			} else {
				output = fmt.Sprintf("{ \"error\": \"2\", \"message\":\"Table %s filling error\"}", tableName)
			}
		}
	} else {
		output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Incorrect access token. Token %s incorrect\"}", accessToken)
	}

	fmt.Fprint(w, output)
}

/**
 * получить файл с серийниками
 */
func GetSerialsHandler(w http.ResponseWriter, r *http.Request) {
	var output string

	vars := mux.Vars(r)
	accessToken := vars["access_token"]
	rangeId := vars["range_id"]

	tableName := vars["table_name"]
	serialFormat := vars["serials_format"]

	if accessToken == model.ACCESS_TOKEN {
		serials, err := services.GetSerialsRange(tableName, rangeId)

		if err != nil {
			output = fmt.Sprintf("{ \"error\": \"2\", \"message\":\"An error occurred when obtaining the serial numbers in the table %s and rage %s\"}", tableName, rangeId)
		} else {
			dealerId, err := strconv.Atoi(vars["dealer_id"])
			if err != nil {
				output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Incorrect dealer id. Dealer id %s incorrect\"}", vars["dealer_id"])
			} else {
				resultString, err := utils.FormatSerials(serials, serialFormat)
				if err != nil {
					output = fmt.Sprintf("{ \"error\": \"2\", \"message\":\"Format %s unsupported\"}", serialFormat)
				} else {
					if last := len(resultString) - 1; last >= 0 {
						resultString = resultString[:last]
					}

					downloadFileName := tableName + "_" + utils.GetDealerName(dealerId) + "_" + rangeId + ".csv"

					w.Header().Set("Content-Type", "text/csv")
					w.Header().Set("Content-Disposition", "attachment; filename="+downloadFileName)
					fmt.Fprint(w, resultString)
				}
			}
		}
	} else {
		output = fmt.Sprintf("{ \"error\": \"1\", \"message\":\"Incorrect access token. Token %s incorrect\"}", accessToken)
	}

	fmt.Fprint(w, output)
}

/**
 * активация серийника
 */
func ActivateSerialsHandler(w http.ResponseWriter, r *http.Request) {
	result := true

	vars := mux.Vars(r)
	tableName := vars["table_name"]
	serialKey := vars["serial_key"]

	canActivate, tryCount := services.SerialCheck(tableName, serialKey)

	if canActivate {
		result = services.SerialUpdate(tableName, tryCount, serialKey)
	} else {
		result = false
	}

	fmt.Fprintf(w, "{is_activated:%t}", result)
}

/**
 * сбросить серийник
 */
func ResetSerialsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["table_name"]
	serialKey := vars["serial_key"]

	res, rowsAffected := services.ResetSerial(tableName, serialKey)

	if res == true {
		if rowsAffected > 0 {
			fmt.Fprint(w, "Ключ успешно сброщен")
		} else {
			fmt.Fprint(w, "0 полей обновлено, что то не так")
		}
	} else {
		fmt.Fprint(w, "Произощла ошибка, ключ не сброшен")
	}
}

/**
 * информация о серийнике
 */
func AboutSerialsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["table_name"]
	serialKey := vars["serial_key"]

	result, activatedCount, maxActivation, activatedTime := services.AboutSerial(tableName, serialKey)

	if result == true {
		fmt.Fprintf(w, "Количество активаций: %d\nВремя последней активации: %s\nМаксимальное количество активаций: %d", activatedCount, activatedTime, maxActivation)
	} else {
		fmt.Fprint(w, "в данной таблице, не существует заданный серийный ключ")
	}
}

/******************************************************MISC*****************************************************************/

/**
 * отправляем почту
 */
func SendMailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.FormValue("name") == "" || r.FormValue("email") == "" || r.FormValue("message") == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "error")

		return
	}

	body := fmt.Sprintf("Имя отправителя: %s,\nПочта отправителя: %s,\nСообщение отправителя: %s", r.FormValue("name"), r.FormValue("email"), r.FormValue("message"))

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@unimedia.uz")
	m.SetHeader("To", "info@unimedia.uz")
	m.SetHeader("Subject", "Отправлено из формы сайта.")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.yandex.ru", 465, "noreply@unimedia.uz", "1q2w3e4r5t")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%s", "complete")
}
