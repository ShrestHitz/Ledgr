-- Users table (for auth)
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Accounts table (Cash, Savings, Allowance, etc.)
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "name" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "currency" varchar NOT NULL DEFAULT 'INR',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Categories table
CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "type" varchar NOT NULL -- 'income' or 'expense'
);

-- Entries table (every transaction on an account)
CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "category_id" bigint,
  "amount" bigint NOT NULL, -- can be negative (expense) or positive (income)
  "note" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Transfers table (moving money between your own accounts)
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL CHECK ("amount" > 0),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Budgets table (monthly limit per category)
CREATE TABLE "budgets" (
  "id" bigserial PRIMARY KEY,
  "category_id" bigint NOT NULL,
  "monthly_limit" bigint NOT NULL,
  "month" int NOT NULL, -- 1-12
  "year" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Savings goals table
CREATE TABLE "savings_goals" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "target_amount" bigint NOT NULL,
  "current_amount" bigint NOT NULL DEFAULT 0,
  "target_date" date,
  "linked_account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Recurring payments table
CREATE TABLE "recurring_payments" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "frequency" varchar NOT NULL, -- 'monthly' or 'weekly'
  "next_due_date" date NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Foreign Keys

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "entries" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "budgets" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "savings_goals" ADD FOREIGN KEY ("linked_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "recurring_payments" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

-- Indexes (for common query patterns)

CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "entries" ("category_id");
CREATE INDEX ON "entries" ("created_at");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "budgets" ("category_id", "month", "year");
CREATE INDEX ON "recurring_payments" ("next_due_date");