package main

import (
	"log"
	"net/http"
	"samdriver/dungeon/llm"
)

func main() {
	log.Println("Starting.")

	prepareRoutes()

	log.Fatal(http.ListenAndServe(":8087", nil))
}

func prepareRoutes() {
	http.HandleFunc("/input", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Processing input.")

		err := llm.ReceiveInputHandler(writer, request)
		if err != nil {
			log.Println("Error processing input:", err)
		}
	})

	// Static files.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/static"))))

	// Index, with 404 for everything else.
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("URL requested:", request.URL.Path)

		if request.URL.Path != "/" {
			http.NotFound(writer, request)
			return
		}

		http.ServeFile(writer, request, "public/index.html")
	})
}
