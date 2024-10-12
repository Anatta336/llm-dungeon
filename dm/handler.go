package dm

import (
	"encoding/json"
	"log"
	"net/http"
	"samdriver/dungeon/world"
)

type inputRequest struct {
	Content string `json:"content"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func ReceiveInputHandler(
	state world.State,
	writer http.ResponseWriter,
	request *http.Request,
) (*world.State, error) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return nil, http.ErrNotSupported
	}

	var input inputRequest
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		onError(writer, err)
		return nil, err
	}

	dmResult, err := Process(state, input.Content)
	if err != nil {
		onError(writer, err)
		return nil, err
	}

	log.Println("Adjudicate:", dmResult.RawAdjudicate)
	log.Println("Encode:", dmResult.RawEncode)

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(dmResult)
	if err != nil {
		onError(writer, err)
		return nil, err
	}

	return &dmResult.OutputState, nil
}

func onError(writer http.ResponseWriter, err error) {
	errorResponse := &errorResponse{
		Error: err.Error(),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(errorResponse)
}
