-- 000001_init_schema.down.sql
-- Rollback: drop everything in reverse dependency order.

DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS price_categories;
DROP TABLE IF EXISTS event_sessions;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS event_categories;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS venue_halls;
DROP TABLE IF EXISTS venues;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS seating_type;
DROP TYPE IF EXISTS event_session_status;
DROP TYPE IF EXISTS event_status;
DROP TYPE IF EXISTS ticket_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS user_role;
