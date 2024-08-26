-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.post_likes
(
    author_id integer NOT NULL,
    post_id integer NOT NULL,
    is_like boolean NOT NULL DEFAULT true,
    created_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    CONSTRAINT pk_post_like PRIMARY KEY (author_id, post_id),
	CONSTRAINT fk_user_post_like FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT fk_post_like FOREIGN KEY (post_id)
        REFERENCES public.posts (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.post_likes
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.post_likes;
-- +goose StatementEnd
