package services

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"crypto/sha1"
	"hash"
	"github.com/korasdor/colarit/utils"
	"github.com/korasdor/colarit/model"
)

var (
	db    *sql.DB
	mHash hash.Hash
)

func InitDb() {
	var err error

	db, err = sql.Open("mysql", model.DB_SOURCE_NAME)
	db.SetMaxIdleConns(0)

	if err != nil {
		model.DbSuccess = "error"
		utils.PrintOutput(err.Error())
	} else {

		model.DbSuccess = "success"
		utils.PrintOutput("db success")
	}

}

func CreateSerialsTable(tableName string) bool {
	stmt, err := db.Prepare("CREATE TABLE " + tableName + "(" +
		"serial_id INT NOT NULL AUTO_INCREMENT," +
		"serial_activated SMALLINT NOT NULL," +
		"max_activations SMALLINT NOT NULL," +
		"serial_activated_time VARCHAR(30)," +
		"dealer_id INT NOT NULL," +
		"range_id INT NOT NULL," +
		"serial_key VARCHAR(200) NOT NULL," +
		"PRIMARY KEY ( serial_id ))")

	defer stmt.Close()

	if err != nil {
		utils.PrintOutput(err.Error())

		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		utils.PrintOutput(err.Error())

		return false
	}

	return true
}

func FillSerialsTable(tableName string, serials []string, dealerId string, rangeId string) bool {
	var values string

	for i := 0; i < len(serials); i++ {
		mHash = sha1.New()
		mHash.Write([]byte(serials[i]))

		serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

		if i < len(serials)-1 {
			values += fmt.Sprintf("('%s',0,3,NULL,%s, %s),", serialHash, dealerId, rangeId)
		} else {
			values += fmt.Sprintf("('%s',0,3,NULL,%s, %s)", serialHash, dealerId, rangeId)
		}
	}

	stmt, err := db.Prepare("INSERT INTO " + tableName + "(serial_key,serial_activated,max_activations,serial_activated_time,dealer_id,range_id) VALUES " + values)
	defer stmt.Close()

	if err != nil {
		utils.PrintOutput(err.Error())

		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		utils.PrintOutput(err.Error())

		return false
	}

	return true
}

func SerialCheck(tableName string, key string) (bool, int) {
	result := false
	serialActivated := -1

	mHash = sha1.New()
	mHash.Write([]byte(key))
	serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

	rows, err := db.Query("SELECT serial_activated,serial_key,max_activations FROM "+tableName+" WHERE serial_key=?", serialHash)
	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	defer rows.Close()

	if rows.Next() == true {
		result = true
		var serialKey string
		var maxActivations int

		err := rows.Scan(&serialActivated, &serialKey, &maxActivations)
		if err != nil {
			result = false
			utils.PrintOutput(err.Error())
		} else {
			if serialActivated < maxActivations {
				result = true
				serialActivated++
			} else {
				result = false
			}
		}
	}

	return result, serialActivated
}

func SerialUpdate(tableName string, tryCount int, key string) bool {
	result := true

	stmt, err := db.Prepare("UPDATE " + tableName + " SET serial_activated=?, serial_activated_time=? WHERE serial_key=?")

	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	defer stmt.Close()

	activatedTime := time.Now().Format(time.RFC3339)
	mHash = sha1.New()
	mHash.Write([]byte(key))
	serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

	_, err = stmt.Exec(tryCount, activatedTime, serialHash)
	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	return result
}

func ResetSerial(tableName string, key string) (bool, int64) {
	result := true

	stmt, err := db.Prepare("UPDATE " + tableName + " SET serial_activated=0, serial_activated_time=NULL WHERE serial_key=?")

	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	defer stmt.Close()

	mHash = sha1.New()
	mHash.Write([]byte(key))
	serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

	out, err := stmt.Exec(serialHash)
	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	rowsAffected, _ := out.RowsAffected()

	return result, rowsAffected
}

func AboutSerial(tableName string, key string) (bool, int, int, string) {
	result := false

	var (
		serialActivated     int
		maxActivations      int
		serialActivatedTime string
	)

	mHash = sha1.New()
	mHash.Write([]byte(key))
	serialHash := fmt.Sprintf("%x", mHash.Sum(nil))

	rows, err := db.Query("SELECT serial_activated, max_activations, COALESCE(serial_activated_time, '') as serial_activated_time FROM "+tableName+" WHERE serial_key=?", serialHash)
	if err != nil {
		result = false
		utils.PrintOutput(err.Error())
	}

	defer rows.Close()

	if rows.Next() == true {
		result = true

		err := rows.Scan(&serialActivated, &maxActivations, &serialActivatedTime)
		if err != nil {
			result = false
			utils.PrintOutput(err.Error())
		} else {
			result = true
		}
	}

	return result, serialActivated, maxActivations, serialActivatedTime
}

func CloseDb() {
	db.Close()
}
