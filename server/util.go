package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
)

/////////////////////////////////////////////////////////////

func getRequestVar(varname string, r *http.Request) string {
	return chi.URLParam(r, varname)
}

func writeJSON(w http.ResponseWriter, contype string, content interface{}) *serverError {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	
	// if content == nil {
	// 	content = make([]string, 0)
	// }

	fmt.Printf("content: %+v", content)

	encodedContent, err := json.Marshal(content)
	if err != nil {
		return serverErrorInternal(err, ErrMsgEncoding)
	}
	writeResponse(w, contype, encodedContent)
	return nil
}

func writeResponse(w http.ResponseWriter, contype string, encodedContent []byte) {
	w.Header().Set("Content-Type", contype)
	w.WriteHeader(http.StatusOK)
	w.Write(encodedContent)
}

func writeError(w http.ResponseWriter, status int, msg string) {

	w.WriteHeader(status)

	result, err := json.Marshal(struct {
		Status      int    `json:"status"`
		Description string `json:"description"`
	}{
		Status:      status,
		Description: msg,
	})

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error fallback: unable to marshal error message: %v", msg)))
		return
	}

	w.Write(result)
}
