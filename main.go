// this is not sanctioned

package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
)

func main() {
	address := flag.String("address", ":8001", "Address to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			w.WriteHeader(200)
			fmt.Fprintf(w, "OK")

		})

	log.Println("Server Starting. Listening on", *address)
	http.ListenAndServe(*address, nil)
}
