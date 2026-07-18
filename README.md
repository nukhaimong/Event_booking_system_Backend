# GoTickets

Event ticket booking API built with Go.

## Features

- **User Authentication** — Register/login with JWT tokens (7-day expiry)
- **Event Management** — Create, update, and list events with photo uploads via Cloudinary
- **Ticket Booking** — Book tickets with real-time availability tracking
- **Stripe Payments** — Checkout sessions for ticket purchases with webhook-based confirmation
- **Image Uploads** — Cloudinary integration for event photos

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.26 |
| Framework | Echo v5 |
| ORM | GORM |
| Database | PostgreSQL |
| Payments | Stripe |
| Image Storage | Cloudinary |
| Auth | JWT (HS256) |

## Project Structure

```
cmd/
  main.go                 — Entry point
internal/
  auth/jwt.go             — JWT generation & validation
  config/                 — Env loading, DB connection, Cloudinary setup
  domain/
    user/                 — Auth, user registration
    event/                — Event CRUD with photo uploads
    bookings/             — Ticket booking & payment flow
  httpResponse/           — Standardized error responses
  middlewares/            — Auth middleware
  payment/                — Stripe service & webhook handler
  server/http.go          — Echo server setup & route registration
  utils/                  — Helpers (e.g., get user ID from context)
```

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL
- [Stripe account](https://stripe.com) (for payment processing)
- [Cloudinary account](https://cloudinary.com) (for image uploads)

### Setup

1. Clone the repo

   ```bash
   git clone https://github.com/your-username/go-ticket.git
   cd go-ticket
   ```

2. Create a `.env` file in the project root:

   ```env
   PORT=8080
   DSN=postgres://user:password@localhost:5432/gotickets?sslmode=disable
   JWT_SECRET=your-secret-key
   FRONTEND_URL=http://localhost:3000
   STRIPE_SECRET_KEY=sk_test_...
   STRIPE_WEBHOOK_SECRET=whsec_...
   STRIPE_SUCCESS_URL=http://localhost:3000/success
   STRIPE_CANCEL_URL=http://localhost:3000/cancel
   ```

   > `CLOUDINARY_URL` is auto-read by the Cloudinary SDK. Set it in your environment or `.env`:
   > `CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name`

3. Run the server

   ```bash
   go run ./cmd/main.go
   ```

   Or use [Air](https://github.com/air-verse/air) for hot reload:

   ```bash
   air
   ```

4. Test Stripe webhooks locally

   ```bash
   stripe listen --forward-to localhost:8080/webhook/stripe
   ```

## API Endpoints

All endpoints are prefixed with `/api/v1` unless noted.

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | No | Register a new user |
| POST | `/api/v1/auth/login` | No | Login, returns JWT |
| GET | `/api/v1/auth/me` | Yes | Get current user |
| GET | `/api/v1/events` | No | List all events |
| GET | `/api/v1/events/:id` | No | Get event by ID |
| GET | `/api/v1/events/my-events` | Yes | Get current user's events |
| POST | `/api/v1/events/create` | Yes | Create event (multipart/form-data) |
| PATCH | `/api/v1/events/:id` | Yes | Update event (multipart/form-data, owner only) |
| POST | `/api/v1/bookings` | Yes | Create booking, returns Stripe checkout URL |
| GET | `/api/v1/bookings` | Yes | Get current user's bookings |
| GET | `/api/v1/bookings/:id` | Yes | Get booking by ID |
| POST | `/webhook/stripe` | No | Stripe webhook (payment confirmation) |

> **Note:** Event create/update use `multipart/form-data` (not JSON) because they accept a `photo` file upload.

## Booking Flow

1. `POST /api/v1/bookings` — creates a pending booking and returns a Stripe checkout URL
2. User completes payment on Stripe
3. Stripe fires `checkout.session.completed` webhook → booking confirmed, tickets deducted
4. If payment expires, `checkout.session.expired` webhook deletes the booking
