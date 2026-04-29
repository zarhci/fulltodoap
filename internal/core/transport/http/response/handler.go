package core_response_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/zarhci/fulltodoap/internal/core/errors"
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

func (h *ResponseHandler) JSONResponse(
	response any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("write JSON response", zap.Error(err))
	}
}

func (h *ResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *ResponseHandler) ErrorResponse(err error, msg string) {

	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)

}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)

}

func (h *ResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
