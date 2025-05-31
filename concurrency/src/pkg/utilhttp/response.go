package utilhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseOk[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusOK, msg)
}

func ResponseCreated[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusCreated, msg)
}

func ResponseAccepted[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusAccepted, msg)
}

func ResponseNoContent[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusNoContent, msg)
}

func ResponseUnauthorized[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusUnauthorized, msg)
}

func ResponseForbidden[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusForbidden, msg)
}

func ResponseConflict[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusConflict, msg)
}

func ResponseInternalServerError[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusInternalServerError, msg)
}

func ResponseBadRequest[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusBadRequest, msg)
}

func ResponseNotFound[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusNotFound, msg) // Corrected status code
}

func ResponseMethodNotAllowed[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusMethodNotAllowed, msg)
}

func ResponseNotImplemented[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusNotImplemented, msg)
}

func ResponseServiceUnavailable[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusServiceUnavailable, msg)
}

func ResponseGatewayTimeout[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusGatewayTimeout, msg)
}

func ResponseTooManyRequests[T any](w http.ResponseWriter, msg T) {
	writeResponse(w, http.StatusTooManyRequests, msg)
}

func writeResponse[T any](w http.ResponseWriter, statusCode int, msg T) {
	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // Use direct string
	w.WriteHeader(statusCode)
	w.Write([]byte(res))
}
