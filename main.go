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
	"math/rand"
)


	//TYPES

	type ActionFunction func(string, map[string]interface{}) // func(entityID,input:'json')


	// INFO
	
	var user string
	var project string
	var version string
	var TUTO_APP_ID string
	var TUTO_APP_SECRET string
	var httpURL string
	var instanceID string
	var ngrok_url string
	var actions map[string]ActionFunction

	var sysou bool


func sendSuccess(requestID string, output string){
	var jsonStr = []byte(output)
	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/actions/" + requestID + "/success", bytes.NewBuffer(jsonStr))
	if err != nil {panic(err)}
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);
	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {panic(err)}
    defer resp.Body.Close()
    if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("Send success from action request Id :"+requestID)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
}



func Light(requestID string, input map[string]interface{}) {
	fmt.Println("Light On !! (start)")
	fmt.Println("My name is ", input["name"])
	fmt.Println("My height is ", input["height"])

	var nice string
	if(input["nice"].(bool)){
		nice = "YES"
	} else {
		nice ="NO"
	}

	fmt.Println("Am I nice ?! ", nice)
  	sendSuccess(requestID,"{}")
}

func LightOff(requestID string, input map[string]interface{}) {
	fmt.Println("Light Off !!")
  	sendSuccess(requestID,"{}")
}

func Explosion(requestID string, input map[string]interface{}) {
	fmt.Println("EXPLOSION !! (start)")
	fmt.Println("My charge is ", input["charge"])
	sendSuccess(requestID,"{}")
}

func Enter(requestID string, entityID string, input map[string]interface{}) {
	fmt.Println("Someone is entering in the room")
	knowledge := getEntityKnowledge(entityID)

	people := knowledge["people"].(float64) + 1

	output := `{"result":`+strconv.FormatInt(int64(people),10)+`}`

	sendSuccess(requestID,output)
}

func Exit(requestID string, entityID string, input map[string]interface{}) {
	fmt.Println("Someone is exiting the room")
	knowledge := getEntityKnowledge(entityID)

	people := knowledge["people"].(float64) - 1

	output := `{"result":`+strconv.FormatInt(int64(people),10)+`}`

	sendSuccess(requestID,output)
}

func Cancel() {
	fmt.Println("(cancel)")
}

func getEntityKnowledge(entityID string) map[string]interface{}{
	request, err := http.NewRequest("GET", httpURL + "/" + instanceID + "/entities/" + entityID + "/knowledge", nil)
	if err != nil {panic(err)}
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {panic(err)}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if(sysou){
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
	var f interface{}
    err = json.Unmarshal(body,&f)
   	m := f.(map[string]interface{})
   	return m
}

func UpdateEntityKnowledge(entityID string, know Knowledge) {
	jsonStr, err := json.Marshal(know)

	var test Knowledge 
	err = json.Unmarshal(jsonStr,&test)
	fmt.Println("TEST = ",test)

	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/entities/" + entityID + "/knowledge", bytes.NewBuffer(jsonStr))
	if err != nil {panic(err)}
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);
	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {panic(err)}
    defer resp.Body.Close()
    if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	   	fmt.Println("KNOWLEDGE CHANGED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}
}


type Knowledge struct {
	Destination string
	Value interface {}
}

func UpdateGlobalKnowledge(know Knowledge){
	jsonStr, err := json.Marshal(know)


	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/knowledge", bytes.NewBuffer(jsonStr))
	if err != nil {panic(err)}
	request.Header.Set("content-type", "application/json; charset=utf-8");
	request.Header.Set("accept", "");
	request.Header.Set("X-Craft-Ai-App-Id", TUTO_APP_ID);
	request.Header.Set("X-Craft-Ai-App-Secret", TUTO_APP_SECRET);
	client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {panic(err)}
    defer resp.Body.Close()
    if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}

}

func InitRoute(){
	for key, fnc := range actions {
		key := key
		fnc := fnc
    	http.HandleFunc("/actions/"+key+"/start", 
    		func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
		        	panic(err)
		    	}
		    	var f interface{}
		    	err = json.Unmarshal(body,&f)
		   		m := f.(map[string]interface{})

		   		requestID64 := m["requestId"].(float64)	
		   		requestID := strconv.FormatInt(int64(requestID64), 10)

		   		input := m["input"].(map[string]interface{})

		   		//fmt.Println("Response from the entity : ",m["entityId"].(float64))
				fnc(requestID,input)
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
	body, err := ioutil.ReadAll(resp.Body)

    var f interface{}
    err = json.Unmarshal(body,&f)
    m := f.(map[string]interface{})

    tunnels := m["tunnels"].([]interface{})
    tunnels_0 := tunnels[0].(map[string]interface{})

    ngrok_url = tunnels_0["public_url"].(string)
}


// Create an instance
func CreateInstance(user string, project string, version string){
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
    body, _ := ioutil.ReadAll(resp.Body)
    instanceID = string(body);

    if(sysou){
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
} 

//Create an entity
func CreateEntity(){
    var jsonStr = []byte(`{"behavior":"main.bt","knowledge":{"people":0}}`)
  	request, err := http.NewRequest("PUT", httpURL + "/" + instanceID + "/entities", bytes.NewBuffer(jsonStr))
  	fmt.Println("HTTP CREATION ENTITY:", request.URL)

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
    if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
}


func UpdateCraft(){
	var jsonStr = []byte(`{"time":"0.5","ts":` + strconv.FormatInt(time.Now().Unix(),10) +`}`)
  	request, err := http.NewRequest("POST", httpURL + "/" + instanceID + "/update", bytes.NewBuffer(jsonStr))
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
   	if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
}

//Register an action 
func RegisterAction(actionID string){
	var jsonStr = []byte(`{"name":"`+ actionID +`","url":"`+ngrok_url+`/actions/`+ actionID +`/"}`)
  	request, err := http.NewRequest("PUT", httpURL + "/" + instanceID + "/actions", bytes.NewBuffer(jsonStr))
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
    if(sysou){
	    body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	}
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

func RandomHumanInteraction(){
	c := time.Tick(1000 * time.Millisecond)
    for _ = range c {
    	r := rand.Intn(10)
    	knowledge := getEntityKnowledge("0") // Get the entity Knowledge to modify the number of people in the room
    	people := knowledge["people"].(float64)
    	if r <4 {
    		fmt.Println("A man is entering")
    		people++;
    	} else {
    		if(people>0){
    			fmt.Println("A man is exiting")
    			people--;
    		}
    	}

    	fmt.Println("Total of men : ",people)

    	know := Knowledge{"people",int64(people)}
    	fmt.Println("***********************************",know)
    	UpdateEntityKnowledge("0",know)
    	fmt.Println("")
    }
}


func main() {

	go initServer()

	// INIT INFO

	user = "Totoketchup"
	project = "GoAI"
	version = "master"
	TUTO_APP_ID = "MgYR67znmswahjlzW4MY"
	TUTO_APP_SECRET = "qjoenA2CXNdco1VXAQncOLCS7zpW9uqeuFNxGXtu"
	httpURL = "https://runtime.craft.ai/api/v1/" + user + "/" + project + "/" + version

	// PARAM
	sysou = false

	actions = map[string]ActionFunction {"light": Light, "explosion": Explosion, "lightoff": LightOff}



	go InitNgrok()

	fmt.Println("")
	CreateInstance(user, project, version)
	fmt.Println("")

	CreateEntity();
	fmt.Println("")

	RegisterAction("light")
	RegisterAction("enter")
	RegisterAction("exit")
	RegisterAction("explosion")
	RegisterAction("lightoff")

	HandleCTRL_C()

	go RandomHumanInteraction();

    c := time.Tick(500 * time.Millisecond)
    for _ = range c {
    	UpdateCraft()
    	fmt.Println("")
    }

	DeleteInstance()
	fmt.Println("")
	

}



