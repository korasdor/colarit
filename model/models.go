package model

import "os"

var (
	DbSuccess       string
	SerialKeyLength int = 10
)

const (
	DB_SOURCE_NAME string = "bacdeca17611c1:a9e99d39@tcp(166.62.10.138:3306)/colar_db"
	//DB_SOURCE_NAME string = "bacdeca17611c1:a9e99d39@tcp(eu-cdbr-west-01.cleardb.com:3306)/heroku_a4f8eaa86b016fe"
)

var (
	ACCESS_TOKEN string = os.Getenv("ACCESS_TOKEN")
)