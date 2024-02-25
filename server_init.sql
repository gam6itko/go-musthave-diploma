# CREATE TABLE IF NOT EXISTS public.counter
# (
#     "name" varchar NOT NULL,
#     value  bigint  NOT NULL DEFAULT 0,
#     CONSTRAINT counter_pk PRIMARY KEY ("name")
# );
#
# CREATE TABLE IF NOT EXISTS public.gauge
# (
#     "key" varchar          NOT NULL,
#     value double precision NULL,
#     CONSTRAINT gauge_pk PRIMARY KEY ("key")
# );



CREATE TABLE IF NOT EXISTS public.user
(
    id       BIGSERIAL PRIMARY KEY,
    login varchar,
    password varchar
);

CREATE TABLE IF NOT EXISTS public.order
(
    id  int PRIMARY KEY,
    user_id BIGINT,
    status smallint,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES "user"(id)
);
