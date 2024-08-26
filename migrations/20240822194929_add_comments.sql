-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.comments
(
    id serial,
    parent_id integer,
    post_id integer NOT NULL,
    author_id integer NOT NULL,
    body text NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    upd_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    CONSTRAINT pk_comment PRIMARY KEY (id),
    CONSTRAINT fk_post_comment FOREIGN KEY (post_id)
        REFERENCES public.posts (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID,
	CONSTRAINT fk_post_post FOREIGN KEY (post_id)
        REFERENCES public.comments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT fk_user_comment FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE IF EXISTS public.comments
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.comments;
-- +goose StatementEnd
