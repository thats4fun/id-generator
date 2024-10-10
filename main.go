package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thats4fun/id-generator/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <getid|freeid> [ID]\n", os.Args[0])
	}

	store, err := internal.NewIDStore()
	if err != nil {
		log.Println(err)
		return
	}
	defer store.Close()

	switch os.Args[1] {
	case "getid":
		id := store.GetId()
		fmt.Println(id)
	case "freeid":
		if len(os.Args) != 3 {
			log.Println("Usage: freeid <ID>")
			return
		}
		id := os.Args[2]
		err := store.FreeId(id)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("ID %s freed\n", id)
	default:
		log.Printf("Unknown command: %s\n", os.Args[1])
	}
}
