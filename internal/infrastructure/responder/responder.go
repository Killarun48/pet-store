package responder

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Code    int    `json:"code" example:"200"`
	Type    string `json:"type" example:"any"`
	Message string `json:"message,omitempty" example:"user deleted"`
}

type Responder interface {
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorNotFound(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, message string)
}

type Respond struct{}

func NewResponder() Responder {
	return &Respond{}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(Response{
		Code:    http.StatusBadRequest,
		Type:    "unknown",
		Message: err.Error(),
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}

func (r *Respond) ErrorNotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(Response{
		Code:    http.StatusNotFound,
		Type:    "unknown",
		Message: err.Error(),
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}

func (r *Respond) Success(w http.ResponseWriter, message string) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(Response{
		Code:    http.StatusOK,
		Type:    "unknown",
		Message: message,
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}

/* func (r *Respond) Success(w http.ResponseWriter, message string) {
	//w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusOK,
		Type:    "unknown",
		Message: message,
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
} */
