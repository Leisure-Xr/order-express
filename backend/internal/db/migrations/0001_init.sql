CREATE TABLE IF NOT EXISTS schema_migrations (
  version    TEXT PRIMARY KEY,
  applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id         TEXT PRIMARY KEY,
  username   TEXT NOT NULL UNIQUE,
  name       TEXT NOT NULL,
  password   TEXT NOT NULL,
  role       TEXT NOT NULL DEFAULT 'customer',
  phone      TEXT,
  avatar     TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
  id         TEXT PRIMARY KEY,
  name       TEXT NOT NULL,
  icon       TEXT,
  image      TEXT,
  sort_order INTEGER NOT NULL DEFAULT 0,
  status     TEXT NOT NULL DEFAULT 'active',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dishes (
  id               TEXT PRIMARY KEY,
  category_id      TEXT NOT NULL,
  name             TEXT NOT NULL,
  description      TEXT NOT NULL,
  price            DOUBLE PRECISION NOT NULL,
  original_price   DOUBLE PRECISION,
  image            TEXT NOT NULL,
  images           TEXT,
  status           TEXT NOT NULL DEFAULT 'on_sale',
  options          TEXT NOT NULL DEFAULT '[]',
  tags             TEXT NOT NULL DEFAULT '[]',
  preparation_time INTEGER,
  created_at       TEXT NOT NULL,
  updated_at       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_dishes_category_id ON dishes(category_id);

CREATE TABLE IF NOT EXISTS orders (
  id                      TEXT PRIMARY KEY,
  order_number            TEXT NOT NULL UNIQUE,
  type                    TEXT NOT NULL,
  status                  TEXT NOT NULL DEFAULT 'pending',
  table_id                TEXT,
  items                   TEXT NOT NULL DEFAULT '[]',
  subtotal                DOUBLE PRECISION NOT NULL DEFAULT 0,
  delivery_fee            DOUBLE PRECISION NOT NULL DEFAULT 0,
  discount                DOUBLE PRECISION NOT NULL DEFAULT 0,
  total                   DOUBLE PRECISION NOT NULL DEFAULT 0,
  remarks                 TEXT,
  payment_method          TEXT NOT NULL DEFAULT 'cash',
  payment_status          TEXT NOT NULL DEFAULT 'unpaid',
  payment_paid_at         TEXT,
  delivery_address        TEXT,
  contact_phone           TEXT,
  estimated_delivery_time TEXT,
  status_history          TEXT NOT NULL DEFAULT '[]',
  created_at              TEXT NOT NULL,
  updated_at              TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

CREATE TABLE IF NOT EXISTS payments (
  payment_id TEXT PRIMARY KEY,
  order_id   TEXT NOT NULL,
  method     TEXT NOT NULL,
  status     TEXT NOT NULL DEFAULT 'processing',
  amount     DOUBLE PRECISION NOT NULL,
  created_at TEXT NOT NULL,
  paid_at    TEXT
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);

CREATE TABLE IF NOT EXISTS tables (
  id               TEXT PRIMARY KEY,
  number           TEXT NOT NULL UNIQUE,
  seats            INTEGER NOT NULL DEFAULT 4,
  status           TEXT NOT NULL DEFAULT 'available',
  current_order_id TEXT,
  qr_code_url      TEXT,
  area             TEXT,
  created_at       TEXT NOT NULL,
  updated_at       TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tables_current_order_id ON tables(current_order_id);

CREATE TABLE IF NOT EXISTS store_info (
  id          INTEGER PRIMARY KEY,
  name        TEXT NOT NULL,
  address     TEXT NOT NULL,
  phone       TEXT NOT NULL,
  logo        TEXT NOT NULL,
  description TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS business_hours (
  id          SERIAL PRIMARY KEY,
  day_of_week INTEGER NOT NULL UNIQUE,
  open_time   TEXT NOT NULL DEFAULT '',
  close_time  TEXT NOT NULL DEFAULT '',
  is_closed   BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS delivery_settings (
  id                     INTEGER PRIMARY KEY,
  enabled                BOOLEAN NOT NULL DEFAULT TRUE,
  minimum_order          DOUBLE PRECISION NOT NULL DEFAULT 0,
  delivery_fee           DOUBLE PRECISION NOT NULL DEFAULT 0,
  free_delivery_threshold DOUBLE PRECISION NOT NULL DEFAULT 0,
  estimated_minutes      INTEGER NOT NULL DEFAULT 30,
  delivery_radius        DOUBLE PRECISION NOT NULL DEFAULT 5,
  updated_at             TEXT NOT NULL
);

