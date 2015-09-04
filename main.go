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


//JSON STRUCTRURE
	type JSON map[string]interface {}
	type ActionFunction func(string, JSON) // func(entityID,input:'json')
	
	// Structure used for setting an knowledge in the creation of an entity
	// Change this structure if you want another initialization
	type KnowledgeJSON struct {
		People int `json:"people"`
		Light bool `json:"light"`
	} 

	//Structure used for creating a new entity
	type Entity struct {
		Behavior string `json:"behavior"`
		Knowledge KnowledgeJSON `json:"knowledge"`
	}

	//Structure used for registering an action
	type ActionJSON struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	type DetectPeopleOutput struct {
		People int `json:"peopleValue"`
	}

	type LightOutput struct {
		Light bool `json:"light"`
	}

	

	// INFO
	var TUTO_APP_ID string
	var TUTO_APP_SECRET string
	var httpURL string
	var instanceID string
	var actions map[string]ActionFunction

	var sysou bool

	// ENVIRONEMENT VARIABLES
	var people int
	var light bool


// REQUEST

func Request(requestType string, url string,JSONbody interface{}) []byte {
	jsonStr, err := json.Marshal(JSONbody)
	request, err := http.NewRequest(requestType, url, bytes.NewBuffer(jsonStr))
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
    	fmt.Println("HTTP REQUEST :", request.URL)
    	fmt.Println("JSON REQUEST :", JSONbody)
	    fmt.Println("response Status:", resp.Status)
	    fmt.Println("response Headers:", resp.Header)
	    fmt.Println("response Body:", string(body))
	   	fmt.Println("")
	}
	return body 
}

// Creates an instance
func CreateInstance(user string, project string, version string){
	body := Request("PUT",httpURL + "?" + "scope=app","")
	var f interface{}
    err := json.Unmarshal(body,&f)
    if err != nil {panic(err)}
   	m := f.(map[string]interface {})
   	instance := m["instance"].(map[string]interface {})
	instanceID = instance["instance_id"].(string);
}

// Creates an entity and return its ID
func CreateEntity() string{
	subm := KnowledgeJSON{People:0,Light:false}
	m := Entity{Behavior:"main.bt",Knowledge:subm}
	body := Request("PUT",httpURL + "/" + instanceID + "/entities",m)
		var f interface{}
    err := json.Unmarshal(body,&f)
    if err != nil {panic(err)}
   	j := f.(map[string]interface {})
   	entity := j["entity"].(map[string]interface {})
	entityId := strconv.FormatInt(int64(entity["id"].(float64)), 10)
	return entityId
}

//Register an action 
func RegisterAction(actionID string,ngrok_url string){
	m := ActionJSON{Name:actionID,Url:ngrok_url+"/actions/"+ actionID +"/"}
	Request("PUT", httpURL + "/" + instanceID + "/actions",m)
}

func sendSuccess(requestID string, output interface{}){
	Request("POST", httpURL + "/" + instanceID + "/actions/" + requestID + "/success", output)
}

func sendCancel(requestID string){
	Request("POST", httpURL + "/" + instanceID + "/actions/" + requestID + "/cancelation", "")
}

func UpdateEntityKnowledge(entityID string, know JSON) {
	Request("POST",httpURL + "/" + instanceID + "/entities/" + entityID + "/knowledge", know)
}

func UpdateGlobalKnowledge(know JSON){
	Request("POST", httpURL + "/" + instanceID + "/knowledge", know)
}

func UpdateCraft(){
	fmt.Println("UPDATE PEOPLE ", people)
	Request("POST",httpURL + "/" + instanceID + "/update",`{"time":"0.5","ts":` + strconv.FormatInt(time.Now().Unix(),10) +`}`)
}

// Delete an instance with its ID
func DeleteInstance(){
    Request("DELETE",httpURL + "/" + instanceID, "")
}

func getEntityKnowledge(entityID string) JSON{
   	body :=	Request("GET", httpURL + "/" + instanceID + "/entities/" + entityID + "/knowledge", "")
   	var f interface{}
    err := json.Unmarshal(body,&f)
    if err != nil {panic(err)}
   	m := f.(map[string]interface{})
   	return m
}


// ACTIONS FUNCTION

func Light(requestID string, input JSON) {
	fmt.Println("Light On !!")
	light = true
	m := LightOutput{Light:light}
  	sendSuccess(requestID,m)
}

func LightOff(requestID string, input JSON) {
	fmt.Println("Light Off !!")
	light = false
	m := LightOutput{Light:light}
  	sendSuccess(requestID,m)
}

func DetectPeople(requestID string, input JSON) {	
	m := DetectPeopleOutput{People:people}
	sendSuccess(requestID,m)
}

func Cancel(requestID string) {
	fmt.Println("(cancel)")
	sendCancel(requestID)
}


func initServer(){
	InitRoute()

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}


func InitNgrok() string{
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

    return tunnels_0["public_url"].(string)
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

/*func RandomHumanInteraction(entityID string){
	c := time.Tick(3000 * time.Millisecond)
    for _ = range c {
    	r := rand.Intn(10)
    	knowledge := getEntityKnowledge(entityID) // Get the entity Knowledge to modify the number of people in the room
    	people = int64(knowledge["people"].(float64))
    	if r <4 {
    		fmt.Println("A man is entering")
    		people++;
    	} else {
    		if(people>0){
    			fmt.Println("A man is exiting")
    			people--;
    		}
    	}
    }
}*/

func InitRoute(){
	for key, fnc := range actions {
		key := key 
		fnc := fnc // <- Memorize the action name and its function, 
    	http.HandleFunc("/actions/"+key+"/start", 
    		func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {panic(err)}

		    	var f interface{}
		    	err = json.Unmarshal(body,&f)
		   		m := f.(map[string]interface{})

		   		requestIDf64 := m["requestId"].(float64)	
		   		requestID := strconv.FormatInt(int64(requestIDf64), 10)
		   		input := m["input"].(map[string]interface{})
				fnc(requestID,input)
			})
		
		http.HandleFunc("/actions/"+key+"/cancel", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {panic(err)}

	    	var f interface{}
	    	err = json.Unmarshal(body,&f)
	   		m := f.(map[string]interface{})

	   		requestIDf64 := m["requestId"].(float64)	
	   		requestID := strconv.FormatInt(int64(requestIDf64), 10)
			Cancel(requestID)
		})


		fmt.Println("The action "+key+" is routed with the function : ", fnc)
	}


	http.HandleFunc("/people", 
    	func(w http.ResponseWriter, r *http.Request) {
    		body, err := ioutil.ReadAll(r.Body)
			if err != nil {panic(err)}

		    var f interface{}
		    err = json.Unmarshal(body,&f)
		   	m := f.(map[string]interface{})

		   	people, _ = strconv.Atoi(m["people"].(string))

		   	fmt.Println("PEOPLE = ",people)
    })

    http.HandleFunc("/get/light", 
    	func(w http.ResponseWriter, r *http.Request) {
    		w.Header().Set("content-type", "application/json; charset=utf-8");

    		lightJSON := make(JSON)
    		lightJSON["light"] = light

    		jsonStr, _ := json.Marshal(lightJSON)

    		w.Write(jsonStr)
    })

    http.HandleFunc("/post/light", 
    	func(w http.ResponseWriter, r *http.Request) {
    		body, err := ioutil.ReadAll(r.Body)
			if err != nil {panic(err)}

		    var f interface{}
		    err = json.Unmarshal(body,&f)
		   	m := f.(map[string]interface{})

		   	light, _ = m["light"].(bool)

		   	lightJSON := make(JSON)
    		lightJSON["light"] = light
		   	UpdateEntityKnowledge("0",lightJSON)
    })

    http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, r.URL.Path[1:])
	})

}


func main() {

	go initServer()

	// PARAM
	sysou = false

	// INIT INFO
	user := "Totoketchup"
	project := "GoAI"
	version := "master"
	TUTO_APP_ID = "MgYR67znmswahjlzW4MY"
	TUTO_APP_SECRET = "qjoenA2CXNdco1VXAQncOLCS7zpW9uqeuFNxGXtu"
	httpURL = "https://runtime.craft.ai/api/v1/" + user + "/" + project + "/" + version
	actions = map[string]ActionFunction {"light": Light, "detectPeople": DetectPeople, "lightoff": LightOff}
	

	ngrok_url := InitNgrok()

	CreateInstance(user, project, version)
	fmt.Println("")

	entityID := CreateEntity();
	fmt.Println("EntityID = ",entityID)

	for key, _ := range actions {
		RegisterAction(key,ngrok_url)
	}

	HandleCTRL_C()

	//go RandomHumanInteraction(entityID);

    c := time.Tick(100 * time.Millisecond)
    for _ = range c {
    	UpdateCraft()
    }

	DeleteInstance()
	fmt.Println("")
}



