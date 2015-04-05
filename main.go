// this is not sanctioned

package main

import (
	"flag"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
)

func main() {
	address := flag.String("address", ":8001", "Address to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			w.Header().Add("Cache-Control", "public, max-age=600")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Max-Age", "600")
			w.Header().Add("Access-Control-Allow-Headers", "accept, custom")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, POST, DELETE")

			if r.Method == "OPTIONS" {

				log.Println("OPTIONS")

				w.WriteHeader(200)
				fmt.Fprintf(w, "OK")

			} else {

				log.Println(r.URL)
				body, err := ioutil.ReadAll(r.Body);

				if err != nil {
					log.Println(err)
				}

				log.Println(string(body))

				w.WriteHeader(200)
				fmt.Fprintf(w, "OK")

			}

		})

	http.ListenAndServe(*address, nil)
}
