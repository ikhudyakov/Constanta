CREATE TABLE transactions (
		id bigint NOT NULL DEFAULT nextval('id_sequence'::regclass),
    userid bigint,
    useremail character varying(50) COLLATE pg_catalog."default",
    amount real,
    currency character(3) COLLATE pg_catalog."default",
    initdate character varying COLLATE pg_catalog."default",
    moddate character varying COLLATE pg_catalog."default",
    status character varying(7) COLLATE pg_catalog."default",
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
	);
