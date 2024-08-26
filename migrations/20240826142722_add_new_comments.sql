-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.new_comments
(
    subsription_id integer NOT NULL,
    comment_id integer NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    CONSTRAINT pk_new_comment PRIMARY KEY (subsription_id, comment_id),
    CONSTRAINT fk_comment_new_comment FOREIGN KEY (comment_id)
        REFERENCES public.comments (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT fk_post_new_comment FOREIGN KEY (subsription_id)
        REFERENCES public.subscriptions (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.new_comments
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.new_comments;
-- +goose StatementEnd
