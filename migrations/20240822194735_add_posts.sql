-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.posts
(
    id serial,
    title text NOT NULL,
    author_id integer NOT NULL,
    body text NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    upd_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    is_commented boolean NOT NULL DEFAULT true,
    CONSTRAINT pk_post PRIMARY KEY (id),
    CONSTRAINT fk_user_post FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID
);


ALTER TABLE IF EXISTS public.posts
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.posts;
-- +goose StatementEnd
