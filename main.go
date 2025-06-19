package main


// go mod init github.com/Madhav-M01/mangodb
//go get -u github.com/gorilla/mux
//go get go.mongodb.org/mongo-driver/v2/mongo
	

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Madhav-M01/mangodb/router"
)

func main() {
	fmt.Println("Starting server at port 4000...")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":4000", r))
}
