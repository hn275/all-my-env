package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/handlers"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	// refresh variable counter every second
	go database.RefreshVariableCounter()

	fmt.Print(
		`
 _______   ________   ___      ___ ___  ___  ___  ___  ________     
|\  ___ \ |\   ___  \|\  \    /  /|\  \|\  \|\  \|\  \|\   __  \    
\ \   __/|\ \  \\ \  \ \  \  /  / \ \  \\\  \ \  \\\  \ \  \|\ /_   
 \ \  \_|/_\ \  \\ \  \ \  \/  / / \ \   __  \ \  \\\  \ \   __  \  
  \ \  \_|\ \ \  \\ \  \ \    / /   \ \  \ \  \ \  \\\  \ \  \|\  \ 
   \ \_______\ \__\\ \__\ \__/ /     \ \__\ \__\ \_______\ \_______\
    \|_______|\|__| \|__|\|__|/       \|__|\|__|\|_______|\|_______|
`)

	r := handlers.New()
	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
