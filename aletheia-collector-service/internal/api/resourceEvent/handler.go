package resourceEvent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Handler – слой HTTP для ResourceEvent.
type Handler struct {
	useCase UseCase
}

// NewHandler – конструктор Handler.
func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc}
}

// SendResourceEvent – HTTP-эндпоинт для приёма ResourceEvent.
func (h *Handler) SendResourceEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ResourceResponse{
			Success:      false,
			ErrorMessage: "Only POST is allowed",
		})
		return
	}

	var evt ResourceEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		errMsg := fmt.Sprintf("failed to parse resource event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusBadRequest, ResourceResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	userID := r.Header.Get("X-User-Id")

	if err := h.useCase.ProcessResourceEvent(evt, userID); err != nil {
		if err.Error() == rateLimitError {
			//rate limit exceeded for userID: %s and serviceName: %s\n available %d requests in min", userID, evt.ServiceName, uc.maxRequestsPerMinute)
			writeJSON(w, http.StatusBadRequest, ResourceResponse{Success: false,
				ErrorMessage: fmt.Sprintf("rate limit exceeded for userID: %s and serviceName: %s\n available %s requests in min", userID, evt.ServiceName, os.Getenv("RATE_LIMIT_MAX_REQUESTS_PER_MINUTE")),
			})
		}
		errMsg := fmt.Sprintf("failed to process resource event: %v", err)
		log.Println(errMsg)
		writeJSON(w, http.StatusInternalServerError, ResourceResponse{Success: false, ErrorMessage: errMsg})
		return
	}

	writeJSON(w, http.StatusOK, ResourceResponse{Success: true})
}

func writeJSON(w http.ResponseWriter, status int, resp ResourceResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
