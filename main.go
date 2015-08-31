package main

import (
	"log"
	"fmt"
	"net/http"
	//"io/ioutil"
	"net/http/httputil"
	"bytes"

)

func initServer(){
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func main() {
	go initServer()

	user := "Totoketchup"
	project := "GoAI"
	version := "master"
	TUTO_APP_ID := "MgYR67znmswahjlzW4MY"
	TUTO_APP_SECRET := "qjoenA2CXNdco1VXAQncOLCS7zpW9uqeuFNxGXtu"

    // CREATE SIMULATION

	httpURL := "https://runtime.craft.ai/api/v1/" + user + "/" + project + "/" + version
	request, err := http.NewRequest("PUT",httpURL + "?" + "scope=app", nil)
	if err != nil {
        panic(err)
    }
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);

	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
	body, err := httputil.DumpResponse(resp, true)
    log.Print("Body :",string(body))

    simID := string(body);

    // CREATE ENTITY

    var jsonStr = []byte(`{"behavior":"main.bt","knowledge":""}`)

  	request, err = http.NewRequest("PUT", httpURL + "/" + simID + "/entities", bytes.NewBuffer(jsonStr))
	if err != nil {
        panic(err)
    }
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);

    resp, err = client.Do(request)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
	body, err = httputil.DumpResponse(resp, true)
    log.Print("Body :",string(body))

}

