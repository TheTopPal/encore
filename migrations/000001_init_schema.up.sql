-- 000001_init_schema.up.sql
-- Initial database schema for Encore concert ticket platform.
--
-- Order matters: enums first, then tables from independent to dependent.

-- ==================== ENUMS ====================

CREATE TYPE user_role AS ENUM ('buyer', 'admin', 'organizer');
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'cancelled', 'refunded');
CREATE TYPE ticket_status AS ENUM ('active', 'used', 'cancelled');
CREATE TYPE event_status AS ENUM ('draft', 'published', 'cancelled');
CREATE TYPE event_session_status AS ENUM ('scheduled', 'selling', 'sold_out', 'cancelled');
CREATE TYPE seating_type AS ENUM ('assigned', 'general_admission');

-- ==================== USERS ====================

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    first_name    VARCHAR NOT NULL,
    last_name     VARCHAR NOT NULL,
    username      VARCHAR NOT NULL UNIQUE,
    role          user_role NOT NULL DEFAULT 'buyer',
    phone         VARCHAR,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==================== VENUES ====================

CREATE TABLE venues (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR NOT NULL,
    address      TEXT NOT NULL,
    city         VARCHAR NOT NULL,
    organizer_id UUID NOT NULL REFERENCES users (id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE venue_halls (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id     UUID NOT NULL REFERENCES venues (id),
    name         VARCHAR NOT NULL,
    seating_type seating_type NOT NULL,
    capacity     INT NOT NULL
);

CREATE TABLE seats (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hall_id     UUID NOT NULL REFERENCES venue_halls (id),
    row_number  VARCHAR NOT NULL,
    seat_number INT NOT NULL,
    section     VARCHAR,

    UNIQUE (hall_id, row_number, seat_number)
);

-- ==================== EVENTS ====================

CREATE TABLE event_categories (
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR NOT NULL,
    description     TEXT,
    category_id     UUID NOT NULL REFERENCES event_categories (id),
    organizer_id    UUID NOT NULL REFERENCES users (id),
    image_url       TEXT,
    age_restriction SMALLINT,
    status          event_status NOT NULL DEFAULT 'draft',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE event_sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID NOT NULL REFERENCES events (id),
    hall_id    UUID NOT NULL REFERENCES venue_halls (id),
    starts_at  TIMESTAMPTZ NOT NULL,
    ends_at    TIMESTAMPTZ NOT NULL,
    status     event_session_status NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE price_categories (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id     UUID NOT NULL REFERENCES event_sessions (id),
    name           VARCHAR NOT NULL,
    price          DECIMAL NOT NULL,
    total_quantity INT,
    section        VARCHAR
);

-- ==================== ORDERS & TICKETS ====================

CREATE TABLE orders (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users (id),
    status       order_status NOT NULL DEFAULT 'pending',
    total_amount DECIMAL NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tickets (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id          UUID NOT NULL REFERENCES orders (id),
    session_id        UUID NOT NULL REFERENCES event_sessions (id),
    price_category_id UUID NOT NULL REFERENCES price_categories (id),
    seat_id           UUID REFERENCES seats (id),
    price             DECIMAL NOT NULL,
    status            ticket_status NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (session_id, seat_id)
);
