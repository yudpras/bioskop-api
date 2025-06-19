CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "email" varchar,
  "username" varchar,
  "password" varchar,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "cities" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "cinemas" (
  "id" SERIAL PRIMARY KEY,
  "cities_id" int,
  "name" varchar,
  "address" varchar,
  "phone" varchar,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "studios" (
  "id" SERIAL PRIMARY KEY,
  "cinema_id" int,
  "name" varchar,
  "capasity" int,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "movies" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "price" int,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "schedules" (
  "id" SERIAL PRIMARY KEY,
  "studio_id" int,
  "movie_id" int,
  "start_time" time,
  "end_time" time,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "seats" (
  "id" SERIAL PRIMARY KEY,
  "schedule_id" int,
  "lock_by_user_id" int,
  "name" varchar,
  "status" varchar,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "transactions" (
  "id" SERIAL PRIMARY KEY,
  "movie_id" int,
  "user_id" int,
  "total_amount" int,
  "status" varchar,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "detail_transactions" (
  "id" SERIAL PRIMARY KEY,
  "transaction_id" int,
  "seat_id" int,
  "ticket_id" varchar UNIQUE,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "payments" (
  "id" SERIAL PRIMARY KEY,
  "transaction_id" int,
  "user_id" int,
  "payment_id" varchar UNIQUE,
  "payment_method" varchar,
  "payment_date" timestamp,
  "amount" int,
  "status" varchar,
  "log_json" text,
  "create_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

ALTER TABLE "cinemas" ADD FOREIGN KEY ("cities_id") REFERENCES "cities" ("id");

ALTER TABLE "studios" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id");

ALTER TABLE "schedules" ADD FOREIGN KEY ("studio_id") REFERENCES "studios" ("id");

ALTER TABLE "schedules" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id");

ALTER TABLE "seats" ADD FOREIGN KEY ("schedule_id") REFERENCES "schedules" ("id");

ALTER TABLE "seats" ADD FOREIGN KEY ("lock_by_user_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "detail_transactions" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "detail_transactions" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
