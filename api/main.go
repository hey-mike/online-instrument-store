package main

import (
	"log"

	"example.com/online-store/src/routes/router"
)

func main() {
	log.Println("Main log....")
	router.RunAPI(":9090")
}
