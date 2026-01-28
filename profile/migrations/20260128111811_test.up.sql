CREATE TABLE IF NOT EXISTS public.profiles
(
    id uuid COLLATE pg_catalog."default" NOT NULL,
    user_id text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    gender text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    date_created date COLLATE pg_catalog."default" NOT NULL,
    photo_path text COLLATE pg_catalog."default",
    CONSTRAINT id PRIMARY KEY (id),
    CONSTRAINT unique_profile UNIQUE (user_id)
);