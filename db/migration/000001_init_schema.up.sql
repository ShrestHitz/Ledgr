-- ==========================================
-- Ledgr — Init Schema (UP)
-- ==========================================

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "name" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "currency" varchar NOT NULL DEFAULT 'INR',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar,
  "name" varchar NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "category_id" bigint,
  "amount" bigint NOT NULL,
  "note" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL CHECK ("amount" > 0),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "budgets" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "category_id" bigint NOT NULL,
  "monthly_limit" bigint NOT NULL,
  "month" int NOT NULL CHECK ("month" BETWEEN 1 AND 12),
  "year" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "savings_goals" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "target_amount" bigint NOT NULL,
  "current_amount" bigint NOT NULL DEFAULT 0,
  "target_date" date,
  "linked_account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "recurring_payments" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "frequency" varchar NOT NULL CHECK ("frequency" IN ('daily', 'weekly', 'monthly', 'yearly')),
  "next_due_date" date NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Foreign Keys
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "categories" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "entries" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "budgets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "budgets" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");
ALTER TABLE "savings_goals" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "savings_goals" ADD FOREIGN KEY ("linked_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "recurring_payments" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "recurring_payments" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

-- Indexes
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "categories" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "entries" ("category_id");
CREATE INDEX ON "entries" ("created_at");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "budgets" ("owner");
CREATE INDEX ON "budgets" ("category_id", "month", "year");
CREATE INDEX ON "savings_goals" ("owner");
CREATE INDEX ON "recurring_payments" ("owner");
CREATE INDEX ON "recurring_payments" ("next_due_date");

COMMENT ON COLUMN "entries"."amount" IS 'negative = expense, positive = income; stored in smallest unit (paise)';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive; stored in smallest unit (paise)';
COMMENT ON COLUMN "budgets"."monthly_limit" IS 'stored in smallest unit (paise)';
COMMENT ON COLUMN "categories"."owner" IS 'NULL means system/global category';
