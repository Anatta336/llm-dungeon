package main

import (
	"log"
	"net/http"
	"samdriver/dungeon/config"
	"samdriver/dungeon/dm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	log.Println("Starting.")

	prepareRoutes()

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, nil))
}

func prepareRoutes() {
	http.HandleFunc("POST /input", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Processing input.")

		err := dm.ReceiveInputHandler(writer, request)
		if err != nil {
			log.Println("Error processing input:", err)
		}
	})

	// Static files.
	http.Handle("GET /static/{file}",
		http.StripPrefix("/static/", http.FileServer(http.Dir("public/static"))),
	)

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
