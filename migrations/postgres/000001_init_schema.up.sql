CREATE TABLE "accounts"
(
    "id"         uuid PRIMARY KEY,
    "cpf"        char(11)    NOT NULL,
    "name"       varchar     NOT NULL,
    "secret"     varchar     NOT NULL,
    "balance"    bigint      NOT NULL DEFAULT (0),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "accounts" ("cpf");