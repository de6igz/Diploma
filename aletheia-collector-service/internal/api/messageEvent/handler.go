package messageEvent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Handler – отвечает за приём HTTP-запросов и передачу в UseCase.
type Handler struct {
	useCase UseCase
}

// NewHandler – конструктор HTTP-обработчика.
func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc}
}

// SendMessageEvent – метод-эндпоинт, например, POST /api/v1/message
func (h *Handler) SendMessageEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, MessageResponse{
			Success:      false,
			ErrorMessage: "Only POST is allowed",
		})
		return
	}

	var evt MessageEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		errMsg := fmt.Sprintf("failed to parse message event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusBadRequest, MessageResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	// Если нужен userID из заголовков
	userID := r.Header.Get("X-User-Id")

	if err := h.useCase.ProcessMessageEvent(evt, userID); err != nil {
		errMsg := fmt.Sprintf("failed to process message event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusInternalServerError, MessageResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{Success: true})
}

// writeJSON – хелпер для ответа в формате JSON.
func writeJSON(w http.ResponseWriter, status int, resp MessageResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
