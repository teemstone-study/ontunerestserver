package app

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBInfo struct {
	host   string
	port   string
	user   string
	pwd    string
	dbname string
	wPort  string
	mIPPort string	
}

// Rectangle 를 반환하는 함수를 만들었다.
func CreateDBInfo(envfilepath string) *DBInfo {
	//db := DBInfo{}
	db := new(DBInfo)
	db.loadDBConf(envfilepath)
	return db
	//return &db
}

func (db *DBInfo) loadDBConf(envfilepath string) {

	err := godotenv.Load(envfilepath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.host = os.Getenv("DBHost")
	db.port = os.Getenv("DBPort")
	db.user = os.Getenv("DBUser")
	db.pwd = os.Getenv("DBPwd")
	db.dbname = os.Getenv("DBName")
	db.wPort = os.Getenv("WPort")
	db.mIPPort = os.Getenv("MIPPort")

	log.Println(db.GetDBConnString())
}

func (db *DBInfo) GetDBConnString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		db.host, db.user, db.pwd, db.dbname, db.port)
}

func (db *DBInfo) GetWPort() string {
	return db.wPort
}

func (db *DBInfo) GetManagerConnString() string {
	return db.mIPPort
}
