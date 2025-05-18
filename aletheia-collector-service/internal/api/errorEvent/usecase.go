package errorEvent

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
)

type UseCase interface {
	ProcessErrorEvent(evt ErrorEvent, userID string) error
}

type useCaseImpl struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) ProcessErrorEvent(evt ErrorEvent, userID string) error {
	key, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate key for error event: %v", err)
	}

	value, err := BuildKafkaMessage(evt, userID)
	if err != nil {
		return fmt.Errorf("failed to build message for error event: %v", err)
	}

	partition, offset, err := uc.repo.SendErrorEvent([]byte(key), value)
	if err != nil {
		return err
	}

	fmt.Printf("[ErrorEvent] Sent to topic partition=%d, offset=%d\n", partition, offset)
	return nil
}
