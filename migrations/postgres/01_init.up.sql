/*
 * Copyright (c) 2022 Fajar Laksono Inc. All Rights Reserved.
 */

CREATE TABLE IF NOT EXISTS users(
    id UUID NOT NULL,
    username VARCHAR(32) NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL, 
    PRIMARY KEY(id),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);