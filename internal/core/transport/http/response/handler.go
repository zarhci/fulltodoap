package core_response_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/zarhci/fulltodoap/internal/core/logger"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	rw  http.ResponseWriter
	log *core_logger.Logger
}

func NewResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		rw:  rw,
		log: log,
	}
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statusCode)
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}
	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("failed to write panic response", zap.Error(err))
	}

}
