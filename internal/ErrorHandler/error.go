package ErrorHandler

import (
	"encoding/json"
	"hot-coffee/models"
	"net/http"
)

func Error(w http.ResponseWriter, ErrorText string, code int) {
	Error := models.Error{
		Code:    code,
		Message: ErrorText,
	}
	jsondata, err := json.MarshalIndent(Error, "", "    ")
	if err != nil {
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsondata)
}
