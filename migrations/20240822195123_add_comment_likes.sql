-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.comment_likes
(
    author_id integer NOT NULL,
    comment_id integer NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    is_like boolean NOT NULL DEFAULT true,
    CONSTRAINT pk_user_comment PRIMARY KEY (author_id, comment_id),
    CONSTRAINT fk_user_comment_like FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT fk_comment_like FOREIGN KEY (comment_id)
        REFERENCES public.comments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.comment_likes
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.comment_likes;
-- +goose StatementEnd
