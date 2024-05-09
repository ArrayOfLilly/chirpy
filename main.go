package main

import (
	"fmt"
	"net/http"
	"os"
	)



var ServeMux = http.NewServeMux()

type Server struct {
	Addr string
	Handler http.Handler
}

func (s Server) ListenAndServe() {
    http.ListenAndServe(s.Addr, s.Handler)
}

func getServer() Server {
	return Server{
		Addr: ":8080", 
		Handler: ServeMux,
	}
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Error: too many arguments")
		fmt.Println("Usage: chirpy")
		return
	}

	fmt.Printf("Server is starting...\n")

	s := getServer()

	s.ListenAndServe()
}

