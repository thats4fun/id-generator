package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thats4fun/dasda/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <getid|freeid> [ID]\n", os.Args[0])
	}

	store, err := internal.NewIDStore()
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	switch os.Args[1] {
	case "getid":
		id := store.GetId()
		fmt.Println(id)
	case "freeid":
		if len(os.Args) != 3 {
			log.Fatal("Usage: freeid <ID>")
		}
		id := os.Args[2]
		err := store.FreeId(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID %s freed\n", id)
	default:
		log.Fatalf("Unknown command: %s\n", os.Args[1])
	}
}
