package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerChirpsRead(w http.ResponseWriter, req *http.Request) {
	file, _ := os.Open(cfg.DB)
	defer file.Close()
	

}



