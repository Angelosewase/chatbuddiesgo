package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJson(res http.ResponseWriter, req *http.Request, data interface{}, code int) error {
	resdata, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to parse the data %v", err)
	}

	res.WriteHeader(code)
	res.Write(resdata)
	return nil
}

func RespondWithError(res http.ResponseWriter, req *http.Request, code int, errorMessage error) error {
	if code < 500 {
		res.WriteHeader(code)
		http.Error(res, errorMessage.Error(), http.StatusBadRequest)
	}

	res.WriteHeader(code)
	http.Error(res, errorMessage.Error(), http.StatusInternalServerError)
	return nil

}
