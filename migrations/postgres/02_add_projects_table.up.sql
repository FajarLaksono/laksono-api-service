/*
 * Copyright (c) 2022 Fajar Laksono Inc. All Rights Reserved.
 */

CREATE TABLE IF NOT EXISTS projects(
    id UUID NOT NULL,
    name VARCHAR(32) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_overlapping BOOLEAN DEFAULT FALSE,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS index_projects_start_date ON projects USING btree (start_date);
CREATE INDEX IF NOT EXISTS index_projects_end_date ON projects USING btree (end_date);
