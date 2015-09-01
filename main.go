package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"os"
	"os/signal"
	"time"
	"strconv"
	//"os/exec"
	"encoding/json"
)


	// INFO
	
	var user string
	var project string
	var version string
	var TUTO_APP_ID string
	var TUTO_APP_SECRET string
	var httpURL string
	var instanceID string
	var ngrok_url string

func sendSuccess(requestID string){
	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/actions/" + requestID + "/success", nil)
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

    fmt.Println("Send success from action request Id :"+requestID)
}



func Light(requestID string) {
	fmt.Println("Light On !! (start)")
  	sendSuccess(requestID)
}

func Cancel() {
	fmt.Println("(cancel)")
}

func Explosion(requestID string) {
	fmt.Println("EXPLOSION !! (start)")
	sendSuccess(requestID)
}




func InitRoute(){

	actions := map[string]func(string) {"light":Light, "explosion":Explosion}

	for key, fnc := range actions {

    	http.HandleFunc("/actions/"+key+"/start", func(w http.ResponseWriter, r *http.Request) {	
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
	        	panic(err)
	    	}
	    	var f interface{}
	    	err = json.Unmarshal(body,&f)
	   		m := f.(map[string]interface{})

	   		requestID64 := m["requestId"].(float64)	
	   		requestID := strconv.FormatInt(int64(requestID64), 10)

	   		fmt.Println("Response from the entity : ",m["entityId"].(float64))
			fnc(requestID)
		})
		
		http.HandleFunc("/actions/"+key+"/cancel", func(w http.ResponseWriter, r *http.Request) {
			Cancel()
		})


		fmt.Println("The action "+key+" is routed with the function : ", fnc)
	 }

}


func initServer(){
	InitRoute()

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func InitNgrok(){
	/*fmt.Println("Launch Ngrok :")
	parts := strings.Fields("ngrok http 8000")
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
        panic(err)
    }
    fmt.Println("Ngrok launched")
    fmt.Println("%s", out)
	time.Sleep(1500 * time.Millisecond)*/

	request, err := http.NewRequest("GET","http://127.0.0.1:4040/api/tunnels", nil)
	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))

    var f interface{}
    err = json.Unmarshal(body,&f)
    m := f.(map[string]interface{})

    tunnels := m["tunnels"].([]interface{})
    tunnels_0 := tunnels[0].(map[string]interface{})

    ngrok_url = tunnels_0["public_url"].(string)

}


// Create an instance
func CreateInstance(){
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
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))

    instanceID = string(body);

} 

//Create an entity
func CreateEntity(){
    var jsonStr = []byte(`{"behavior":"main.bt","knowledge":{}}`)
  	request, err := http.NewRequest("PUT", httpURL + "/" + instanceID + "/entities", bytes.NewBuffer(jsonStr))
  	//fmt.Println("********HTTP URL********" , httpURL)
  	//fmt.Println("********ENTITY********" , request.URL)
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
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))
}


func UpdateCraft(){
	var jsonStr = []byte(`{"time":"0.5","ts":` + strconv.FormatInt(time.Now().Unix(),10) +`}`)
  	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/update", bytes.NewBuffer(jsonStr))
  	//fmt.Println("********HTTP URL********" , httpURL)
  	//fmt.Println("********Update********", request.URL)
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

/*    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))*/
}

//Register an action 
func RegisterAction(actionID string){
	var jsonStr = []byte(`{"name":"`+ actionID +`","url":"`+ngrok_url+`/actions/`+ actionID +`/"}`)
  	request, err := http.NewRequest("PUT", httpURL + "/" + instanceID + "/actions", bytes.NewBuffer(jsonStr))
  	//fmt.Println("********HTTP URL********" , httpURL)
  	//fmt.Println("********ENTITY********" , request.URL)
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
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))
}


// Delete an instance with its ID
func DeleteInstance(){
	request, err := http.NewRequest("DELETE",httpURL + "/" + instanceID, nil)
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
	body, err := ioutil.ReadAll(resp.Body)
    fmt.Println("Body :",string(body))
}

func HandleCTRL_C(){
	// CTRL C HANDLING
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
	    for _ = range signalChan {
	        fmt.Println("\nReceived an interrupt, stopping instance...\n")
	        DeleteInstance()
	        panic("ERREUR")
	    }
	}()
}


func main() {

	go initServer()

	c := time.Tick(500 * time.Millisecond)
    for _ = range c {
    }

	// INIT INFO
/*
	user = "Totoketchup"
	project = "GoAI"
	version = "master"
	TUTO_APP_ID = "MgYR67znmswahjlzW4MY"
	TUTO_APP_SECRET = "qjoenA2CXNdco1VXAQncOLCS7zpW9uqeuFNxGXtu"
	httpURL = "https://runtime.craft.ai/api/v1/" + user + "/" + project + "/" + version

	go InitNgrok()

	fmt.Println("")
	CreateInstance()
	fmt.Println("")

	CreateEntity();
	fmt.Println("")

	RegisterAction("light")
	RegisterAction("explosion")

	HandleCTRL_C()

	c := time.Tick(500 * time.Millisecond)
    for _ = range c {
    	UpdateCraft()
    	//fmt.Println("")
    }

	DeleteInstance()
	fmt.Println("")
	
*/
}



