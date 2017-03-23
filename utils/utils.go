package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
	"fmt"
)

//CheckError if err is not nil, then log the detail and return true, otherwise return false
func CheckError(err error, desc string) (ret bool) {
	if err != nil {
		log.Error(fmt.Sprintf("Error[%s]: %s", desc, err.Error()))
		return true
	}
	return false
}

//CheckFileIsExist return true when filename exists, otherwise return false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
