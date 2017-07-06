package main

import (
	"os"
	"io/ioutil"
	"net/http"
	"github.com/Sirupsen/logrus"
	"strconv"
	"time"
)

func checkError(err error) {
	if err != nil {
		logrus.Errorf("Fatal Error: %v", err)
		os.Exit(1)
	}
}

func getServiceIndex() (index int){
	// Simple function to return index of service //
	resp, err := http.Get("http://rancher-metadata/latest/self/container/service_index")

	for err != nil  {
		// Assuming error is a non nil value then retry //
		time.Sleep(60 * time.Second)
		resp, err = http.Get("http://rancher-metadata/latest/self/container/service_index")
	}
	defer resp.Body.Close();
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	index, err = strconv.Atoi(string(body))
	checkError(err)
	return

}
