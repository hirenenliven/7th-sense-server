-- Migration: parts, their processes, applicable machines and sub programs
--
-- A part belongs to a customer (user_id) and is made in one or more processes. Each
-- process runs one program on any of the customer's applicable machines and may call
-- sub programs. Times are whole seconds.
--
--   parts                  Part No., name, count; production_time = sum of its processes
--   part_processes         one row per process; production_time = setup + cycle + loading
--                          + the cycle times of its sub programs
--   part_process_machines  machines a process may run on (no soft delete, like
--                          user_machine_access: removing one deletes the row)
--   part_sub_programs      sub programs of a process (has_sub_program = true)
--
-- clientele_id is the customer's client the part is made for. There is no clientele
-- table yet, so the column has no foreign key; add it with the Clientele module.
--
-- drawing_file is the stored name of the uploaded drawing PDF (UPLOAD_DIR/part-drawings/),
-- drawing_name the name it was uploaded with.
--
-- Rollback:
--   DROP TABLE IF EXISTS part_sub_programs, part_process_machines, part_processes, parts;

BEGIN;

CREATE TABLE IF NOT EXISTS parts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    clientele_id INTEGER,
    part_no VARCHAR(100) NOT NULL,
    part_name VARCHAR(200) NOT NULL,
    part_count INTEGER NOT NULL DEFAULT 1 CHECK (part_count > 0),
    production_time INTEGER NOT NULL DEFAULT 0 CHECK (production_time >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A customer cannot have the same Part No. and name twice (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS parts_user_part_key
    ON parts (user_id, LOWER(part_no), LOWER(part_name)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_parts_user_id ON parts(user_id);
CREATE INDEX IF NOT EXISTS idx_parts_is_deleted ON parts(is_deleted);

CREATE TRIGGER update_parts_updated_at BEFORE UPDATE ON parts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS part_processes (
    id SERIAL PRIMARY KEY,
    part_id INTEGER NOT NULL REFERENCES parts(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    process VARCHAR(100) NOT NULL,
    program_name VARCHAR(100) NOT NULL,
    drawing_no VARCHAR(100) NOT NULL,
    drawing_revision_no VARCHAR(100) NOT NULL,
    drawing_file VARCHAR(100) NOT NULL,
    drawing_name VARCHAR(255) NOT NULL,
    accuracy_id INTEGER NOT NULL,
    setup_time INTEGER NOT NULL DEFAULT 0 CHECK (setup_time >= 0),
    cycle_time INTEGER NOT NULL CHECK (cycle_time > 0),
    loading_time INTEGER NOT NULL DEFAULT 0 CHECK (loading_time >= 0),
    production_time INTEGER NOT NULL DEFAULT 0 CHECK (production_time >= 0),
    has_sub_program BOOLEAN NOT NULL DEFAULT false,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_part_processes_part_id ON part_processes(part_id);
CREATE INDEX IF NOT EXISTS idx_part_processes_is_deleted ON part_processes(is_deleted);

CREATE TRIGGER update_part_processes_updated_at BEFORE UPDATE ON part_processes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS part_process_machines (
    id SERIAL PRIMARY KEY,
    process_id INTEGER NOT NULL REFERENCES part_processes(id) ON DELETE CASCADE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT part_process_machines_process_machine_key UNIQUE (process_id, machine_id)
);

-- process_id lookups use the unique index above; this one serves "which parts run on machine X"
CREATE INDEX IF NOT EXISTS idx_part_process_machines_machine_id ON part_process_machines(machine_id);

CREATE TRIGGER update_part_process_machines_updated_at BEFORE UPDATE ON part_process_machines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS part_sub_programs (
    id SERIAL PRIMARY KEY,
    process_id INTEGER NOT NULL REFERENCES part_processes(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    name VARCHAR(200) NOT NULL,
    cycle_time INTEGER NOT NULL DEFAULT 0 CHECK (cycle_time >= 0),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_part_sub_programs_process_id ON part_sub_programs(process_id);
CREATE INDEX IF NOT EXISTS idx_part_sub_programs_is_deleted ON part_sub_programs(is_deleted);

CREATE TRIGGER update_part_sub_programs_updated_at BEFORE UPDATE ON part_sub_programs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
