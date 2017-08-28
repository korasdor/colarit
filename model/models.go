package model

import "os"

var (
	DbSuccess       string
	SerialKeyLength int = 10

	//DB_SOURCE_NAME string = "bacdeca17611c1:a9e99d39@tcp(eu-cdbr-west-01.cleardb.com:3306)/heroku_a4f8eaa86b016fe"
	DB_SOURCE_NAME string = os.Getenv("MYSQ_DATA_SOURCE_NAME")
	ACCESS_TOKEN   string = os.Getenv("ACCESS_TOKEN")
)
