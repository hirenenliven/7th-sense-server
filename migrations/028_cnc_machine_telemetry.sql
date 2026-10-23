-- Migration: CNC machine telemetry (docs/cnc-telemetry-collector-plan.md)
--
-- Replaces 028_add_machine_telementry.sql, 029_machine_variables_key.sql and
-- 030_machine_variables_id.sql. Part 1 removes everything those created (on a fresh database
-- it removes nothing); part 2 creates it again in its final form. Run it once, on any database:
--   psql -U postgres -d 7th_sense -f migrations/028_cnc_machine_telementry.sql
--
-- Data loss: every row in the tables below is deleted (telemetry, live status, mode log,
-- events and per-machine overrides in machine_variables). Variables and modes are seeded again.
-- controllers.cnc_engine is kept as it is (only empty values are filled from the brand name).
--
-- Written by the CNC collectors that the API server starts and supervises itself
-- (cmd/cnc-collector, one process per shard; no Windows service, script or scheduler).
--
-- Configuration (edited by the superadmin; soft delete except machine_variables):
--   machine_modes            the business modes of the old system (machines_modes), same ids
--   cnc_variables            every variable that can be read (old mm_variables), same ids
--   cnc_variable_addresses   where a variable is read, per controller family / version:
--                            a field of the DLL snapshot, a PLC node / address, or an NC address
--   machine_variables        per machine: switch a variable off or override its address / limits;
--                            id is the variable id, key (machine_id, id)
--   controllers.cnc_engine   which DLL engine talks to the controller (1 Fanuc, 2 Siemens, ...)
--
-- Data (no soft delete: append-only or one row per key, written only by the collectors):
--   machine_variable_values  one row per machine-variable: latest numeric value (value only)
--                            upserted on every reading, then used to detect changes; text variables not stored
--   machine_telementry       time series, one row per stored value (recorded_at, machine, variable,
--                            value for numeric variables only); the latest value is its newest row.
--                            TimescaleDB hypertable, compressed after 1 day, kept 1 year
--   machine_live_status     one row per machine: connection, current mode, collector health
--   machine_mode_log        one row per period a machine stays in one mode (machine_modes)
--   machine_events          connect, disconnect, DLL errors, alarms, collector start / stop
--
-- Notifications (the collectors LISTEN, so nothing is polled):
--   machines_changed    '<machine id>'            a machine is added, changed or removed
--   cnc_config_changed  'variables' | 'machine:<id>'  the variable configuration changed
--
-- Rollback: part 1 of this file, plus
--   ALTER TABLE controllers DROP CONSTRAINT IF EXISTS controllers_cnc_engine_check,
--       DROP COLUMN IF EXISTS cnc_engine;

BEGIN;

-- ============================================================================================
-- Part 1: remove what 028 / 029 / 030 created (any version of them)
-- ============================================================================================

DROP TRIGGER IF EXISTS notify_machines_changed ON machines;
DROP TRIGGER IF EXISTS notify_cnc_config_changed ON controllers;   -- first version of 028 only

-- Children before parents. machine_telementry is the first draft of 028.
-- Dropping the hypertable also drops its chunks and its retention / compression jobs.
DROP TABLE IF EXISTS
    machine_events,
    machine_mode_log,
    machine_live_status,
    machine_variable_live,                       -- removed: the latest value is read from machine_telementry
    machine_telementry,
    machine_variable_values,
    machine_variables,
    cnc_variable_addresses,
    cnc_variables,
    machine_modes;

DROP FUNCTION IF EXISTS notify_machines_changed();
DROP FUNCTION IF EXISTS notify_cnc_config_changed();

-- ============================================================================================
-- Part 2: create
-- ============================================================================================

-- TimescaleDB when it is installed and preloaded. Without it machine_telementry is a plain
-- table: fine for development, too slow for production (the collector logs an error at start).
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_available_extensions WHERE name = 'timescaledb')
       AND current_setting('shared_preload_libraries', true) LIKE '%timescaledb%' THEN
        CREATE EXTENSION IF NOT EXISTS timescaledb;
    ELSE
        RAISE NOTICE 'timescaledb not available: machine_telementry will be a plain table';
    END IF;
END
$$;

-- ============================================================================================
-- Configuration
-- ============================================================================================

-- Business modes, with the ids of the old machines_modes table (fixed list, like roles).
-- mode_for: 'machine' = set by the collector from the machine data, 'user' = set by an operator.
CREATE TABLE machine_modes (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mode_for VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT machine_modes_mode_for_check CHECK (mode_for IN ('machine', 'user'))
);

INSERT INTO machine_modes (id, name, mode_for) VALUES
    (1, 'Setup', 'machine'),
    (2, 'Loading / Unloading Time', 'user'),
    (3, 'In Production', 'machine'),
    (4, 'Hold', 'machine'),
    (5, 'Idle', 'machine'),
    (6, 'Downtime', 'machine'),
    (7, 'Planned Downtime', 'user'),
    (8, 'Unfinished Cycle', 'machine'),
    (9, 'Offline', 'machine'),
    (10, 'Resume Cycle', 'machine'),
    (11, 'End Cycle', 'machine'),
    (12, 'Rapid Mode', 'machine'),
    (13, 'Leave Planned Downtime', 'user')
ON CONFLICT (id) DO NOTHING;

-- machines.mode_id (old mm_machines.iModeId) holds one of these ids; it gets no foreign key
-- here because rows copied from the old database are not checked yet.

-- Which DLL engine talks to a controller: 0 simulator, 1 Fanuc, 2 Siemens, 3 Syntec,
-- 4 Mitsubishi. NULL = not set (machines with that controller are not collected).
ALTER TABLE controllers ADD COLUMN IF NOT EXISTS cnc_engine SMALLINT;
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'controllers_cnc_engine_check') THEN
        ALTER TABLE controllers ADD CONSTRAINT controllers_cnc_engine_check
            CHECK (cnc_engine BETWEEN 0 AND 4);
    END IF;
END
$$;

-- First guess from the brand name; the superadmin corrects it on the controller form
UPDATE controllers c SET cnc_engine = CASE
        WHEN b.name ILIKE '%fanuc%' THEN 1
        WHEN b.name ILIKE '%siemens%' OR c.name ILIKE '%sinumerik%' THEN 2
        WHEN b.name ILIKE '%syntec%' THEN 3
        WHEN b.name ILIKE '%mitsubishi%' THEN 4
    END
FROM brands b
WHERE b.id = c.brand_id AND c.cnc_engine IS NULL;

-- Every variable that can be read. ids are the old mm_variables.iVarId (new ones from 200).
--   data_type        number | bool | text (bool is stored as 0 / 1)
--   category         state | process | health | alarm
--   store_mode       old isStore: 0 not stored (read for the mode rules / events only),
--                    1 every change, 2 one sample every min_interval_ms (sensors),
--                    3 changes past the deadband, at most one per min_interval_ms (analog values)
--   deadband         smallest change that is stored (absolute, in the variable's unit);
--                    NULL = any change
--   min_interval_ms  old iInterval; 0 = no limit
CREATE TABLE cnc_variables (
    id INTEGER PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    data_type VARCHAR(10) NOT NULL DEFAULT 'number',
    unit VARCHAR(20),
    category VARCHAR(20) NOT NULL DEFAULT 'process',
    store_mode SMALLINT NOT NULL DEFAULT 1,
    deadband DOUBLE PRECISION,
    min_interval_ms INTEGER NOT NULL DEFAULT 0,
    remarks VARCHAR(250),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT cnc_variables_data_type_check CHECK (data_type IN ('number', 'bool', 'text')),
    CONSTRAINT cnc_variables_category_check CHECK (category IN ('state', 'process', 'health', 'alarm')),
    CONSTRAINT cnc_variables_store_mode_check CHECK (store_mode BETWEEN 0 AND 3),
    CONSTRAINT cnc_variables_deadband_check CHECK (deadband IS NULL OR deadband >= 0),
    CONSTRAINT cnc_variables_min_interval_check CHECK (min_interval_ms >= 0)
);

CREATE UNIQUE INDEX cnc_variables_code_key
    ON cnc_variables (LOWER(code)) WHERE is_deleted = false;
CREATE INDEX idx_cnc_variables_is_deleted ON cnc_variables(is_deleted);

CREATE TRIGGER update_cnc_variables_updated_at BEFORE UPDATE ON cnc_variables
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Where a variable is read, per engine and controller version (old vVarPathv1 = 828D,
-- vVarPathv2 = 840D):
--   source 'snapshot'  a field of the DLL's CncMasterSnapshot, e.g. 'SpindleLoadPercetnage'
--                      (read on every DLL callback; works with the current DLL)
--          'plc'       a PLC node id / address, e.g. 'ns=2;s=/Random/Plc/DB9004.DBD8:REAL[1]'
--          'nc'        an NC variable address, e.g. 'ns=2;s=/Channel/State/actFeedRateIpo[U1,1]'
--                      ('plc' / 'nc' need a DLL read function; until then they stay unused)
--   controller_version_type  1 = 828D, 2 = 840D (machines.controller_version_type), NULL = all
-- The collector picks, per machine and variable: the row for the machine's version before the
-- NULL row, then 'snapshot' before 'nc' before 'plc', skipping sources it cannot read.
CREATE TABLE cnc_variable_addresses (
    id SERIAL PRIMARY KEY,
    variable_id INTEGER NOT NULL REFERENCES cnc_variables(id),
    cnc_engine SMALLINT NOT NULL,
    controller_version_type SMALLINT,
    source VARCHAR(10) NOT NULL,
    address VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT cnc_variable_addresses_engine_check CHECK (cnc_engine BETWEEN 0 AND 4),
    CONSTRAINT cnc_variable_addresses_source_check CHECK (source IN ('snapshot', 'plc', 'nc'))
);

CREATE UNIQUE INDEX cnc_variable_addresses_key
    ON cnc_variable_addresses (variable_id, cnc_engine, COALESCE(controller_version_type, 0), source)
    WHERE is_deleted = false;
CREATE INDEX idx_cnc_variable_addresses_is_deleted ON cnc_variable_addresses(is_deleted);

CREATE TRIGGER update_cnc_variable_addresses_updated_at BEFORE UPDATE ON cnc_variable_addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Per machine changes to the defaults. No row = the variable is read with its defaults.
-- A NULL column keeps the default; source and address are overridden together.
-- id is the variable id (cnc_variables.id); the key is the machine + the variable. No soft
-- delete (same exception as user_machine_access): deleting the row returns the variable to its
-- defaults.
CREATE TABLE machine_variables (
    id INTEGER NOT NULL REFERENCES cnc_variables(id),
    machine_id INTEGER NOT NULL REFERENCES machines(id),
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    source VARCHAR(10),
    address VARCHAR(255),
    deadband DOUBLE PRECISION,
    min_interval_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (machine_id, id),
    CONSTRAINT machine_variables_source_check CHECK (source IS NULL OR source IN ('snapshot', 'plc', 'nc')),
    CONSTRAINT machine_variables_address_check CHECK ((source IS NULL) = (address IS NULL)),
    CONSTRAINT machine_variables_deadband_check CHECK (deadband IS NULL OR deadband >= 0),
    CONSTRAINT machine_variables_min_interval_check CHECK (min_interval_ms IS NULL OR min_interval_ms >= 0)
);

CREATE TRIGGER update_machine_variables_updated_at BEFORE UPDATE ON machine_variables
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Seed: the old mm_variables list (Siemens OPC UA node ids; same id, name, status, isStore,
-- iInterval) plus the values the DLL snapshot already delivers.
CREATE TEMP TABLE seed_cnc_variables (
    id INTEGER, code VARCHAR(50), name VARCHAR(100), data_type VARCHAR(10), unit VARCHAR(20),
    category VARCHAR(20), store_mode SMALLINT, deadband DOUBLE PRECISION, min_interval_ms INTEGER,
    remarks VARCHAR(250), status VARCHAR(20), path_828d VARCHAR(255), path_840d VARCHAR(255)
) ON COMMIT DROP;

INSERT INTO seed_cnc_variables VALUES
    (1, 'ProgLevel', 'Program level', 'number', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/ProgramInfo/displProgLevel', 'ns=2;s=/Channel/ProgramInfo/displProgLevel'),
    (2, 'cycleStart', 'Cycle start', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB3300.DBX3.0', 'ns=2;s=/Random/Plc/DB21.DBX35.0'),
    (3, 'cycleStop', 'Cycle stop', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB3300.DBX3.2', 'ns=2;s=/Random/Plc/DB21.DBX35.2'),
    (4, 'cycleCondition', 'Cycle condition', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB3300.DBX1.5', 'ns=2;s=/Random/Plc/DB21.DBX33.5'),
    (5, 'rapidMode', 'Rapid mode (G0)', 'bool', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/State/G0Mode', 'ns=2;s=/Channel/State/G0Mode'),
    (7, 'protectionLevel5', 'Protection level 5', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB2600.DBX0.7', 'ns=2;s=/Random/Plc/DB10.DBX56.6'),
    (8, 'feedStopActive', 'Feed stop active', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Channel/ProgramModification/feedStopActive', 'ns=2;s=/Channel/ProgramModification/feedStopActive'),
    (9, 'actBlock', 'Active block', 'text', NULL, 'process', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Channel/ProgramInfo/actBlock', 'ns=2;s=/Channel/ProgramInfo/actBlock'),
    (10, 'feedRateOvr', 'Feed rate override', 'number', '%', 'process', 3, NULL, 0, 'Dashboard And Chart', 'active',
        'ns=2;s=/Channel/GeometricAxis/feedRateOvr', 'ns=2;s=/Channel/GeometricAxis/feedRateOvr'),
    (11, 'searchRunMode', 'Search run mode', 'number', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Channel/State/searchRunMode', 'ns=2;s=/Channel/State/searchRunMode'),
    (12, 'nckAlarm', 'NCK alarm', 'bool', NULL, 'alarm', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB2700.DBX3.0', 'ns=2;s=/Random/Plc/DB10.DBX109.0'),
    (13, 'allAlarmNo', 'All alarm numbers', 'number', NULL, 'alarm', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Nck/AlarmEvent/textIndex', 'ns=2;s=/Nck/AlarmEvent/textIndex'),
    (14, 'opMode', 'Operating mode (0 JOG, 1 MDA, 2 AUTO)', 'number', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Bag/State/opMode', 'ns=2;s=/Bag/State/opMode'),
    (15, 'progName', 'Program name', 'text', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/ProgramInfo/progName', 'ns=2;s=/Channel/ProgramInfo/progName'),
    (16, 'speedOvr', 'Spindle override', 'number', '%', 'process', 3, NULL, 0, 'Only Dashboard', 'active',
        'ns=2;s=/Channel/Spindle/speedOvr', 'ns=2;s=/Channel/Spindle/speedOvr'),
    (17, 'actToolIdent', 'Active tool ident', 'text', NULL, 'process', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/State/actToolIdent', 'ns=2;s=/Channel/State/actToolIdent'),
    (19, 'chanStatus', 'Channel status', 'number', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Channel/State/chanStatus[u1]', 'ns=2;s=/Channel/State/chanStatus[u1]'),
    (21, 'noProRun', 'No program run', 'bool', NULL, 'state', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/DB3300.DBX3.4', 'ns=2;s=/Random/Plc/DB21.DBX35.4'),
    (23, 'driveReady', 'Drive ready', 'bool', NULL, 'state', 0, NULL, 0,
        '=1 when drive ready  =0 when drive not ready', 'inactive',
        'ns=2;s=/Random/Plc/DB2700.DBX2.6', 'ns=2;s=/Random/Plc/DB10.DBX108.6'),
    (24, 'OEMAlarmMCDown', 'OEM alarm machine down', 'bool', NULL, 'alarm', 0, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Random/Plc/M300.0', 'ns=2;s=/Random/Plc/M15.0'),
    -- Old row read the alarm "from code" (no node); now the DLL snapshot's AlarmNo
    (25, 'Alarm', 'Alarm number (0 = none)', 'number', NULL, 'alarm', 1, NULL, 0, 'DLL snapshot AlarmNo', 'active',
        NULL, NULL),
    (26, 'SpindleLoad', 'Spindle load', 'number', '%', 'process', 3, 2, 1000, 'Only Dashboard', 'active',
        'ns=2;s=/Channel/Spindle/driveLoad[1,1]', 'ns=2;s=/Channel/Spindle/driveLoad[1,1]'),
    (27, 'FeedRev', 'Feed rate unit', 'number', NULL, 'process', 0, NULL, 0, 'Only Dashboard', 'active',
        'ns=2;s=/Channel/State/feedRateIpoUnit', 'ns=2;s=/Channel/State/feedRateIpoUnit'),
    (28, 'SpindleRpm', 'Spindle speed', 'number', 'rpm', 'process', 1, NULL, 0, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/Spindle/actSpeed[1,1]', 'ns=2;s=/Channel/Spindle/actSpeed[1,1]'),
    (31, 'HeadstockTemp', 'Headstock temperature', 'number', '°C', 'health', 2, NULL, 0, 'Senser Variable', 'active',
        'ns=2;s=/Random/Plc/DB9006.DBD8', 'ns=2;s=/Random/Plc/DB9006.DBD8'),
    (32, 'axPowerX', 'X axis power', 'number', NULL, 'process', 0, NULL, 0, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,1]', 'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,1]'),
    (33, 'axPowerY', 'Y axis power', 'number', NULL, 'process', 0, NULL, 0, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,2]', 'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,2]'),
    (34, 'axPowerZ', 'Z axis power', 'number', NULL, 'process', 0, NULL, 0, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,3]', 'ns=2;s=/Channel/MachineAxis/aaPowerSmooth[U1,3]'),
    (35, 'axLoadZ', 'Z axis load', 'number', '%', 'process', 3, 2, 1000, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,3]', 'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,3]'),
    (36, 'axLoadY', 'Y axis load', 'number', '%', 'process', 3, 2, 1000, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,2]', 'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,2]'),
    (37, 'axLoadX', 'X axis load', 'number', '%', 'process', 3, 2, 1000, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,1]', 'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,1]'),
    (39, 'hydPressure', 'Hydraulic system pressure', 'number', 'bar', 'health', 2, NULL, 60000, 'System Pressure-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD8:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD8:REAL[1]'),
    (40, 'hydTemp', 'Tank oil temperature', 'number', '°C', 'health', 2, NULL, 300000, 'Tank Oil Temperature-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD0:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD0:REAL[1]'),
    (41, 'elePaTemp', 'Electrical panel temperature', 'number', '°C', 'health', 2, NULL, 300000, 'Electrical Panel Tempurature-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD108:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD108:REAL[1]'),
    (42, 'PnePress', 'Pneumatic pressure', 'number', 'bar', 'health', 2, NULL, 60000, 'Pneumatic Pressure-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD12:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD12:REAL[1]'),
    (43, 'CTSPressure', 'CTS pressure', 'number', 'bar', 'health', 3, NULL, 30000, 'CTS Pressure-Real Time-Coolant Flow Pressure Status', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD16:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD16:REAL[1]'),
    (44, 'SpVibr', 'Spindle vibration', 'number', NULL, 'health', 3, NULL, 1000, 'Spindle Vibration-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD28:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD28:REAL[1]'),
    (45, 'SPBeTemp', 'Spindle bearing temperature', 'number', '°C', 'health', 2, NULL, 300000, 'Spindle Bearing Temparture-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD4:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD4:REAL[1]'),
    (46, 'SPMotorTemp', 'Spindle motor temperature', 'number', '°C', 'health', 2, NULL, 120000, 'Spindle Motor Temparture-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD124:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD124:REAL[1]'),
    (47, 'ChillerTemp', 'Chiller temperature', 'number', '°C', 'health', 2, NULL, 300000, 'Chiller Temperature-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD112:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD112:REAL[1]'),
    (48, 'CFBRate', 'Chiller flow bearing rate', 'number', NULL, 'health', 2, NULL, 60000, 'Chiller Flow Bearing Rate-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD20:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD20:REAL[1]'),
    (49, 'CFMRate', 'Chiller flow motor rate', 'number', NULL, 'health', 2, NULL, 60000, 'Chiller Flow Motor Rate-Real Time', 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD24:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD24:REAL[1]'),
    (50, 'MMode', 'Main mode (PLC)', 'number', NULL, 'state', 1, NULL, 0, 'Main Mode', 'active',
        'ns=2;s=/Random/Plc/DB9011.DBB8', 'ns=2;s=/Random/Plc/DB9011.DBB8'),
    (51, 'SMode', 'Sub mode (PLC)', 'number', NULL, 'state', 1, NULL, 0, 'Sub MODE', 'active',
        'ns=2;s=/Random/Plc/DB9011.DBB9', 'ns=2;s=/Random/Plc/DB9011.DBB9'),
    (52, 'DownRe', 'Downtime reason (PLC)', 'number', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9011.DBW4', 'ns=2;s=/Random/Plc/DB9011.DBW4'),
    (53, 'TaskList', 'Task list (PLC)', 'number', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9007.DBW2', 'ns=2;s=/Random/Plc/DB9007.DBW2'),
    (54, 'actFeedRate', 'Actual feed rate', 'number', 'mm/min', 'process', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/State/actFeedRateIpo[U1,1]', 'ns=2;s=/Channel/State/actFeedRateIpo[U1,1]'),
    -- Old row was inactive; the DLL snapshot now delivers it (ProcStatus)
    (55, 'proStatus', 'Program status (0 Reset, 1 Stop, 2 Active, 3 Waiting, 4 Interrupted)', 'number', NULL, 'state', 1, NULL, 0, NULL, 'active',
        'ns=2;s=/Channel/State/progStatus', 'ns=2;s=/Channel/State/progStatus'),
    (56, 'cuttingTime', 'Cutting time', 'number', 's', 'process', 1, NULL, 0, NULL, 'inactive',
        'ns=2;s=/Channel/ChannelDiagnose/cuttingTime', 'ns=2;s=/Channel/ChannelDiagnose/cuttingTime'),
    (57, 'Hyd_act_pre', 'Turning hydraulic actual pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD8:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD8:REAL[1]'),
    (58, 'Hyd_set_pre', 'Turning hydraulic set pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD52:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD52:REAL[1]'),
    (59, 'Hyd_rst_pre', 'Turning hydraulic reset pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD104:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD104:REAL[1]'),
    (60, 'CTS_act_pre', 'Turning CTS actual pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD16:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD16:REAL[1]'),
    (61, 'CTS_high_pre', 'Turning CTS high pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD68:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD68:REAL[1]'),
    (62, 'CTS_low_pre', 'Turning CTS low pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD72:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD72:REAL[1]'),
    (63, 'Chuck_id_act_pre', 'Turning chuck ID actual pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD128:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD128:REAL[1]'),
    (64, 'chuck_id_set_pre', 'Turning chuck ID set pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD132:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD132:REAL[1]'),
    (65, 'chuck_id_rest_pre', 'Turning chuck ID reset pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD136:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD136:REAL[1]'),
    (66, 'chuck_od_act_pre', 'Turning chuck OD actual pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD140:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD140:REAL[1]'),
    (67, 'chuck_od_set_pre', 'Turning chuck OD set pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD144:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD144:REAL[1]'),
    (68, 'chuck_od_rest_pre', 'Turning chuck OD reset pressure', 'number', 'bar', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD148:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD148:REAL[1]'),
    (69, 'ele_panel_set_temp', 'Electrical panel set temperature', 'number', '°C', 'health', 2, NULL, 60000, NULL, 'active',
        'ns=2;s=/Random/Plc/DB9004.DBD116:REAL[1]', 'ns=2;s=/Random/Plc/DB9004.DBD116:REAL[1]'),
    (132, 'axLoadA', 'A axis load', 'number', '%', 'process', 3, 2, 1000, 'Only For Dashboard', 'active',
        'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,5]', 'ns=2;s=/Channel/MachineAxis/aaLoadSmooth[U1,5]'),
    -- New: values the DLL snapshot delivers that the old list did not have
    (200, 'isEmergency', 'Emergency stop', 'bool', NULL, 'alarm', 1, NULL, 0, 'DLL snapshot IsEmergency', 'active', NULL, NULL),
    (201, 'partCount', 'Part count', 'number', NULL, 'process', 1, NULL, 0, 'DLL snapshot ProgramPartCount', 'active', NULL, NULL),
    (202, 'requiredPartCount', 'Required part count', 'number', NULL, 'process', 1, NULL, 0, 'DLL snapshot RequiredPartCount', 'active', NULL, NULL),
    (203, 'toolNumber', 'Tool number', 'number', NULL, 'process', 1, NULL, 0, 'DLL snapshot ToolNumber', 'active', NULL, NULL),
    (204, 'curProgLine', 'Current program line', 'number', NULL, 'process', 1, NULL, 0, 'DLL snapshot CurProgLine', 'active', NULL, NULL),
    (205, 'totalProgLines', 'Total program lines', 'number', NULL, 'process', 1, NULL, 0, 'DLL snapshot TotalProgLines', 'active', NULL, NULL),
    (206, 'progPath', 'Active program path', 'text', NULL, 'state', 1, NULL, 0, 'DLL snapshot ActiveProgPath', 'active', NULL, NULL),
    (207, 'alarmMsg', 'Alarm message', 'text', NULL, 'alarm', 1, NULL, 0, 'DLL snapshot AlarmMsg', 'active', NULL, NULL);

INSERT INTO cnc_variables (id, code, name, data_type, unit, category, store_mode, deadband,
                           min_interval_ms, remarks, status)
SELECT id, code, name, data_type, unit, category, store_mode, deadband, min_interval_ms,
       NULLIF(remarks, ''), status
FROM seed_cnc_variables
ON CONFLICT (id) DO NOTHING;

-- Siemens (engine 2) node ids: one row for all versions when 828D and 840D agree,
-- otherwise one row per version. /Random/Plc/ and /Plc/ nodes are PLC, the rest NC.
INSERT INTO cnc_variable_addresses (variable_id, cnc_engine, controller_version_type, source, address)
SELECT s.id, 2, v.version_type,
       CASE WHEN v.address LIKE '%/Plc/%' THEN 'plc' ELSE 'nc' END, v.address
FROM seed_cnc_variables s
CROSS JOIN LATERAL (VALUES
    (CASE WHEN s.path_828d = s.path_840d THEN NULL ELSE 1 END::SMALLINT, s.path_828d),
    (2::SMALLINT, CASE WHEN s.path_828d = s.path_840d THEN NULL ELSE s.path_840d END)
) AS v(version_type, address)
WHERE v.address IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM cnc_variable_addresses a
                  WHERE a.variable_id = s.id AND a.cnc_engine = 2 AND a.is_deleted = false
                    AND COALESCE(a.controller_version_type, 0) = COALESCE(v.version_type, 0)
                    AND a.source = CASE WHEN v.address LIKE '%/Plc/%' THEN 'plc' ELSE 'nc' END);

-- Fields of the DLL snapshot (docs/cnc-communication.md §6), for every engine the DLL has.
-- Field names keep the DLL's spelling.
INSERT INTO cnc_variable_addresses (variable_id, cnc_engine, controller_version_type, source, address)
SELECT f.variable_id, e.engine, NULL, 'snapshot', f.field
FROM (VALUES
    (10, 'FeedRatePercetnage'),
    (14, 'MachineMode'),
    (15, 'ProgName'),
    (16, 'SpindeOverridePercentage'),
    (25, 'AlarmNo'),
    (26, 'SpindleLoadPercetnage'),
    (28, 'SpindeOverrideActual'),
    (35, 'AxisLoadZPercetnage'),
    (36, 'AxisLoadYPercetnage'),
    (37, 'AxisLoadXPercetnage'),
    (54, 'FeedRateActual'),
    (55, 'ProcStatus'),
    (200, 'IsEmergency'),
    (201, 'ProgramPartCount'),
    (202, 'RequiredPartCount'),
    (203, 'ToolNumber'),
    (204, 'CurProgLine'),
    (205, 'TotalProgLines'),
    (206, 'ActiveProgPath'),
    (207, 'AlarmMsg')
) AS f(variable_id, field)
CROSS JOIN (VALUES (0::SMALLINT), (1::SMALLINT), (2::SMALLINT), (3::SMALLINT)) AS e(engine)
WHERE NOT EXISTS (SELECT 1 FROM cnc_variable_addresses a
                  WHERE a.variable_id = f.variable_id AND a.cnc_engine = e.engine
                    AND a.controller_version_type IS NULL AND a.source = 'snapshot'
                    AND a.is_deleted = false);

-- ============================================================================================
-- Data
-- ============================================================================================

-- One row per machine-variable: latest value (value for number/bool variables only).
-- Upserted on every reading, then used to detect changes before inserting into time-series.
-- Text variables are not cached here (only in time-series). No id column (costs 8 bytes + sequence per row, never queried).
-- No foreign keys (machines are soft-deleted, so ON DELETE never fires). updated_at tracks when value changed.
CREATE TABLE machine_variable_values (
    machine_id integer NOT NULL,
    variable_id integer NOT NULL,
    value double precision,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (machine_id, variable_id)
);

-- Time series: one row per stored value, value for number / bool variables only.
-- Only changed values are inserted (detected by comparing with machine_variable_values).
-- Text variables are not stored. The latest value of a variable is its newest row (the index below).
-- No id column (it would cost 8 bytes and a sequence call per row and is never queried) and no foreign keys
-- (a lookup per inserted row; machines are soft-deleted, so ON DELETE never fires). recorded_at = server time the value arrived.
CREATE TABLE machine_telementry (
    machine_id integer NOT NULL,
    variable_id integer NOT NULL,
    value double precision,
    recorded_at timestamp with time zone NOT NULL
);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb') THEN
        PERFORM create_hypertable('machine_telementry', 'recorded_at',
            chunk_time_interval => INTERVAL '6 hours',
            create_default_indexes => false, if_not_exists => true);
        PERFORM add_retention_policy('machine_telementry', INTERVAL '365 days', if_not_exists => true);
        -- Compression is not in the Apache-licensed TimescaleDB build; keep going without it
        BEGIN
            EXECUTE 'ALTER TABLE machine_telementry SET (timescaledb.compress, '
                'timescaledb.compress_segmentby = ''machine_id, variable_id'', '
                'timescaledb.compress_orderby = ''recorded_at DESC'')';
            PERFORM add_compression_policy('machine_telementry', INTERVAL '1 day', if_not_exists => true);
        EXCEPTION WHEN OTHERS THEN
            RAISE NOTICE 'machine_telementry compression not enabled: %', SQLERRM;
        END;
    END IF;
END
$$;

CREATE INDEX idx_machine_telementry_machine_variable_ts
    ON machine_telementry (machine_id, variable_id, recorded_at DESC);

-- One row per collected machine: connection, current mode and collector health
CREATE TABLE machine_live_status (
    machine_id INTEGER PRIMARY KEY REFERENCES machines(id) ON DELETE CASCADE,
    collector_shard SMALLINT,                 -- collector process that owns the machine
    is_connected BOOLEAN NOT NULL DEFAULT false,
    mode_id INTEGER REFERENCES machine_modes(id),
    connected_since TIMESTAMPTZ,
    last_seen_at TIMESTAMPTZ,                 -- server clock, last DLL callback
    last_change_at TIMESTAMPTZ,               -- server clock, last stored value
    sample_ts TIMESTAMPTZ,                    -- DLL clock of the last callback (clock drift check)
    missing_nodes TEXT,                       -- nodes the controller does not have (DLL GcodeGroup)
    callback_count BIGINT NOT NULL DEFAULT 0, -- DLL callbacks since the collector started
    stored_count BIGINT NOT NULL DEFAULT 0,   -- values stored since the collector started
    dropped_count BIGINT NOT NULL DEFAULT 0,  -- values dropped (write buffer full)
    error_count BIGINT NOT NULL DEFAULT 0,    -- DLL error callbacks since the collector started
    last_error TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- One row per period a machine stays in one mode; ended_at NULL = current period
CREATE TABLE machine_mode_log (
    id BIGSERIAL PRIMARY KEY,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    mode_id INTEGER NOT NULL REFERENCES machine_modes(id),
    program_name VARCHAR(64),                 -- at the start of the period
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    duration_sec INTEGER,                     -- set when the period closes
    CONSTRAINT machine_mode_log_period_check CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX idx_machine_mode_log_machine_started ON machine_mode_log (machine_id, started_at DESC);

-- At most one open period per machine
CREATE UNIQUE INDEX machine_mode_log_open_key
    ON machine_mode_log (machine_id) WHERE ended_at IS NULL;

CREATE TABLE machine_events (
    id BIGSERIAL PRIMARY KEY,
    machine_id INTEGER REFERENCES machines(id) ON DELETE CASCADE, -- NULL for collector events
    collector_shard SMALLINT,
    event_type VARCHAR(30) NOT NULL,
    code INTEGER,                             -- DLL error code, CNCConnect result, alarm number or count
    message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT machine_events_event_type_check CHECK (event_type IN (
        'collector_start', 'collector_stop', 'values_dropped', 'config_error',
        'connect_ok', 'connect_failed', 'disconnected', 'stale', 'reconnect', 'dll_error',
        'alarm_on', 'alarm_off', 'emergency_on', 'emergency_off'))
);

CREATE INDEX idx_machine_events_machine_created ON machine_events (machine_id, created_at DESC);
CREATE INDEX idx_machine_events_type_created ON machine_events (event_type, created_at DESC);

-- ============================================================================================
-- Notifications
-- ============================================================================================

-- A machine is added or removed, or a field the collector uses changes. The collectors
-- connect to every active, not deleted machine; they check the machine again on each notification.
CREATE OR REPLACE FUNCTION notify_machines_changed() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('machines_changed', OLD.id::text);
        RETURN OLD;
    END IF;

    IF TG_OP = 'INSERT' THEN
        PERFORM pg_notify('machines_changed', NEW.id::text);
        RETURN NEW;
    END IF;

    IF (NEW.is_deleted, NEW.status, NEW.live_ip, NEW.port, NEW.username, NEW.password,
        NEW.brand_id, NEW.controller_id, NEW.controller_version_type)
       IS DISTINCT FROM
       (OLD.is_deleted, OLD.status, OLD.live_ip, OLD.port, OLD.username, OLD.password,
        OLD.brand_id, OLD.controller_id, OLD.controller_version_type) THEN
        PERFORM pg_notify('machines_changed', NEW.id::text);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_machines_changed AFTER INSERT OR UPDATE OR DELETE ON machines
    FOR EACH ROW EXECUTE FUNCTION notify_machines_changed();

-- The variable configuration changes. PostgreSQL sends the same payload once per transaction,
-- so a bulk edit is one notification. Controller cnc_engine changes are manual admin operations
-- and do not require collector notification (collectors reload config on startup or on
-- cnc_config_changed for variables / machine-specific overrides).
CREATE OR REPLACE FUNCTION notify_cnc_config_changed() RETURNS trigger AS $$
DECLARE
    row_data RECORD;
BEGIN
    IF TG_OP = 'DELETE' THEN
        row_data := OLD;
    ELSE
        row_data := NEW;
    END IF;

    IF TG_TABLE_NAME = 'machine_variables' THEN
        PERFORM pg_notify('cnc_config_changed', 'machine:' || row_data.machine_id);
    ELSE
        PERFORM pg_notify('cnc_config_changed', 'variables');
    END IF;
    RETURN row_data;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_cnc_config_changed AFTER INSERT OR UPDATE OR DELETE ON cnc_variables
    FOR EACH ROW EXECUTE FUNCTION notify_cnc_config_changed();

CREATE TRIGGER notify_cnc_config_changed AFTER INSERT OR UPDATE OR DELETE ON cnc_variable_addresses
    FOR EACH ROW EXECUTE FUNCTION notify_cnc_config_changed();

CREATE TRIGGER notify_cnc_config_changed AFTER INSERT OR UPDATE OR DELETE ON machine_variables
    FOR EACH ROW EXECUTE FUNCTION notify_cnc_config_changed();

COMMIT;
