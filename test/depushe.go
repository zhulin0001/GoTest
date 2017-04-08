package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Info struct {
	Password string
	UserType string `json:"userType"`
}

type Result struct {
	Code int32
	Data map[string]interface{}
}

func Depushe() {
	url := ""
	index := 1
	errorCount := 0
	fileName := "result.txt"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	for {
		url := fmt.Sprintf(url, 10000+index)
		resp, err := http.Get(url)
		CheckError(err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		CheckError(err)
		req := new(Result)
		err = json.Unmarshal(body, &req)
		CheckError(err)
		index = index + 1
		if req.Data["userType"] == nil {
			if errorCount >= 5 {
				break
			}
			errorCount = errorCount + 1
			continue
		}
		fmt.Println(string(body))
		dstFile.WriteString(string(body) + "\n")
	}
}

//CheckError if err is not nil, then log the detail and return true, otherwise return false
func CheckError(err error) (ret bool) {
	ret = false
	if err != nil {
		fmt.Printf(fmt.Sprintf("Error[]: %s", err.Error()))
		ret = true
		os.Exit(1)
	}
	return
}
