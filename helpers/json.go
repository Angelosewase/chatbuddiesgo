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

func RespondWithError(res http.ResponseWriter, req *http.Request, code int, errorMessage error)  {
	if code < 500 {
		http.Error(res, errorMessage.Error(), http.StatusBadRequest)
	}else{
		http.Error(res, errorMessage.Error(), http.StatusInternalServerError)
	}



}


