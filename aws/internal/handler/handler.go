package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	appErrors "github.com/quimnieto/akeneo-serverless-connector/aws/internal/errors"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/logger"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/models"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/processor"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/publisher"
)

type LambdaHandler struct {
	processor processor.EventProcessor
	publisher publisher.SNSPublisher
	logger    logger.Logger
}

func NewLambdaHandler(proc processor.EventProcessor, pub publisher.SNSPublisher, log logger.Logger) *LambdaHandler {
	return &LambdaHandler{
		processor: proc,
		publisher: pub,
		logger:    log,
	}
}

func (h *LambdaHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqLogger := h.logger
	if req.RequestContext.RequestID != "" {
		reqLogger = reqLogger.WithRequestID(req.RequestContext.RequestID)
	}

	payload, decodeErr := h.decodeBody(req.Body, req.IsBase64Encoded)
	if decodeErr != nil {
		appErr := appErrors.ErrInvalidPayload.WithDetails(decodeErr.Error())
		reqLogger.Error("failed to decode request body", decodeErr, nil)
		return h.errorResponse(http.StatusBadRequest, appErr), nil
	}

	if len(payload) == 0 {
		appErr := appErrors.ErrInvalidPayload.WithDetails("request body is empty")
		reqLogger.Error("empty request body", appErr, nil)
		return h.errorResponse(http.StatusBadRequest, appErr), nil
	}

	event, parseErr := h.processor.Parse(payload)
	if parseErr != nil {
		appErr := appErrors.ErrInvalidPayload.WithDetails(parseErr.Error())
		reqLogger.Error("failed to parse payload", parseErr, nil)
		return h.errorResponse(http.StatusBadRequest, appErr), nil
	}

	validateErr := h.processor.Validate(event)
	if validateErr != nil {
		appErr := appErrors.ErrMissingField.WithDetails(validateErr.Error())
		reqLogger.Error("payload validation failed", validateErr, map[string]interface{}{
			"event_id": event.EventID,
		})
		return h.errorResponse(http.StatusBadRequest, appErr), nil
	}

	publishErr := h.publisher.Publish(ctx, event)
	if publishErr != nil {
		appErr := appErrors.ErrSNSPublishFailed.WithDetails(publishErr.Error())
		reqLogger.Error("failed to publish event", publishErr, map[string]interface{}{
			"event_id":   event.EventID,
			"event_type": event.EventType,
		})
		return h.errorResponse(http.StatusInternalServerError, appErr), nil
	}

	reqLogger.Info("event published", map[string]interface{}{
		"event_id":   event.EventID,
		"event_type": event.EventType,
	})

	return h.successResponse(event), nil
}

func (h *LambdaHandler) decodeBody(body string, encoded bool) ([]byte, error) {
	if !encoded {
		return []byte(body), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (h *LambdaHandler) successResponse(event *models.AkeneoEvent) events.APIGatewayProxyResponse {
	payload := map[string]string{
		"status":   "accepted",
		"event_id": event.EventID,
	}

	body, _ := json.Marshal(payload)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (h *LambdaHandler) errorResponse(status int, appErr *appErrors.AppError) events.APIGatewayProxyResponse {
	payload := map[string]string{
		"error_code": appErr.Code,
		"message":    appErr.Message,
	}

	if appErr.Details != "" {
		payload["details"] = appErr.Details
	}

	body, _ := json.Marshal(payload)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
