package v1

import (
	"DevelopsToday/internal/controller/http/v1/cat"
	"DevelopsToday/internal/controller/http/v1/mission"
	"DevelopsToday/internal/controller/http/v1/target"
	"DevelopsToday/pkg/logger"
)

type V1 struct {
	cat     *cat.Handler
	mission *mission.Handler
	target  *target.Handler
	l       logger.Interface `json:"l,omitempty"`
}
