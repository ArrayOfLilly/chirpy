package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type message struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }

	type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        Error string `json:"error"`
		Valid bool `json:"valid"`
	}

	respBody := returnVals{
        Error: "",
		Valid: true,
    }
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	decoder := json.NewDecoder(req.Body)
	msg := message{}
	err := decoder.Decode(&msg)
	if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		respBody.Error = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
		return
    }

	

	// msg is a struct with data populated successfully
	if len(msg.Body) > 140 {
		log.Printf("Message length exceeds the permitted length: %d", len(msg.Body))
		respBody.Error = "Message length exceeds the permitted length"
		respBody.Valid = false
		w.WriteHeader(400)
		
		_, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			respBody.Error = fmt.Sprintf("%v", err)
			w.WriteHeader(500)
			return
		}
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		respBody.Error = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
		return
	}

	// valid chirp
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dat))
}



    