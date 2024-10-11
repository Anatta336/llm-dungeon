package dm

import (
	"encoding/json"
	"log"
	"net/http"
)

type inputRequest struct {
	Content string `json:"content"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func ReceiveInputHandler(writer http.ResponseWriter, request *http.Request) error {

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return http.ErrNotSupported
	}

	var input inputRequest
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		onError(writer, err)
		return err
	}

	// TODO: Provide live world state to the LLMs.

	dmResult, err := Process(input.Content)
	if err != nil {
		onError(writer, err)
		return err
	}

	log.Println("Adjudicate:", dmResult.RawAdjudicate)
	log.Println("Encode:", dmResult.RawActionEncode)

	// TODO: Apply encoded actions to the game state.

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(dmResult)
	if err != nil {
		onError(writer, err)
		return err
	}

	return nil
}

func onError(writer http.ResponseWriter, err error) {
	errorResponse := &errorResponse{
		Error: err.Error(),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(errorResponse)
}
