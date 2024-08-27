-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    id serial,
    login text NOT NULL,
    create_date timestamp WITHOUT time zone NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    name text NOT NULL,
    surname text NOT NULL,
    phone text NOT NULL,
    email text NOT NULL,
    CONSTRAINT pk_user PRIMARY KEY (id),
    CONSTRAINT u_login UNIQUE (login),
    CONSTRAINT u_phone UNIQUE (phone),
    CONSTRAINT u_email UNIQUE (email)
);

ALTER TABLE IF EXISTS public.users
    OWNER to pguser;

CREATE UNIQUE INDEX ind_login
ON public.users (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
