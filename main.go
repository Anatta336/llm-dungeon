package main

import (
	"log"
	"net/http"
	"samdriver/dungeon/llm"
)

func main() {
	// file, err := os.OpenFile("dungeon.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening log file: %v", err)
	// }
	// log.SetOutput(file)
	log.Println("Starting.")
	// defer file.Close()

	prepareRoutes()

	log.Fatal(http.ListenAndServe(":8087", nil))
}

func testCommand() {
	command := "push the bed towards the door."

	var msg llm.UserMessage = llm.UserMessage{
		Content: command,
	}
	result, err := msg.DmProcess()
	if err != nil {
		log.Fatalf("error categorising message: %v", err)
	}

	log.Println("Message:", command)
	log.Println("Result:", result)
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
