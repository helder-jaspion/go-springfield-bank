CREATE TABLE "transfers"
(
    "id"                     uuid PRIMARY KEY,
    "account_origin_id"      uuid        NOT NULL,
    "account_destination_id" uuid        NOT NULL,
    "amount"                 bigint      NOT NULL,
    "created_at"             timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("account_origin_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("account_destination_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "transfers" ("account_origin_id");

CREATE INDEX ON "transfers" ("account_destination_id");

CREATE INDEX ON "transfers" ("account_origin_id", "account_destination_id");