# Settings API

Restaurant settings are scoped to the current single-restaurant deployment and are persisted in PostgreSQL when `DATABASE_URL` is configured. The development fallback uses in-memory storage.

## Endpoints

- `GET /api/v1/settings` returns the current configuration.
- `PUT /api/v1/settings` replaces the restaurant name, timezone, operational-alerts preference, and auto-advance preference.

Both endpoints are protected by the API's configured authentication middleware.

## Persistence

Migration `0008_settings.sql` creates the singleton `restaurant_settings` row. The server supplies defaults when no row has been saved yet.

`auto_advance_tickets` is intentionally a persisted preference only. The automation scheduler does not consume it until that capability is introduced in a later milestone.
