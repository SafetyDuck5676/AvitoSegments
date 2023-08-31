-- Adminer 4.8.1 PostgreSQL 15.3 (Debian 15.3-1.pgdg120+1) dump

\connect "postgres";

CREATE SEQUENCE eventlog_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."eventlog" (
    "id" integer DEFAULT nextval('eventlog_id_seq') NOT NULL,
    "action" character varying(10) NOT NULL,
    "created_at" timestamp NOT NULL,
    "user_id" integer NOT NULL,
    "segment_id" integer NOT NULL,
    CONSTRAINT "eventlog_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE SEQUENCE segment_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."segment" (
    "id" integer DEFAULT nextval('segment_id_seq') NOT NULL,
    "slug" character varying(255) NOT NULL,
    CONSTRAINT "segment_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE TABLE "public"."segment_user" (
    "segment_id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "created_at" timestamp NOT NULL,
    "ttl" timestamp
) WITH (oids = false);


CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "user_name" character varying NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "users" ("id", "user_name") VALUES
(1,	'Ooti'),
(2,	'test'),
(3,	'test2'),
(4,	'test3'),
(5,	'test4'),
(6,	'test5'),
(7,	'test6'),
(8,	'test7'),
(9,	'test8'),
(10,	'test9');

ALTER TABLE ONLY "public"."segment_user" ADD CONSTRAINT "segment_user_segment_id_fkey" FOREIGN KEY (segment_id) REFERENCES segment(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."segment_user" ADD CONSTRAINT "segment_user_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE;

-- 2023-08-30 19:42:32.346069+00