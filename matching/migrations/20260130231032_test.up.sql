CREATE TABLE IF NOT EXISTS public.likes
(
    user_id text COLLATE pg_catalog."default" NOT NULL,
    liked_id text COLLATE pg_catalog."default" NOT NULL,
    UNIQUE (user_id, liked_id)
);