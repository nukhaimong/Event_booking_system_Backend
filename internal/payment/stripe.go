package payment

import (
	"fmt"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type CreateCheckoutRequest struct {
	BookingID uint  `json:"booking_id"`
	UserID    uint  `json:"user_id"`
	EventID   uint  `json:"event_id"`
	Amount    int64 `json:"amount"` // in cents
	Quantity  int   `json:"quantity"`
}

type CheckoutResponse struct {
	URL       string `json:"url"`
	SessionID string `json:"session_id"`
	BookingID uint   `json:"booking_id"`
}

type StripeService struct {
	secretKey  string
	successURL string
	cancelURL  string
}

func NewStripeService(successURL, cancelURL, secretKey string) *StripeService {
	if secretKey == "" {
		println("⚠️ WARNING: STRIPE_SECRET_KEY is empty!")
	}

	return &StripeService{
		secretKey:  secretKey,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
}

func (s *StripeService) CreateCheckoutSession(req *CreateCheckoutRequest) (*CheckoutResponse, error) {
	if s.secretKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY not set")
	}

	stripe.Key = s.secretKey

	// Amount in cents
	amountInCents := req.Amount * 100

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("Event Tickets x%d", req.Quantity)),
					},
					UnitAmount: stripe.Int64(amountInCents),
				},
				Quantity: stripe.Int64(int64(req.Quantity)),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}", s.successURL)),
		CancelURL:  stripe.String(fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}", s.cancelURL)),
		Metadata: map[string]string{
			"booking_id": fmt.Sprintf("%d", req.BookingID),
			"user_id":    fmt.Sprintf("%d", req.UserID),
			"event_id":   fmt.Sprintf("%d", req.EventID),
		},
	}

	session, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &CheckoutResponse{
		URL:       session.URL,
		SessionID: session.ID,
		BookingID: req.BookingID,
	}, nil
}
