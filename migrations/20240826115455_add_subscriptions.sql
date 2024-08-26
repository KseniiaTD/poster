-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.subscriptions
(
    id serial,
    post_id integer NOT NULL,
    user_id integer NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    upd_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    CONSTRAINT pk_subscription PRIMARY KEY (id) ,
    CONSTRAINT fk_user_subscr FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT fk_post_subscr FOREIGN KEY (post_id)
        REFERENCES public.posts (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

CREATE UNIQUE INDEX sub_post_user 
ON public.subscriptions (post_id, user_id);

ALTER TABLE IF EXISTS public.subscriptions
    OWNER to pguser;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.subscriptions;
-- +goose StatementEnd
