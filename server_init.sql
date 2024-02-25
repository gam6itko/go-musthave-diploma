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
