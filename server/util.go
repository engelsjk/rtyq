package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CrunchyData/pg_featureserv/api"
	"github.com/go-chi/chi"
)

func getRequestVar(varname string, r *http.Request) string {
	return chi.URLParam(r, varname)
}

func writeJSON(w http.ResponseWriter, contype string, content interface{}) *serverError {
	encodedContent, err := json.Marshal(content)
	if err != nil {
		return serverErrorInternal(err, api.ErrMsgEncoding)
	}
	writeResponse(w, contype, encodedContent)
	return nil
}

func writeResponse(w http.ResponseWriter, contype string, encodedContent []byte) {
	w.Header().Set("Content-Type", contype)
	w.WriteHeader(http.StatusOK)
	w.Write(encodedContent)
}

func writeError(w http.ResponseWriter, code string, msg string, status int) {

	w.WriteHeader(status)

	result, err := json.Marshal(struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}{
		Code:        code,
		Description: msg,
	})

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error fallback: unable to marshal error message: %v", msg)))
		return
	}

	w.Write(result)
}
