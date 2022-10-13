-- +migrate Up
-- Table: public.user

CREATE TABLE IF NOT EXISTS public.user
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    username character varying(45) UNIQUE NOT NULL ,
    password_hash_hint character varying(8) NOT NULL,
    check_hash character varying(128) NOT NULL ,
    customer_id integer UNIQUE NOT NULL,
    CONSTRAINT user_id PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.user
    OWNER to postgres;

INSERT INTO public.user(
    username, password_hash_hint, check_hash, customer_id)
VALUES ('Derek', '4cbad12e', '296fd6d505f3ddf41f550a754a27541d754295fe1c125f7805e349f1d94d5330', 1);

CREATE TABLE IF NOT EXISTS public.jwt
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    user_id integer unique,
    jwt text,
    expiration_time timestamp,
    CONSTRAINT jwt_id PRIMARY KEY (id),
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES public.user (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.jwt
    OWNER to postgres;

-- +migrate Down
DROP TABLE IF EXISTS public.jwt;
DROP TABLE IF EXISTS public.user;