package common

import (
	"log"
	"os"
	"path/filepath"
)

func GetServiceName(str string) string {
	return filepath.Base(filepath.Dir(str))
}

func LogToFile() {
	file, _ := os.OpenFile(filepath.Join(filepath.Dir(os.Args[0]), "log"), os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	log.SetOutput(file)
}

func FileIsExist(FileName string) bool {
	_, err := os.Stat(FileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetRunDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return dir
}
