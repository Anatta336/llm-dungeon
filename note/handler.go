package note

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNoteHandler(writer http.ResponseWriter, request *http.Request) {

	idStr := request.URL.Path[len("/note/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		onError(writer, err)
		return
	}

	note, err := find(id)

	if err != nil {
		onError(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(note)

	if err != nil {
		onError(writer, err)
		return
	}
}

func GetAllNotesHandler(writer http.ResponseWriter, request *http.Request) {
	notes, err := all()

	if err != nil {
		onError(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(notes)

	if err != nil {
		onError(writer, err)
		return
	}
}

func onError(writer http.ResponseWriter, err error) {
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(errorResponse)
}
