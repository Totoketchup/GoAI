package main

import (
	"log"
	"fmt"
	"net/http"
)

type String string

type Struct struct {
    Greeting string
    Punct    string
    Who      string
}


func (s String) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	fmt.Fprint(w, "<-- String response!")
}

func (s Struct) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	fmt.Fprint(w, "<-- Struct response!")
}

func main() {
	http.Handle("/string", String("Im a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	user := "Totoketchup"
	project := "Test"
	version := "master"
	TUTO_APP_ID := "0"
	TUTO_APP_SECRET := "0"


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
    fmt.Println("response Body:", resp.Body)
}