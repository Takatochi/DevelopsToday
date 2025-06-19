package v1

import (
	"DevelopsToday/internal/controller/http/v1/cat"
	"DevelopsToday/pkg/logger"
)

type V1 struct {
	cat     *cat.Handler
	mission *cat.Handler
	target  *cat.Handler
	l       logger.Interface `json:"l,omitempty"`
}
