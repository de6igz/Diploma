package errorEvent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) SendErrorEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Success:      false,
			ErrorMessage: "Only POST is allowed",
		})
		return
	}

	var evt ErrorEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		errMsg := fmt.Sprintf("failed to parse error event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	userID := r.Header.Get("X-User-Id")

	if err := h.useCase.ProcessErrorEvent(evt, userID); err != nil {
		errMsg := fmt.Sprintf("failed to process error event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	writeJSON(w, http.StatusOK, ErrorResponse{Success: true})
}

func writeJSON(w http.ResponseWriter, status int, resp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
