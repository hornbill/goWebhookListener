package main

import (
	"encoding/json" //-- For JSON decoding
	"log"           //-- Log errors
	"net/http"      //-- HTTP server
)

//---- Constants ---- //
const authKey string = "123456" //-- Auth Key Constant
const version string = "0.0.1"  //-- Version
const port string = "9000"      //-- Port
//---- Structures ---- //

//-- Structure for JSON Response on Entity Update
type webhookJSONRespStruct struct {
	OnEntityEvent struct {
		ActionBy         string `json:"actionBy"`
		ActionByName     string `json:"actionByName"`
		CallingSessionID string `json:"callingSessionId"`
		Entity           string `json:"entity"`
		EventSource      string `json:"eventSource"`
		EventTime        string `json:"eventTime"`
		Record           struct {
			HPkID string `json:"h_pk_id"`
		} `json:"record"`
	} `json:"onEntityEvent"`
}

//-- Main Function
func main() {

	log.Println("Hornbill Webhook Listner V", version)
	log.Println("Listening on Port:", port)
	//-- Run WebhookCatcher when the url :9000/api is called
	http.HandleFunc("/api", webhookCatcher)
	//-- Run HTTP server on port 9000
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

//-- Function called from /api
func webhookCatcher(w http.ResponseWriter, r *http.Request) {
	//-- Log Request
	log.Println("Incomming Request")
	//-- Check if Auth Key is Valid
	boolKeyIsValid := checkAuthKey(r)

	//-- if Not Throw error
	if !boolKeyIsValid {
		throwError("Invalid or Missing Key", w)
		return
	}
	//-- Try and Decode JSON From Wehbook
	boolJSONProcess := processJSON(r, w)

	//-- If JSON is not decoded correctly then throw error
	if !boolJSONProcess {
		throwError("Unable to Process JSON Response", w)
		return
	}
	//-- All Good
	w.Write([]byte("Success"))
}

//-- Function to check if ?key matched the authKey set
func checkAuthKey(r *http.Request) bool {
	//-- Get URL Parameter 'key'
	key := r.URL.Query().Get("key")

	//-- Check if we have something
	if len(key) != 0 {
		//-- Validate key
		if key == authKey {
			return true
		}
		return false
	}
	return false
}

//-- Decode JSON Response from Webhook
func processJSON(r *http.Request, w http.ResponseWriter) bool {
	log.Println("Process JSON")
	//-- Decode JSON from Request Body
	//-- This only works if the webhook is set to JSON format and not XMLformat
	decoder := json.NewDecoder(r.Body)
	//-- t now becombes a stucture based on the webhookJSONRespStruct struct
	var t webhookJSONRespStruct
	//-- Catch Errors
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Error: ", err)
		return false
	}
	//-- OutputEvent Source
	log.Println("Action Name ", t.OnEntityEvent.EventSource)
	return true
}

//-- Any non 200 reponse to a webook will cause it to fail
func throwError(s string, w http.ResponseWriter) {
	http.Error(w, s, 500)
}
