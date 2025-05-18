// Package interfaces
// @tg version=v1.0.0
// @tg backend=event-service
// @tg title=`Event & Project Service API`
// @tg desc=`API для получения событий, данных пользователя, проектов и правил`
// @tg servers=`http://localhost;local`
package interfaces

import (
	"context"

	"aletheia-public-api/interfaces/types/v1" // поправьте под свой реальный import path
)

// App
// @tg http-server log metrics
// @tg http-prefix=v1
type App interface {

	// GetMe
	// @tg summary=`Получить данные о пользователе`
	// @tg desc=`Возвращает информацию о текущем пользователе`
	// @tg http-method=GET
	// @tg http-path=/me
	// @tg http-header=X-User-Id
	GetMe(ctx context.Context) (resp v1.MeResponse, err error)
}
