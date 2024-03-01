package exception

import (
	"github.com/go-errors/errors"
	"net/http"
)

// AppException Модель данных для описания ошибки
type AppException struct {
	Code    int
	Error   error
	Context map[string]any
}

func NewAppException(code int, err error, context map[string]any) *AppException {
	return &AppException{code, err, context}
}

func NewValidationAppException(context map[string]any) *AppException {
	return &AppException{http.StatusUnprocessableEntity, errors.New("VALIDATION ERROR"), context}
}
