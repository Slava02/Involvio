BEGIN;

CREATE TABLE "space" (
                         "id" integer PRIMARY KEY,
                         "name" varchar,
                         "description" varchar,
                         "tags" jsonb
);

CREATE TABLE "category_space" (
                                  "category_id" integer,
                                  "space_id" integer,
                                  PRIMARY KEY ("category_id", "space_id")
);

CREATE TABLE "user" (
                        "id" integer PRIMARY KEY,
                        "first_name" varchar,
                        "last_name" varchar,
                        "username" varchar,
                        "photo_url" varchar,
                        "auth_date" timestamp
);

CREATE TABLE "user_space" (
                              "user_id" integer,
                              "space_id" integer,
                              "tags" jsonb,
                              PRIMARY KEY ("user_id", "space_id")
);

CREATE TABLE "user_event" (
                              "user_id" integer,
                              "event_id" integer,
                              PRIMARY KEY ("user_id", "event_id")
);

CREATE TABLE "event" (
                         "id" integer PRIMARY KEY,
                         "name" varchar,
                         "description" varchar,
                         "begin_date" timestamp,
                         "end_date" timestamp
);

ALTER TABLE "user_space" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_space" ADD FOREIGN KEY ("space_id") REFERENCES "space" ("id");

ALTER TABLE "user_event" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_event" ADD FOREIGN KEY ("event_id") REFERENCES "event" ("id");

ALTER TABLE "category_space" ADD FOREIGN KEY ("space_id") REFERENCES "space" ("id");




COMMIT;