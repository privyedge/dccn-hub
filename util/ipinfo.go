package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)


type IPInfo struct {
	IP       string
	Hostname string
	City     string
	Country  string
	Loc      string
	Org      string
}

func GetLatLng(ip string) (string, string)  {
	reader := strings.NewReader(``)
	url := "https://ipinfo.io/"+ip+"?token=05afd766593f88"
	WriteLog(url)
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		fmt.Print(err.Error())
		return "", ""
	}
	client := &http.Client{}
	resp, err := client.Do(request)

	contents, err := ioutil.ReadAll(resp.Body)

	var ipinfo IPInfo

	err = json.Unmarshal(contents, &ipinfo)

	values :=  strings.Split(ipinfo.Loc, ",")


	log :=fmt.Sprintf("loc lat: %s lng: %s \n", values[0], values[1])
	WriteLog(log)
	return values[0], values[1]
}


