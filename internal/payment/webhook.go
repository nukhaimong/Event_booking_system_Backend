package payment

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type WebhookHandler struct {
	secret         string
	bookingService BookingService
}

type BookingService interface {
	HandlePaymentSuccess(bookingID string) error
	HandlePaymentExpired(bookingID string) error
}

func NewWebhookHandler(service BookingService, webhookSecret string) *WebhookHandler {
	return &WebhookHandler{
		secret:         webhookSecret,
		bookingService: service,
	}
}

func (w *WebhookHandler) HandleWebhook(c *echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
	}

	signature := c.Request().Header.Get("Stripe-Signature")
	if signature == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing signature"})
	}

	if w.secret == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Webhook secret not configured"})
	}

	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		w.secret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid signature",
			"details": err.Error(),
		})
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse session"})
		}

		bookingID := session.Metadata["booking_id"]

		if w.bookingService != nil {
			if err := w.bookingService.HandlePaymentSuccess(bookingID); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse session"})
		}

		bookingID := session.Metadata["booking_id"]

		if w.bookingService != nil {
			if err := w.bookingService.HandlePaymentExpired(bookingID); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
