package utils

import (
	"os"

	"fmt"

	log "github.com/sirupsen/logrus"
)

//CheckError if err is not nil, then log the detail and return true, otherwise return false
func CheckError(err error, desc string) (ret bool) {
	ret = false
	if err != nil {
		log.Error(fmt.Sprintf("Error[%s]: %s", desc, err.Error()))
		ret = true
	}
	return
}

//CheckFileIsExist return true when filename exists, otherwise return false
func CheckFileIsExist(filename string) (exist bool) {
	exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return
}
