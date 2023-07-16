package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	f, err := os.ReadFile("./mock.json")
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./mock.json", j, 0777); err != nil {
		log.Fatal(err)
	}
}
