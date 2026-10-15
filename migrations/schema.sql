--
-- PostgreSQL database dump
--

\restrict v1DMVH0erDHUDQG7KifHiGyRZ0ExeMfDqOTfZd3w8RIvRVhL6sLVBaalPdv3FZ4

-- Dumped from database version 18.6
-- Dumped by pg_dump version 18.6

-- Started on 2026-10-15 12:11:20

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 304 (class 1255 OID 116754)
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 230 (class 1259 OID 116630)
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    user_id integer,
    address_type character varying(20) DEFAULT 'home'::character varying,
    address_line1 character varying(255),
    address_line2 character varying(255),
    city character varying(100),
    pin_code character varying(20),
    contact_person character varying(100),
    contact_phone character varying(20),
    is_default boolean DEFAULT false,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    country_name character varying(100),
    state_name character varying(100)
);


ALTER TABLE public.addresses OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 116629)
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.addresses_id_seq OWNER TO postgres;

--
-- TOC entry 5718 (class 0 OID 0)
-- Dependencies: 229
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- TOC entry 238 (class 1259 OID 116732)
-- Name: auth_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.auth_tokens (
    id integer NOT NULL,
    user_id integer,
    token character varying(255) NOT NULL,
    token_type character varying(20) NOT NULL,
    device_type character varying(50),
    device_info character varying(255),
    expires_at timestamp without time zone,
    last_used_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.auth_tokens OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 116731)
-- Name: auth_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.auth_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.auth_tokens_id_seq OWNER TO postgres;

--
-- TOC entry 5719 (class 0 OID 0)
-- Dependencies: 237
-- Name: auth_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.auth_tokens_id_seq OWNED BY public.auth_tokens.id;


--
-- TOC entry 240 (class 1259 OID 116773)
-- Name: branches; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.branches (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100),
    address text,
    customer_id integer,
    gst_number character varying(30),
    pan_number character varying(30),
    color_code character varying(10),
    status character varying(20) DEFAULT 'active'::character varying,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    user_city_id integer
);


ALTER TABLE public.branches OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 116772)
-- Name: branches_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.branches_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.branches_id_seq OWNER TO postgres;

--
-- TOC entry 5720 (class 0 OID 0)
-- Dependencies: 239
-- Name: branches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.branches_id_seq OWNED BY public.branches.id;


--
-- TOC entry 244 (class 1259 OID 116866)
-- Name: brands; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.brands (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50),
    description text,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.brands OWNER TO postgres;

--
-- TOC entry 5721 (class 0 OID 0)
-- Dependencies: 244
-- Name: TABLE brands; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.brands IS 'Machine brands - Added in migration 005';


--
-- TOC entry 243 (class 1259 OID 116865)
-- Name: brands_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.brands_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.brands_id_seq OWNER TO postgres;

--
-- TOC entry 5722 (class 0 OID 0)
-- Dependencies: 243
-- Name: brands_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.brands_id_seq OWNED BY public.brands.id;


--
-- TOC entry 290 (class 1259 OID 117957)
-- Name: breaks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.breaks (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(100) NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    shift_id integer,
    CONSTRAINT breaks_check CHECK ((end_time <> start_time))
);


ALTER TABLE public.breaks OWNER TO postgres;

--
-- TOC entry 289 (class 1259 OID 117956)
-- Name: breaks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.breaks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.breaks_id_seq OWNER TO postgres;

--
-- TOC entry 5723 (class 0 OID 0)
-- Dependencies: 289
-- Name: breaks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.breaks_id_seq OWNED BY public.breaks.id;


--
-- TOC entry 242 (class 1259 OID 116849)
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50),
    description text,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- TOC entry 5724 (class 0 OID 0)
-- Dependencies: 242
-- Name: TABLE categories; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.categories IS 'Product categories - Added in migration 005';


--
-- TOC entry 241 (class 1259 OID 116848)
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- TOC entry 5725 (class 0 OID 0)
-- Dependencies: 241
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 254 (class 1259 OID 117048)
-- Name: category_machines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.category_machines (
    id integer NOT NULL,
    category_id integer,
    machine_id integer,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.category_machines OWNER TO postgres;

--
-- TOC entry 253 (class 1259 OID 117047)
-- Name: category_machines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.category_machines_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.category_machines_id_seq OWNER TO postgres;

--
-- TOC entry 5726 (class 0 OID 0)
-- Dependencies: 253
-- Name: category_machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.category_machines_id_seq OWNED BY public.category_machines.id;


--
-- TOC entry 286 (class 1259 OID 117842)
-- Name: clienteles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clienteles (
    id integer NOT NULL,
    user_id integer NOT NULL,
    client_code character varying(50) NOT NULL,
    name character varying(250) NOT NULL,
    email character varying(255),
    contact_number character varying(20) NOT NULL,
    address character varying(250) NOT NULL,
    country_name character varying(100) CONSTRAINT clienteles_country_code_not_null NOT NULL,
    state_name character varying(100) CONSTRAINT clienteles_state_code_not_null NOT NULL,
    city character varying(100) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.clienteles OWNER TO postgres;

--
-- TOC entry 285 (class 1259 OID 117841)
-- Name: clienteles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clienteles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clienteles_id_seq OWNER TO postgres;

--
-- TOC entry 5727 (class 0 OID 0)
-- Dependencies: 285
-- Name: clienteles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clienteles_id_seq OWNED BY public.clienteles.id;


--
-- TOC entry 234 (class 1259 OID 116689)
-- Name: company_profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.company_profiles (
    id integer NOT NULL,
    user_id integer,
    gst_number character varying(30),
    pan_number character varying(30),
    website character varying(255),
    finix_id integer,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    company_name character varying(150)
);


ALTER TABLE public.company_profiles OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 116688)
-- Name: company_profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.company_profiles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.company_profiles_id_seq OWNER TO postgres;

--
-- TOC entry 5728 (class 0 OID 0)
-- Dependencies: 233
-- Name: company_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.company_profiles_id_seq OWNED BY public.company_profiles.id;


--
-- TOC entry 250 (class 1259 OID 116928)
-- Name: controller_products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.controller_products (
    id integer NOT NULL,
    controller_id integer,
    product_id integer,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.controller_products OWNER TO postgres;

--
-- TOC entry 5729 (class 0 OID 0)
-- Dependencies: 250
-- Name: TABLE controller_products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.controller_products IS 'Controller-Product relationships - Added in migration 005';


--
-- TOC entry 249 (class 1259 OID 116927)
-- Name: controller_products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.controller_products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.controller_products_id_seq OWNER TO postgres;

--
-- TOC entry 5730 (class 0 OID 0)
-- Dependencies: 249
-- Name: controller_products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.controller_products_id_seq OWNED BY public.controller_products.id;


--
-- TOC entry 248 (class 1259 OID 116906)
-- Name: controllers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.controllers (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    brand_id integer,
    description text,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.controllers OWNER TO postgres;

--
-- TOC entry 5731 (class 0 OID 0)
-- Dependencies: 248
-- Name: TABLE controllers; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.controllers IS 'Machine controllers - Added in migration 005';


--
-- TOC entry 247 (class 1259 OID 116905)
-- Name: controllers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.controllers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.controllers_id_seq OWNER TO postgres;

--
-- TOC entry 5732 (class 0 OID 0)
-- Dependencies: 247
-- Name: controllers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.controllers_id_seq OWNED BY public.controllers.id;


--
-- TOC entry 222 (class 1259 OID 116513)
-- Name: countries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.countries (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(3) NOT NULL,
    phone_code character varying(10) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.countries OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 116512)
-- Name: countries_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.countries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.countries_id_seq OWNER TO postgres;

--
-- TOC entry 5733 (class 0 OID 0)
-- Dependencies: 221
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.countries_id_seq OWNED BY public.countries.id;


--
-- TOC entry 292 (class 1259 OID 118011)
-- Name: holidays; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.holidays (
    id integer NOT NULL,
    user_id integer NOT NULL,
    holiday_type character varying(10) NOT NULL,
    name character varying(100),
    start_date date,
    end_date date,
    week_day smallint,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT holidays_holiday_type_check CHECK (((holiday_type)::text = ANY ((ARRAY['holiday'::character varying, 'week_off'::character varying])::text[]))),
    CONSTRAINT holidays_kind_check CHECK (((((holiday_type)::text = 'holiday'::text) AND (name IS NOT NULL) AND (start_date IS NOT NULL) AND (end_date IS NOT NULL) AND (end_date >= start_date) AND (week_day IS NULL)) OR (((holiday_type)::text = 'week_off'::text) AND ((week_day >= 1) AND (week_day <= 7)) AND (start_date IS NULL) AND (end_date IS NULL))))
);


ALTER TABLE public.holidays OWNER TO postgres;

--
-- TOC entry 291 (class 1259 OID 118010)
-- Name: holidays_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.holidays_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.holidays_id_seq OWNER TO postgres;

--
-- TOC entry 5734 (class 0 OID 0)
-- Dependencies: 291
-- Name: holidays_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.holidays_id_seq OWNED BY public.holidays.id;


--
-- TOC entry 252 (class 1259 OID 116968)
-- Name: machines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.machines (
    id integer NOT NULL,
    user_id integer,
    machine_type_id integer,
    product_id integer,
    name character varying(50) NOT NULL,
    serial_number character varying(50) NOT NULL,
    live_ip character varying(20) NOT NULL,
    port integer NOT NULL,
    username character varying(15) NOT NULL,
    password character varying(255) NOT NULL,
    local_ip character varying(20) NOT NULL,
    brand_id integer,
    controller_id integer,
    controller_version character varying(10) NOT NULL,
    accuracy_id integer,
    line_id integer,
    installation_date date NOT NULL,
    bill_date date,
    warranty_date date,
    amc character varying(50),
    is_insured boolean DEFAULT false NOT NULL,
    online_status boolean DEFAULT false NOT NULL,
    maintenance_required boolean DEFAULT false NOT NULL,
    hourly_cost numeric(10,2) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_monitor boolean DEFAULT false NOT NULL,
    online_date date,
    offline_date timestamp with time zone,
    mode_id integer,
    controller_version_type integer,
    health_ids character varying(100),
    longitude character varying(50),
    latitude character varying(50),
    cycle_analysis boolean DEFAULT false NOT NULL,
    machine_system_type character varying(20) DEFAULT 'NewPlc'::character varying NOT NULL,
    security_mode character varying(20) DEFAULT 'Basic256Sha256'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    model_id integer
);


ALTER TABLE public.machines OWNER TO postgres;

--
-- TOC entry 5735 (class 0 OID 0)
-- Dependencies: 252
-- Name: TABLE machines; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.machines IS 'Machine inventory and details - Added in migration 002';


--
-- TOC entry 251 (class 1259 OID 116967)
-- Name: machines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.machines_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.machines_id_seq OWNER TO postgres;

--
-- TOC entry 5736 (class 0 OID 0)
-- Dependencies: 251
-- Name: machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.machines_id_seq OWNED BY public.machines.id;


--
-- TOC entry 274 (class 1259 OID 117647)
-- Name: models; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.models (
    id integer NOT NULL,
    controller_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.models OWNER TO postgres;

--
-- TOC entry 273 (class 1259 OID 117646)
-- Name: models_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.models_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.models_id_seq OWNER TO postgres;

--
-- TOC entry 5737 (class 0 OID 0)
-- Dependencies: 273
-- Name: models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.models_id_seq OWNED BY public.models.id;


--
-- TOC entry 282 (class 1259 OID 117783)
-- Name: part_process_machines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.part_process_machines (
    id integer NOT NULL,
    process_id integer NOT NULL,
    machine_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.part_process_machines OWNER TO postgres;

--
-- TOC entry 281 (class 1259 OID 117782)
-- Name: part_process_machines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.part_process_machines_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.part_process_machines_id_seq OWNER TO postgres;

--
-- TOC entry 5738 (class 0 OID 0)
-- Dependencies: 281
-- Name: part_process_machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.part_process_machines_id_seq OWNED BY public.part_process_machines.id;


--
-- TOC entry 280 (class 1259 OID 117738)
-- Name: part_processes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.part_processes (
    id integer NOT NULL,
    part_id integer NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    process character varying(100) NOT NULL,
    program_name character varying(100) NOT NULL,
    drawing_no character varying(100) NOT NULL,
    drawing_revision_no character varying(100) NOT NULL,
    drawing_file character varying(100) NOT NULL,
    drawing_name character varying(255) NOT NULL,
    accuracy_id integer NOT NULL,
    setup_time integer DEFAULT 0 NOT NULL,
    cycle_time integer NOT NULL,
    loading_time integer DEFAULT 0 NOT NULL,
    production_time integer DEFAULT 0 NOT NULL,
    has_sub_program boolean DEFAULT false NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT part_processes_cycle_time_check CHECK ((cycle_time > 0)),
    CONSTRAINT part_processes_loading_time_check CHECK ((loading_time >= 0)),
    CONSTRAINT part_processes_production_time_check CHECK ((production_time >= 0)),
    CONSTRAINT part_processes_setup_time_check CHECK ((setup_time >= 0))
);


ALTER TABLE public.part_processes OWNER TO postgres;

--
-- TOC entry 279 (class 1259 OID 117737)
-- Name: part_processes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.part_processes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.part_processes_id_seq OWNER TO postgres;

--
-- TOC entry 5739 (class 0 OID 0)
-- Dependencies: 279
-- Name: part_processes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.part_processes_id_seq OWNED BY public.part_processes.id;


--
-- TOC entry 284 (class 1259 OID 117809)
-- Name: part_sub_programs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.part_sub_programs (
    id integer NOT NULL,
    process_id integer NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    name character varying(200) NOT NULL,
    cycle_time integer DEFAULT 0 NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT part_sub_programs_cycle_time_check CHECK ((cycle_time >= 0))
);


ALTER TABLE public.part_sub_programs OWNER TO postgres;

--
-- TOC entry 283 (class 1259 OID 117808)
-- Name: part_sub_programs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.part_sub_programs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.part_sub_programs_id_seq OWNER TO postgres;

--
-- TOC entry 5740 (class 0 OID 0)
-- Dependencies: 283
-- Name: part_sub_programs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.part_sub_programs_id_seq OWNED BY public.part_sub_programs.id;


--
-- TOC entry 278 (class 1259 OID 117706)
-- Name: parts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.parts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    clientele_id integer,
    part_no character varying(100) NOT NULL,
    part_name character varying(200) NOT NULL,
    part_count integer DEFAULT 1 NOT NULL,
    production_time integer DEFAULT 0 NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT parts_part_count_check CHECK ((part_count > 0)),
    CONSTRAINT parts_production_time_check CHECK ((production_time >= 0))
);


ALTER TABLE public.parts OWNER TO postgres;

--
-- TOC entry 277 (class 1259 OID 117705)
-- Name: parts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.parts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.parts_id_seq OWNER TO postgres;

--
-- TOC entry 5741 (class 0 OID 0)
-- Dependencies: 277
-- Name: parts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.parts_id_seq OWNED BY public.parts.id;


--
-- TOC entry 246 (class 1259 OID 116883)
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50),
    category_id integer,
    image character varying(255) DEFAULT 'no_image_available.png'::character varying,
    description text,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.products OWNER TO postgres;

--
-- TOC entry 5742 (class 0 OID 0)
-- Dependencies: 246
-- Name: TABLE products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.products IS 'Machine products - Added in migration 005';


--
-- TOC entry 245 (class 1259 OID 116882)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- TOC entry 5743 (class 0 OID 0)
-- Dependencies: 245
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 226 (class 1259 OID 116542)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 116541)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- TOC entry 5744 (class 0 OID 0)
-- Dependencies: 225
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 288 (class 1259 OID 117890)
-- Name: shifts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.shifts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(100) NOT NULL,
    shift_type character varying(10) NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,
    working_days integer[] DEFAULT '{1,2,3,4,5,6,7}'::integer[] NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT shifts_check CHECK ((end_time <> start_time)),
    CONSTRAINT shifts_shift_type_check CHECK (((shift_type)::text = ANY ((ARRAY['day'::character varying, 'night'::character varying])::text[])))
);


ALTER TABLE public.shifts OWNER TO postgres;

--
-- TOC entry 287 (class 1259 OID 117889)
-- Name: shifts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.shifts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.shifts_id_seq OWNER TO postgres;

--
-- TOC entry 5745 (class 0 OID 0)
-- Dependencies: 287
-- Name: shifts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.shifts_id_seq OWNED BY public.shifts.id;


--
-- TOC entry 224 (class 1259 OID 116527)
-- Name: states; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.states (
    id integer NOT NULL,
    country_id integer,
    name character varying(100) NOT NULL,
    code character varying(10),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.states OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 116526)
-- Name: states_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.states_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.states_id_seq OWNER TO postgres;

--
-- TOC entry 5746 (class 0 OID 0)
-- Dependencies: 223
-- Name: states_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.states_id_seq OWNED BY public.states.id;


--
-- TOC entry 270 (class 1259 OID 117332)
-- Name: user_cell; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_cell (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_line_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_cell OWNER TO postgres;

--
-- TOC entry 269 (class 1259 OID 117331)
-- Name: user_cell_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_cell_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_cell_id_seq OWNER TO postgres;

--
-- TOC entry 5747 (class 0 OID 0)
-- Dependencies: 269
-- Name: user_cell_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_cell_id_seq OWNED BY public.user_cell.id;


--
-- TOC entry 262 (class 1259 OID 117182)
-- Name: user_city; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_city (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_state_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_city OWNER TO postgres;

--
-- TOC entry 261 (class 1259 OID 117181)
-- Name: user_city_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_city_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_city_id_seq OWNER TO postgres;

--
-- TOC entry 5748 (class 0 OID 0)
-- Dependencies: 261
-- Name: user_city_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_city_id_seq OWNED BY public.user_city.id;


--
-- TOC entry 258 (class 1259 OID 117129)
-- Name: user_country; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_country (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_country OWNER TO postgres;

--
-- TOC entry 257 (class 1259 OID 117128)
-- Name: user_country_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_country_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_country_id_seq OWNER TO postgres;

--
-- TOC entry 5749 (class 0 OID 0)
-- Dependencies: 257
-- Name: user_country_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_country_id_seq OWNED BY public.user_country.id;


--
-- TOC entry 256 (class 1259 OID 117105)
-- Name: user_filters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_filters (
    id integer NOT NULL,
    user_id integer NOT NULL,
    filter_name character varying(30) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_filters OWNER TO postgres;

--
-- TOC entry 255 (class 1259 OID 117104)
-- Name: user_filters_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_filters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_filters_id_seq OWNER TO postgres;

--
-- TOC entry 5750 (class 0 OID 0)
-- Dependencies: 255
-- Name: user_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_filters_id_seq OWNED BY public.user_filters.id;


--
-- TOC entry 268 (class 1259 OID 117302)
-- Name: user_line; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_line (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_shop_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_line OWNER TO postgres;

--
-- TOC entry 267 (class 1259 OID 117301)
-- Name: user_line_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_line_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_line_id_seq OWNER TO postgres;

--
-- TOC entry 5751 (class 0 OID 0)
-- Dependencies: 267
-- Name: user_line_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_line_id_seq OWNED BY public.user_line.id;


--
-- TOC entry 264 (class 1259 OID 117242)
-- Name: user_location; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_location (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    branch_id integer NOT NULL
);


ALTER TABLE public.user_location OWNER TO postgres;

--
-- TOC entry 263 (class 1259 OID 117241)
-- Name: user_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_location_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_location_id_seq OWNER TO postgres;

--
-- TOC entry 5752 (class 0 OID 0)
-- Dependencies: 263
-- Name: user_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_location_id_seq OWNED BY public.user_location.id;


--
-- TOC entry 276 (class 1259 OID 117680)
-- Name: user_machine_access; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_machine_access (
    id integer NOT NULL,
    user_id integer NOT NULL,
    machine_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_machine_access OWNER TO postgres;

--
-- TOC entry 275 (class 1259 OID 117679)
-- Name: user_machine_access_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_machine_access_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_machine_access_id_seq OWNER TO postgres;

--
-- TOC entry 5753 (class 0 OID 0)
-- Dependencies: 275
-- Name: user_machine_access_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_machine_access_id_seq OWNED BY public.user_machine_access.id;


--
-- TOC entry 236 (class 1259 OID 116710)
-- Name: user_machines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_machines (
    id integer NOT NULL,
    user_id integer,
    machine_id integer NOT NULL,
    assigned_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.user_machines OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 116709)
-- Name: user_machines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_machines_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_machines_id_seq OWNER TO postgres;

--
-- TOC entry 5754 (class 0 OID 0)
-- Dependencies: 235
-- Name: user_machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_machines_id_seq OWNED BY public.user_machines.id;


--
-- TOC entry 232 (class 1259 OID 116663)
-- Name: user_profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_profiles (
    id integer NOT NULL,
    user_id integer,
    designation character varying(100),
    grade character varying(50),
    birth_date date,
    joining_date date,
    blood_group character varying(5),
    profile_picture character varying(255),
    contact_person character varying(100),
    contact_phone character varying(20),
    theme character varying(20) DEFAULT 'dark_version'::character varying,
    timezone character varying(50) DEFAULT 'Asia/Kolkata'::character varying,
    language character varying(10) DEFAULT 'en'::character varying,
    start_day_time time without time zone,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    filters text[] DEFAULT '{}'::text[]
);


ALTER TABLE public.user_profiles OWNER TO postgres;

--
-- TOC entry 231 (class 1259 OID 116662)
-- Name: user_profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_profiles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_profiles_id_seq OWNER TO postgres;

--
-- TOC entry 5755 (class 0 OID 0)
-- Dependencies: 231
-- Name: user_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_profiles_id_seq OWNED BY public.user_profiles.id;


--
-- TOC entry 266 (class 1259 OID 117272)
-- Name: user_shop; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_shop (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_location_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_shop OWNER TO postgres;

--
-- TOC entry 265 (class 1259 OID 117271)
-- Name: user_shop_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_shop_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_shop_id_seq OWNER TO postgres;

--
-- TOC entry 5756 (class 0 OID 0)
-- Dependencies: 265
-- Name: user_shop_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_shop_id_seq OWNED BY public.user_shop.id;


--
-- TOC entry 260 (class 1259 OID 117152)
-- Name: user_state; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_state (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_country_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_state OWNER TO postgres;

--
-- TOC entry 259 (class 1259 OID 117151)
-- Name: user_state_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_state_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_state_id_seq OWNER TO postgres;

--
-- TOC entry 5757 (class 0 OID 0)
-- Dependencies: 259
-- Name: user_state_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_state_id_seq OWNED BY public.user_state.id;


--
-- TOC entry 272 (class 1259 OID 117362)
-- Name: user_supervisor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_supervisor (
    id integer NOT NULL,
    user_id integer NOT NULL,
    user_cell_id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_supervisor OWNER TO postgres;

--
-- TOC entry 271 (class 1259 OID 117361)
-- Name: user_supervisor_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_supervisor_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_supervisor_id_seq OWNER TO postgres;

--
-- TOC entry 5758 (class 0 OID 0)
-- Dependencies: 271
-- Name: user_supervisor_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_supervisor_id_seq OWNED BY public.user_supervisor.id;


--
-- TOC entry 228 (class 1259 OID 116580)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    mobile character varying(20),
    role_id integer NOT NULL,
    branch_id integer,
    company_id integer,
    status character varying(20) DEFAULT 'active'::character varying,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 116579)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5759 (class 0 OID 0)
-- Dependencies: 227
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 5045 (class 2604 OID 116633)
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- TOC entry 5066 (class 2604 OID 116735)
-- Name: auth_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens ALTER COLUMN id SET DEFAULT nextval('public.auth_tokens_id_seq'::regclass);


--
-- TOC entry 5069 (class 2604 OID 116776)
-- Name: branches id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches ALTER COLUMN id SET DEFAULT nextval('public.branches_id_seq'::regclass);


--
-- TOC entry 5078 (class 2604 OID 116869)
-- Name: brands id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.brands ALTER COLUMN id SET DEFAULT nextval('public.brands_id_seq'::regclass);


--
-- TOC entry 5188 (class 2604 OID 117960)
-- Name: breaks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.breaks ALTER COLUMN id SET DEFAULT nextval('public.breaks_id_seq'::regclass);


--
-- TOC entry 5074 (class 2604 OID 116852)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 5106 (class 2604 OID 117051)
-- Name: category_machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines ALTER COLUMN id SET DEFAULT nextval('public.category_machines_id_seq'::regclass);


--
-- TOC entry 5177 (class 2604 OID 117845)
-- Name: clienteles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clienteles ALTER COLUMN id SET DEFAULT nextval('public.clienteles_id_seq'::regclass);


--
-- TOC entry 5059 (class 2604 OID 116692)
-- Name: company_profiles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles ALTER COLUMN id SET DEFAULT nextval('public.company_profiles_id_seq'::regclass);


--
-- TOC entry 5091 (class 2604 OID 116931)
-- Name: controller_products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products ALTER COLUMN id SET DEFAULT nextval('public.controller_products_id_seq'::regclass);


--
-- TOC entry 5087 (class 2604 OID 116909)
-- Name: controllers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers ALTER COLUMN id SET DEFAULT nextval('public.controllers_id_seq'::regclass);


--
-- TOC entry 5034 (class 2604 OID 116516)
-- Name: countries id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries ALTER COLUMN id SET DEFAULT nextval('public.countries_id_seq'::regclass);


--
-- TOC entry 5193 (class 2604 OID 118014)
-- Name: holidays id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.holidays ALTER COLUMN id SET DEFAULT nextval('public.holidays_id_seq'::regclass);


--
-- TOC entry 5094 (class 2604 OID 116971)
-- Name: machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines ALTER COLUMN id SET DEFAULT nextval('public.machines_id_seq'::regclass);


--
-- TOC entry 5145 (class 2604 OID 117650)
-- Name: models id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models ALTER COLUMN id SET DEFAULT nextval('public.models_id_seq'::regclass);


--
-- TOC entry 5168 (class 2604 OID 117786)
-- Name: part_process_machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_process_machines ALTER COLUMN id SET DEFAULT nextval('public.part_process_machines_id_seq'::regclass);


--
-- TOC entry 5159 (class 2604 OID 117741)
-- Name: part_processes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_processes ALTER COLUMN id SET DEFAULT nextval('public.part_processes_id_seq'::regclass);


--
-- TOC entry 5171 (class 2604 OID 117812)
-- Name: part_sub_programs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_sub_programs ALTER COLUMN id SET DEFAULT nextval('public.part_sub_programs_id_seq'::regclass);


--
-- TOC entry 5152 (class 2604 OID 117709)
-- Name: parts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parts ALTER COLUMN id SET DEFAULT nextval('public.parts_id_seq'::regclass);


--
-- TOC entry 5082 (class 2604 OID 116886)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 5038 (class 2604 OID 116545)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 5182 (class 2604 OID 117893)
-- Name: shifts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shifts ALTER COLUMN id SET DEFAULT nextval('public.shifts_id_seq'::regclass);


--
-- TOC entry 5036 (class 2604 OID 116530)
-- Name: states id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states ALTER COLUMN id SET DEFAULT nextval('public.states_id_seq'::regclass);


--
-- TOC entry 5137 (class 2604 OID 117335)
-- Name: user_cell id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cell ALTER COLUMN id SET DEFAULT nextval('public.user_cell_id_seq'::regclass);


--
-- TOC entry 5121 (class 2604 OID 117185)
-- Name: user_city id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_city ALTER COLUMN id SET DEFAULT nextval('public.user_city_id_seq'::regclass);


--
-- TOC entry 5113 (class 2604 OID 117132)
-- Name: user_country id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_country ALTER COLUMN id SET DEFAULT nextval('public.user_country_id_seq'::regclass);


--
-- TOC entry 5109 (class 2604 OID 117108)
-- Name: user_filters id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters ALTER COLUMN id SET DEFAULT nextval('public.user_filters_id_seq'::regclass);


--
-- TOC entry 5133 (class 2604 OID 117305)
-- Name: user_line id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_line ALTER COLUMN id SET DEFAULT nextval('public.user_line_id_seq'::regclass);


--
-- TOC entry 5125 (class 2604 OID 117245)
-- Name: user_location id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_location ALTER COLUMN id SET DEFAULT nextval('public.user_location_id_seq'::regclass);


--
-- TOC entry 5149 (class 2604 OID 117683)
-- Name: user_machine_access id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machine_access ALTER COLUMN id SET DEFAULT nextval('public.user_machine_access_id_seq'::regclass);


--
-- TOC entry 5063 (class 2604 OID 116713)
-- Name: user_machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines ALTER COLUMN id SET DEFAULT nextval('public.user_machines_id_seq'::regclass);


--
-- TOC entry 5051 (class 2604 OID 116666)
-- Name: user_profiles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles ALTER COLUMN id SET DEFAULT nextval('public.user_profiles_id_seq'::regclass);


--
-- TOC entry 5129 (class 2604 OID 117275)
-- Name: user_shop id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_shop ALTER COLUMN id SET DEFAULT nextval('public.user_shop_id_seq'::regclass);


--
-- TOC entry 5117 (class 2604 OID 117155)
-- Name: user_state id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_state ALTER COLUMN id SET DEFAULT nextval('public.user_state_id_seq'::regclass);


--
-- TOC entry 5141 (class 2604 OID 117365)
-- Name: user_supervisor id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_supervisor ALTER COLUMN id SET DEFAULT nextval('public.user_supervisor_id_seq'::regclass);


--
-- TOC entry 5040 (class 2604 OID 116583)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 5650 (class 0 OID 116630)
-- Dependencies: 230
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.addresses (id, user_id, address_type, address_line1, address_line2, city, pin_code, contact_person, contact_phone, is_default, is_deleted, deleted_at, created_at, updated_at, country_name, state_name) FROM stdin;
2	23	work	3879 Shelia Cape	85145 S Washington Avenue	Rockwall	27947-5205	Kolby Skiles	+81394556979	t	f	\N	2026-10-14 10:34:33.413564	2026-10-14 10:34:33.413564	Japan	Hokkaido
3	35	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:00.916939	2026-10-14 12:07:00.916939	India	Gujarat
5	37	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:01.062939	2026-10-14 12:07:01.062939	India	Gujarat
6	38	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:01.136321	2026-10-14 12:07:01.136321	India	Gujarat
4	36	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	t	2026-10-14 12:07:01.13979+05:30	2026-10-14 12:07:00.988792	2026-10-14 12:07:01.13979	India	Gujarat
7	45	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:08.64863	2026-10-14 12:07:08.64863	India	Gujarat
9	47	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:08.783353	2026-10-14 12:07:08.783353	India	Gujarat
10	48	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:07:08.853665	2026-10-14 12:07:08.853665	India	Gujarat
8	46	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	t	2026-10-14 12:07:08.857401+05:30	2026-10-14 12:07:08.715634	2026-10-14 12:07:08.857401	India	Gujarat
11	49	work	5302 Consuelo View	37132 Yasmine Station	Santa Maria	88753	Carol Schultz	+13512223332	t	f	\N	2026-10-14 12:14:25.190123	2026-10-14 12:14:25.190123	United States	South Carolina
12	50	work	966 Haag Turnpike	14227 Stoltenberg Wells	Irving	51619-0825	Bernadette Runolfsdottir	+9715522717507	t	f	\N	2026-10-14 12:30:26.127054	2026-10-14 12:30:26.127054	United Arab Emirates	Sharjah
13	58	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:41:49.2054	2026-10-14 12:41:49.2054	India	Gujarat
15	60	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:41:49.363436	2026-10-14 12:41:49.363436	India	Gujarat
16	61	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 12:41:49.435932	2026-10-14 12:41:49.435932	India	Gujarat
14	59	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	t	2026-10-14 12:41:49.440609+05:30	2026-10-14 12:41:49.288914	2026-10-14 12:41:49.440609	India	Gujarat
17	102	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 16:09:28.932166	2026-10-14 16:09:28.932166	India	Gujarat
19	104	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 16:09:29.078007	2026-10-14 16:09:29.078007	India	Gujarat
20	105	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-14 16:09:29.152617	2026-10-14 16:09:29.152617	India	Gujarat
18	103	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	t	2026-10-14 16:09:29.157727+05:30	2026-10-14 16:09:29.004149	2026-10-14 16:09:29.157727	India	Gujarat
1	22	work	315 Emanuel Via	374 Gleichner Hills	Dublin	65893	Edison Stiedemann	+442896779618	t	f	\N	2026-10-14 10:26:06.7819	2026-10-15 09:56:20.930348	United Kingdom	Kensington and Chelsea
21	138	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-15 10:20:48.24788	2026-10-15 10:20:48.24788	India	Gujarat
23	140	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-15 10:20:48.399613	2026-10-15 10:20:48.399613	India	Gujarat
24	141	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	f	\N	2026-10-15 10:20:48.475144	2026-10-15 10:20:48.475144	India	Gujarat
22	139	work	Plot 12	\N	Rajkot	360001	Member	+919876543210	t	t	2026-10-15 10:20:48.479162+05:30	2026-10-15 10:20:48.322469	2026-10-15 10:20:48.479162	India	Gujarat
\.


--
-- TOC entry 5658 (class 0 OID 116732)
-- Dependencies: 238
-- Data for Name: auth_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.auth_tokens (id, user_id, token, token_type, device_type, device_info, expires_at, last_used_at, created_at, is_deleted, deleted_at) FROM stdin;
12	1	fcaa77af29697378698e5571ead1fa7ed0d57c68e2938587a30b88ec8a47a29e	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 22:12:54.924825	2026-10-14 10:12:55.001454	2026-10-14 10:12:54.924825	f	\N
9	18	7072ac9400216afdbfbf68cd4eda034eac7f9e73bf60330b3a2ce6c0b94df3ed	session	web	\N	2026-11-13 09:38:36.189374	2026-10-14 09:38:36.191205	2026-10-14 09:38:36.189374	f	\N
8	19	469b41c69cda8b887660b69a0c299895645e62e051c84e66a357dfb96984b10c	session	web	\N	2026-10-14 21:38:36.18725	2026-10-14 09:38:36.19063	2026-10-14 09:38:36.18725	t	2026-10-14 09:38:36.192558+05:30
14	1	2cdf00261e87e669cd382f3a43fca485bed0aff09eadffd926e4e2f1e7a3fe60	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 22:37:43.849997	2026-10-14 11:12:37.138049	2026-10-14 10:37:43.849997	f	\N
34	89	718ee654eff5741cd7634727a9c49e522e5e0699e895a9004580356b543824ae	session	web	\N	2026-10-15 04:09:28.281877	2026-10-14 16:09:28.282949	2026-10-14 16:09:28.281877	f	\N
5	1	06b953fce8f05a9f47347c58fe4a037b300948d0ba72a04a1850a8781a772f84	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 21:27:35.558493	2026-10-14 10:02:29.232511	2026-10-14 09:27:35.558493	t	2026-10-14 10:02:29.233492+05:30
13	1	d9e071666120b19d586db0bf1cce28b396f4369b6f2fa7f6a2f37bd06950c31c	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 22:15:48.938892	2026-10-14 10:37:26.713539	2026-10-14 10:15:48.938892	t	2026-10-14 10:37:40.368274+05:30
3	1	4f752910c5e7290a2f85d099a43605f254319a79536cee84c05fa605c0270b89	session	web	curl/8.14.1	2026-10-14 06:05:37.469379	2026-10-13 18:05:38.532844	2026-10-13 18:05:37.469379	t	2026-10-13 18:05:44.578833+05:30
4	1	09014a461321e4239704171e399d771b251483ffc7085365cb035cf2dc8f63c9	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 06:18:11.862948	2026-10-13 18:21:46.370041	2026-10-13 18:18:11.862948	t	2026-10-14 10:12:30.404709+05:30
2	1	12feea19ea5333c76fbf09f74eef5a7cb9e239f0e1842c7112f84e4343c8244f	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36	2026-10-14 05:50:00.818717	2026-10-13 18:22:07.10556	2026-10-13 17:50:00.818717	t	2026-10-14 10:12:30.404709+05:30
7	14	3001694a31162cacc517aaec53a605a766cc271d2ee235949a6f604c5e14aab4	session	web	\N	2026-11-13 09:37:37.593484	2026-10-14 09:37:37.600736	2026-10-14 09:37:37.593484	f	\N
6	15	7ec30b85be9a8aee87d67441637b5db3578dd54eefdb30fbfbc4355c34d3225d	session	web	\N	2026-10-14 21:37:37.587298	2026-10-14 09:37:37.599006	2026-10-14 09:37:37.587298	t	2026-10-14 09:37:37.605344+05:30
10	1	84c35637a4b50eb546b20b3ae930d6d3df9b4ea6c8233e832b2e54d9268eac8d	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 21:49:30.076821	2026-10-14 09:51:07.429134	2026-10-14 09:49:30.076821	t	2026-10-14 10:12:30.404709+05:30
11	1	eb42dd1495d5ccebb06acb03bdbfe34b162537d9eb0ac6d70f7510e50b4c7450	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 22:04:27.652845	2026-10-14 10:10:21.177119	2026-10-14 10:04:27.652845	t	2026-10-14 10:12:30.404709+05:30
21	30	3ce70093d41d68b02ce0bc0008b5be9e33f1ad6799a9f5739deadd1829caf28f	session	web	\N	2026-10-15 00:07:00.519418	2026-10-14 12:07:00.520134	2026-10-14 12:07:00.519418	f	\N
18	24	c478758cc685b68b1161d580762edb7cf49f47f70ba8a04ea90ede174f8e5584	session	web	\N	2026-11-13 11:57:29.956008	2026-10-14 11:57:29.959058	2026-10-14 11:57:29.956008	f	\N
17	25	f7396bbcd53908eee7120a200966bf53c40fc48e565d8f0f6b42ec1f7c10abd1	session	web	\N	2026-10-14 23:57:29.953529	2026-10-14 11:57:29.957845	2026-10-14 11:57:29.953529	t	2026-10-14 11:57:29.961115+05:30
30	1	7a875600be4d58db386433132e382279ff134f8f7bbcc057ddee467abd854458	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 02:30:30.119766	2026-10-14 18:38:18.329083	2026-10-14 14:30:30.119766	f	\N
20	28	4c986526e707dfdc945897c1131574118a2f1ec685113a980310eef628ec5daa	session	web	\N	2026-11-13 12:07:00.516183	2026-10-14 12:07:00.522682	2026-10-14 12:07:00.516183	f	\N
19	29	2fb1cdc256b1d4aeda1b399ba78bc20a29adaefe89ed28b9c3a27f17aea671d7	session	web	\N	2026-10-15 00:07:00.513764	2026-10-14 12:07:00.520992	2026-10-14 12:07:00.513764	t	2026-10-14 12:07:00.524074+05:30
15	22	e010aec1782f57b1390cd38b9b69c2da9bddf58a37ae8dea4c31d04d390eca2c	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 22:40:14.519254	2026-10-14 12:32:24.878728	2026-10-14 10:40:14.519254	t	2026-10-14 12:32:31.734892+05:30
24	41	e7a581d151992bf83fae44f93ab26cde4582a694dcf53ddf1d35ff7d65af8dd0	session	web	\N	2026-10-15 00:07:08.458744	2026-10-14 12:07:08.459263	2026-10-14 12:07:08.458744	f	\N
23	39	ae86ef6165c82146b9609c11d61baecade8b0b0e2c135f9b24d44425549bc5b4	session	web	\N	2026-11-13 12:07:08.455808	2026-10-14 12:07:08.461308	2026-10-14 12:07:08.455808	f	\N
22	40	3e56020486524d857e9d204d5997713a8c45827572b7b158a804e996f78a2368	session	web	\N	2026-10-15 00:07:08.450426	2026-10-14 12:07:08.45997	2026-10-14 12:07:08.450426	t	2026-10-14 12:07:08.462528+05:30
16	1	7ac90e682b86f5fbe3f8738dacd4b967c9f9fceddf0b3bb59518b1f7d15ee9de	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-14 23:20:21.409821	2026-10-14 12:14:34.398303	2026-10-14 11:20:21.409821	t	2026-10-14 12:14:38.180281+05:30
29	53	09d399fbd653f3805f86b4117242f45ccf517270129c01fac67cf6c4e82ad272	session	web	\N	2026-10-15 00:41:48.800459	2026-10-14 12:41:48.800935	2026-10-14 12:41:48.800459	f	\N
25	1	6d2821c6c3555acf56fb293898ca4ae19535b62f24fa9b1c802ade6fe04b4bf8	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 00:14:39.983091	2026-10-14 14:27:31.681589	2026-10-14 12:14:39.983091	f	\N
28	51	9bc54caaefd044e584386721680260862d51fdd27502ed0be585edc99bef5bc2	session	web	\N	2026-11-13 12:41:48.797745	2026-10-14 12:41:48.804391	2026-10-14 12:41:48.797745	f	\N
27	52	8f0917306aec95a0a2d925c8dac47ec43be61c7f21d80dce8a3a9467fbf24875	session	web	\N	2026-10-15 00:41:48.795597	2026-10-14 12:41:48.801708	2026-10-14 12:41:48.795597	t	2026-10-14 12:41:48.806936+05:30
31	1	43a1cdd80995d9dd083f067e8a5cf6714dadb8a967b0efc67b556063886104c4	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-11-13 16:03:51.46574	2026-10-14 17:22:38.259133	2026-10-14 16:03:51.46574	f	\N
36	1	e19fc5d68244181fefc6df7da62430e614a3d5e4636b69383f5192fb6bbd02a3	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 21:35:12.019213	2026-10-15 09:37:31.6674	2026-10-15 09:35:12.019213	f	\N
32	88	cc135cd524eadf1017ed1986e45d04cd2529995d98f56e4736de68a431d4b9c2	session	web	\N	2026-10-15 04:09:28.276386	2026-10-14 16:09:28.283594	2026-10-14 16:09:28.276386	t	2026-10-14 16:09:28.286223+05:30
33	87	968a7877a0a14fa882aae6e5321d9561217af1b2607f44599031a9c32dc85249	session	web	\N	2026-11-13 16:09:28.279358	2026-10-14 16:09:28.285015	2026-10-14 16:09:28.279358	f	\N
26	50	ae9140a46af2116c582e60b7e402b7293ba4753b90f1c8fe23eb1d3e6ae83ce8	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 00:32:35.469419	2026-10-14 16:24:25.079527	2026-10-14 12:32:35.469419	f	\N
39	125	7d42aff15f2ea3d6da3285d9cfa29789929b50fbfdfca222daf7399250faed15	session	web	\N	2026-10-15 22:20:47.48813	2026-10-15 10:20:47.488702	2026-10-15 10:20:47.48813	f	\N
38	123	179015c3eb85c9d399572fdc3da2274fd3cea9254c44cf3e30f82bba0d7a25f4	session	web	\N	2026-11-14 10:20:47.485095	2026-10-15 10:20:47.490917	2026-10-15 10:20:47.485095	f	\N
37	124	c2814f979f601cc3684097379db6b166e0eb6ef4791c3b3ec3a2a1bb4c73ceff	session	web	\N	2026-10-15 22:20:47.474828	2026-10-15 10:20:47.489488	2026-10-15 10:20:47.474828	t	2026-10-15 10:20:47.492223+05:30
35	1	e4b9f118a8c04b4c8deb59d678b9af11aa75c58c9b81380aa44bf75723e7d059	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 21:29:35.297494	2026-10-15 11:53:42.423237	2026-10-15 09:29:35.297494	f	\N
41	1	b2ca60563c5e91b81a70157df461c7f2df28f5c260e2379c64687a67e4ee87f2	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 23:54:07.770709	2026-10-15 12:06:28.462113	2026-10-15 11:54:07.770709	f	\N
40	1	89427e2a7147b8659e0b2b37822560e93f7d5c78d6662a7fc490754504af2455	session	web	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36	2026-10-15 23:42:43.492062	2026-10-15 11:55:13.932406	2026-10-15 11:42:43.492062	f	\N
\.


--
-- TOC entry 5660 (class 0 OID 116773)
-- Dependencies: 240
-- Data for Name: branches; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.branches (id, code, name, email, address, customer_id, gst_number, pan_number, color_code, status, is_deleted, deleted_at, created_at, updated_at, user_city_id) FROM stdin;
34	M2140857059	Branch 2140857059	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:28:02.186017+05:30	2026-10-14 15:28:02.15702	2026-10-14 15:28:02.15702	28
33	M214085705	Branch 214085705	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:28:02.190044+05:30	2026-10-14 15:28:02.145988	2026-10-14 15:28:02.145988	27
4	RAJKOT1	Metoda	rjt@gmail.com	3210 Lakin Landing	\N	27ABCDE1234F1Z5	CHSDP4933C	\N	active	f	\N	2026-10-14 10:19:18.986442	2026-10-14 10:19:18.986442	\N
6	GUJ002	Surat	surat@example.com	456 Textile Market, Surat	\N	27ABCDE5678F1Z5	BCDEF5678G	#33FF57	active	f	\N	2026-10-14 10:28:37.179377	2026-10-14 10:28:37.179377	\N
7	GUJ003	Vadodara	vadodara@example.com	789 Chemical Zone, Vadodara	\N	27ABCDE9012F1Z5	CDEF9012H	#3357FF	active	f	\N	2026-10-14 10:28:37.179377	2026-10-14 10:28:37.179377	\N
8	GUJ004	Rajkot	rajkot@example.com	321 Engineering Hub, Rajkot	\N	27ABCDE3456F1Z5	DEFG3456I	#FF33A8	active	f	\N	2026-10-14 10:28:37.179377	2026-10-14 10:28:37.179377	\N
9	GUJ005	Gandhinagar	gandhinagar@example.com	654 IT Park, Gandhinagar	\N	27ABCDE7890F1Z5	EFGH7890J	#A833FF	active	f	\N	2026-10-14 10:28:37.179377	2026-10-14 10:28:37.179377	\N
1	M93421215	Branch 93421215	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 09:37:21.063879+05:30	2026-10-14 09:37:20.987928	2026-10-14 09:37:20.987928	\N
2	M229469522	Branch 229469522	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 09:37:53.7365+05:30	2026-10-14 09:37:52.773447	2026-10-14 09:37:52.773447	\N
3	M631445522	Branch 631445522	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 09:38:36.382173+05:30	2026-10-14 09:38:36.327784	2026-10-14 09:38:36.327784	\N
5	GUJ001	Ahmedabad	ahmedabad@example.com	123 Industrial Area, Ahmedabad	\N	27ABCDE1234F1Z5	ABCDE1234F	#FF5733	active	f	\N	2026-10-14 10:28:37.179377	2026-10-14 10:28:37.179377	\N
10	M13336822	Branch 13336822	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 11:57:30.215036+05:30	2026-10-14 11:57:30.149761	2026-10-14 11:57:30.149761	4
11	M73394922	Branch 73394922	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:00.815654+05:30	2026-10-14 12:07:00.752829	2026-10-14 12:07:00.752829	5
13	MB83246528	Branch B83246528	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:01.143652+05:30	2026-10-14 12:07:00.854891	2026-10-14 12:07:00.854891	7
12	MA83246528	Branch A83246528	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:01.148346+05:30	2026-10-14 12:07:00.836484	2026-10-14 12:07:00.836484	6
14	M847312702	Branch 847312702	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:08.559035+05:30	2026-10-14 12:07:08.492769	2026-10-14 12:07:08.492769	8
16	MB857087506	Branch B857087506	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:08.861311+05:30	2026-10-14 12:07:08.588283	2026-10-14 12:07:08.588283	10
15	MA857087506	Branch A857087506	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:07:08.865396+05:30	2026-10-14 12:07:08.575369	2026-10-14 12:07:08.575369	9
17	M899315822	Branch 899315822	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:41:49.092367+05:30	2026-10-14 12:41:49.015764	2026-10-14 12:41:49.015764	11
19	MB911479328	Branch B911479328	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:41:49.445181+05:30	2026-10-14 12:41:49.142476	2026-10-14 12:41:49.142476	13
18	MA911479328	Branch A911479328	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:41:49.450537+05:30	2026-10-14 12:41:49.121047	2026-10-14 12:41:49.121047	12
21	M3279490019	Branch 3279490019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:55:13.592727+05:30	2026-10-14 12:55:13.495923	2026-10-14 12:55:13.495923	15
20	M327949001	Branch 327949001	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:55:13.600455+05:30	2026-10-14 12:55:13.384903	2026-10-14 12:55:13.384903	14
22	M769633603	Branch 769633603	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:55:27.77513+05:30	2026-10-14 12:55:27.718232	2026-10-14 12:55:27.718232	16
24	M7779789049	Branch 7779789049	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:55:27.833529+05:30	2026-10-14 12:55:27.798593	2026-10-14 12:55:27.798593	18
23	M777978904	Branch 777978904	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 12:55:27.83733+05:30	2026-10-14 12:55:27.783983	2026-10-14 12:55:27.783983	17
26	M161897019	Branch 161897019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:27:40.356995+05:30	2026-10-14 15:27:40.275382	2026-10-14 15:27:40.275382	20
25	M16189701	Branch 16189701	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:27:40.363542+05:30	2026-10-14 15:27:40.248811	2026-10-14 15:27:40.248811	19
27	M40066304	Branch 40066304	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:27:40.467176+05:30	2026-10-14 15:27:40.414007	2026-10-14 15:27:40.414007	21
29	M471232059	Branch 471232059	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:27:40.539863+05:30	2026-10-14 15:27:40.488006	2026-10-14 15:27:40.488006	23
28	M47123205	Branch 47123205	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:27:40.545445+05:30	2026-10-14 15:27:40.476499	2026-10-14 15:27:40.476499	22
31	M1932926019	Branch 1932926019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:28:02.062413+05:30	2026-10-14 15:28:01.996979	2026-10-14 15:28:01.996979	25
30	M193292601	Branch 193292601	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:28:02.068032+05:30	2026-10-14 15:28:01.976673	2026-10-14 15:28:01.976673	24
32	M209017604	Branch 209017604	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:28:02.137004+05:30	2026-10-14 15:28:02.095857	2026-10-14 15:28:02.095857	26
36	M1589422019	Branch 1589422019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:41:11.77537+05:30	2026-10-14 15:41:11.690267	2026-10-14 15:41:11.690267	30
35	M158942201	Branch 158942201	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:41:11.78157+05:30	2026-10-14 15:41:11.664866	2026-10-14 15:41:11.664866	29
38	M1785631029	Branch 1785631029	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:41:11.835998+05:30	2026-10-14 15:41:11.807658	2026-10-14 15:41:11.807658	32
37	M178563102	Branch 178563102	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:41:11.840118+05:30	2026-10-14 15:41:11.791659	2026-10-14 15:41:11.791659	31
40	M5730982019	Branch 5730982019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:51:05.964305+05:30	2026-10-14 15:51:05.83948	2026-10-14 15:51:05.83948	34
39	M573098201	Branch 573098201	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:51:05.992795+05:30	2026-10-14 15:51:05.811357	2026-10-14 15:51:05.811357	33
42	M6003831029	Branch 6003831029	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:51:06.108483+05:30	2026-10-14 15:51:06.039495	2026-10-14 15:51:06.039495	36
41	M600383102	Branch 600383102	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 15:51:06.11509+05:30	2026-10-14 15:51:06.014498	2026-10-14 15:51:06.014498	35
44	M8387817179	Branch 8387817179	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:28.479331+05:30	2026-10-14 16:09:28.427141	2026-10-14 16:09:28.427141	38
43	M838781717	Branch 838781717	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:28.484262+05:30	2026-10-14 16:09:28.405764	2026-10-14 16:09:28.405764	37
45	M851444423	Branch 851444423	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:28.558268+05:30	2026-10-14 16:09:28.520696	2026-10-14 16:09:28.520696	39
47	M8567715269	Branch 8567715269	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:28.616256+05:30	2026-10-14 16:09:28.587582	2026-10-14 16:09:28.587582	41
46	M856771526	Branch 856771526	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:28.620304+05:30	2026-10-14 16:09:28.575418	2026-10-14 16:09:28.575418	40
49	MB868524232	Branch B868524232	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:29.163102+05:30	2026-10-14 16:09:28.858336	2026-10-14 16:09:28.858336	43
48	MA868524232	Branch A868524232	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:09:29.168329+05:30	2026-10-14 16:09:28.689281	2026-10-14 16:09:28.689281	42
51	M5216632019	Branch 5216632019	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:26:15.70877+05:30	2026-10-14 16:26:15.586716	2026-10-14 16:26:15.586716	45
50	M521663201	Branch 521663201	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:26:15.71582+05:30	2026-10-14 16:26:15.531511	2026-10-14 16:26:15.531511	44
52	M572189502	Branch 572189502	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:26:15.771428+05:30	2026-10-14 16:26:15.727005	2026-10-14 16:26:15.727005	46
54	M5775058039	Branch 5775058039	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:26:15.815564+05:30	2026-10-14 16:26:15.791157	2026-10-14 16:26:15.791157	48
53	M577505803	Branch 577505803	\N	\N	\N	\N	\N	\N	active	t	2026-10-14 16:26:15.818937+05:30	2026-10-14 16:26:15.780442	2026-10-14 16:26:15.780442	47
56	M7640435179	Branch 7640435179	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:47.781062+05:30	2026-10-15 10:20:47.689741	2026-10-15 10:20:47.689741	52
55	M764043517	Branch 764043517	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:47.788079+05:30	2026-10-15 10:20:47.664028	2026-10-15 10:20:47.664028	51
57	M781488023	Branch 781488023	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:47.858299+05:30	2026-10-15 10:20:47.819327	2026-10-15 10:20:47.819327	53
59	M7870395269	Branch 7870395269	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:47.924504+05:30	2026-10-15 10:20:47.885866	2026-10-15 10:20:47.885866	55
58	M787039526	Branch 787039526	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:47.929044+05:30	2026-10-15 10:20:47.875233	2026-10-15 10:20:47.875233	54
61	MB801317832	Branch B801317832	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:48.483713+05:30	2026-10-15 10:20:48.184806	2026-10-15 10:20:48.184806	57
60	MA801317832	Branch A801317832	\N	\N	\N	\N	\N	\N	active	t	2026-10-15 10:20:48.488433+05:30	2026-10-15 10:20:48.018938	2026-10-15 10:20:48.018938	56
64	32619	Champaign	111252@gmail.com	7831 Runte River	49	28AFEUE5412D5ZM	OWUHW8044K	#570000	active	f	\N	2026-10-15 11:13:05.133698	2026-10-15 11:13:05.133698	59
66	3453432432423	Champaign	111252@gmail.com	7831 Runte River	49	28AFEUE5412D5ZM	OWUHW8044K	#570000	active	f	\N	2026-10-15 11:13:44.992638	2026-10-15 11:13:44.992638	59
\.


--
-- TOC entry 5664 (class 0 OID 116866)
-- Dependencies: 244
-- Data for Name: brands; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.brands (id, name, code, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
185	Brand 847312702	\N	\N	t	2026-10-14 12:07:08.563248+05:30	2026-10-14 12:07:08.474976	2026-10-14 12:07:08.563248
1	Test Brand 138251500 Updated	TC138251500	\N	t	2026-10-12 17:48:44.146918+05:30	2026-10-12 17:48:44.139025	2026-10-12 17:48:44.146918
49	Test Brand 53606100 Updated	TC53606100	\N	t	2026-10-13 15:14:27.086272+05:30	2026-10-13 15:14:27.054052	2026-10-13 15:14:27.086272
4	Test Brand 138251500 Updated	TC138251500	\N	t	2026-10-12 17:48:44.148185+05:30	2026-10-12 17:48:44.147868	2026-10-12 17:48:44.148185
111	Test Brand 74514501 Updated	TC74514501	\N	t	2026-10-14 09:37:20.843512+05:30	2026-10-14 09:37:20.747079	2026-10-14 09:37:20.843512
5	Test Brand Ref 162588400	\N	\N	t	2026-10-12 17:48:44.17353+05:30	2026-10-12 17:48:44.162812	2026-10-12 17:48:44.17353
52	Test Brand 53606100 Updated	TC53606100	\N	t	2026-10-13 15:14:27.087929+05:30	2026-10-13 15:14:27.087522	2026-10-13 15:14:27.087929
6	Test Brand Ref 174094300	\N	\N	t	2026-10-12 17:48:44.175785+05:30	2026-10-12 17:48:44.174538	2026-10-12 17:48:44.175785
53	Test Brand Ref 117393500	\N	\N	t	2026-10-13 15:14:27.143826+05:30	2026-10-13 15:14:27.117505	2026-10-13 15:14:27.143826
114	Test Brand 74514501 Updated	TC74514501	\N	t	2026-10-14 09:37:20.846709+05:30	2026-10-14 09:37:20.846154	2026-10-14 09:37:20.846709
54	Test Brand Ref 144523800	\N	\N	t	2026-10-13 15:14:27.146367+05:30	2026-10-13 15:14:27.144928	2026-10-13 15:14:27.146367
7	Test Brand 672113800 Updated	TC672113800	\N	t	2026-10-12 17:48:50.679916+05:30	2026-10-12 17:48:50.672868	2026-10-12 17:48:50.679916
10	Test Brand 672113800 Updated	TC672113800	\N	t	2026-10-12 17:48:50.681956+05:30	2026-10-12 17:48:50.681304	2026-10-12 17:48:50.681956
115	Test Brand Ref 84711802	\N	\N	t	2026-10-14 09:37:20.888587+05:30	2026-10-14 09:37:20.847422	2026-10-14 09:37:20.888587
11	Test Brand Ref 697801600	\N	\N	t	2026-10-12 17:48:50.711368+05:30	2026-10-12 17:48:50.698191	2026-10-12 17:48:50.711368
12	Test Brand Ref 712340000	\N	\N	t	2026-10-12 17:48:50.714369+05:30	2026-10-12 17:48:50.712745	2026-10-12 17:48:50.714369
55	Test Brand 728818800 Updated	TC728818800	\N	t	2026-10-13 15:48:34.754092+05:30	2026-10-13 15:48:34.729436	2026-10-13 15:48:34.754092
117	Test Brand Ref 89089307	\N	\N	t	2026-10-14 09:37:20.8949+05:30	2026-10-14 09:37:20.891376	2026-10-14 09:37:20.8949
58	Test Brand 728818800 Updated	TC728818800	\N	t	2026-10-13 15:48:34.755478+05:30	2026-10-13 15:48:34.755149	2026-10-13 15:48:34.755478
13	Test Brand 261232500 Updated	TC261232500	\N	t	2026-10-12 17:48:53.272033+05:30	2026-10-12 17:48:53.261771	2026-10-12 17:48:53.272033
116	Test Brand Ref 88938906	\N	\N	t	2026-10-14 09:37:20.895828+05:30	2026-10-14 09:37:20.89004	2026-10-14 09:37:20.895828
16	Test Brand 261232500 Updated	TC261232500	\N	t	2026-10-12 17:48:53.273392+05:30	2026-10-12 17:48:53.273057	2026-10-12 17:48:53.273392
59	Test Brand Ref 782521100	\N	\N	t	2026-10-13 15:48:34.803249+05:30	2026-10-13 15:48:34.782612	2026-10-13 15:48:34.803249
17	Test Brand Ref 285218700	\N	\N	t	2026-10-12 17:48:53.296425+05:30	2026-10-12 17:48:53.285774	2026-10-12 17:48:53.296425
18	Test Brand Ref 297068900	\N	\N	t	2026-10-12 17:48:53.298916+05:30	2026-10-12 17:48:53.297581	2026-10-12 17:48:53.298916
60	Test Brand Ref 804018600	\N	\N	t	2026-10-13 15:48:34.805781+05:30	2026-10-13 15:48:34.804335	2026-10-13 15:48:34.805781
118	Test Brand Ref 89657208	\N	\N	t	2026-10-14 09:37:20.898784+05:30	2026-10-14 09:37:20.896709	2026-10-14 09:37:20.898784
62	Other brand 394910600	\N	\N	t	2026-10-13 16:26:29.514863+05:30	2026-10-13 16:26:29.406804	2026-10-13 16:26:29.514863
19	Test Brand 906849900 Updated	TC906849900	\N	t	2026-10-12 17:48:55.914131+05:30	2026-10-12 17:48:55.907655	2026-10-12 17:48:55.914131
61	Brand 394910600	\N	\N	t	2026-10-13 16:26:29.51529+05:30	2026-10-13 16:26:29.40501	2026-10-13 16:26:29.51529
22	Test Brand 906849900 Updated	TC906849900	\N	t	2026-10-12 17:48:55.915584+05:30	2026-10-12 17:48:55.915222	2026-10-12 17:48:55.915584
119	Test Brand Ref 90113811	\N	\N	t	2026-10-14 09:37:20.919213+05:30	2026-10-14 09:37:20.901408	2026-10-14 09:37:20.919213
23	Test Brand Ref 928165900	\N	\N	t	2026-10-12 17:48:55.938048+05:30	2026-10-12 17:48:55.928399	2026-10-12 17:48:55.938048
24	Test Brand Ref 939142900	\N	\N	t	2026-10-12 17:48:55.942397+05:30	2026-10-12 17:48:55.939378	2026-10-12 17:48:55.942397
63	Fanuc	\N	\N	f	\N	2026-10-13 17:51:53.39891	2026-10-13 17:51:53.39891
25	Test Brand 601793700 Updated	TC601793700	\N	t	2026-10-13 11:20:45.610841+05:30	2026-10-13 11:20:45.602418	2026-10-13 11:20:45.610841
64	Siemens	\N	\N	f	\N	2026-10-13 17:51:53.39891	2026-10-13 17:51:53.39891
28	Test Brand 601793700 Updated	TC601793700	\N	t	2026-10-13 11:20:45.613197+05:30	2026-10-13 11:20:45.612789	2026-10-13 11:20:45.613197
65	Heidenhain	\N	\N	f	\N	2026-10-13 17:51:53.39891	2026-10-13 17:51:53.39891
29	Test Brand Ref 626772900	\N	\N	t	2026-10-13 11:20:45.637674+05:30	2026-10-13 11:20:45.627099	2026-10-13 11:20:45.637674
66	Mitsubishi	\N	\N	f	\N	2026-10-13 17:51:53.39891	2026-10-13 17:51:53.39891
30	Test Brand Ref 638294600	\N	\N	t	2026-10-13 11:20:45.640009+05:30	2026-10-13 11:20:45.638794	2026-10-13 11:20:45.640009
31	Test Brand 605661400 Updated	TC605661400	\N	t	2026-10-13 13:25:36.616594+05:30	2026-10-13 13:25:36.606453	2026-10-13 13:25:36.616594
67	Test Brand 206659600 Updated	TC206659600	\N	t	2026-10-13 18:00:14.21509+05:30	2026-10-13 18:00:14.207455	2026-10-13 18:00:14.21509
34	Test Brand 605661400 Updated	TC605661400	\N	t	2026-10-13 13:25:36.619665+05:30	2026-10-13 13:25:36.619122	2026-10-13 13:25:36.619665
35	Test Brand Ref 636772500	\N	\N	t	2026-10-13 13:25:36.648981+05:30	2026-10-13 13:25:36.636994	2026-10-13 13:25:36.648981
70	Test Brand 206659600 Updated	TC206659600	\N	t	2026-10-13 18:00:14.21665+05:30	2026-10-13 18:00:14.216184	2026-10-13 18:00:14.21665
36	Test Brand Ref 650172700	\N	\N	t	2026-10-13 13:25:36.652559+05:30	2026-10-13 13:25:36.650626	2026-10-13 13:25:36.652559
71	Test Brand Ref 239722300	\N	\N	t	2026-10-13 18:00:14.269001+05:30	2026-10-13 18:00:14.240057	2026-10-13 18:00:14.269001
37	Test Brand 302674900 Updated	TC302674900	\N	t	2026-10-13 13:25:45.309436+05:30	2026-10-13 13:25:45.303508	2026-10-13 13:25:45.309436
73	Test Brand Ref 270968700	\N	\N	t	2026-10-13 18:00:14.273825+05:30	2026-10-13 18:00:14.271305	2026-10-13 18:00:14.273825
40	Test Brand 302674900 Updated	TC302674900	\N	t	2026-10-13 13:25:45.310569+05:30	2026-10-13 13:25:45.31028	2026-10-13 13:25:45.310569
72	Test Brand Ref 269419000	\N	\N	t	2026-10-13 18:00:14.274751+05:30	2026-10-13 18:00:14.269931	2026-10-13 18:00:14.274751
41	Test Brand Ref 321222900	\N	\N	t	2026-10-13 13:25:45.330205+05:30	2026-10-13 13:25:45.32171	2026-10-13 13:25:45.330205
42	Test Brand Ref 330976600	\N	\N	t	2026-10-13 13:25:45.332265+05:30	2026-10-13 13:25:45.331161	2026-10-13 13:25:45.332265
74	Test Brand Ref 274932800	\N	\N	t	2026-10-13 18:00:14.277192+05:30	2026-10-13 18:00:14.2753	2026-10-13 18:00:14.277192
75	Test Brand Ref 278468100	\N	\N	t	2026-10-13 18:00:14.291796+05:30	2026-10-13 18:00:14.278619	2026-10-13 18:00:14.291796
43	Test Brand 61015300 Updated	TC61015300	\N	t	2026-10-13 14:57:56.069073+05:30	2026-10-13 14:57:56.061807	2026-10-13 14:57:56.069073
46	Test Brand 61015300 Updated	TC61015300	\N	t	2026-10-13 14:57:56.071041+05:30	2026-10-13 14:57:56.070513	2026-10-13 14:57:56.071041
76	Test Brand Ref 291690900	\N	\N	t	2026-10-13 18:00:14.293075+05:30	2026-10-13 18:00:14.292264	2026-10-13 18:00:14.293075
47	Test Brand Ref 87293500	\N	\N	t	2026-10-13 14:57:56.101766+05:30	2026-10-13 14:57:56.087818	2026-10-13 14:57:56.101766
48	Test Brand Ref 102595300	\N	\N	t	2026-10-13 14:57:56.1046+05:30	2026-10-13 14:57:56.102909	2026-10-13 14:57:56.1046
78	Test Brand 562404002 Updated	TC562404002	\N	t	2026-10-13 18:01:25.632977+05:30	2026-10-13 18:01:25.625011	2026-10-13 18:01:25.632977
81	Test Brand 562404002 Updated	TC562404002	\N	t	2026-10-13 18:01:25.634959+05:30	2026-10-13 18:01:25.634479	2026-10-13 18:01:25.634959
82	Test Brand Ref 564871308	\N	\N	t	2026-10-13 18:01:25.671071+05:30	2026-10-13 18:01:25.649146	2026-10-13 18:01:25.671071
84	Test Brand Ref 567262513	\N	\N	t	2026-10-13 18:01:25.675639+05:30	2026-10-13 18:01:25.673196	2026-10-13 18:01:25.675639
83	Test Brand Ref 567153212	\N	\N	t	2026-10-13 18:01:25.676565+05:30	2026-10-13 18:01:25.671969	2026-10-13 18:01:25.676565
85	Test Brand Ref 567663914	\N	\N	t	2026-10-13 18:01:25.679855+05:30	2026-10-13 18:01:25.67719	2026-10-13 18:01:25.679855
86	Test Brand Ref 568127917	\N	\N	t	2026-10-13 18:01:25.697407+05:30	2026-10-13 18:01:25.68182	2026-10-13 18:01:25.697407
88	Test Brand Ref 569834520	\N	\N	t	2026-10-13 18:01:25.717201+05:30	2026-10-13 18:01:25.698442	2026-10-13 18:01:25.717201
87	Test Brand Ref 569781219	\N	\N	t	2026-10-13 18:01:25.717667+05:30	2026-10-13 18:01:25.697967	2026-10-13 18:01:25.717667
89	Test Brand 966478102 Updated	TC966478102	\N	t	2026-10-13 18:01:39.672678+05:30	2026-10-13 18:01:39.665576	2026-10-13 18:01:39.672678
92	Test Brand 966478102 Updated	TC966478102	\N	t	2026-10-13 18:01:39.67429+05:30	2026-10-13 18:01:39.673863	2026-10-13 18:01:39.67429
93	Test Brand Ref 968643908	\N	\N	t	2026-10-13 18:01:39.708872+05:30	2026-10-13 18:01:39.686964	2026-10-13 18:01:39.708872
95	Test Brand Ref 971057213	\N	\N	t	2026-10-13 18:01:39.713207+05:30	2026-10-13 18:01:39.711027	2026-10-13 18:01:39.713207
94	Test Brand Ref 970954612	\N	\N	t	2026-10-13 18:01:39.714107+05:30	2026-10-13 18:01:39.709759	2026-10-13 18:01:39.714107
96	Test Brand Ref 971524614	\N	\N	t	2026-10-13 18:01:39.717498+05:30	2026-10-13 18:01:39.715568	2026-10-13 18:01:39.717498
97	Test Brand Ref 971846917	\N	\N	t	2026-10-13 18:01:39.732039+05:30	2026-10-13 18:01:39.718815	2026-10-13 18:01:39.732039
99	Test Brand Ref 973238120	\N	\N	t	2026-10-13 18:01:39.740587+05:30	2026-10-13 18:01:39.732895	2026-10-13 18:01:39.740587
98	Test Brand Ref 973238119	\N	\N	t	2026-10-13 18:01:39.740871+05:30	2026-10-13 18:01:39.732545	2026-10-13 18:01:39.740871
100	Test Brand 715317003 Updated	TC715317003	\N	t	2026-10-13 18:01:47.160797+05:30	2026-10-13 18:01:47.153763	2026-10-13 18:01:47.160797
103	Test Brand 715317003 Updated	TC715317003	\N	t	2026-10-13 18:01:47.162176+05:30	2026-10-13 18:01:47.161786	2026-10-13 18:01:47.162176
104	Test Brand Ref 717503409	\N	\N	t	2026-10-13 18:01:47.19759+05:30	2026-10-13 18:01:47.175648	2026-10-13 18:01:47.19759
106	Test Brand Ref 719956514	\N	\N	t	2026-10-13 18:01:47.202312+05:30	2026-10-13 18:01:47.199833	2026-10-13 18:01:47.202312
105	Test Brand Ref 719798313	\N	\N	t	2026-10-13 18:01:47.203337+05:30	2026-10-13 18:01:47.198537	2026-10-13 18:01:47.203337
107	Test Brand Ref 720315415	\N	\N	t	2026-10-13 18:01:47.206431+05:30	2026-10-13 18:01:47.203885	2026-10-13 18:01:47.206431
108	Test Brand Ref 720784618	\N	\N	t	2026-10-13 18:01:47.225215+05:30	2026-10-13 18:01:47.208362	2026-10-13 18:01:47.225215
110	Test Brand Ref 722616921	\N	\N	t	2026-10-13 18:01:47.235157+05:30	2026-10-13 18:01:47.226578	2026-10-13 18:01:47.235157
109	Test Brand Ref 722564620	\N	\N	t	2026-10-13 18:01:47.235436+05:30	2026-10-13 18:01:47.225999	2026-10-13 18:01:47.235436
186	Other brand 847312702	\N	\N	t	2026-10-14 12:07:08.56291+05:30	2026-10-14 12:07:08.476141	2026-10-14 12:07:08.56291
121	Test Brand Ref 91962114	\N	\N	t	2026-10-14 09:37:20.933767+05:30	2026-10-14 09:37:20.920152	2026-10-14 09:37:20.933767
120	Test Brand Ref 91962113	\N	\N	t	2026-10-14 09:37:20.934129+05:30	2026-10-14 09:37:20.919729	2026-10-14 09:37:20.934129
123	Other brand 93421215	\N	\N	t	2026-10-14 09:37:21.072829+05:30	2026-10-14 09:37:20.949927	2026-10-14 09:37:21.072829
122	Brand 93421215	\N	\N	t	2026-10-14 09:37:21.079633+05:30	2026-10-14 09:37:20.949576	2026-10-14 09:37:21.079633
124	Test Brand Ref 166674301	\N	\N	t	2026-10-14 09:37:31.714734+05:30	2026-10-14 09:37:31.668126	2026-10-14 09:37:31.714734
190	Other brand B857087506	\N	\N	t	2026-10-14 12:07:08.863264+05:30	2026-10-14 12:07:08.585953	2026-10-14 12:07:08.863264
126	Test Brand Ref 171984406	\N	\N	t	2026-10-14 09:37:31.726038+05:30	2026-10-14 09:37:31.720251	2026-10-14 09:37:31.726038
125	Test Brand Ref 171828505	\N	\N	t	2026-10-14 09:37:31.729525+05:30	2026-10-14 09:37:31.718781	2026-10-14 09:37:31.729525
189	Brand B857087506	\N	\N	t	2026-10-14 12:07:08.863537+05:30	2026-10-14 12:07:08.585739	2026-10-14 12:07:08.863537
127	Test Brand Ref 172990307	\N	\N	t	2026-10-14 09:37:31.737941+05:30	2026-10-14 09:37:31.730498	2026-10-14 09:37:31.737941
188	Other brand A857087506	\N	\N	t	2026-10-14 12:07:08.868157+05:30	2026-10-14 12:07:08.572855	2026-10-14 12:07:08.868157
128	Test Brand Ref 173831310	\N	\N	t	2026-10-14 09:37:31.766685+05:30	2026-10-14 09:37:31.738761	2026-10-14 09:37:31.766685
187	Brand A857087506	\N	\N	t	2026-10-14 12:07:08.868423+05:30	2026-10-14 12:07:08.572605	2026-10-14 12:07:08.868423
129	Test Brand 815311403 Updated	TC815311403	\N	t	2026-10-14 09:37:38.231717+05:30	2026-10-14 09:37:38.157867	2026-10-14 09:37:38.231717
132	Test Brand 815311403 Updated	TC815311403	\N	t	2026-10-14 09:37:38.294645+05:30	2026-10-14 09:37:38.283384	2026-10-14 09:37:38.294645
191	Test Brand 884581003 Updated	TC884581003	\N	t	2026-10-14 12:41:48.865379+05:30	2026-10-14 12:41:48.846381	2026-10-14 12:41:48.865379
133	Test Brand Ref 848741109	\N	\N	t	2026-10-14 09:37:38.609449+05:30	2026-10-14 09:37:38.488209	2026-10-14 09:37:38.609449
194	Test Brand 884581003 Updated	TC884581003	\N	t	2026-10-14 12:41:48.866698+05:30	2026-10-14 12:41:48.866347	2026-10-14 12:41:48.866698
135	Test Brand Ref 861533014	\N	\N	t	2026-10-14 09:37:39.998409+05:30	2026-10-14 09:37:38.61554	2026-10-14 09:37:39.998409
134	Test Brand Ref 861068913	\N	\N	t	2026-10-14 09:37:41.355936+05:30	2026-10-14 09:37:38.611027	2026-10-14 09:37:41.355936
136	Test Brand Ref 253905115	\N	\N	t	2026-10-14 09:37:46.49702+05:30	2026-10-14 09:37:42.539611	2026-10-14 09:37:46.49702
195	Test Brand Ref 889583909	\N	\N	t	2026-10-14 12:41:48.937882+05:30	2026-10-14 12:41:48.895998	2026-10-14 12:41:48.937882
137	Test Brand Ref 721624618	\N	\N	t	2026-10-14 09:37:51.49057+05:30	2026-10-14 09:37:47.217351	2026-10-14 09:37:51.49057
139	Test Brand Ref 168632821	\N	\N	t	2026-10-14 09:37:52.245191+05:30	2026-10-14 09:37:51.686638	2026-10-14 09:37:52.245191
138	Test Brand Ref 155984520	\N	\N	t	2026-10-14 09:37:52.262182+05:30	2026-10-14 09:37:51.560221	2026-10-14 09:37:52.262182
197	Test Brand Ref 893991614	\N	\N	t	2026-10-14 12:41:48.94328+05:30	2026-10-14 12:41:48.940267	2026-10-14 12:41:48.94328
196	Test Brand Ref 893835413	\N	\N	t	2026-10-14 12:41:48.944875+05:30	2026-10-14 12:41:48.938813	2026-10-14 12:41:48.944875
141	Other brand 229469522	\N	\N	t	2026-10-14 09:37:53.969355+05:30	2026-10-14 09:37:52.371637	2026-10-14 09:37:53.969355
140	Brand 229469522	\N	\N	t	2026-10-14 09:37:53.998175+05:30	2026-10-14 09:37:52.340279	2026-10-14 09:37:53.998175
198	Test Brand Ref 894557615	\N	\N	t	2026-10-14 12:41:48.948303+05:30	2026-10-14 12:41:48.945961	2026-10-14 12:41:48.948303
199	Test Brand Ref 894911518	\N	\N	t	2026-10-14 12:41:48.979018+05:30	2026-10-14 12:41:48.949642	2026-10-14 12:41:48.979018
142	Test Brand 621898603 Updated	TC621898603	\N	t	2026-10-14 09:38:36.226351+05:30	2026-10-14 09:38:36.219783	2026-10-14 09:38:36.226351
145	Test Brand 621898603 Updated	TC621898603	\N	t	2026-10-14 09:38:36.227665+05:30	2026-10-14 09:38:36.227343	2026-10-14 09:38:36.227665
146	Test Brand Ref 624292409	\N	\N	t	2026-10-14 09:38:36.273387+05:30	2026-10-14 09:38:36.2435	2026-10-14 09:38:36.273387
201	Test Brand Ref 897993821	\N	\N	t	2026-10-14 12:41:48.992687+05:30	2026-10-14 12:41:48.980009	2026-10-14 12:41:48.992687
200	Test Brand Ref 897943020	\N	\N	t	2026-10-14 12:41:48.993006+05:30	2026-10-14 12:41:48.979574	2026-10-14 12:41:48.993006
148	Test Brand Ref 627498114	\N	\N	t	2026-10-14 09:38:36.279363+05:30	2026-10-14 09:38:36.275468	2026-10-14 09:38:36.279363
147	Test Brand Ref 627375513	\N	\N	t	2026-10-14 09:38:36.28134+05:30	2026-10-14 09:38:36.274234	2026-10-14 09:38:36.28134
149	Test Brand Ref 628177915	\N	\N	t	2026-10-14 09:38:36.284409+05:30	2026-10-14 09:38:36.282092	2026-10-14 09:38:36.284409
150	Test Brand Ref 628534618	\N	\N	t	2026-10-14 09:38:36.301725+05:30	2026-10-14 09:38:36.285596	2026-10-14 09:38:36.301725
203	Other brand 899315822	\N	\N	t	2026-10-14 12:41:49.095518+05:30	2026-10-14 12:41:48.995288	2026-10-14 12:41:49.095518
202	Brand 899315822	\N	\N	t	2026-10-14 12:41:49.095801+05:30	2026-10-14 12:41:48.995009	2026-10-14 12:41:49.095801
152	Test Brand Ref 630237621	\N	\N	t	2026-10-14 09:38:36.314335+05:30	2026-10-14 09:38:36.302688	2026-10-14 09:38:36.314335
151	Test Brand Ref 630186520	\N	\N	t	2026-10-14 09:38:36.314625+05:30	2026-10-14 09:38:36.30227	2026-10-14 09:38:36.314625
154	Other brand 631445522	\N	\N	t	2026-10-14 09:38:36.386788+05:30	2026-10-14 09:38:36.316678	2026-10-14 09:38:36.386788
153	Brand 631445522	\N	\N	t	2026-10-14 09:38:36.387111+05:30	2026-10-14 09:38:36.316404	2026-10-14 09:38:36.387111
207	Other brand B911479328	\N	\N	t	2026-10-14 12:41:49.447665+05:30	2026-10-14 12:41:49.138209	2026-10-14 12:41:49.447665
206	Brand B911479328	\N	\N	t	2026-10-14 12:41:49.447953+05:30	2026-10-14 12:41:49.137488	2026-10-14 12:41:49.447953
155	Test Brand 622003 Updated	TC622003	\N	t	2026-10-14 11:57:30.026823+05:30	2026-10-14 11:57:30.008617	2026-10-14 11:57:30.026823
205	Other brand A911479328	\N	\N	t	2026-10-14 12:41:49.453949+05:30	2026-10-14 12:41:49.116831	2026-10-14 12:41:49.453949
158	Test Brand 622003 Updated	TC622003	\N	t	2026-10-14 11:57:30.028251+05:30	2026-10-14 11:57:30.027869	2026-10-14 11:57:30.028251
204	Brand A911479328	\N	\N	t	2026-10-14 12:41:49.454296+05:30	2026-10-14 12:41:49.116564	2026-10-14 12:41:49.454296
159	Test Brand Ref 5399309	\N	\N	t	2026-10-14 11:57:30.085625+05:30	2026-10-14 11:57:30.054518	2026-10-14 11:57:30.085625
161	Test Brand Ref 8736014	\N	\N	t	2026-10-14 11:57:30.091446+05:30	2026-10-14 11:57:30.087942	2026-10-14 11:57:30.091446
160	Test Brand Ref 8623713	\N	\N	t	2026-10-14 11:57:30.092173+05:30	2026-10-14 11:57:30.086464	2026-10-14 11:57:30.092173
162	Test Brand Ref 9228215	\N	\N	t	2026-10-14 11:57:30.09454+05:30	2026-10-14 11:57:30.09265	2026-10-14 11:57:30.09454
163	Test Brand Ref 9539818	\N	\N	t	2026-10-14 11:57:30.120749+05:30	2026-10-14 11:57:30.095788	2026-10-14 11:57:30.120749
211	Other brand 3279490019	\N	\N	t	2026-10-14 12:55:13.597173+05:30	2026-10-14 12:55:13.490421	2026-10-14 12:55:13.597173
210	Brand 3279490019	\N	\N	t	2026-10-14 12:55:13.597811+05:30	2026-10-14 12:55:13.489945	2026-10-14 12:55:13.597811
165	Test Brand Ref 12153221	\N	\N	t	2026-10-14 11:57:30.132854+05:30	2026-10-14 11:57:30.12163	2026-10-14 11:57:30.132854
164	Test Brand Ref 12098820	\N	\N	t	2026-10-14 11:57:30.133169+05:30	2026-10-14 11:57:30.121205	2026-10-14 11:57:30.133169
209	Other brand 327949001	\N	\N	t	2026-10-14 12:55:13.602917+05:30	2026-10-14 12:55:13.31296	2026-10-14 12:55:13.602917
208	Brand 327949001	\N	\N	t	2026-10-14 12:55:13.603201+05:30	2026-10-14 12:55:13.299833	2026-10-14 12:55:13.603201
167	Other brand 13336822	\N	\N	t	2026-10-14 11:57:30.218055+05:30	2026-10-14 11:57:30.135549	2026-10-14 11:57:30.218055
166	Brand 13336822	\N	\N	t	2026-10-14 11:57:30.218435+05:30	2026-10-14 11:57:30.135234	2026-10-14 11:57:30.218435
212	Test Brand Ref 765262502	\N	\N	t	2026-10-14 12:55:27.695979+05:30	2026-10-14 12:55:27.654634	2026-10-14 12:55:27.695979
168	Test Brand 56441303 Updated	TC56441303	\N	t	2026-10-14 12:07:00.584121+05:30	2026-10-14 12:07:00.565733	2026-10-14 12:07:00.584121
214	Other brand 769633603	\N	\N	t	2026-10-14 12:55:27.778007+05:30	2026-10-14 12:55:27.701941	2026-10-14 12:55:27.778007
171	Test Brand 56441303 Updated	TC56441303	\N	t	2026-10-14 12:07:00.586556+05:30	2026-10-14 12:07:00.585686	2026-10-14 12:07:00.586556
213	Brand 769633603	\N	\N	t	2026-10-14 12:55:27.778276+05:30	2026-10-14 12:55:27.701237	2026-10-14 12:55:27.778276
172	Test Brand Ref 61743809	\N	\N	t	2026-10-14 12:07:00.662703+05:30	2026-10-14 12:07:00.617571	2026-10-14 12:07:00.662703
174	Test Brand Ref 66465114	\N	\N	t	2026-10-14 12:07:00.669501+05:30	2026-10-14 12:07:00.665163	2026-10-14 12:07:00.669501
173	Test Brand Ref 66307213	\N	\N	t	2026-10-14 12:07:00.670772+05:30	2026-10-14 12:07:00.663613	2026-10-14 12:07:00.670772
175	Test Brand Ref 67107715	\N	\N	t	2026-10-14 12:07:00.676172+05:30	2026-10-14 12:07:00.67145	2026-10-14 12:07:00.676172
176	Test Brand Ref 68091318	\N	\N	t	2026-10-14 12:07:00.71697+05:30	2026-10-14 12:07:00.681381	2026-10-14 12:07:00.71697
218	Other brand 7779789049	\N	\N	t	2026-10-14 12:55:27.835563+05:30	2026-10-14 12:55:27.796356	2026-10-14 12:55:27.835563
217	Brand 7779789049	\N	\N	t	2026-10-14 12:55:27.835815+05:30	2026-10-14 12:55:27.796121	2026-10-14 12:55:27.835815
178	Test Brand Ref 71863721	\N	\N	t	2026-10-14 12:07:00.73347+05:30	2026-10-14 12:07:00.719037	2026-10-14 12:07:00.73347
177	Test Brand Ref 71813420	\N	\N	t	2026-10-14 12:07:00.73385+05:30	2026-10-14 12:07:00.718457	2026-10-14 12:07:00.73385
216	Other brand 777978904	\N	\N	t	2026-10-14 12:55:27.839465+05:30	2026-10-14 12:55:27.781467	2026-10-14 12:55:27.839465
215	Brand 777978904	\N	\N	t	2026-10-14 12:55:27.839718+05:30	2026-10-14 12:55:27.781238	2026-10-14 12:55:27.839718
180	Other brand 73394922	\N	\N	t	2026-10-14 12:07:00.818882+05:30	2026-10-14 12:07:00.736077	2026-10-14 12:07:00.818882
179	Brand 73394922	\N	\N	t	2026-10-14 12:07:00.819186+05:30	2026-10-14 12:07:00.735801	2026-10-14 12:07:00.819186
184	Other brand B83246528	\N	\N	t	2026-10-14 12:07:01.145904+05:30	2026-10-14 12:07:00.851683	2026-10-14 12:07:01.145904
183	Brand B83246528	\N	\N	t	2026-10-14 12:07:01.146179+05:30	2026-10-14 12:07:00.851458	2026-10-14 12:07:01.146179
182	Other brand A83246528	\N	\N	t	2026-10-14 12:07:01.151508+05:30	2026-10-14 12:07:00.834168	2026-10-14 12:07:01.151508
181	Brand A83246528	\N	\N	t	2026-10-14 12:07:01.151876+05:30	2026-10-14 12:07:00.83393	2026-10-14 12:07:01.151876
222	Other brand 161897019	\N	\N	t	2026-10-14 15:27:40.361055+05:30	2026-10-14 15:27:40.271994	2026-10-14 15:27:40.361055
221	Brand 161897019	\N	\N	t	2026-10-14 15:27:40.361453+05:30	2026-10-14 15:27:40.271691	2026-10-14 15:27:40.361453
220	Other brand 16189701	\N	\N	t	2026-10-14 15:27:40.36686+05:30	2026-10-14 15:27:40.19897	2026-10-14 15:27:40.36686
219	Brand 16189701	\N	\N	t	2026-10-14 15:27:40.36733+05:30	2026-10-14 15:27:40.185945	2026-10-14 15:27:40.36733
223	Test Brand Ref 36865503	\N	\N	t	2026-10-14 15:27:40.400373+05:30	2026-10-14 15:27:40.370242	2026-10-14 15:27:40.400373
224	Brand 40066304	\N	\N	t	2026-10-14 15:27:40.470048+05:30	2026-10-14 15:27:40.402742	2026-10-14 15:27:40.470048
225	Other brand 40066304	\N	\N	t	2026-10-14 15:27:40.469777+05:30	2026-10-14 15:27:40.408606	2026-10-14 15:27:40.469777
285	Other brand 5216632019	\N	\N	t	2026-10-14 16:26:15.713072+05:30	2026-10-14 16:26:15.580656	2026-10-14 16:26:15.713072
284	Brand 5216632019	\N	\N	t	2026-10-14 16:26:15.713553+05:30	2026-10-14 16:26:15.580022	2026-10-14 16:26:15.713553
283	Other brand 521663201	\N	\N	t	2026-10-14 16:26:15.720005+05:30	2026-10-14 16:26:15.289943	2026-10-14 16:26:15.720005
229	Other brand 471232059	\N	\N	t	2026-10-14 15:27:40.543418+05:30	2026-10-14 15:27:40.485331	2026-10-14 15:27:40.543418
228	Brand 471232059	\N	\N	t	2026-10-14 15:27:40.543738+05:30	2026-10-14 15:27:40.48509	2026-10-14 15:27:40.543738
227	Other brand 47123205	\N	\N	t	2026-10-14 15:27:40.547969+05:30	2026-10-14 15:27:40.47388	2026-10-14 15:27:40.547969
226	Brand 47123205	\N	\N	t	2026-10-14 15:27:40.548222+05:30	2026-10-14 15:27:40.47363	2026-10-14 15:27:40.548222
282	Brand 521663201	\N	\N	t	2026-10-14 16:26:15.720301+05:30	2026-10-14 16:26:15.270989	2026-10-14 16:26:15.720301
287	Other brand 572189502	\N	\N	t	2026-10-14 16:26:15.773962+05:30	2026-10-14 16:26:15.724019	2026-10-14 16:26:15.773962
286	Brand 572189502	\N	\N	t	2026-10-14 16:26:15.774273+05:30	2026-10-14 16:26:15.723764	2026-10-14 16:26:15.774273
291	Other brand 5775058039	\N	\N	t	2026-10-14 16:26:15.817496+05:30	2026-10-14 16:26:15.788054	2026-10-14 16:26:15.817496
290	Brand 5775058039	\N	\N	t	2026-10-14 16:26:15.817724+05:30	2026-10-14 16:26:15.787756	2026-10-14 16:26:15.817724
289	Other brand 577505803	\N	\N	t	2026-10-14 16:26:15.820723+05:30	2026-10-14 16:26:15.777482	2026-10-14 16:26:15.820723
288	Brand 577505803	\N	\N	t	2026-10-14 16:26:15.820914+05:30	2026-10-14 16:26:15.777166	2026-10-14 16:26:15.820914
233	Other brand 1932926019	\N	\N	t	2026-10-14 15:28:02.06589+05:30	2026-10-14 15:28:01.994343	2026-10-14 15:28:02.06589
232	Brand 1932926019	\N	\N	t	2026-10-14 15:28:02.066217+05:30	2026-10-14 15:28:01.994044	2026-10-14 15:28:02.066217
231	Other brand 193292601	\N	\N	t	2026-10-14 15:28:02.070549+05:30	2026-10-14 15:28:01.942887	2026-10-14 15:28:02.070549
230	Brand 193292601	\N	\N	t	2026-10-14 15:28:02.070818+05:30	2026-10-14 15:28:01.941294	2026-10-14 15:28:02.070818
234	Test Brand Ref 207266003	\N	\N	t	2026-10-14 15:28:02.089731+05:30	2026-10-14 15:28:02.073098	2026-10-14 15:28:02.089731
292	Test Brand 754258003 Updated	TC754258003	\N	t	2026-10-15 10:20:47.562986+05:30	2026-10-15 10:20:47.543245	2026-10-15 10:20:47.562986
236	Other brand 209017604	\N	\N	t	2026-10-14 15:28:02.139642+05:30	2026-10-14 15:28:02.09202	2026-10-14 15:28:02.139642
235	Brand 209017604	\N	\N	t	2026-10-14 15:28:02.139916+05:30	2026-10-14 15:28:02.091791	2026-10-14 15:28:02.139916
295	Test Brand 754258003 Updated	TC754258003	\N	t	2026-10-15 10:20:47.564463+05:30	2026-10-15 10:20:47.564026	2026-10-15 10:20:47.564463
296	Test Brand Ref 759715909	\N	\N	t	2026-10-15 10:20:47.631582+05:30	2026-10-15 10:20:47.597329	2026-10-15 10:20:47.631582
240	Other brand 2140857059	\N	\N	t	2026-10-14 15:28:02.188332+05:30	2026-10-14 15:28:02.154491	2026-10-14 15:28:02.188332
239	Brand 2140857059	\N	\N	t	2026-10-14 15:28:02.188627+05:30	2026-10-14 15:28:02.154256	2026-10-14 15:28:02.188627
238	Other brand 214085705	\N	\N	t	2026-10-14 15:28:02.191905+05:30	2026-10-14 15:28:02.143232	2026-10-14 15:28:02.191905
237	Brand 214085705	\N	\N	t	2026-10-14 15:28:02.192099+05:30	2026-10-14 15:28:02.143005	2026-10-14 15:28:02.192099
298	Test Brand Ref 763315614	\N	\N	t	2026-10-15 10:20:47.635794+05:30	2026-10-15 10:20:47.633523	2026-10-15 10:20:47.635794
297	Test Brand Ref 763209413	\N	\N	t	2026-10-15 10:20:47.636698+05:30	2026-10-15 10:20:47.632447	2026-10-15 10:20:47.636698
244	Other brand 1589422019	\N	\N	t	2026-10-14 15:41:11.779356+05:30	2026-10-14 15:41:11.686997	2026-10-14 15:41:11.779356
243	Brand 1589422019	\N	\N	t	2026-10-14 15:41:11.779731+05:30	2026-10-14 15:41:11.686627	2026-10-14 15:41:11.779731
242	Other brand 158942201	\N	\N	t	2026-10-14 15:41:11.783656+05:30	2026-10-14 15:41:11.620247	2026-10-14 15:41:11.783656
241	Brand 158942201	\N	\N	t	2026-10-14 15:41:11.783886+05:30	2026-10-14 15:41:11.606853	2026-10-14 15:41:11.783886
299	Test Brand Ref 763684815	\N	\N	t	2026-10-15 10:20:47.639921+05:30	2026-10-15 10:20:47.637164	2026-10-15 10:20:47.639921
248	Other brand 1785631029	\N	\N	t	2026-10-14 15:41:11.83829+05:30	2026-10-14 15:41:11.804998	2026-10-14 15:41:11.83829
247	Brand 1785631029	\N	\N	t	2026-10-14 15:41:11.838593+05:30	2026-10-14 15:41:11.804696	2026-10-14 15:41:11.838593
246	Other brand 178563102	\N	\N	t	2026-10-14 15:41:11.84291+05:30	2026-10-14 15:41:11.787772	2026-10-14 15:41:11.84291
245	Brand 178563102	\N	\N	t	2026-10-14 15:41:11.843132+05:30	2026-10-14 15:41:11.787367	2026-10-14 15:41:11.843132
303	Other brand 7640435179	\N	\N	t	2026-10-15 10:20:47.785774+05:30	2026-10-15 10:20:47.686857	2026-10-15 10:20:47.785774
302	Brand 7640435179	\N	\N	t	2026-10-15 10:20:47.78614+05:30	2026-10-15 10:20:47.686532	2026-10-15 10:20:47.78614
301	Other brand 764043517	\N	\N	t	2026-10-15 10:20:47.790511+05:30	2026-10-15 10:20:47.642675	2026-10-15 10:20:47.790511
252	Other brand 5730982019	\N	\N	t	2026-10-14 15:51:05.986069+05:30	2026-10-14 15:51:05.834625	2026-10-14 15:51:05.986069
251	Brand 5730982019	\N	\N	t	2026-10-14 15:51:05.986886+05:30	2026-10-14 15:51:05.833777	2026-10-14 15:51:05.986886
250	Other brand 573098201	\N	\N	t	2026-10-14 15:51:05.996497+05:30	2026-10-14 15:51:05.782569	2026-10-14 15:51:05.996497
249	Brand 573098201	\N	\N	t	2026-10-14 15:51:05.996846+05:30	2026-10-14 15:51:05.740046	2026-10-14 15:51:05.996846
300	Brand 764043517	\N	\N	t	2026-10-15 10:20:47.790832+05:30	2026-10-15 10:20:47.64244	2026-10-15 10:20:47.790832
304	Test Brand Ref 779250819	\N	\N	t	2026-10-15 10:20:47.805744+05:30	2026-10-15 10:20:47.7931	2026-10-15 10:20:47.805744
256	Other brand 6003831029	\N	\N	t	2026-10-14 15:51:06.112418+05:30	2026-10-14 15:51:06.033718	2026-10-14 15:51:06.112418
255	Brand 6003831029	\N	\N	t	2026-10-14 15:51:06.112816+05:30	2026-10-14 15:51:06.033395	2026-10-14 15:51:06.112816
254	Other brand 600383102	\N	\N	t	2026-10-14 15:51:06.121717+05:30	2026-10-14 15:51:06.006347	2026-10-14 15:51:06.121717
253	Brand 600383102	\N	\N	t	2026-10-14 15:51:06.121971+05:30	2026-10-14 15:51:06.005669	2026-10-14 15:51:06.121971
306	Test Brand Ref 780626722	\N	\N	t	2026-10-15 10:20:47.814716+05:30	2026-10-15 10:20:47.806534	2026-10-15 10:20:47.814716
305	Test Brand Ref 780558821	\N	\N	t	2026-10-15 10:20:47.814978+05:30	2026-10-15 10:20:47.806165	2026-10-15 10:20:47.814978
257	Test Brand 832013803 Updated	TC832013803	\N	t	2026-10-14 16:09:28.328396+05:30	2026-10-14 16:09:28.320565	2026-10-14 16:09:28.328396
260	Test Brand 832013803 Updated	TC832013803	\N	t	2026-10-14 16:09:28.329724+05:30	2026-10-14 16:09:28.329363	2026-10-14 16:09:28.329724
308	Other brand 781488023	\N	\N	t	2026-10-15 10:20:47.860817+05:30	2026-10-15 10:20:47.816933	2026-10-15 10:20:47.860817
261	Test Brand Ref 835719509	\N	\N	t	2026-10-14 16:09:28.379696+05:30	2026-10-14 16:09:28.357328	2026-10-14 16:09:28.379696
307	Brand 781488023	\N	\N	t	2026-10-15 10:20:47.861116+05:30	2026-10-15 10:20:47.816663	2026-10-15 10:20:47.861116
263	Test Brand Ref 838115414	\N	\N	t	2026-10-14 16:09:28.383594+05:30	2026-10-14 16:09:28.381424	2026-10-14 16:09:28.383594
262	Test Brand Ref 838011713	\N	\N	t	2026-10-14 16:09:28.384497+05:30	2026-10-14 16:09:28.380454	2026-10-14 16:09:28.384497
264	Test Brand Ref 838469315	\N	\N	t	2026-10-14 16:09:28.386762+05:30	2026-10-14 16:09:28.384986	2026-10-14 16:09:28.386762
312	Other brand 7870395269	\N	\N	t	2026-10-15 10:20:47.927136+05:30	2026-10-15 10:20:47.883139	2026-10-15 10:20:47.927136
311	Brand 7870395269	\N	\N	t	2026-10-15 10:20:47.927463+05:30	2026-10-15 10:20:47.882808	2026-10-15 10:20:47.927463
268	Other brand 8387817179	\N	\N	t	2026-10-14 16:09:28.482642+05:30	2026-10-14 16:09:28.423827	2026-10-14 16:09:28.482642
267	Brand 8387817179	\N	\N	t	2026-10-14 16:09:28.482913+05:30	2026-10-14 16:09:28.423419	2026-10-14 16:09:28.482913
266	Other brand 838781717	\N	\N	t	2026-10-14 16:09:28.488364+05:30	2026-10-14 16:09:28.389743	2026-10-14 16:09:28.488364
265	Brand 838781717	\N	\N	t	2026-10-14 16:09:28.488632+05:30	2026-10-14 16:09:28.389512	2026-10-14 16:09:28.488632
310	Other brand 787039526	\N	\N	t	2026-10-15 10:20:47.931196+05:30	2026-10-15 10:20:47.872457	2026-10-15 10:20:47.931196
269	Test Brand Ref 849088319	\N	\N	t	2026-10-14 16:09:28.503365+05:30	2026-10-14 16:09:28.491202	2026-10-14 16:09:28.503365
309	Brand 787039526	\N	\N	t	2026-10-15 10:20:47.931471+05:30	2026-10-15 10:20:47.872157	2026-10-15 10:20:47.931471
271	Test Brand Ref 850367622	\N	\N	t	2026-10-14 16:09:28.513924+05:30	2026-10-14 16:09:28.504089	2026-10-14 16:09:28.513924
270	Test Brand Ref 850367621	\N	\N	t	2026-10-14 16:09:28.514808+05:30	2026-10-14 16:09:28.503769	2026-10-14 16:09:28.514808
273	Other brand 851444423	\N	\N	t	2026-10-14 16:09:28.560709+05:30	2026-10-14 16:09:28.518033	2026-10-14 16:09:28.560709
272	Brand 851444423	\N	\N	t	2026-10-14 16:09:28.561006+05:30	2026-10-14 16:09:28.517775	2026-10-14 16:09:28.561006
316	Other brand B801317832	\N	\N	t	2026-10-15 10:20:48.486073+05:30	2026-10-15 10:20:48.025031	2026-10-15 10:20:48.486073
315	Brand B801317832	\N	\N	t	2026-10-15 10:20:48.486343+05:30	2026-10-15 10:20:48.024764	2026-10-15 10:20:48.486343
314	Other brand A801317832	\N	\N	t	2026-10-15 10:20:48.490328+05:30	2026-10-15 10:20:48.014605	2026-10-15 10:20:48.490328
277	Other brand 8567715269	\N	\N	t	2026-10-14 16:09:28.618827+05:30	2026-10-14 16:09:28.584846	2026-10-14 16:09:28.618827
276	Brand 8567715269	\N	\N	t	2026-10-14 16:09:28.619073+05:30	2026-10-14 16:09:28.584576	2026-10-14 16:09:28.619073
275	Other brand 856771526	\N	\N	t	2026-10-14 16:09:28.622144+05:30	2026-10-14 16:09:28.570617	2026-10-14 16:09:28.622144
274	Brand 856771526	\N	\N	t	2026-10-14 16:09:28.622342+05:30	2026-10-14 16:09:28.570194	2026-10-14 16:09:28.622342
313	Brand A801317832	\N	\N	t	2026-10-15 10:20:48.49052+05:30	2026-10-15 10:20:48.014411	2026-10-15 10:20:48.49052
281	Other brand B868524232	\N	\N	t	2026-10-14 16:09:29.165824+05:30	2026-10-14 16:09:28.697143	2026-10-14 16:09:29.165824
280	Brand B868524232	\N	\N	t	2026-10-14 16:09:29.166133+05:30	2026-10-14 16:09:28.696886	2026-10-14 16:09:29.166133
279	Other brand A868524232	\N	\N	t	2026-10-14 16:09:29.170912+05:30	2026-10-14 16:09:28.686946	2026-10-14 16:09:29.170912
278	Brand A868524232	\N	\N	t	2026-10-14 16:09:29.171152+05:30	2026-10-14 16:09:28.686666	2026-10-14 16:09:29.171152
\.


--
-- TOC entry 5710 (class 0 OID 117957)
-- Dependencies: 290
-- Data for Name: breaks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.breaks (id, user_id, name, start_time, end_time, status, is_deleted, deleted_at, created_at, updated_at, shift_id) FROM stdin;
6	96	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-14 16:09:28.629059	2026-10-14 16:09:28.629059	\N
9	96	Snack	02:05:00	02:20:00	active	f	\N	2026-10-14 16:09:28.640049	2026-10-14 16:09:28.640049	\N
10	97	Other	01:00:00	01:10:00	active	f	\N	2026-10-14 16:09:28.641122	2026-10-14 16:09:28.641122	\N
7	96	N-Tea	03:00:00	03:15:00	active	t	2026-10-14 16:09:28.657955+05:30	2026-10-14 16:09:28.637627	2026-10-14 16:09:28.657955	\N
8	96	Lunch	13:00:00	13:30:00	active	t	2026-10-14 16:09:28.661971+05:30	2026-10-14 16:09:28.638895	2026-10-14 16:09:28.661971	\N
3	85	End	05:50:00	06:00:00	active	f	\N	2026-10-14 16:08:20.200574	2026-10-14 16:25:27.182204	3
5	86	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-14 16:08:20.200574	2026-10-14 16:25:27.182204	4
4	85	Start	22:00:00	22:10:00	active	f	\N	2026-10-14 16:08:20.200574	2026-10-14 16:25:27.182204	3
1	22	989	19:27:00	21:29:00	active	f	\N	2026-10-14 16:08:20.200574	2026-10-14 16:25:27.182204	5
15	111	Snack	02:05:00	02:20:00	active	f	\N	2026-10-14 16:26:15.837176	2026-10-14 16:26:15.837176	\N
16	112	Other	01:00:00	01:10:00	active	f	\N	2026-10-14 16:26:15.838358	2026-10-14 16:26:15.838358	\N
12	111	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-14 16:26:15.827706	2026-10-14 16:26:15.860479	\N
13	111	N-Tea	03:00:00	03:15:00	active	f	\N	2026-10-14 16:26:15.834888	2026-10-14 16:26:15.860479	\N
18	111	Late tea	04:00:00	04:10:00	active	t	2026-10-14 16:26:15.859315+05:30	2026-10-14 16:26:15.85665	2026-10-14 16:26:15.860479	\N
14	111	Lunch	13:00:00	13:30:00	active	f	\N	2026-10-14 16:26:15.836094	2026-10-14 16:26:15.864865	\N
22	115	Snack	02:05:00	02:20:00	active	f	\N	2026-10-14 16:26:31.016278	2026-10-14 16:26:31.016278	\N
23	116	Other	01:00:00	01:10:00	active	f	\N	2026-10-14 16:26:31.017848	2026-10-14 16:26:31.017848	\N
19	115	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-14 16:26:31.002761	2026-10-14 16:26:31.039865	\N
20	115	N-Tea	03:00:00	03:15:00	active	f	\N	2026-10-14 16:26:31.013267	2026-10-14 16:26:31.039865	\N
25	115	Late tea	04:00:00	04:10:00	active	t	2026-10-14 16:26:31.039027+05:30	2026-10-14 16:26:31.037619	2026-10-14 16:26:31.039865	\N
21	115	Lunch	13:00:00	13:30:00	active	f	\N	2026-10-14 16:26:31.014783	2026-10-14 16:26:31.043887	\N
29	117	Snack	02:05:00	02:20:00	active	f	\N	2026-10-14 16:26:41.23764	2026-10-14 16:26:41.23764	\N
30	118	Other	01:00:00	01:10:00	active	f	\N	2026-10-14 16:26:41.239163	2026-10-14 16:26:41.239163	\N
35	132	N-Tea	03:00:00	03:15:00	active	f	\N	2026-10-15 10:20:47.954936	2026-10-15 10:20:47.979512	\N
26	117	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-14 16:26:41.224633	2026-10-14 16:26:41.262064	\N
27	117	N-Tea	03:00:00	03:15:00	active	f	\N	2026-10-14 16:26:41.234587	2026-10-14 16:26:41.262064	\N
32	117	Late tea	04:00:00	04:10:00	active	t	2026-10-14 16:26:41.260743+05:30	2026-10-14 16:26:41.259266	2026-10-14 16:26:41.262064	\N
28	117	Lunch	13:00:00	13:30:00	active	f	\N	2026-10-14 16:26:41.236154	2026-10-14 16:26:41.265472	\N
33	22	HELLO	10:30:00	11:30:00	active	f	\N	2026-10-14 16:30:06.507946	2026-10-14 16:30:06.507946	5
2	22	Lunch	05:30:00	05:56:00	inactive	f	\N	2026-10-14 16:08:20.200574	2026-10-14 16:35:08.276958	5
37	132	Snack	02:05:00	02:20:00	active	f	\N	2026-10-15 10:20:47.957409	2026-10-15 10:20:47.957409	\N
38	133	Other	01:00:00	01:10:00	active	f	\N	2026-10-15 10:20:47.95871	2026-10-15 10:20:47.95871	\N
40	132	Late tea	04:00:00	04:10:00	active	t	2026-10-15 10:20:47.978218+05:30	2026-10-15 10:20:47.976816	2026-10-15 10:20:47.979512	\N
34	132	Dinner	23:30:00	00:15:00	active	f	\N	2026-10-15 10:20:47.935963	2026-10-15 10:20:47.979512	\N
36	132	Lunch	13:00:00	13:30:00	active	f	\N	2026-10-15 10:20:47.956243	2026-10-15 10:20:47.982587	\N
\.


--
-- TOC entry 5662 (class 0 OID 116849)
-- Dependencies: 242
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, code, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Category 121986100 Updated	TC121986100	\N	t	2026-10-12 17:48:44.136062+05:30	2026-10-12 17:48:44.124319	2026-10-12 17:48:44.136062
51	Test Category 173660100 Updated	TC173660100	\N	t	2026-10-13 18:00:14.204968+05:30	2026-10-13 18:00:14.175274	2026-10-13 18:00:14.204968
4	Test Category 121986100 Updated	TC121986100	\N	t	2026-10-12 17:48:44.138245+05:30	2026-10-12 17:48:44.137075	2026-10-12 17:48:44.138245
5	Test Category Ref 148465800	\N	\N	t	2026-10-12 17:48:44.160127+05:30	2026-10-12 17:48:44.148591	2026-10-12 17:48:44.160127
54	Test Category 173660100 Updated	TC173660100	\N	t	2026-10-13 18:00:14.206549+05:30	2026-10-13 18:00:14.206132	2026-10-13 18:00:14.206549
55	Test Category Ref 216880400	\N	\N	t	2026-10-13 18:00:14.237417+05:30	2026-10-13 18:00:14.21719	2026-10-13 18:00:14.237417
6	Test Category 662000300 Updated	TC662000300	\N	t	2026-10-12 17:48:50.67064+05:30	2026-10-12 17:48:50.663036	2026-10-12 17:48:50.67064
9	Test Category 662000300 Updated	TC662000300	\N	t	2026-10-12 17:48:50.672009+05:30	2026-10-12 17:48:50.671633	2026-10-12 17:48:50.672009
10	Test Category Ref 682313000	\N	\N	t	2026-10-12 17:48:50.695098+05:30	2026-10-12 17:48:50.682786	2026-10-12 17:48:50.695098
56	Test Category 561037401 Updated	TC561037401	\N	t	2026-10-13 18:01:25.622396+05:30	2026-10-13 18:01:25.612086	2026-10-13 18:01:25.622396
59	Test Category 561037401 Updated	TC561037401	\N	t	2026-10-13 18:01:25.623963+05:30	2026-10-13 18:01:25.623549	2026-10-13 18:01:25.623963
11	Test Category 250434800 Updated	TC250434800	\N	t	2026-10-12 17:48:53.259031+05:30	2026-10-12 17:48:53.251486	2026-10-12 17:48:53.259031
14	Test Category 250434800 Updated	TC250434800	\N	t	2026-10-12 17:48:53.260799+05:30	2026-10-12 17:48:53.260328	2026-10-12 17:48:53.260799
60	Test Category Ref 563496103	\N	\N	t	2026-10-13 18:01:25.646686+05:30	2026-10-13 18:01:25.635534	2026-10-13 18:01:25.646686
15	Test Category Ref 273755100	\N	\N	t	2026-10-12 17:48:53.2836+05:30	2026-10-12 17:48:53.27392	2026-10-12 17:48:53.2836
16	Test Category 895651600 Updated	TC895651600	\N	t	2026-10-12 17:48:55.90516+05:30	2026-10-12 17:48:55.896477	2026-10-12 17:48:55.90516
61	Test Category 965147001 Updated	TC965147001	\N	t	2026-10-13 18:01:39.663056+05:30	2026-10-13 18:01:39.65306	2026-10-13 18:01:39.663056
19	Test Category 895651600 Updated	TC895651600	\N	t	2026-10-12 17:48:55.906783+05:30	2026-10-12 17:48:55.906387	2026-10-12 17:48:55.906783
20	Test Category Ref 915794800	\N	\N	t	2026-10-12 17:48:55.926216+05:30	2026-10-12 17:48:55.91606	2026-10-12 17:48:55.926216
64	Test Category 965147001 Updated	TC965147001	\N	t	2026-10-13 18:01:39.664538+05:30	2026-10-13 18:01:39.664117	2026-10-13 18:01:39.664538
65	Test Category Ref 967432703	\N	\N	t	2026-10-13 18:01:39.684486+05:30	2026-10-13 18:01:39.674898	2026-10-13 18:01:39.684486
21	Test Category 582972500 Updated	TC582972500	\N	t	2026-10-13 11:20:45.600236+05:30	2026-10-13 11:20:45.5841	2026-10-13 11:20:45.600236
24	Test Category 582972500 Updated	TC582972500	\N	t	2026-10-13 11:20:45.601671+05:30	2026-10-13 11:20:45.60127	2026-10-13 11:20:45.601671
25	Test Category Ref 613390200	\N	\N	t	2026-10-13 11:20:45.625141+05:30	2026-10-13 11:20:45.613693	2026-10-13 11:20:45.625141
66	Test Category 714280902 Updated	TC714280902	\N	t	2026-10-13 18:01:47.15139+05:30	2026-10-13 18:01:47.143883	2026-10-13 18:01:47.15139
69	Test Category 714280902 Updated	TC714280902	\N	t	2026-10-13 18:01:47.152969+05:30	2026-10-13 18:01:47.152556	2026-10-13 18:01:47.152969
26	Test Category 592941600 Updated	TC592941600	\N	t	2026-10-13 13:25:36.603961+05:30	2026-10-13 13:25:36.593886	2026-10-13 13:25:36.603961
29	Test Category 592941600 Updated	TC592941600	\N	t	2026-10-13 13:25:36.605616+05:30	2026-10-13 13:25:36.605178	2026-10-13 13:25:36.605616
70	Test Category Ref 716261604	\N	\N	t	2026-10-13 18:01:47.173279+05:30	2026-10-13 18:01:47.162814	2026-10-13 18:01:47.173279
30	Test Category Ref 620100000	\N	\N	t	2026-10-13 13:25:36.633599+05:30	2026-10-13 13:25:36.620434	2026-10-13 13:25:36.633599
31	Test Category 291934700 Updated	TC291934700	\N	t	2026-10-13 13:25:45.300931+05:30	2026-10-13 13:25:45.292518	2026-10-13 13:25:45.300931
71	Test Category 763789402 Updated	TC763789402	\N	t	2026-10-14 09:37:37.742662+05:30	2026-10-14 09:37:37.639941	2026-10-14 09:37:37.742662
34	Test Category 291934700 Updated	TC291934700	\N	t	2026-10-13 13:25:45.302597+05:30	2026-10-13 13:25:45.302155	2026-10-13 13:25:45.302597
35	Test Category Ref 310467400	\N	\N	t	2026-10-13 13:25:45.3198+05:30	2026-10-13 13:25:45.311016	2026-10-13 13:25:45.3198
74	Test Category 763789402 Updated	TC763789402	\N	t	2026-10-14 09:37:38.13528+05:30	2026-10-14 09:37:38.126422	2026-10-14 09:37:38.13528
75	Test Category Ref 829975504	\N	\N	t	2026-10-14 09:37:38.474028+05:30	2026-10-14 09:37:38.300115	2026-10-14 09:37:38.474028
36	Test Category 47782400 Updated	TC47782400	\N	t	2026-10-13 14:57:56.059579+05:30	2026-10-13 14:57:56.048584	2026-10-13 14:57:56.059579
39	Test Category 47782400 Updated	TC47782400	\N	t	2026-10-13 14:57:56.061087+05:30	2026-10-13 14:57:56.060681	2026-10-13 14:57:56.061087
40	Test Category Ref 71207900	\N	\N	t	2026-10-13 14:57:56.082683+05:30	2026-10-13 14:57:56.071705	2026-10-13 14:57:56.082683
76	Test Category 620695402 Updated	TC620695402	\N	t	2026-10-14 09:38:36.217314+05:30	2026-10-14 09:38:36.207774	2026-10-14 09:38:36.217314
79	Test Category 620695402 Updated	TC620695402	\N	t	2026-10-14 09:38:36.218744+05:30	2026-10-14 09:38:36.21834	2026-10-14 09:38:36.218744
41	Test Category 42296300 Updated	TC42296300	\N	t	2026-10-13 15:14:27.051811+05:30	2026-10-13 15:14:27.042744	2026-10-13 15:14:27.051811
44	Test Category 42296300 Updated	TC42296300	\N	t	2026-10-13 15:14:27.053139+05:30	2026-10-13 15:14:27.052787	2026-10-13 15:14:27.053139
80	Test Category Ref 622788604	\N	\N	t	2026-10-14 09:38:36.240715+05:30	2026-10-14 09:38:36.228345	2026-10-14 09:38:36.240715
45	Test Category Ref 88364700	\N	\N	t	2026-10-13 15:14:27.115251+05:30	2026-10-13 15:14:27.088558	2026-10-13 15:14:27.115251
46	Test Category 710835100 Updated	TC710835100	\N	t	2026-10-13 15:48:34.727186+05:30	2026-10-13 15:48:34.711484	2026-10-13 15:48:34.727186
81	Test Category 997618802 Updated	TC997618802	\N	t	2026-10-14 11:57:30.003079+05:30	2026-10-14 11:57:29.978967	2026-10-14 11:57:30.003079
49	Test Category 710835100 Updated	TC710835100	\N	t	2026-10-13 15:48:34.728723+05:30	2026-10-13 15:48:34.728296	2026-10-13 15:48:34.728723
50	Test Category Ref 755702700	\N	\N	t	2026-10-13 15:48:34.779807+05:30	2026-10-13 15:48:34.755933	2026-10-13 15:48:34.779807
84	Test Category 997618802 Updated	TC997618802	\N	t	2026-10-14 11:57:30.005631+05:30	2026-10-14 11:57:30.005085	2026-10-14 11:57:30.005631
85	Test Category Ref 2833804	\N	\N	t	2026-10-14 11:57:30.051957+05:30	2026-10-14 11:57:30.028849	2026-10-14 11:57:30.051957
86	Test Category 54064502 Updated	TC54064502	\N	t	2026-10-14 12:07:00.561975+05:30	2026-10-14 12:07:00.541146	2026-10-14 12:07:00.561975
89	Test Category 54064502 Updated	TC54064502	\N	t	2026-10-14 12:07:00.564135+05:30	2026-10-14 12:07:00.563361	2026-10-14 12:07:00.564135
90	Test Category Ref 58690204	\N	\N	t	2026-10-14 12:07:00.613864+05:30	2026-10-14 12:07:00.587366	2026-10-14 12:07:00.613864
91	Test Category 881750402 Updated	TC881750402	\N	t	2026-10-14 12:41:48.843949+05:30	2026-10-14 12:41:48.818546	2026-10-14 12:41:48.843949
94	Test Category 881750402 Updated	TC881750402	\N	t	2026-10-14 12:41:48.845482+05:30	2026-10-14 12:41:48.844954	2026-10-14 12:41:48.845482
95	Test Category Ref 886697504	\N	\N	t	2026-10-14 12:41:48.89176+05:30	2026-10-14 12:41:48.8675	2026-10-14 12:41:48.89176
96	Test Category 829638802 Updated	TC829638802	\N	t	2026-10-14 16:09:28.318577+05:30	2026-10-14 16:09:28.296749	2026-10-14 16:09:28.318577
99	Test Category 829638802 Updated	TC829638802	\N	t	2026-10-14 16:09:28.319845+05:30	2026-10-14 16:09:28.319527	2026-10-14 16:09:28.319845
100	Test Category Ref 833006004	\N	\N	t	2026-10-14 16:09:28.353365+05:30	2026-10-14 16:09:28.330184	2026-10-14 16:09:28.353365
101	Test Category 750248402 Updated	TC750248402	\N	t	2026-10-15 10:20:47.541185+05:30	2026-10-15 10:20:47.503311	2026-10-15 10:20:47.541185
104	Test Category 750248402 Updated	TC750248402	\N	t	2026-10-15 10:20:47.542415+05:30	2026-10-15 10:20:47.542053	2026-10-15 10:20:47.542415
105	Test Category Ref 756454004	\N	\N	t	2026-10-15 10:20:47.593334+05:30	2026-10-15 10:20:47.565071	2026-10-15 10:20:47.593334
\.


--
-- TOC entry 5674 (class 0 OID 117048)
-- Dependencies: 254
-- Data for Name: category_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.category_machines (id, category_id, machine_id, is_deleted, deleted_at, created_at) FROM stdin;
\.


--
-- TOC entry 5706 (class 0 OID 117842)
-- Dependencies: 286
-- Data for Name: clienteles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clienteles (id, user_id, client_code, name, email, contact_number, address, country_name, state_name, city, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
29	107	CL-521663201	Client CL-521663201	client.cl-521663201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:26:15.608256	2026-10-14 16:26:15.608256
5	68	OWN-16189701	Client OWN-16189701	client.own-16189701@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	f	\N	2026-10-14 15:27:40.30874	2026-10-14 15:27:40.30874
1	67	CL-16189701	Renamed 16189701	client.cl-16189701@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	inactive	t	2026-10-14 15:27:40.354329+05:30	2026-10-14 15:27:40.279384	2026-10-14 15:27:40.354329
9	72	NE-193292601	Client NE-193292601	\N	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	f	\N	2026-10-14 15:28:02.010173	2026-10-14 15:28:02.010173
30	106	NE-521663201	Client NE-521663201	\N	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:26:15.610604	2026-10-14 16:26:15.610604
10	73	OWN-193292601	Client OWN-193292601	client.own-193292601@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	f	\N	2026-10-14 15:28:02.019429	2026-10-14 15:28:02.019429
6	72	CL-193292601	Renamed 193292601	client.cl-193292601@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	inactive	t	2026-10-14 15:28:02.059675+05:30	2026-10-14 15:28:01.999683	2026-10-14 15:28:02.059675
11	49	78966	Lockman	153734@gmail.com	+814073094509	93214 Odessa Knolls	JP	35	Las Vegas	active	f	\N	2026-10-14 15:37:43.685303	2026-10-14 15:37:43.685303
8	73	CL-193292601	Client CL-193292601	client.cl-193292601@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	t	2026-10-14 15:37:48.565435+05:30	2026-10-14 15:28:02.008568	2026-10-14 15:37:48.565435
3	23	CL-16189701	Client CL-16189701	client.cl-16189701@example.com	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	t	2026-10-14 15:37:50.34705+05:30	2026-10-14 15:27:40.296104	2026-10-14 15:37:50.34705
4	67	NE-16189701	Client NE-16189701	\N	+919876543210	Plot 12, GIDC	IN	GJ	Rajkot	active	t	2026-10-14 15:37:52.059229+05:30	2026-10-14 15:27:40.298098	2026-10-14 15:37:52.059229
15	77	NE-158942201	Client NE-158942201	\N	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:41:11.707785	2026-10-14 15:41:11.707785
16	78	OWN-158942201	Client OWN-158942201	client.own-158942201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:41:11.717534	2026-10-14 15:41:11.717534
12	77	CL-158942201	Renamed 158942201	client.cl-158942201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	inactive	t	2026-10-14 15:41:11.772851+05:30	2026-10-14 15:41:11.69395	2026-10-14 15:41:11.772851
19	82	CL-573098201	Client CL-573098201	client.cl-573098201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:51:05.862141	2026-10-14 15:51:05.862141
20	81	NE-573098201	Client NE-573098201	\N	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:51:05.866151	2026-10-14 15:51:05.866151
31	107	OWN-521663201	Client OWN-521663201	client.own-521663201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:26:15.622181	2026-10-14 16:26:15.622181
21	82	OWN-573098201	Client OWN-573098201	client.own-573098201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:51:05.889271	2026-10-14 15:51:05.889271
17	81	CL-573098201	Renamed 573098201	client.cl-573098201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	inactive	t	2026-10-14 15:51:05.959319+05:30	2026-10-14 15:51:05.846086	2026-10-14 15:51:05.959319
24	91	CL-838781717	Client CL-838781717	client.cl-838781717@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:09:28.438497	2026-10-14 16:09:28.438497
25	90	NE-838781717	Client NE-838781717	\N	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:09:28.440065	2026-10-14 16:09:28.440065
27	106	CL-521663201	Renamed 521663201	client.cl-521663201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	inactive	t	2026-10-14 16:26:15.705883+05:30	2026-10-14 16:26:15.590682	2026-10-14 16:26:15.705883
26	91	OWN-838781717	Client OWN-838781717	client.own-838781717@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 16:09:28.44773	2026-10-14 16:09:28.44773
22	90	CL-838781717	Renamed 838781717	client.cl-838781717@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	inactive	t	2026-10-14 16:09:28.477221+05:30	2026-10-14 16:09:28.429655	2026-10-14 16:09:28.477221
14	23	CL-158942201	Client CL-158942201	client.cl-158942201@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-14 15:41:11.705651	2026-10-14 17:16:25.464558
32	23	53914-6710	Kiehn-Moore	172252@gmail.com	+441697084890	7798 Union Avenue	United Kingdom	Powys	Raleigh	active	f	\N	2026-10-14 17:22:58.271185	2026-10-14 17:22:58.271185
33	49	08556-9746	Okuneva	174506@gmail.com	+9718176753154	344 Russel Park	United Arab Emirates	Dubai	North Little Rock	active	f	\N	2026-10-14 17:45:08.66453	2026-10-14 17:45:08.66453
34	22	21954-8480	Kassulke	174645@gmail.com	+15686960695	528 Lionel Junctions	United States	Vermont	Greenwood	active	f	\N	2026-10-14 17:46:50.23219	2026-10-14 17:46:50.23219
35	22	29482-4065	Green	174719@gmail.com	+9716814554001	9940 Telly Ridges	United Arab Emirates	Ras al Khaimah	Daytona Beach	active	f	\N	2026-10-14 17:47:27.419318	2026-10-14 17:47:27.419318
36	22	53936-8264	Flatley	174816@gmail.com	+811279125792	1330 W Maple Street	Japan	Yamanashi	San Angelo	active	f	\N	2026-10-14 17:48:19.947263	2026-10-14 17:48:19.947263
39	127	CL-764043517	Client CL-764043517	client.cl-764043517@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-15 10:20:47.706919	2026-10-15 10:20:47.706919
40	126	NE-764043517	Client NE-764043517	\N	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-15 10:20:47.708821	2026-10-15 10:20:47.708821
41	127	OWN-764043517	Client OWN-764043517	client.own-764043517@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	active	f	\N	2026-10-15 10:20:47.725501	2026-10-15 10:20:47.725501
37	126	CL-764043517	Renamed 764043517	client.cl-764043517@example.com	+919876543210	Plot 12, GIDC	India	Gujarat	Rajkot	inactive	t	2026-10-15 10:20:47.778784+05:30	2026-10-15 10:20:47.693345	2026-10-15 10:20:47.778784
\.


--
-- TOC entry 5654 (class 0 OID 116689)
-- Dependencies: 234
-- Data for Name: company_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.company_profiles (id, user_id, gst_number, pan_number, website, finix_id, is_deleted, deleted_at, created_at, updated_at, company_name) FROM stdin;
2	23	09LMHKR6500K5ZT	QLTBB5174O	https://example.com/eshknhhe	99458498	f	\N	2026-10-14 10:34:33.413564	2026-10-14 10:34:33.413564	Cartwright Inc
3	35	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:00.916939	2026-10-14 12:07:00.916939	Member adm83246528
5	37	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:01.062939	2026-10-14 12:07:01.062939	Member sup83246528
6	38	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:01.136321	2026-10-14 12:07:01.136321	Member sa83246528
4	36	24AAACE1234F1Z5	AAACE1234F	\N	\N	t	2026-10-14 12:07:01.13979+05:30	2026-10-14 12:07:00.988792	2026-10-14 12:07:01.13979	Member op83246528
7	45	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:08.64863	2026-10-14 12:07:08.64863	Member adm857087506
9	47	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:08.783353	2026-10-14 12:07:08.783353	Member sup857087506
10	48	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:07:08.853665	2026-10-14 12:07:08.853665	Member sa857087506
8	46	24AAACE1234F1Z5	AAACE1234F	\N	\N	t	2026-10-14 12:07:08.857401+05:30	2026-10-14 12:07:08.715634	2026-10-14 12:07:08.857401	Member op857087506
11	49	27SFLQS3642W2ZB	QNKFS6343F	https://example.com/xqhi3cg7	87058404	f	\N	2026-10-14 12:14:25.190123	2026-10-14 12:14:25.190123	Bogan - Nikolaus
12	50	07OTMTM3107U7ZL	ZUVLP6399C	https://example.com/voogoysp	50458847	f	\N	2026-10-14 12:30:26.127054	2026-10-14 12:30:26.127054	Jakubowski Group
13	58	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:41:49.2054	2026-10-14 12:41:49.2054	Member adm911479328
15	60	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:41:49.363436	2026-10-14 12:41:49.363436	Member sup911479328
16	61	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 12:41:49.435932	2026-10-14 12:41:49.435932	Member sa911479328
14	59	24AAACE1234F1Z5	AAACE1234F	\N	\N	t	2026-10-14 12:41:49.440609+05:30	2026-10-14 12:41:49.288914	2026-10-14 12:41:49.440609	Member op911479328
17	102	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 16:09:28.932166	2026-10-14 16:09:28.932166	Member adm868524232
19	104	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 16:09:29.078007	2026-10-14 16:09:29.078007	Member sup868524232
20	105	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-14 16:09:29.152617	2026-10-14 16:09:29.152617	Member sa868524232
18	103	24AAACE1234F1Z5	AAACE1234F	\N	\N	t	2026-10-14 16:09:29.157727+05:30	2026-10-14 16:09:29.004149	2026-10-14 16:09:29.157727	Member op868524232
1	22	19IMUEL4295E4ZK	JYOFF6545C	https://example.com/cq2kaifn	12345678	f	\N	2026-10-14 10:26:06.7819	2026-10-15 09:56:20.930348	Dare - Kirlin
21	138	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-15 10:20:48.24788	2026-10-15 10:20:48.24788	Member adm801317832
23	140	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-15 10:20:48.399613	2026-10-15 10:20:48.399613	Member sup801317832
24	141	24AAACE1234F1Z5	AAACE1234F	\N	\N	f	\N	2026-10-15 10:20:48.475144	2026-10-15 10:20:48.475144	Member sa801317832
22	139	24AAACE1234F1Z5	AAACE1234F	\N	\N	t	2026-10-15 10:20:48.479162+05:30	2026-10-15 10:20:48.322469	2026-10-15 10:20:48.479162	Member op801317832
\.


--
-- TOC entry 5670 (class 0 OID 116928)
-- Dependencies: 250
-- Data for Name: controller_products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.controller_products (id, controller_id, product_id, is_deleted, deleted_at, created_at) FROM stdin;
\.


--
-- TOC entry 5668 (class 0 OID 116906)
-- Dependencies: 248
-- Data for Name: controllers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.controllers (id, name, brand_id, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
151	FANUC 1	63	\N	f	\N	2026-10-14 10:36:36.123204	2026-10-14 10:36:36.123204
152	SIEMENS - 1	64	\N	f	\N	2026-10-14 10:36:46.657914	2026-10-14 10:36:46.657914
153	Test Controller 5454610 Updated	159	\N	t	2026-10-14 11:57:30.082354+05:30	2026-10-14 11:57:30.055559	2026-10-14 11:57:30.082354
155	Test Controller 5454610 Updated	159	\N	t	2026-10-14 11:57:30.08525+05:30	2026-10-14 11:57:30.084128	2026-10-14 11:57:30.08525
206	Other controller 3279490019	211	\N	t	2026-10-14 12:55:13.595953+05:30	2026-10-14 12:55:13.491441	2026-10-14 12:55:13.595953
205	Controller 3279490019	210	\N	t	2026-10-14 12:55:13.596392+05:30	2026-10-14 12:55:13.490885	2026-10-14 12:55:13.596392
157	Test Controller 8623712	161	\N	t	2026-10-14 11:57:30.091094+05:30	2026-10-14 11:57:30.089691	2026-10-14 11:57:30.091094
156	Test Controller 8623712	160	\N	t	2026-10-14 11:57:30.091889+05:30	2026-10-14 11:57:30.087128	2026-10-14 11:57:30.091889
204	Other controller 327949001	209	\N	t	2026-10-14 12:55:13.60234+05:30	2026-10-14 12:55:13.339233	2026-10-14 12:55:13.60234
158	Test Controller 9282916 840D	162	\N	t	2026-10-14 11:57:30.09419+05:30	2026-10-14 11:57:30.093202	2026-10-14 11:57:30.09419
203	Controller 327949001	208	\N	t	2026-10-14 12:55:13.602632+05:30	2026-10-14 12:55:13.313781	2026-10-14 12:55:13.602632
160	Other Model Controller 9539817	163	\N	t	2026-10-14 11:57:30.119895+05:30	2026-10-14 11:57:30.097054	2026-10-14 11:57:30.119895
159	Model Controller 9539817	163	\N	t	2026-10-14 11:57:30.120415+05:30	2026-10-14 11:57:30.096336	2026-10-14 11:57:30.120415
208	Other Model Controller 765262501	212	\N	t	2026-10-14 12:55:27.694929+05:30	2026-10-14 12:55:27.663145	2026-10-14 12:55:27.694929
207	Model Controller 765262501	212	\N	t	2026-10-14 12:55:27.695547+05:30	2026-10-14 12:55:27.660518	2026-10-14 12:55:27.695547
162	Spare Controller 12098819	164	\N	t	2026-10-14 11:57:30.131973+05:30	2026-10-14 11:57:30.122472	2026-10-14 11:57:30.131973
161	Renamed Controller 12098819	164	\N	t	2026-10-14 11:57:30.132485+05:30	2026-10-14 11:57:30.122021	2026-10-14 11:57:30.132485
164	Other controller 13336822	167	\N	t	2026-10-14 11:57:30.217293+05:30	2026-10-14 11:57:30.136215	2026-10-14 11:57:30.217293
163	Controller 13336822	166	\N	t	2026-10-14 11:57:30.217574+05:30	2026-10-14 11:57:30.135933	2026-10-14 11:57:30.217574
210	Other controller 769633603	214	\N	t	2026-10-14 12:55:27.777397+05:30	2026-10-14 12:55:27.703651	2026-10-14 12:55:27.777397
209	Controller 769633603	213	\N	t	2026-10-14 12:55:27.77767+05:30	2026-10-14 12:55:27.703065	2026-10-14 12:55:27.77767
165	Test Controller 61794610 Updated	172	\N	t	2026-10-14 12:07:00.65904+05:30	2026-10-14 12:07:00.618842	2026-10-14 12:07:00.65904
167	Test Controller 61794610 Updated	172	\N	t	2026-10-14 12:07:00.662266+05:30	2026-10-14 12:07:00.661146	2026-10-14 12:07:00.662266
169	Test Controller 66307212	174	\N	t	2026-10-14 12:07:00.668984+05:30	2026-10-14 12:07:00.666686	2026-10-14 12:07:00.668984
168	Test Controller 66307212	173	\N	t	2026-10-14 12:07:00.670292+05:30	2026-10-14 12:07:00.664272	2026-10-14 12:07:00.670292
214	Other controller 7779789049	218	\N	t	2026-10-14 12:55:27.835005+05:30	2026-10-14 12:55:27.7968	2026-10-14 12:55:27.835005
170	Test Controller 67168516 840D	175	\N	t	2026-10-14 12:07:00.674911+05:30	2026-10-14 12:07:00.672293	2026-10-14 12:07:00.674911
213	Controller 7779789049	217	\N	t	2026-10-14 12:55:27.835311+05:30	2026-10-14 12:55:27.796569	2026-10-14 12:55:27.835311
212	Other controller 777978904	216	\N	t	2026-10-14 12:55:27.838957+05:30	2026-10-14 12:55:27.781975	2026-10-14 12:55:27.838957
172	Other Model Controller 68041017	176	\N	t	2026-10-14 12:07:00.715799+05:30	2026-10-14 12:07:00.686922	2026-10-14 12:07:00.715799
171	Model Controller 68041017	176	\N	t	2026-10-14 12:07:00.716531+05:30	2026-10-14 12:07:00.682926	2026-10-14 12:07:00.716531
211	Controller 777978904	215	\N	t	2026-10-14 12:55:27.839219+05:30	2026-10-14 12:55:27.781676	2026-10-14 12:55:27.839219
174	Spare Controller 71813419	177	\N	t	2026-10-14 12:07:00.731294+05:30	2026-10-14 12:07:00.720475	2026-10-14 12:07:00.731294
173	Renamed Controller 71813419	177	\N	t	2026-10-14 12:07:00.732996+05:30	2026-10-14 12:07:00.7196	2026-10-14 12:07:00.732996
176	Other controller 73394922	180	\N	t	2026-10-14 12:07:00.818114+05:30	2026-10-14 12:07:00.739493	2026-10-14 12:07:00.818114
175	Controller 73394922	179	\N	t	2026-10-14 12:07:00.818477+05:30	2026-10-14 12:07:00.73906	2026-10-14 12:07:00.818477
218	Other controller 161897019	222	\N	t	2026-10-14 15:27:40.360106+05:30	2026-10-14 15:27:40.272779	2026-10-14 15:27:40.360106
217	Controller 161897019	221	\N	t	2026-10-14 15:27:40.360537+05:30	2026-10-14 15:27:40.272305	2026-10-14 15:27:40.360537
216	Other controller 16189701	220	\N	t	2026-10-14 15:27:40.365898+05:30	2026-10-14 15:27:40.212762	2026-10-14 15:27:40.365898
215	Controller 16189701	219	\N	t	2026-10-14 15:27:40.366467+05:30	2026-10-14 15:27:40.199703	2026-10-14 15:27:40.366467
180	Other controller B83246528	184	\N	t	2026-10-14 12:07:01.145316+05:30	2026-10-14 12:07:00.852192	2026-10-14 12:07:01.145316
179	Controller B83246528	183	\N	t	2026-10-14 12:07:01.145599+05:30	2026-10-14 12:07:00.851924	2026-10-14 12:07:01.145599
178	Other controller A83246528	182	\N	t	2026-10-14 12:07:01.150815+05:30	2026-10-14 12:07:00.834629	2026-10-14 12:07:01.150815
177	Controller A83246528	181	\N	t	2026-10-14 12:07:01.151138+05:30	2026-10-14 12:07:00.834388	2026-10-14 12:07:01.151138
182	Other controller 847312702	186	\N	t	2026-10-14 12:07:08.562172+05:30	2026-10-14 12:07:08.479533	2026-10-14 12:07:08.562172
181	Controller 847312702	185	\N	t	2026-10-14 12:07:08.562546+05:30	2026-10-14 12:07:08.476646	2026-10-14 12:07:08.562546
220	Other Model Controller 36865502	223	\N	t	2026-10-14 15:27:40.398986+05:30	2026-10-14 15:27:40.373112	2026-10-14 15:27:40.398986
219	Model Controller 36865502	223	\N	t	2026-10-14 15:27:40.399639+05:30	2026-10-14 15:27:40.371661	2026-10-14 15:27:40.399639
186	Other controller B857087506	190	\N	t	2026-10-14 12:07:08.862778+05:30	2026-10-14 12:07:08.586465	2026-10-14 12:07:08.862778
185	Controller B857087506	189	\N	t	2026-10-14 12:07:08.863022+05:30	2026-10-14 12:07:08.586167	2026-10-14 12:07:08.863022
184	Other controller A857087506	188	\N	t	2026-10-14 12:07:08.867521+05:30	2026-10-14 12:07:08.57345	2026-10-14 12:07:08.867521
183	Controller A857087506	187	\N	t	2026-10-14 12:07:08.867844+05:30	2026-10-14 12:07:08.573109	2026-10-14 12:07:08.867844
222	Other controller 40066304	225	\N	t	2026-10-14 15:27:40.469184+05:30	2026-10-14 15:27:40.411267	2026-10-14 15:27:40.469184
221	Controller 40066304	224	\N	t	2026-10-14 15:27:40.469481+05:30	2026-10-14 15:27:40.41087	2026-10-14 15:27:40.469481
187	Test Controller 889635110 Updated	195	\N	t	2026-10-14 12:41:48.933816+05:30	2026-10-14 12:41:48.897101	2026-10-14 12:41:48.933816
189	Test Controller 889635110 Updated	195	\N	t	2026-10-14 12:41:48.937419+05:30	2026-10-14 12:41:48.936217	2026-10-14 12:41:48.937419
191	Test Controller 893835412	197	\N	t	2026-10-14 12:41:48.942734+05:30	2026-10-14 12:41:48.940856	2026-10-14 12:41:48.942734
190	Test Controller 893835412	196	\N	t	2026-10-14 12:41:48.944417+05:30	2026-10-14 12:41:48.939404	2026-10-14 12:41:48.944417
226	Other controller 471232059	229	\N	t	2026-10-14 15:27:40.542741+05:30	2026-10-14 15:27:40.485789	2026-10-14 15:27:40.542741
192	Test Controller 894608716 840D	198	\N	t	2026-10-14 12:41:48.947827+05:30	2026-10-14 12:41:48.946633	2026-10-14 12:41:48.947827
225	Controller 471232059	228	\N	t	2026-10-14 15:27:40.543027+05:30	2026-10-14 15:27:40.485542	2026-10-14 15:27:40.543027
224	Other controller 47123205	227	\N	t	2026-10-14 15:27:40.547215+05:30	2026-10-14 15:27:40.474345	2026-10-14 15:27:40.547215
194	Other Model Controller 894911517	199	\N	t	2026-10-14 12:41:48.978119+05:30	2026-10-14 12:41:48.951146	2026-10-14 12:41:48.978119
193	Model Controller 894911517	199	\N	t	2026-10-14 12:41:48.97865+05:30	2026-10-14 12:41:48.95027	2026-10-14 12:41:48.97865
223	Controller 47123205	226	\N	t	2026-10-14 15:27:40.547645+05:30	2026-10-14 15:27:40.474108	2026-10-14 15:27:40.547645
196	Spare Controller 897943019	200	\N	t	2026-10-14 12:41:48.991863+05:30	2026-10-14 12:41:48.980883	2026-10-14 12:41:48.991863
195	Renamed Controller 897943019	200	\N	t	2026-10-14 12:41:48.99232+05:30	2026-10-14 12:41:48.98043	2026-10-14 12:41:48.99232
198	Other controller 899315822	203	\N	t	2026-10-14 12:41:49.094886+05:30	2026-10-14 12:41:48.996006	2026-10-14 12:41:49.094886
197	Controller 899315822	202	\N	t	2026-10-14 12:41:49.095148+05:30	2026-10-14 12:41:48.995687	2026-10-14 12:41:49.095148
230	Other controller 1932926019	233	\N	t	2026-10-14 15:28:02.065073+05:30	2026-10-14 15:28:01.994946	2026-10-14 15:28:02.065073
229	Controller 1932926019	232	\N	t	2026-10-14 15:28:02.065405+05:30	2026-10-14 15:28:01.99459	2026-10-14 15:28:02.065405
228	Other controller 193292601	231	\N	t	2026-10-14 15:28:02.070002+05:30	2026-10-14 15:28:01.948253	2026-10-14 15:28:02.070002
227	Controller 193292601	230	\N	t	2026-10-14 15:28:02.070292+05:30	2026-10-14 15:28:01.94404	2026-10-14 15:28:02.070292
202	Other controller B911479328	207	\N	t	2026-10-14 12:41:49.447056+05:30	2026-10-14 12:41:49.138936	2026-10-14 12:41:49.447056
201	Controller B911479328	206	\N	t	2026-10-14 12:41:49.44736+05:30	2026-10-14 12:41:49.138591	2026-10-14 12:41:49.44736
200	Other controller A911479328	205	\N	t	2026-10-14 12:41:49.453186+05:30	2026-10-14 12:41:49.117497	2026-10-14 12:41:49.453186
199	Controller A911479328	204	\N	t	2026-10-14 12:41:49.45359+05:30	2026-10-14 12:41:49.117133	2026-10-14 12:41:49.45359
232	Other Model Controller 207266002	234	\N	t	2026-10-14 15:28:02.088879+05:30	2026-10-14 15:28:02.074872	2026-10-14 15:28:02.088879
231	Model Controller 207266002	234	\N	t	2026-10-14 15:28:02.089332+05:30	2026-10-14 15:28:02.073992	2026-10-14 15:28:02.089332
234	Other controller 209017604	236	\N	t	2026-10-14 15:28:02.138951+05:30	2026-10-14 15:28:02.09369	2026-10-14 15:28:02.138951
233	Controller 209017604	235	\N	t	2026-10-14 15:28:02.139367+05:30	2026-10-14 15:28:02.092286	2026-10-14 15:28:02.139367
282	Other controller 5216632019	285	\N	t	2026-10-14 16:26:15.711988+05:30	2026-10-14 16:26:15.581914	2026-10-14 16:26:15.711988
281	Controller 5216632019	284	\N	t	2026-10-14 16:26:15.712529+05:30	2026-10-14 16:26:15.581212	2026-10-14 16:26:15.712529
280	Other controller 521663201	283	\N	t	2026-10-14 16:26:15.71938+05:30	2026-10-14 16:26:15.313366	2026-10-14 16:26:15.71938
238	Other controller 2140857059	240	\N	t	2026-10-14 15:28:02.187872+05:30	2026-10-14 15:28:02.155015	2026-10-14 15:28:02.187872
237	Controller 2140857059	239	\N	t	2026-10-14 15:28:02.188085+05:30	2026-10-14 15:28:02.154694	2026-10-14 15:28:02.188085
236	Other controller 214085705	238	\N	t	2026-10-14 15:28:02.191517+05:30	2026-10-14 15:28:02.143713	2026-10-14 15:28:02.191517
235	Controller 214085705	237	\N	t	2026-10-14 15:28:02.191717+05:30	2026-10-14 15:28:02.143453	2026-10-14 15:28:02.191717
279	Controller 521663201	282	\N	t	2026-10-14 16:26:15.719706+05:30	2026-10-14 16:26:15.291093	2026-10-14 16:26:15.719706
284	Other controller 572189502	287	\N	t	2026-10-14 16:26:15.773306+05:30	2026-10-14 16:26:15.724676	2026-10-14 16:26:15.773306
242	Other controller 1589422019	244	\N	t	2026-10-14 15:41:11.777948+05:30	2026-10-14 15:41:11.68771	2026-10-14 15:41:11.777948
241	Controller 1589422019	243	\N	t	2026-10-14 15:41:11.778935+05:30	2026-10-14 15:41:11.687291	2026-10-14 15:41:11.778935
240	Other controller 158942201	242	\N	t	2026-10-14 15:41:11.783161+05:30	2026-10-14 15:41:11.634122	2026-10-14 15:41:11.783161
239	Controller 158942201	241	\N	t	2026-10-14 15:41:11.7834+05:30	2026-10-14 15:41:11.621165	2026-10-14 15:41:11.7834
283	Controller 572189502	286	\N	t	2026-10-14 16:26:15.773653+05:30	2026-10-14 16:26:15.724322	2026-10-14 16:26:15.773653
246	Other controller 1785631029	248	\N	t	2026-10-14 15:41:11.83778+05:30	2026-10-14 15:41:11.805606	2026-10-14 15:41:11.83778
245	Controller 1785631029	247	\N	t	2026-10-14 15:41:11.838036+05:30	2026-10-14 15:41:11.805267	2026-10-14 15:41:11.838036
244	Other controller 178563102	246	\N	t	2026-10-14 15:41:11.842378+05:30	2026-10-14 15:41:11.788667	2026-10-14 15:41:11.842378
243	Controller 178563102	245	\N	t	2026-10-14 15:41:11.842638+05:30	2026-10-14 15:41:11.788216	2026-10-14 15:41:11.842638
288	Other controller 5775058039	291	\N	t	2026-10-14 16:26:15.81705+05:30	2026-10-14 16:26:15.788737	2026-10-14 16:26:15.81705
287	Controller 5775058039	290	\N	t	2026-10-14 16:26:15.817296+05:30	2026-10-14 16:26:15.788341	2026-10-14 16:26:15.817296
286	Other controller 577505803	289	\N	t	2026-10-14 16:26:15.820317+05:30	2026-10-14 16:26:15.778139	2026-10-14 16:26:15.820317
250	Other controller 5730982019	252	\N	t	2026-10-14 15:51:05.971484+05:30	2026-10-14 15:51:05.835585	2026-10-14 15:51:05.971484
249	Controller 5730982019	251	\N	t	2026-10-14 15:51:05.972443+05:30	2026-10-14 15:51:05.835083	2026-10-14 15:51:05.972443
248	Other controller 573098201	250	\N	t	2026-10-14 15:51:05.995832+05:30	2026-10-14 15:51:05.789338	2026-10-14 15:51:05.995832
247	Controller 573098201	249	\N	t	2026-10-14 15:51:05.996158+05:30	2026-10-14 15:51:05.783829	2026-10-14 15:51:05.996158
285	Controller 577505803	288	\N	t	2026-10-14 16:26:15.820509+05:30	2026-10-14 16:26:15.777832	2026-10-14 16:26:15.820509
254	Other controller 6003831029	256	\N	t	2026-10-14 15:51:06.111505+05:30	2026-10-14 15:51:06.034875	2026-10-14 15:51:06.111505
253	Controller 6003831029	255	\N	t	2026-10-14 15:51:06.112075+05:30	2026-10-14 15:51:06.034214	2026-10-14 15:51:06.112075
252	Other controller 600383102	254	\N	t	2026-10-14 15:51:06.121223+05:30	2026-10-14 15:51:06.008566	2026-10-14 15:51:06.121223
251	Controller 600383102	253	\N	t	2026-10-14 15:51:06.121467+05:30	2026-10-14 15:51:06.006715	2026-10-14 15:51:06.121467
289	Test Controller 759767510 Updated	296	\N	t	2026-10-15 10:20:47.62746+05:30	2026-10-15 10:20:47.598173	2026-10-15 10:20:47.62746
291	Test Controller 759767510 Updated	296	\N	t	2026-10-15 10:20:47.63113+05:30	2026-10-15 10:20:47.630039	2026-10-15 10:20:47.63113
255	Test Controller 835770910 Updated	261	\N	t	2026-10-14 16:09:28.37666+05:30	2026-10-14 16:09:28.358825	2026-10-14 16:09:28.37666
257	Test Controller 835770910 Updated	261	\N	t	2026-10-14 16:09:28.379352+05:30	2026-10-14 16:09:28.378421	2026-10-14 16:09:28.379352
293	Test Controller 763209412	298	\N	t	2026-10-15 10:20:47.635396+05:30	2026-10-15 10:20:47.634029	2026-10-15 10:20:47.635396
259	Test Controller 838011712	263	\N	t	2026-10-14 16:09:28.383183+05:30	2026-10-14 16:09:28.381855	2026-10-14 16:09:28.383183
258	Test Controller 838011712	262	\N	t	2026-10-14 16:09:28.384146+05:30	2026-10-14 16:09:28.380893	2026-10-14 16:09:28.384146
292	Test Controller 763209412	297	\N	t	2026-10-15 10:20:47.636332+05:30	2026-10-15 10:20:47.632901	2026-10-15 10:20:47.636332
260	Test Controller 838520416 840D	264	\N	t	2026-10-14 16:09:28.386434+05:30	2026-10-14 16:09:28.38543	2026-10-14 16:09:28.386434
294	Test Controller 763735916 840D	299	\N	t	2026-10-15 10:20:47.639613+05:30	2026-10-15 10:20:47.637604	2026-10-15 10:20:47.639613
264	Other controller 8387817179	268	\N	t	2026-10-14 16:09:28.481945+05:30	2026-10-14 16:09:28.42461	2026-10-14 16:09:28.481945
263	Controller 8387817179	267	\N	t	2026-10-14 16:09:28.482225+05:30	2026-10-14 16:09:28.424213	2026-10-14 16:09:28.482225
262	Other controller 838781717	266	\N	t	2026-10-14 16:09:28.487827+05:30	2026-10-14 16:09:28.390436	2026-10-14 16:09:28.487827
261	Controller 838781717	265	\N	t	2026-10-14 16:09:28.488112+05:30	2026-10-14 16:09:28.390119	2026-10-14 16:09:28.488112
266	Other Model Controller 849088318	269	\N	t	2026-10-14 16:09:28.502673+05:30	2026-10-14 16:09:28.49248	2026-10-14 16:09:28.502673
265	Model Controller 849088318	269	\N	t	2026-10-14 16:09:28.503093+05:30	2026-10-14 16:09:28.49171	2026-10-14 16:09:28.503093
298	Other controller 7640435179	303	\N	t	2026-10-15 10:20:47.784785+05:30	2026-10-15 10:20:47.687424	2026-10-15 10:20:47.784785
297	Controller 7640435179	302	\N	t	2026-10-15 10:20:47.785207+05:30	2026-10-15 10:20:47.687139	2026-10-15 10:20:47.785207
296	Other controller 764043517	301	\N	t	2026-10-15 10:20:47.789838+05:30	2026-10-15 10:20:47.643336	2026-10-15 10:20:47.789838
268	Spare Controller 850316820	270	\N	t	2026-10-14 16:09:28.51294+05:30	2026-10-14 16:09:28.504864	2026-10-14 16:09:28.51294
267	Renamed Controller 850316820	270	\N	t	2026-10-14 16:09:28.513474+05:30	2026-10-14 16:09:28.504498	2026-10-14 16:09:28.513474
295	Controller 764043517	300	\N	t	2026-10-15 10:20:47.790197+05:30	2026-10-15 10:20:47.643006	2026-10-15 10:20:47.790197
270	Other controller 851444423	273	\N	t	2026-10-14 16:09:28.560124+05:30	2026-10-14 16:09:28.518758	2026-10-14 16:09:28.560124
269	Controller 851444423	272	\N	t	2026-10-14 16:09:28.560442+05:30	2026-10-14 16:09:28.518361	2026-10-14 16:09:28.560442
300	Other Model Controller 779250818	304	\N	t	2026-10-15 10:20:47.805028+05:30	2026-10-15 10:20:47.794336	2026-10-15 10:20:47.805028
299	Model Controller 779250818	304	\N	t	2026-10-15 10:20:47.80545+05:30	2026-10-15 10:20:47.793713	2026-10-15 10:20:47.80545
274	Other controller 8567715269	277	\N	t	2026-10-14 16:09:28.618232+05:30	2026-10-14 16:09:28.585482	2026-10-14 16:09:28.618232
273	Controller 8567715269	276	\N	t	2026-10-14 16:09:28.6185+05:30	2026-10-14 16:09:28.585119	2026-10-14 16:09:28.6185
272	Other controller 856771526	275	\N	t	2026-10-14 16:09:28.621738+05:30	2026-10-14 16:09:28.571321	2026-10-14 16:09:28.621738
271	Controller 856771526	274	\N	t	2026-10-14 16:09:28.621925+05:30	2026-10-14 16:09:28.570986	2026-10-14 16:09:28.621925
302	Spare Controller 780558820	305	\N	t	2026-10-15 10:20:47.813958+05:30	2026-10-15 10:20:47.807311	2026-10-15 10:20:47.813958
301	Renamed Controller 780558820	305	\N	t	2026-10-15 10:20:47.814367+05:30	2026-10-15 10:20:47.806913	2026-10-15 10:20:47.814367
278	Other controller B868524232	281	\N	t	2026-10-14 16:09:29.165198+05:30	2026-10-14 16:09:28.697682	2026-10-14 16:09:29.165198
277	Controller B868524232	280	\N	t	2026-10-14 16:09:29.165492+05:30	2026-10-14 16:09:28.697401	2026-10-14 16:09:29.165492
276	Other controller A868524232	279	\N	t	2026-10-14 16:09:29.170328+05:30	2026-10-14 16:09:28.68749	2026-10-14 16:09:29.170328
275	Controller A868524232	278	\N	t	2026-10-14 16:09:29.170651+05:30	2026-10-14 16:09:28.687227	2026-10-14 16:09:29.170651
304	Other controller 781488023	308	\N	t	2026-10-15 10:20:47.860145+05:30	2026-10-15 10:20:47.817457	2026-10-15 10:20:47.860145
303	Controller 781488023	307	\N	t	2026-10-15 10:20:47.860506+05:30	2026-10-15 10:20:47.817171	2026-10-15 10:20:47.860506
308	Other controller 7870395269	312	\N	t	2026-10-15 10:20:47.926478+05:30	2026-10-15 10:20:47.883719	2026-10-15 10:20:47.926478
307	Controller 7870395269	311	\N	t	2026-10-15 10:20:47.926806+05:30	2026-10-15 10:20:47.883405	2026-10-15 10:20:47.926806
306	Other controller 787039526	310	\N	t	2026-10-15 10:20:47.930678+05:30	2026-10-15 10:20:47.873067	2026-10-15 10:20:47.930678
305	Controller 787039526	309	\N	t	2026-10-15 10:20:47.93092+05:30	2026-10-15 10:20:47.872744	2026-10-15 10:20:47.93092
312	Other controller B801317832	316	\N	t	2026-10-15 10:20:48.485584+05:30	2026-10-15 10:20:48.182868	2026-10-15 10:20:48.485584
311	Controller B801317832	315	\N	t	2026-10-15 10:20:48.485828+05:30	2026-10-15 10:20:48.182455	2026-10-15 10:20:48.485828
310	Other controller A801317832	314	\N	t	2026-10-15 10:20:48.489908+05:30	2026-10-15 10:20:48.015045	2026-10-15 10:20:48.489908
309	Controller A801317832	313	\N	t	2026-10-15 10:20:48.490114+05:30	2026-10-15 10:20:48.014803	2026-10-15 10:20:48.490114
\.


--
-- TOC entry 5642 (class 0 OID 116513)
-- Dependencies: 222
-- Data for Name: countries; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.countries (id, name, code, phone_code, created_at) FROM stdin;
1	India	IND	+91	2026-10-12 10:38:40.418244
2	United States	USA	+1	2026-10-12 10:38:40.418244
\.


--
-- TOC entry 5712 (class 0 OID 118011)
-- Dependencies: 292
-- Data for Name: holidays; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.holidays (id, user_id, holiday_type, name, start_date, end_date, week_day, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
5	99	week_off	\N	\N	\N	3	active	f	\N	2026-10-14 16:09:28.674553	2026-10-14 16:09:28.674553
2	98	holiday	Holi	2026-03-04	2026-03-04	\N	inactive	f	\N	2026-10-14 16:09:28.671231	2026-10-14 16:09:28.676496
1	98	holiday	Diwali	2026-11-08	2026-11-10	\N	active	t	2026-10-14 16:09:28.678279+05:30	2026-10-14 16:09:28.66592	2026-10-14 16:09:28.678279
10	114	week_off	\N	\N	\N	3	active	f	\N	2026-10-14 16:26:15.877696	2026-10-14 16:26:15.877696
7	113	holiday	Holi	2026-03-04	2026-03-04	\N	inactive	f	\N	2026-10-14 16:26:15.874627	2026-10-14 16:26:15.879189
6	113	holiday	Diwali	2026-11-08	2026-11-10	\N	active	t	2026-10-14 16:26:15.880799+05:30	2026-10-14 16:26:15.870102	2026-10-14 16:26:15.880799
11	22	week_off	GENERAL WEEK OFF	\N	\N	5	active	f	\N	2026-10-14 16:30:34.76002	2026-10-14 16:30:34.76002
16	135	week_off	\N	\N	\N	3	active	f	\N	2026-10-15 10:20:48.003025	2026-10-15 10:20:48.003025
13	134	holiday	Holi	2026-03-04	2026-03-04	\N	inactive	f	\N	2026-10-15 10:20:48.000327	2026-10-15 10:20:48.004807
12	134	holiday	Diwali	2026-11-08	2026-11-10	\N	active	t	2026-10-15 10:20:48.006583+05:30	2026-10-15 10:20:47.987766	2026-10-15 10:20:48.006583
\.


--
-- TOC entry 5672 (class 0 OID 116968)
-- Dependencies: 252
-- Data for Name: machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.machines (id, user_id, machine_type_id, product_id, name, serial_number, live_ip, port, username, password, local_ip, brand_id, controller_id, controller_version, accuracy_id, line_id, installation_date, bill_date, warranty_date, amc, is_insured, online_status, maintenance_required, hourly_cost, status, is_monitor, online_date, offline_date, mode_id, controller_version_type, health_ids, longitude, latitude, cycle_analysis, machine_system_type, security_mode, is_deleted, deleted_at, created_at, updated_at, model_id) FROM stdin;
37	64	\N	\N	VMC-769633603-B	SN-VMC-769633603-B	103.81.119.106	4840	operator	secret123	192.168.1.10	213	209	V04.07	1	16	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:55:27.770631+05:30	2026-10-14 12:55:27.736653	2026-10-14 12:55:27.770631	82
35	62	\N	\N	VMC-B-327949001	SN-VMC-B-327949001	103.81.119.106	4840	operator	secret123	192.168.1.10	208	203	V04.07	1	14	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:27.241145+05:30	2026-10-14 12:55:13.484007	2026-10-14 15:17:27.241145	75
36	63	\N	\N	VMC-O-327949001	SN-VMC-O-327949001	103.81.119.106	4840	operator	secret123	192.168.1.10	210	205	V04.07	1	15	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:28.865916+05:30	2026-10-14 12:55:13.4999	2026-10-14 15:17:28.865916	77
34	62	\N	\N	VMC-A-327949001	SN-VMC-A-327949001	103.81.119.106	4840	operator	secret123	192.168.1.10	208	203	V04.07	1	14	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:30.569678+05:30	2026-10-14 12:55:13.44253	2026-10-14 15:17:30.569678	75
1	\N	\N	\N	In use 569781218	In use 569781218	10.0.0.1	4840	user	secret	10.0.0.2	87	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-13 18:01:25.715056+05:30	2026-10-13 18:01:25.702766	2026-10-14 09:44:42.178231	\N
2	\N	\N	\N	In use 973238118	In use 973238118	10.0.0.1	4840	user	secret	10.0.0.2	98	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-13 18:01:39.738791+05:30	2026-10-13 18:01:39.734563	2026-10-14 09:44:42.178231	\N
3	\N	\N	\N	In use 722564619	In use 722564619	10.0.0.1	4840	user	secret	10.0.0.2	109	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-13 18:01:47.233409+05:30	2026-10-13 18:01:47.228952	2026-10-14 09:44:42.178231	\N
4	\N	\N	\N	In use 91908212	In use 91908212	10.0.0.1	4840	user	secret	10.0.0.2	120	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 09:37:20.931863+05:30	2026-10-14 09:37:20.923628	2026-10-14 09:44:42.178231	\N
5	13	\N	\N	VMC-93421215-B	SN-VMC-93421215-B	103.81.119.106	4840	operator	secret123	192.168.1.10	122	\N	V04.07	1	1	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 09:37:21.05651+05:30	2026-10-14 09:37:21.021758	2026-10-14 09:44:42.178231	\N
10	\N	\N	\N	In use 630186519	In use 630186519	10.0.0.1	4840	user	secret	10.0.0.2	151	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 09:38:36.311179+05:30	2026-10-14 09:38:36.305329	2026-10-14 09:44:42.178231	\N
11	20	\N	\N	VMC-631445522-B	SN-VMC-631445522-B	103.81.119.106	4840	operator	secret123	192.168.1.10	153	\N	V04.07	1	3	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 09:38:36.3781+05:30	2026-10-14 09:38:36.347903	2026-10-14 09:44:42.178231	\N
7	\N	\N	\N	In use 155984519	In use 155984519	10.0.0.1	4840	user	secret	10.0.0.2	138	\N	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 09:37:52.133363+05:30	2026-10-14 09:37:51.973139	2026-10-14 09:44:42.178231	\N
8	16	\N	\N	VMC-229469522-B	SN-VMC-229469522-B	103.81.119.106	4840	operator	secret123	192.168.1.10	140	\N	V04.07	1	2	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 09:37:53.358382+05:30	2026-10-14 09:37:52.996745	2026-10-14 09:44:42.178231	\N
14	\N	\N	\N	In use 12098819	In use 12098819	10.0.0.1	4840	user	secret	10.0.0.2	164	161	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 11:57:30.130895+05:30	2026-10-14 11:57:30.124428	2026-10-14 11:57:30.130895	46
15	26	\N	\N	VMC-13336822-B	SN-VMC-13336822-B	103.81.119.106	4840	operator	secret123	192.168.1.10	166	163	V04.07	1	4	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 11:57:30.20959+05:30	2026-10-14 11:57:30.174185	2026-10-14 11:57:30.20959	47
13	23	\N	\N	Goodwin	17	238.32.88.28	9878	ENV	12345678	166.134.176.94	63	151	v.123	3	\N	2025-10-27	2026-11-24	2027-01-09	Quidem delectatio amiculum arca defaeco vito arcus	t	f	t	2131.00	active	f	\N	\N	\N	2	\N	58	56	f	OldPlc	Basic256Sha256	f	\N	2026-10-14 10:46:47.455094	2026-10-14 12:06:54.193703	41
17	\N	\N	\N	In use 71813419	In use 71813419	10.0.0.1	4840	user	secret	10.0.0.2	177	173	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:00.727797+05:30	2026-10-14 12:07:00.721662	2026-10-14 12:07:00.727797	52
18	31	\N	\N	VMC-73394922-B	SN-VMC-73394922-B	103.81.119.106	4840	operator	secret123	192.168.1.10	179	175	V04.07	1	5	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:00.810584+05:30	2026-10-14 12:07:00.777598	2026-10-14 12:07:00.810584	53
22	34	\N	\N	UA-B83246528-0	SN-UA-B83246528-0	103.81.119.106	4840	operator	secret123	192.168.1.10	183	179	V04.07	1	7	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:01.142507+05:30	2026-10-14 12:07:00.858351	2026-10-14 12:07:01.142507	57
21	33	\N	\N	UA-A83246528-1	SN-UA-A83246528-1	103.81.119.106	4840	operator	secret123	192.168.1.10	181	177	V04.07	1	6	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:01.146697+05:30	2026-10-14 12:07:00.849116	2026-10-14 12:07:01.146697	55
20	33	\N	\N	UA-A83246528-0	SN-UA-A83246528-0	103.81.119.106	4840	operator	secret123	192.168.1.10	181	177	V04.07	1	6	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:01.147032+05:30	2026-10-14 12:07:00.839221	2026-10-14 12:07:01.147032	55
23	42	\N	\N	VMC-847312702-B	SN-VMC-847312702-B	103.81.119.106	4840	operator	secret123	192.168.1.10	185	181	V04.07	1	8	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:08.554426+05:30	2026-10-14 12:07:08.513104	2026-10-14 12:07:08.554426	59
27	44	\N	\N	UA-B857087506-0	SN-UA-B857087506-0	103.81.119.106	4840	operator	secret123	192.168.1.10	189	185	V04.07	1	10	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:08.860223+05:30	2026-10-14 12:07:08.590604	2026-10-14 12:07:08.860223	63
26	43	\N	\N	UA-A857087506-1	SN-UA-A857087506-1	103.81.119.106	4840	operator	secret123	192.168.1.10	187	183	V04.07	1	9	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:08.864012+05:30	2026-10-14 12:07:08.583373	2026-10-14 12:07:08.864012	61
25	43	\N	\N	UA-A857087506-0	SN-UA-A857087506-0	103.81.119.106	4840	operator	secret123	192.168.1.10	187	183	V04.07	1	9	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:07:08.864348+05:30	2026-10-14 12:07:08.57789	2026-10-14 12:07:08.864348	61
28	\N	\N	\N	In use 897943019	In use 897943019	10.0.0.1	4840	user	secret	10.0.0.2	200	195	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 12:41:48.990714+05:30	2026-10-14 12:41:48.98236	2026-10-14 12:41:48.990714	68
29	54	\N	\N	VMC-899315822-B	SN-VMC-899315822-B	103.81.119.106	4840	operator	secret123	192.168.1.10	202	197	V04.07	1	11	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:41:49.085011+05:30	2026-10-14 12:41:49.047621	2026-10-14 12:41:49.085011	69
33	57	\N	\N	UA-B911479328-0	SN-UA-B911479328-0	103.81.119.106	4840	operator	secret123	192.168.1.10	206	201	V04.07	1	13	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:41:49.443688+05:30	2026-10-14 12:41:49.144854	2026-10-14 12:41:49.443688	73
32	56	\N	\N	UA-A911479328-1	SN-UA-A911479328-1	103.81.119.106	4840	operator	secret123	192.168.1.10	204	199	V04.07	1	12	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:41:49.448562+05:30	2026-10-14 12:41:49.134317	2026-10-14 12:41:49.448562	71
31	56	\N	\N	UA-A911479328-0	SN-UA-A911479328-0	103.81.119.106	4840	operator	secret123	192.168.1.10	204	199	V04.07	1	12	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 12:41:49.448981+05:30	2026-10-14 12:41:49.126232	2026-10-14 12:41:49.448981	71
41	66	\N	\N	VMC-O-777978904	SN-VMC-O-777978904	103.81.119.106	4840	operator	secret123	192.168.1.10	217	213	V04.07	1	18	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:15.00191+05:30	2026-10-14 12:55:27.800654	2026-10-14 15:17:15.00191	86
40	65	\N	\N	VMC-B-777978904	SN-VMC-B-777978904	103.81.119.106	4840	operator	secret123	192.168.1.10	215	211	V04.07	1	17	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:18.248941+05:30	2026-10-14 12:55:27.793032	2026-10-14 15:17:18.248941	84
39	65	\N	\N	VMC-A-777978904	SN-VMC-A-777978904	103.81.119.106	4840	operator	secret123	192.168.1.10	215	211	V04.07	1	17	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:17:19.936905+05:30	2026-10-14 12:55:27.786436	2026-10-14 15:17:19.936905	84
42	22	\N	\N	Huels	19	154.129.172.33	6157	cust_98	12345678	103.67.208.223	63	151	v23.23	2	\N	2026-09-30	\N	\N	Aegrotatio canto sufficio quidem supplanto custodi	t	f	t	878.00	active	f	\N	\N	\N	1	\N	46	77	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:18:14.56921	2026-10-14 15:18:25.937115	41
43	67	\N	\N	VMC-C-16189701	SN-VMC-C-16189701	103.81.119.106	4840	operator	secret123	192.168.1.10	219	215	V04.07	1	19	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:27:40.318031	2026-10-14 15:27:40.318031	88
44	69	\N	\N	VMC-40066304-B	SN-VMC-40066304-B	103.81.119.106	4840	operator	secret123	192.168.1.10	224	221	V04.07	1	21	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:27:40.458798+05:30	2026-10-14 15:27:40.417789	2026-10-14 15:27:40.458798	95
46	70	\N	\N	VMC-A-47123205	SN-VMC-A-47123205	103.81.119.106	4840	operator	secret123	192.168.1.10	226	223	V04.07	1	22	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:27:40.479365	2026-10-14 15:27:40.479365	97
47	70	\N	\N	VMC-B-47123205	SN-VMC-B-47123205	103.81.119.106	4840	operator	secret123	192.168.1.10	226	223	V04.07	1	22	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:27:40.482251	2026-10-14 15:27:40.482251	97
48	71	\N	\N	VMC-O-47123205	SN-VMC-O-47123205	103.81.119.106	4840	operator	secret123	192.168.1.10	228	225	V04.07	1	23	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:27:40.491071	2026-10-14 15:27:40.491071	99
49	72	\N	\N	VMC-C-193292601	SN-VMC-C-193292601	103.81.119.106	4840	operator	secret123	192.168.1.10	230	227	V04.07	1	24	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:28:02.028265	2026-10-14 15:28:02.028265	101
50	74	\N	\N	VMC-209017604-B	SN-VMC-209017604-B	103.81.119.106	4840	operator	secret123	192.168.1.10	235	233	V04.07	1	26	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 15:28:02.129369+05:30	2026-10-14 15:28:02.098914	2026-10-14 15:28:02.129369	108
52	75	\N	\N	VMC-A-214085705	SN-VMC-A-214085705	103.81.119.106	4840	operator	secret123	192.168.1.10	237	235	V04.07	1	27	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:28:02.148729	2026-10-14 15:28:02.148729	110
53	75	\N	\N	VMC-B-214085705	SN-VMC-B-214085705	103.81.119.106	4840	operator	secret123	192.168.1.10	237	235	V04.07	1	27	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:28:02.150963	2026-10-14 15:28:02.150963	110
54	76	\N	\N	VMC-O-214085705	SN-VMC-O-214085705	103.81.119.106	4840	operator	secret123	192.168.1.10	239	237	V04.07	1	28	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:28:02.159478	2026-10-14 15:28:02.159478	112
55	77	\N	\N	VMC-C-158942201	SN-VMC-C-158942201	103.81.119.106	4840	operator	secret123	192.168.1.10	241	239	V04.07	1	29	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:41:11.732482	2026-10-14 15:41:11.732482	114
56	79	\N	\N	VMC-A-178563102	SN-VMC-A-178563102	103.81.119.106	4840	operator	secret123	192.168.1.10	245	243	V04.07	1	31	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:41:11.794921	2026-10-14 15:41:11.794921	118
57	79	\N	\N	VMC-B-178563102	SN-VMC-B-178563102	103.81.119.106	4840	operator	secret123	192.168.1.10	245	243	V04.07	1	31	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:41:11.799328	2026-10-14 15:41:11.799328	118
58	80	\N	\N	VMC-O-178563102	SN-VMC-O-178563102	103.81.119.106	4840	operator	secret123	192.168.1.10	247	245	V04.07	1	32	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:41:11.811194	2026-10-14 15:41:11.811194	120
59	81	\N	\N	VMC-C-573098201	SN-VMC-C-573098201	103.81.119.106	4840	operator	secret123	192.168.1.10	249	247	V04.07	1	33	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:51:05.905405	2026-10-14 15:51:05.905405	122
60	83	\N	\N	VMC-A-600383102	SN-VMC-A-600383102	103.81.119.106	4840	operator	secret123	192.168.1.10	253	251	V04.07	1	35	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:51:06.021294	2026-10-14 15:51:06.021294	126
61	83	\N	\N	VMC-B-600383102	SN-VMC-B-600383102	103.81.119.106	4840	operator	secret123	192.168.1.10	253	251	V04.07	1	35	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:51:06.027833	2026-10-14 15:51:06.027833	126
62	84	\N	\N	VMC-O-600383102	SN-VMC-O-600383102	103.81.119.106	4840	operator	secret123	192.168.1.10	255	253	V04.07	1	36	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 15:51:06.043135	2026-10-14 15:51:06.043135	128
63	90	\N	\N	VMC-C-838781717	SN-VMC-C-838781717	103.81.119.106	4840	operator	secret123	192.168.1.10	265	261	V04.07	1	37	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:09:28.453099	2026-10-14 16:09:28.453099	130
64	\N	\N	\N	In use 850316820	In use 850316820	10.0.0.1	4840	user	secret	10.0.0.2	270	267	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-14 16:09:28.511974+05:30	2026-10-14 16:09:28.506269	2026-10-14 16:09:28.511974	137
65	92	\N	\N	VMC-851444423-B	SN-VMC-851444423-B	103.81.119.106	4840	operator	secret123	192.168.1.10	272	269	V04.07	1	39	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 16:09:28.552072+05:30	2026-10-14 16:09:28.523828	2026-10-14 16:09:28.552072	138
67	94	\N	\N	VMC-A-856771526	SN-VMC-A-856771526	103.81.119.106	4840	operator	secret123	192.168.1.10	274	271	V04.07	1	40	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:09:28.5793	2026-10-14 16:09:28.5793	140
68	94	\N	\N	VMC-B-856771526	SN-VMC-B-856771526	103.81.119.106	4840	operator	secret123	192.168.1.10	274	271	V04.07	1	40	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:09:28.581584	2026-10-14 16:09:28.581584	140
69	95	\N	\N	VMC-O-856771526	SN-VMC-O-856771526	103.81.119.106	4840	operator	secret123	192.168.1.10	276	273	V04.07	1	41	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:09:28.589839	2026-10-14 16:09:28.589839	142
70	100	\N	\N	UA-A868524232-0	SN-UA-A868524232-0	103.81.119.106	4840	operator	secret123	192.168.1.10	278	275	V04.07	1	42	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 16:09:29.167056+05:30	2026-10-14 16:09:28.691661	2026-10-14 16:09:29.167056	144
72	101	\N	\N	UA-B868524232-0	SN-UA-B868524232-0	103.81.119.106	4840	operator	secret123	192.168.1.10	280	277	V04.07	1	43	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 16:09:29.160971+05:30	2026-10-14 16:09:28.870511	2026-10-14 16:09:29.160971	146
71	100	\N	\N	UA-A868524232-1	SN-UA-A868524232-1	103.81.119.106	4840	operator	secret123	192.168.1.10	278	275	V04.07	1	42	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 16:09:29.166708+05:30	2026-10-14 16:09:28.694816	2026-10-14 16:09:29.166708	144
73	106	\N	\N	VMC-C-521663201	SN-VMC-C-521663201	103.81.119.106	4840	operator	secret123	192.168.1.10	282	279	V04.07	1	44	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:26:15.645119	2026-10-14 16:26:15.645119	148
74	108	\N	\N	VMC-572189502-B	SN-VMC-572189502-B	103.81.119.106	4840	operator	secret123	192.168.1.10	286	283	V04.07	1	46	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-14 16:26:15.764633+05:30	2026-10-14 16:26:15.730811	2026-10-14 16:26:15.764633	152
76	109	\N	\N	VMC-A-577505803	SN-VMC-A-577505803	103.81.119.106	4840	operator	secret123	192.168.1.10	288	285	V04.07	1	47	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:26:15.7835	2026-10-14 16:26:15.7835	154
77	109	\N	\N	VMC-B-577505803	SN-VMC-B-577505803	103.81.119.106	4840	operator	secret123	192.168.1.10	288	285	V04.07	1	47	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:26:15.785386	2026-10-14 16:26:15.785386	154
78	110	\N	\N	VMC-O-577505803	SN-VMC-O-577505803	103.81.119.106	4840	operator	secret123	192.168.1.10	290	287	V04.07	1	48	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-14 16:26:15.793924	2026-10-14 16:26:15.793924	156
79	126	\N	\N	VMC-C-764043517	SN-VMC-C-764043517	103.81.119.106	4840	operator	secret123	192.168.1.10	300	295	V04.07	1	49	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-15 10:20:47.730218	2026-10-15 10:20:47.730218	158
80	\N	\N	\N	In use 780558820	In use 780558820	10.0.0.1	4840	user	secret	10.0.0.2	305	301	V1	\N	\N	2026-01-01	\N	\N	\N	f	f	f	0.00	active	f	\N	\N	\N	\N	\N	\N	\N	f	NewPlc	Basic256Sha256	t	2026-10-15 10:20:47.813153+05:30	2026-10-15 10:20:47.808421	2026-10-15 10:20:47.813153	165
81	128	\N	\N	VMC-781488023-B	SN-VMC-781488023-B	103.81.119.106	4840	operator	secret123	192.168.1.10	307	303	V04.07	1	51	2026-01-15	2026-01-10	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-15 10:20:47.851458+05:30	2026-10-15 10:20:47.821525	2026-10-15 10:20:47.851458	166
83	130	\N	\N	VMC-A-787039526	SN-VMC-A-787039526	103.81.119.106	4840	operator	secret123	192.168.1.10	309	305	V04.07	1	52	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-15 10:20:47.877961	2026-10-15 10:20:47.877961	168
84	130	\N	\N	VMC-B-787039526	SN-VMC-B-787039526	103.81.119.106	4840	operator	secret123	192.168.1.10	309	305	V04.07	1	52	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-15 10:20:47.880041	2026-10-15 10:20:47.880041	168
85	131	\N	\N	VMC-O-787039526	SN-VMC-O-787039526	103.81.119.106	4840	operator	secret123	192.168.1.10	311	307	V04.07	1	53	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	f	\N	2026-10-15 10:20:47.88848	2026-10-15 10:20:47.88848	170
88	137	\N	\N	UA-B801317832-0	SN-UA-B801317832-0	103.81.119.106	4840	operator	secret123	192.168.1.10	315	311	V04.07	1	55	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-15 10:20:48.482246+05:30	2026-10-15 10:20:48.187212	2026-10-15 10:20:48.482246	174
87	136	\N	\N	UA-A801317832-1	SN-UA-A801317832-1	103.81.119.106	4840	operator	secret123	192.168.1.10	313	309	V04.07	1	54	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-15 10:20:48.486861+05:30	2026-10-15 10:20:48.0229	2026-10-15 10:20:48.486861	172
86	136	\N	\N	UA-A801317832-0	SN-UA-A801317832-0	103.81.119.106	4840	operator	secret123	192.168.1.10	313	309	V04.07	1	54	2026-01-15	\N	2027-01-15	\N	t	f	f	800.50	active	f	\N	\N	\N	1	1,5,6	70.671297	22.244054	f	NewPlc	Basic256Sha256	t	2026-10-15 10:20:48.487247+05:30	2026-10-15 10:20:48.021478	2026-10-15 10:20:48.487247	172
\.


--
-- TOC entry 5694 (class 0 OID 117647)
-- Dependencies: 274
-- Data for Name: models; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.models (id, controller_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
40	152	09809	f	\N	2026-10-14 10:37:01.390476+05:30	2026-10-14 10:37:01.390476+05:30
41	151	09789897	f	\N	2026-10-14 10:37:06.427052+05:30	2026-10-14 10:37:06.427052+05:30
42	151	Medhurst	f	\N	2026-10-14 10:37:12.794114+05:30	2026-10-14 10:37:19.440001+05:30
39	151	hello	t	2026-10-14 10:37:22.813812+05:30	2026-10-14 10:36:56.0576+05:30	2026-10-14 10:37:22.813812+05:30
104	230	Other model 1932926019	t	2026-10-14 15:28:02.064408+05:30	2026-10-14 15:28:01.995543+05:30	2026-10-14 15:28:02.064408+05:30
45	160	Model 9539817	t	2026-10-14 11:57:30.111515+05:30	2026-10-14 11:57:30.110378+05:30	2026-10-14 11:57:30.111515+05:30
103	229	Model 1932926019	t	2026-10-14 15:28:02.064735+05:30	2026-10-14 15:28:01.995246+05:30	2026-10-14 15:28:02.064735+05:30
43	159	Model 9539817 Updated	t	2026-10-14 11:57:30.118525+05:30	2026-10-14 11:57:30.098086+05:30	2026-10-14 11:57:30.118525+05:30
102	228	Other model 193292601	t	2026-10-14 15:28:02.069324+05:30	2026-10-14 15:28:01.960392+05:30	2026-10-14 15:28:02.069324+05:30
101	227	Model 193292601	t	2026-10-14 15:28:02.069593+05:30	2026-10-14 15:28:01.948883+05:30	2026-10-14 15:28:02.069593+05:30
46	161	Renamed Model 12098819	t	2026-10-14 11:57:30.131468+05:30	2026-10-14 11:57:30.12288+05:30	2026-10-14 11:57:30.131468+05:30
48	164	Other model 13336822	t	2026-10-14 11:57:30.216694+05:30	2026-10-14 11:57:30.136822+05:30	2026-10-14 11:57:30.216694+05:30
47	163	Model 13336822	t	2026-10-14 11:57:30.216941+05:30	2026-10-14 11:57:30.136482+05:30	2026-10-14 11:57:30.216941+05:30
107	232	Model 207266002	t	2026-10-14 15:28:02.080179+05:30	2026-10-14 15:28:02.079144+05:30	2026-10-14 15:28:02.080179+05:30
51	172	Model 68041017	t	2026-10-14 12:07:00.702721+05:30	2026-10-14 12:07:00.700031+05:30	2026-10-14 12:07:00.702721+05:30
105	231	Model 207266002 Updated	t	2026-10-14 15:28:02.087151+05:30	2026-10-14 15:28:02.076406+05:30	2026-10-14 15:28:02.087151+05:30
49	171	Model 68041017 Updated	t	2026-10-14 12:07:00.713486+05:30	2026-10-14 12:07:00.688861+05:30	2026-10-14 12:07:00.713486+05:30
52	173	Renamed Model 71813419	t	2026-10-14 12:07:00.729059+05:30	2026-10-14 12:07:00.721002+05:30	2026-10-14 12:07:00.729059+05:30
109	234	Other model 209017604	t	2026-10-14 15:28:02.138375+05:30	2026-10-14 15:28:02.094325+05:30	2026-10-14 15:28:02.138375+05:30
108	233	Model 209017604	t	2026-10-14 15:28:02.138657+05:30	2026-10-14 15:28:02.093975+05:30	2026-10-14 15:28:02.138657+05:30
54	176	Other model 73394922	t	2026-10-14 12:07:00.817415+05:30	2026-10-14 12:07:00.740237+05:30	2026-10-14 12:07:00.817415+05:30
53	175	Model 73394922	t	2026-10-14 12:07:00.817711+05:30	2026-10-14 12:07:00.739858+05:30	2026-10-14 12:07:00.817711+05:30
58	180	Other model B83246528	t	2026-10-14 12:07:01.144777+05:30	2026-10-14 12:07:00.852724+05:30	2026-10-14 12:07:01.144777+05:30
57	179	Model B83246528	t	2026-10-14 12:07:01.145056+05:30	2026-10-14 12:07:00.852456+05:30	2026-10-14 12:07:01.145056+05:30
56	178	Other model A83246528	t	2026-10-14 12:07:01.15011+05:30	2026-10-14 12:07:00.835177+05:30	2026-10-14 12:07:01.15011+05:30
55	177	Model A83246528	t	2026-10-14 12:07:01.150557+05:30	2026-10-14 12:07:00.834863+05:30	2026-10-14 12:07:01.150557+05:30
113	238	Other model 2140857059	t	2026-10-14 15:28:02.187309+05:30	2026-10-14 15:28:02.155513+05:30	2026-10-14 15:28:02.187309+05:30
112	237	Model 2140857059	t	2026-10-14 15:28:02.187589+05:30	2026-10-14 15:28:02.155261+05:30	2026-10-14 15:28:02.187589+05:30
60	182	Other model 847312702	t	2026-10-14 12:07:08.56133+05:30	2026-10-14 12:07:08.483073+05:30	2026-10-14 12:07:08.56133+05:30
59	181	Model 847312702	t	2026-10-14 12:07:08.561756+05:30	2026-10-14 12:07:08.480122+05:30	2026-10-14 12:07:08.561756+05:30
111	236	Other model 214085705	t	2026-10-14 15:28:02.191075+05:30	2026-10-14 15:28:02.144238+05:30	2026-10-14 15:28:02.191075+05:30
110	235	Model 214085705	t	2026-10-14 15:28:02.191315+05:30	2026-10-14 15:28:02.143961+05:30	2026-10-14 15:28:02.191315+05:30
64	186	Other model B857087506	t	2026-10-14 12:07:08.862301+05:30	2026-10-14 12:07:08.587016+05:30	2026-10-14 12:07:08.862301+05:30
63	185	Model B857087506	t	2026-10-14 12:07:08.862541+05:30	2026-10-14 12:07:08.586737+05:30	2026-10-14 12:07:08.862541+05:30
62	184	Other model A857087506	t	2026-10-14 12:07:08.866875+05:30	2026-10-14 12:07:08.574051+05:30	2026-10-14 12:07:08.866875+05:30
61	183	Model A857087506	t	2026-10-14 12:07:08.867194+05:30	2026-10-14 12:07:08.57375+05:30	2026-10-14 12:07:08.867194+05:30
117	242	Other model 1589422019	t	2026-10-14 15:41:11.777159+05:30	2026-10-14 15:41:11.688513+05:30	2026-10-14 15:41:11.777159+05:30
67	194	Model 894911517	t	2026-10-14 12:41:48.965426+05:30	2026-10-14 12:41:48.964167+05:30	2026-10-14 12:41:48.965426+05:30
116	241	Model 1589422019	t	2026-10-14 15:41:11.777555+05:30	2026-10-14 15:41:11.688136+05:30	2026-10-14 15:41:11.777555+05:30
65	193	Model 894911517 Updated	t	2026-10-14 12:41:48.976546+05:30	2026-10-14 12:41:48.952682+05:30	2026-10-14 12:41:48.976546+05:30
115	240	Other model 158942201	t	2026-10-14 15:41:11.78265+05:30	2026-10-14 15:41:11.646741+05:30	2026-10-14 15:41:11.78265+05:30
114	239	Model 158942201	t	2026-10-14 15:41:11.782895+05:30	2026-10-14 15:41:11.634811+05:30	2026-10-14 15:41:11.782895+05:30
68	195	Renamed Model 897943019	t	2026-10-14 12:41:48.991389+05:30	2026-10-14 12:41:48.981337+05:30	2026-10-14 12:41:48.991389+05:30
70	198	Other model 899315822	t	2026-10-14 12:41:49.094284+05:30	2026-10-14 12:41:48.996622+05:30	2026-10-14 12:41:49.094284+05:30
69	197	Model 899315822	t	2026-10-14 12:41:49.094547+05:30	2026-10-14 12:41:48.996295+05:30	2026-10-14 12:41:49.094547+05:30
121	246	Other model 1785631029	t	2026-10-14 15:41:11.837271+05:30	2026-10-14 15:41:11.80619+05:30	2026-10-14 15:41:11.837271+05:30
120	245	Model 1785631029	t	2026-10-14 15:41:11.837534+05:30	2026-10-14 15:41:11.805887+05:30	2026-10-14 15:41:11.837534+05:30
74	202	Other model B911479328	t	2026-10-14 12:41:49.446454+05:30	2026-10-14 12:41:49.139788+05:30	2026-10-14 12:41:49.446454+05:30
73	201	Model B911479328	t	2026-10-14 12:41:49.446785+05:30	2026-10-14 12:41:49.13928+05:30	2026-10-14 12:41:49.446785+05:30
72	200	Other model A911479328	t	2026-10-14 12:41:49.452541+05:30	2026-10-14 12:41:49.118158+05:30	2026-10-14 12:41:49.452541+05:30
71	199	Model A911479328	t	2026-10-14 12:41:49.452875+05:30	2026-10-14 12:41:49.117794+05:30	2026-10-14 12:41:49.452875+05:30
119	244	Other model 178563102	t	2026-10-14 15:41:11.841847+05:30	2026-10-14 15:41:11.789535+05:30	2026-10-14 15:41:11.841847+05:30
118	243	Model 178563102	t	2026-10-14 15:41:11.842116+05:30	2026-10-14 15:41:11.789147+05:30	2026-10-14 15:41:11.842116+05:30
78	206	Other model 3279490019	t	2026-10-14 12:55:13.595148+05:30	2026-10-14 12:55:13.492171+05:30	2026-10-14 12:55:13.595148+05:30
77	205	Model 3279490019	t	2026-10-14 12:55:13.595546+05:30	2026-10-14 12:55:13.491836+05:30	2026-10-14 12:55:13.595546+05:30
76	204	Other model 327949001	t	2026-10-14 12:55:13.601735+05:30	2026-10-14 12:55:13.359355+05:30	2026-10-14 12:55:13.601735+05:30
75	203	Model 327949001	t	2026-10-14 12:55:13.602029+05:30	2026-10-14 12:55:13.340528+05:30	2026-10-14 12:55:13.602029+05:30
125	250	Other model 5730982019	t	2026-10-14 15:51:05.96898+05:30	2026-10-14 15:51:05.836407+05:30	2026-10-14 15:51:05.96898+05:30
81	208	Model 765262501	t	2026-10-14 12:55:27.680659+05:30	2026-10-14 12:55:27.673401+05:30	2026-10-14 12:55:27.680659+05:30
124	249	Model 5730982019	t	2026-10-14 15:51:05.970369+05:30	2026-10-14 15:51:05.836018+05:30	2026-10-14 15:51:05.970369+05:30
79	207	Model 765262501 Updated	t	2026-10-14 12:55:27.693322+05:30	2026-10-14 12:55:27.6672+05:30	2026-10-14 12:55:27.693322+05:30
123	248	Other model 573098201	t	2026-10-14 15:51:05.995111+05:30	2026-10-14 15:51:05.79656+05:30	2026-10-14 15:51:05.995111+05:30
122	247	Model 573098201	t	2026-10-14 15:51:05.995485+05:30	2026-10-14 15:51:05.789946+05:30	2026-10-14 15:51:05.995485+05:30
83	210	Other model 769633603	t	2026-10-14 12:55:27.776774+05:30	2026-10-14 12:55:27.704896+05:30	2026-10-14 12:55:27.776774+05:30
82	209	Model 769633603	t	2026-10-14 12:55:27.777049+05:30	2026-10-14 12:55:27.704476+05:30	2026-10-14 12:55:27.777049+05:30
87	214	Other model 7779789049	t	2026-10-14 12:55:27.834516+05:30	2026-10-14 12:55:27.797291+05:30	2026-10-14 12:55:27.834516+05:30
86	213	Model 7779789049	t	2026-10-14 12:55:27.834762+05:30	2026-10-14 12:55:27.797026+05:30	2026-10-14 12:55:27.834762+05:30
85	212	Other model 777978904	t	2026-10-14 12:55:27.838445+05:30	2026-10-14 12:55:27.782578+05:30	2026-10-14 12:55:27.838445+05:30
84	211	Model 777978904	t	2026-10-14 12:55:27.838706+05:30	2026-10-14 12:55:27.782259+05:30	2026-10-14 12:55:27.838706+05:30
129	254	Other model 6003831029	t	2026-10-14 15:51:06.110286+05:30	2026-10-14 15:51:06.035925+05:30	2026-10-14 15:51:06.110286+05:30
128	253	Model 6003831029	t	2026-10-14 15:51:06.11107+05:30	2026-10-14 15:51:06.035398+05:30	2026-10-14 15:51:06.11107+05:30
127	252	Other model 600383102	t	2026-10-14 15:51:06.119375+05:30	2026-10-14 15:51:06.010009+05:30	2026-10-14 15:51:06.119375+05:30
126	251	Model 600383102	t	2026-10-14 15:51:06.120908+05:30	2026-10-14 15:51:06.009507+05:30	2026-10-14 15:51:06.120908+05:30
91	218	Other model 161897019	t	2026-10-14 15:27:40.359322+05:30	2026-10-14 15:27:40.27347+05:30	2026-10-14 15:27:40.359322+05:30
90	217	Model 161897019	t	2026-10-14 15:27:40.359715+05:30	2026-10-14 15:27:40.273105+05:30	2026-10-14 15:27:40.359715+05:30
89	216	Other model 16189701	t	2026-10-14 15:27:40.364995+05:30	2026-10-14 15:27:40.226618+05:30	2026-10-14 15:27:40.364995+05:30
88	215	Model 16189701	t	2026-10-14 15:27:40.365464+05:30	2026-10-14 15:27:40.213755+05:30	2026-10-14 15:27:40.365464+05:30
94	220	Model 36865502	t	2026-10-14 15:27:40.382461+05:30	2026-10-14 15:27:40.379664+05:30	2026-10-14 15:27:40.382461+05:30
92	219	Model 36865502 Updated	t	2026-10-14 15:27:40.395781+05:30	2026-10-14 15:27:40.374865+05:30	2026-10-14 15:27:40.395781+05:30
133	264	Other model 8387817179	t	2026-10-14 16:09:28.481206+05:30	2026-10-14 16:09:28.425389+05:30	2026-10-14 16:09:28.481206+05:30
132	263	Model 8387817179	t	2026-10-14 16:09:28.481559+05:30	2026-10-14 16:09:28.424972+05:30	2026-10-14 16:09:28.481559+05:30
96	222	Other model 40066304	t	2026-10-14 15:27:40.468521+05:30	2026-10-14 15:27:40.411913+05:30	2026-10-14 15:27:40.468521+05:30
95	221	Model 40066304	t	2026-10-14 15:27:40.46891+05:30	2026-10-14 15:27:40.411596+05:30	2026-10-14 15:27:40.46891+05:30
131	262	Other model 838781717	t	2026-10-14 16:09:28.487267+05:30	2026-10-14 16:09:28.392752+05:30	2026-10-14 16:09:28.487267+05:30
130	261	Model 838781717	t	2026-10-14 16:09:28.487575+05:30	2026-10-14 16:09:28.390772+05:30	2026-10-14 16:09:28.487575+05:30
100	226	Other model 471232059	t	2026-10-14 15:27:40.542144+05:30	2026-10-14 15:27:40.48625+05:30	2026-10-14 15:27:40.542144+05:30
99	225	Model 471232059	t	2026-10-14 15:27:40.542469+05:30	2026-10-14 15:27:40.486012+05:30	2026-10-14 15:27:40.542469+05:30
98	224	Other model 47123205	t	2026-10-14 15:27:40.546596+05:30	2026-10-14 15:27:40.4749+05:30	2026-10-14 15:27:40.546596+05:30
97	223	Model 47123205	t	2026-10-14 15:27:40.546946+05:30	2026-10-14 15:27:40.474614+05:30	2026-10-14 15:27:40.546946+05:30
136	266	Model 849088318	t	2026-10-14 16:09:28.496211+05:30	2026-10-14 16:09:28.495266+05:30	2026-10-14 16:09:28.496211+05:30
134	265	Model 849088318 Updated	t	2026-10-14 16:09:28.501521+05:30	2026-10-14 16:09:28.493286+05:30	2026-10-14 16:09:28.501521+05:30
137	267	Renamed Model 850316820	t	2026-10-14 16:09:28.512498+05:30	2026-10-14 16:09:28.505121+05:30	2026-10-14 16:09:28.512498+05:30
139	270	Other model 851444423	t	2026-10-14 16:09:28.559499+05:30	2026-10-14 16:09:28.519333+05:30	2026-10-14 16:09:28.559499+05:30
138	269	Model 851444423	t	2026-10-14 16:09:28.559849+05:30	2026-10-14 16:09:28.519056+05:30	2026-10-14 16:09:28.559849+05:30
140	271	Model 856771526	t	2026-10-14 16:09:28.62154+05:30	2026-10-14 16:09:28.571611+05:30	2026-10-14 16:09:28.62154+05:30
143	274	Other model 8567715269	t	2026-10-14 16:09:28.617656+05:30	2026-10-14 16:09:28.586117+05:30	2026-10-14 16:09:28.617656+05:30
142	273	Model 8567715269	t	2026-10-14 16:09:28.617935+05:30	2026-10-14 16:09:28.585777+05:30	2026-10-14 16:09:28.617935+05:30
141	272	Other model 856771526	t	2026-10-14 16:09:28.621325+05:30	2026-10-14 16:09:28.571891+05:30	2026-10-14 16:09:28.621325+05:30
147	278	Other model B868524232	t	2026-10-14 16:09:29.164548+05:30	2026-10-14 16:09:28.698231+05:30	2026-10-14 16:09:29.164548+05:30
146	277	Model B868524232	t	2026-10-14 16:09:29.164891+05:30	2026-10-14 16:09:28.69796+05:30	2026-10-14 16:09:29.164891+05:30
145	276	Other model A868524232	t	2026-10-14 16:09:29.169762+05:30	2026-10-14 16:09:28.688002+05:30	2026-10-14 16:09:29.169762+05:30
144	275	Model A868524232	t	2026-10-14 16:09:29.170084+05:30	2026-10-14 16:09:28.68773+05:30	2026-10-14 16:09:29.170084+05:30
151	282	Other model 5216632019	t	2026-10-14 16:26:15.710994+05:30	2026-10-14 16:26:15.583415+05:30	2026-10-14 16:26:15.710994+05:30
150	281	Model 5216632019	t	2026-10-14 16:26:15.711513+05:30	2026-10-14 16:26:15.582698+05:30	2026-10-14 16:26:15.711513+05:30
149	280	Other model 521663201	t	2026-10-14 16:26:15.718647+05:30	2026-10-14 16:26:15.328923+05:30	2026-10-14 16:26:15.718647+05:30
148	279	Model 521663201	t	2026-10-14 16:26:15.719067+05:30	2026-10-14 16:26:15.316215+05:30	2026-10-14 16:26:15.719067+05:30
153	284	Other model 572189502	t	2026-10-14 16:26:15.772704+05:30	2026-10-14 16:26:15.725389+05:30	2026-10-14 16:26:15.772704+05:30
152	283	Model 572189502	t	2026-10-14 16:26:15.772969+05:30	2026-10-14 16:26:15.725051+05:30	2026-10-14 16:26:15.772969+05:30
157	288	Other model 5775058039	t	2026-10-14 16:26:15.816664+05:30	2026-10-14 16:26:15.789395+05:30	2026-10-14 16:26:15.816664+05:30
156	287	Model 5775058039	t	2026-10-14 16:26:15.816857+05:30	2026-10-14 16:26:15.789053+05:30	2026-10-14 16:26:15.816857+05:30
155	286	Other model 577505803	t	2026-10-14 16:26:15.81994+05:30	2026-10-14 16:26:15.778785+05:30	2026-10-14 16:26:15.81994+05:30
154	285	Model 577505803	t	2026-10-14 16:26:15.820131+05:30	2026-10-14 16:26:15.778472+05:30	2026-10-14 16:26:15.820131+05:30
161	298	Other model 7640435179	t	2026-10-15 10:20:47.78322+05:30	2026-10-15 10:20:47.688061+05:30	2026-10-15 10:20:47.78322+05:30
160	297	Model 7640435179	t	2026-10-15 10:20:47.784318+05:30	2026-10-15 10:20:47.687721+05:30	2026-10-15 10:20:47.784318+05:30
159	296	Other model 764043517	t	2026-10-15 10:20:47.789199+05:30	2026-10-15 10:20:47.653705+05:30	2026-10-15 10:20:47.789199+05:30
158	295	Model 764043517	t	2026-10-15 10:20:47.789534+05:30	2026-10-15 10:20:47.643678+05:30	2026-10-15 10:20:47.789534+05:30
164	300	Model 779250818	t	2026-10-15 10:20:47.798295+05:30	2026-10-15 10:20:47.797233+05:30	2026-10-15 10:20:47.798295+05:30
162	299	Model 779250818 Updated	t	2026-10-15 10:20:47.803783+05:30	2026-10-15 10:20:47.79511+05:30	2026-10-15 10:20:47.803783+05:30
165	301	Renamed Model 780558820	t	2026-10-15 10:20:47.813558+05:30	2026-10-15 10:20:47.807588+05:30	2026-10-15 10:20:47.813558+05:30
167	304	Other model 781488023	t	2026-10-15 10:20:47.859594+05:30	2026-10-15 10:20:47.817999+05:30	2026-10-15 10:20:47.859594+05:30
166	303	Model 781488023	t	2026-10-15 10:20:47.859903+05:30	2026-10-15 10:20:47.817715+05:30	2026-10-15 10:20:47.859903+05:30
171	308	Other model 7870395269	t	2026-10-15 10:20:47.925914+05:30	2026-10-15 10:20:47.884351+05:30	2026-10-15 10:20:47.925914+05:30
170	307	Model 7870395269	t	2026-10-15 10:20:47.92621+05:30	2026-10-15 10:20:47.884061+05:30	2026-10-15 10:20:47.92621+05:30
169	306	Other model 787039526	t	2026-10-15 10:20:47.930216+05:30	2026-10-15 10:20:47.873649+05:30	2026-10-15 10:20:47.930216+05:30
168	305	Model 787039526	t	2026-10-15 10:20:47.930449+05:30	2026-10-15 10:20:47.873337+05:30	2026-10-15 10:20:47.930449+05:30
175	312	Other model B801317832	t	2026-10-15 10:20:48.484992+05:30	2026-10-15 10:20:48.183487+05:30	2026-10-15 10:20:48.484992+05:30
174	311	Model B801317832	t	2026-10-15 10:20:48.485293+05:30	2026-10-15 10:20:48.1832+05:30	2026-10-15 10:20:48.485293+05:30
173	310	Other model A801317832	t	2026-10-15 10:20:48.489484+05:30	2026-10-15 10:20:48.015474+05:30	2026-10-15 10:20:48.489484+05:30
172	309	Model A801317832	t	2026-10-15 10:20:48.489717+05:30	2026-10-15 10:20:48.01526+05:30	2026-10-15 10:20:48.489717+05:30
\.


--
-- TOC entry 5702 (class 0 OID 117783)
-- Dependencies: 282
-- Data for Name: part_process_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.part_process_machines (id, process_id, machine_id, created_at, updated_at) FROM stdin;
11	7	13	2026-10-14 14:52:04.853209	2026-10-14 14:52:04.853209
12	8	42	2026-10-14 15:19:25.249141	2026-10-14 15:19:25.249141
13	9	43	2026-10-14 15:27:40.335897	2026-10-14 15:27:40.335897
15	10	47	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.494874
16	11	47	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.494874
18	12	46	2026-10-14 15:27:40.50556	2026-10-14 15:27:40.50556
19	13	49	2026-10-14 15:28:02.04349	2026-10-14 15:28:02.04349
21	14	53	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.16231
22	15	53	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.16231
24	16	52	2026-10-14 15:28:02.171864	2026-10-14 15:28:02.171864
25	17	55	2026-10-14 15:41:11.75549	2026-10-14 15:41:11.75549
27	18	57	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.815485
28	19	57	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.815485
30	20	56	2026-10-14 15:41:11.823714	2026-10-14 15:41:11.823714
31	21	59	2026-10-14 15:51:05.931791	2026-10-14 15:51:05.931791
33	22	61	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.047795
34	23	61	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.047795
36	24	60	2026-10-14 15:51:06.060619	2026-10-14 15:51:06.060619
37	25	63	2026-10-14 16:09:28.463121	2026-10-14 16:09:28.463121
39	26	68	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.59261
40	27	68	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.59261
42	28	67	2026-10-14 16:09:28.602817	2026-10-14 16:09:28.602817
43	29	73	2026-10-14 16:26:15.68487	2026-10-14 16:26:15.68487
45	30	77	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.796292
46	31	77	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.796292
48	32	76	2026-10-14 16:26:15.804379	2026-10-14 16:26:15.804379
49	33	42	2026-10-14 17:49:08.733003	2026-10-14 17:49:08.733003
50	34	79	2026-10-15 10:20:47.751206	2026-10-15 10:20:47.751206
52	35	84	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.890978
53	36	84	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.890978
55	37	83	2026-10-15 10:20:47.901286	2026-10-15 10:20:47.901286
\.


--
-- TOC entry 5700 (class 0 OID 117738)
-- Dependencies: 280
-- Data for Name: part_processes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.part_processes (id, part_id, sort_order, process, program_name, drawing_no, drawing_revision_no, drawing_file, drawing_name, accuracy_id, setup_time, cycle_time, loading_time, production_time, has_sub_program, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
23	17	2	P2	PRG_P2	DRW-P2	1	b44891a6dcd86a7b57c05c191bff0288.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 15:51:06.060619+05:30	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.060619
22	17	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	b44891a6dcd86a7b57c05c191bff0288.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.060619
7	5	1	5871 Sheryl Tunnel	Runolfsson	836 Lang Vista	52114 Emard Knoll	b1e7e55e99b0661b222c2d6c5563f4f3.pdf	Spindle_run_hour_flow.pdf	2	4953	8553	9224	22730	t	f	\N	2026-10-14 14:52:04.853209	2026-10-14 14:52:04.853209
8	6	1	906 Center Road	Volkman	37085 Jovany Fort	197 River Street	fadd1738e78d97fc7271ec0fe955696c.pdf	Spindle_run_hour_flow.pdf	2	5622	5622	5622	16866	f	f	\N	2026-10-14 15:19:25.249141	2026-10-14 15:19:25.249141
9	7	1	P1	PRG_P1	DRW-P1	1	796bcf561e4261a903d3bb5b8d718792.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:27:40.335897	2026-10-14 15:27:40.335897
11	8	2	P2	PRG_P2	DRW-P2	1	394307a5a2216ea08f19788d104f2aa6.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 15:27:40.50556+05:30	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.50556
10	8	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	394307a5a2216ea08f19788d104f2aa6.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.50556
12	8	2	P3	PRG_P3	DRW-P3	1	394307a5a2216ea08f19788d104f2aa6.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:27:40.50556	2026-10-14 15:27:40.50556
13	10	1	P1	PRG_P1	DRW-P1	1	921dc94d7b4a215a339332d02be80641.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:28:02.04349	2026-10-14 15:28:02.04349
15	11	2	P2	PRG_P2	DRW-P2	1	90856328e5503c9899efe9448faa5dd4.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 15:28:02.171864+05:30	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.171864
14	11	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	90856328e5503c9899efe9448faa5dd4.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.171864
16	11	2	P3	PRG_P3	DRW-P3	1	90856328e5503c9899efe9448faa5dd4.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:28:02.171864	2026-10-14 15:28:02.171864
17	13	1	P1	PRG_P1	DRW-P1	1	5e3a4005cdb5b29b22f6226513997e50.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:41:11.75549	2026-10-14 15:41:11.75549
19	14	2	P2	PRG_P2	DRW-P2	1	85d5d543c82357132237fdbf30da72ca.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 15:41:11.823714+05:30	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.823714
18	14	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	85d5d543c82357132237fdbf30da72ca.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.823714
20	14	2	P3	PRG_P3	DRW-P3	1	85d5d543c82357132237fdbf30da72ca.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:41:11.823714	2026-10-14 15:41:11.823714
21	16	1	P1	PRG_P1	DRW-P1	1	e1eb5b6e84742edfc5b67010dcb5ab15.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:51:05.931791	2026-10-14 15:51:05.931791
24	17	2	P3	PRG_P3	DRW-P3	1	b44891a6dcd86a7b57c05c191bff0288.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 15:51:06.060619	2026-10-14 15:51:06.060619
25	19	1	P1	PRG_P1	DRW-P1	1	4bb2b7d41ad5423746e7a524d240ddb5.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:09:28.463121	2026-10-14 16:09:28.463121
27	20	2	P2	PRG_P2	DRW-P2	1	17f161a12dc8d51b208cee6560d45cd1.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 16:09:28.602817+05:30	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.602817
26	20	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	17f161a12dc8d51b208cee6560d45cd1.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.602817
28	20	2	P3	PRG_P3	DRW-P3	1	17f161a12dc8d51b208cee6560d45cd1.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:09:28.602817	2026-10-14 16:09:28.602817
29	22	1	P1	PRG_P1	DRW-P1	1	507b1be1b6c8ca6bc925cdd473db0ed7.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:26:15.68487	2026-10-14 16:26:15.68487
31	23	2	P2	PRG_P2	DRW-P2	1	4384d5709a32cf5b4223e01187dd1e76.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-14 16:26:15.804379+05:30	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.804379
30	23	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	4384d5709a32cf5b4223e01187dd1e76.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.804379
32	23	2	P3	PRG_P3	DRW-P3	1	4384d5709a32cf5b4223e01187dd1e76.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-14 16:26:15.804379	2026-10-14 16:26:15.804379
34	26	1	P1	PRG_P1	DRW-P1	1	d880db5204db5cfaa5384f5575bb136b.pdf	drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-15 10:20:47.751206	2026-10-15 10:20:47.751206
36	27	2	P2	PRG_P2	DRW-P2	1	45ea59099e95f473f0480c63cdcb826c.pdf	shaft drawing.pdf	1	900	3600	900	8100	t	t	2026-10-15 10:20:47.901286+05:30	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.901286
35	27	1	P1-edited	PRG_P1-edited	DRW-P1-edited	1	45ea59099e95f473f0480c63cdcb826c.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.901286
37	27	2	P3	PRG_P3	DRW-P3	1	45ea59099e95f473f0480c63cdcb826c.pdf	shaft drawing.pdf	1	900	3600	900	5400	f	f	\N	2026-10-15 10:20:47.901286	2026-10-15 10:20:47.901286
33	25	1	2906 Broadway Street	Wehner	586 Dennis Alley	34763 Dejon Bypass	f7ab5ffb86c6a546b465b3af7ef6b964.pdf	7th_sense_workflow.pdf	2	4210	4210	4210	12630	f	f	\N	2026-10-14 17:49:08.733003	2026-10-15 12:01:33.940098
\.


--
-- TOC entry 5704 (class 0 OID 117809)
-- Dependencies: 284
-- Data for Name: part_sub_programs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.part_sub_programs (id, process_id, sort_order, name, cycle_time, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
5	7	1	S_SUB1	0	f	\N	2026-10-14 14:52:04.853209	2026-10-14 14:52:04.853209
6	11	1	S_SUB1	1800	f	\N	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.494874
7	11	2	S_SUB2	900	f	\N	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.494874
8	15	1	S_SUB1	1800	f	\N	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.16231
9	15	2	S_SUB2	900	f	\N	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.16231
10	19	1	S_SUB1	1800	f	\N	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.815485
11	19	2	S_SUB2	900	f	\N	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.815485
12	23	1	S_SUB1	1800	f	\N	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.047795
13	23	2	S_SUB2	900	f	\N	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.047795
14	27	1	S_SUB1	1800	f	\N	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.59261
15	27	2	S_SUB2	900	f	\N	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.59261
16	31	1	S_SUB1	1800	f	\N	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.796292
17	31	2	S_SUB2	900	f	\N	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.796292
18	36	1	S_SUB1	1800	f	\N	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.890978
19	36	2	S_SUB2	900	f	\N	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.890978
\.


--
-- TOC entry 5698 (class 0 OID 117706)
-- Dependencies: 278
-- Data for Name: parts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.parts (id, user_id, clientele_id, part_no, part_name, part_count, production_time, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
5	23	\N	GAXJD3548U	Larkin	18409	22730	active	f	\N	2026-10-14 14:52:04.853209	2026-10-14 14:52:04.853209
6	22	\N	ZGIAH2501V	O'Connell	66341	16866	active	f	\N	2026-10-14 15:19:25.249141	2026-10-14 15:19:25.249141
7	67	1	PN-C-16189701	Clientele part 16189701	1	5400	active	t	2026-10-14 15:27:40.352857+05:30	2026-10-14 15:27:40.335897	2026-10-14 15:27:40.352857
8	70	\N	PN-47123205	Shaft edited 47123205	2	10800	inactive	t	2026-10-14 15:27:40.536658+05:30	2026-10-14 15:27:40.494874	2026-10-14 15:27:40.536658
10	72	6	PN-C-193292601	Clientele part 193292601	1	5400	active	t	2026-10-14 15:28:02.05847+05:30	2026-10-14 15:28:02.04349	2026-10-14 15:28:02.05847
11	75	\N	PN-214085705	Shaft edited 214085705	2	10800	inactive	t	2026-10-14 15:28:02.183806+05:30	2026-10-14 15:28:02.16231	2026-10-14 15:28:02.183806
13	77	12	PN-C-158942201	Clientele part 158942201	1	5400	active	t	2026-10-14 15:41:11.771543+05:30	2026-10-14 15:41:11.75549	2026-10-14 15:41:11.771543
14	79	\N	PN-178563102	Shaft edited 178563102	2	10800	inactive	t	2026-10-14 15:41:11.834186+05:30	2026-10-14 15:41:11.815485	2026-10-14 15:41:11.834186
16	81	17	PN-C-573098201	Clientele part 573098201	1	5400	active	t	2026-10-14 15:51:05.957885+05:30	2026-10-14 15:51:05.931791	2026-10-14 15:51:05.957885
17	83	\N	PN-600383102	Shaft edited 600383102	2	10800	inactive	t	2026-10-14 15:51:06.099264+05:30	2026-10-14 15:51:06.047795	2026-10-14 15:51:06.099264
19	90	22	PN-C-838781717	Clientele part 838781717	1	5400	active	t	2026-10-14 16:09:28.476197+05:30	2026-10-14 16:09:28.463121	2026-10-14 16:09:28.476197
20	94	\N	PN-856771526	Shaft edited 856771526	2	10800	inactive	t	2026-10-14 16:09:28.614167+05:30	2026-10-14 16:09:28.59261	2026-10-14 16:09:28.614167
22	106	27	PN-C-521663201	Clientele part 521663201	1	5400	active	t	2026-10-14 16:26:15.704785+05:30	2026-10-14 16:26:15.68487	2026-10-14 16:26:15.704785
23	109	\N	PN-577505803	Shaft edited 577505803	2	10800	inactive	t	2026-10-14 16:26:15.813744+05:30	2026-10-14 16:26:15.796292	2026-10-14 16:26:15.813744
26	126	37	PN-C-764043517	Clientele part 764043517	1	5400	active	t	2026-10-15 10:20:47.777777+05:30	2026-10-15 10:20:47.751206	2026-10-15 10:20:47.777777
27	130	\N	PN-787039526	Shaft edited 787039526	2	10800	inactive	t	2026-10-15 10:20:47.922535+05:30	2026-10-15 10:20:47.890978	2026-10-15 10:20:47.922535
25	22	36	NBGJW3534G	Rolfson	14530	12630	inactive	f	\N	2026-10-14 17:49:08.733003	2026-10-15 12:01:33.940098
\.


--
-- TOC entry 5666 (class 0 OID 116883)
-- Dependencies: 246
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, code, category_id, image, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Product 148465800 Updated	TC148465800	5	lathe.png	\N	t	2026-10-12 17:48:44.157782+05:30	2026-10-12 17:48:44.149268	2026-10-12 17:48:44.157782
49	Test Product 88364700 Updated	TC88364700	45	lathe.png	\N	t	2026-10-13 15:14:27.113127+05:30	2026-10-13 15:14:27.089429	2026-10-13 15:14:27.113127
4	Test Product 148465800 Updated	TC148465800	5	lathe.png	\N	t	2026-10-12 17:48:44.15966+05:30	2026-10-12 17:48:44.159117	2026-10-12 17:48:44.15966
52	Test Product 88364700 Updated	TC88364700	45	lathe.png	\N	t	2026-10-13 15:14:27.114855+05:30	2026-10-13 15:14:27.114333	2026-10-13 15:14:27.114855
6	Test Product 161014200	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:44.162273+05:30	2026-10-12 17:48:44.161869	2026-10-12 17:48:44.162273
54	Test Product 115845500	\N	\N	no_image_available.png	\N	t	2026-10-13 15:14:27.117006+05:30	2026-10-13 15:14:27.116641	2026-10-13 15:14:27.117006
7	Test Product 682313000 Updated	TC682313000	10	lathe.png	\N	t	2026-10-12 17:48:50.692503+05:30	2026-10-12 17:48:50.683849	2026-10-12 17:48:50.692503
10	Test Product 682313000 Updated	TC682313000	10	lathe.png	\N	t	2026-10-12 17:48:50.694535+05:30	2026-10-12 17:48:50.693926	2026-10-12 17:48:50.694535
12	Test Product 696127300	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:50.697607+05:30	2026-10-12 17:48:50.696914	2026-10-12 17:48:50.697607
55	Test Product 756217500 Updated	TC756217500	50	lathe.png	\N	t	2026-10-13 15:48:34.776296+05:30	2026-10-13 15:48:34.756721	2026-10-13 15:48:34.776296
58	Test Product 756217500 Updated	TC756217500	50	lathe.png	\N	t	2026-10-13 15:48:34.778468+05:30	2026-10-13 15:48:34.77759	2026-10-13 15:48:34.778468
13	Test Product 274269500 Updated	TC274269500	15	lathe.png	\N	t	2026-10-12 17:48:53.281889+05:30	2026-10-12 17:48:53.274804	2026-10-12 17:48:53.281889
16	Test Product 274269500 Updated	TC274269500	15	lathe.png	\N	t	2026-10-12 17:48:53.2832+05:30	2026-10-12 17:48:53.282777	2026-10-12 17:48:53.2832
60	Test Product 780965500	\N	\N	no_image_available.png	\N	t	2026-10-13 15:48:34.781973+05:30	2026-10-13 15:48:34.781401	2026-10-13 15:48:34.781973
18	Test Product 284189700	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:53.285268+05:30	2026-10-12 17:48:53.284917	2026-10-12 17:48:53.285268
61	Model 394910600	\N	\N	no_image_available.png	\N	t	2026-10-13 16:26:29.515733+05:30	2026-10-13 16:26:29.402377	2026-10-13 16:26:29.515733
19	Test Product 916338700 Updated	TC916338700	20	lathe.png	\N	t	2026-10-12 17:48:55.924207+05:30	2026-10-12 17:48:55.916769	2026-10-12 17:48:55.924207
22	Test Product 916338700 Updated	TC916338700	20	lathe.png	\N	t	2026-10-12 17:48:55.925796+05:30	2026-10-12 17:48:55.925295	2026-10-12 17:48:55.925796
63	Test Product 217444000 Updated	TC217444000	55	lathe.png	\N	t	2026-10-13 18:00:14.235226+05:30	2026-10-13 18:00:14.218034	2026-10-13 18:00:14.235226
24	Test Product 927046000	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:55.927792+05:30	2026-10-12 17:48:55.927456	2026-10-12 17:48:55.927792
66	Test Product 217444000 Updated	TC217444000	55	lathe.png	\N	t	2026-10-13 18:00:14.23696+05:30	2026-10-13 18:00:14.236397	2026-10-13 18:00:14.23696
25	Test Product 613950100 Updated	TC613950100	25	lathe.png	\N	t	2026-10-13 11:20:45.623118+05:30	2026-10-13 11:20:45.614442	2026-10-13 11:20:45.623118
68	Test Product 238139900	\N	\N	no_image_available.png	\N	t	2026-10-13 18:00:14.23936+05:30	2026-10-13 18:00:14.238911	2026-10-13 18:00:14.23936
28	Test Product 613950100 Updated	TC613950100	25	lathe.png	\N	t	2026-10-13 11:20:45.624731+05:30	2026-10-13 11:20:45.624358	2026-10-13 11:20:45.624731
30	Test Product 625742900	\N	\N	no_image_available.png	\N	t	2026-10-13 11:20:45.62664+05:30	2026-10-13 11:20:45.626309	2026-10-13 11:20:45.62664
69	Test Product 563547804 Updated	TC563547804	60	lathe.png	\N	t	2026-10-13 18:01:25.644518+05:30	2026-10-13 18:01:25.636366	2026-10-13 18:01:25.644518
31	Test Product 621062000 Updated	TC621062000	30	lathe.png	\N	t	2026-10-13 13:25:36.631412+05:30	2026-10-13 13:25:36.621797	2026-10-13 13:25:36.631412
72	Test Product 563547804 Updated	TC563547804	60	lathe.png	\N	t	2026-10-13 18:01:25.646254+05:30	2026-10-13 18:01:25.645699	2026-10-13 18:01:25.646254
34	Test Product 621062000 Updated	TC621062000	30	lathe.png	\N	t	2026-10-13 13:25:36.63301+05:30	2026-10-13 13:25:36.632474	2026-10-13 13:25:36.63301
36	Test Product 634745200	\N	\N	no_image_available.png	\N	t	2026-10-13 13:25:36.636254+05:30	2026-10-13 13:25:36.635623	2026-10-13 13:25:36.636254
74	Test Product 564769807	\N	\N	no_image_available.png	\N	t	2026-10-13 18:01:25.648639+05:30	2026-10-13 18:01:25.648199	2026-10-13 18:01:25.648639
37	Test Product 310985800 Updated	TC310985800	35	lathe.png	\N	t	2026-10-13 13:25:45.318253+05:30	2026-10-13 13:25:45.311748	2026-10-13 13:25:45.318253
40	Test Product 310985800 Updated	TC310985800	35	lathe.png	\N	t	2026-10-13 13:25:45.319484+05:30	2026-10-13 13:25:45.319116	2026-10-13 13:25:45.319484
75	Test Product 967484204 Updated	TC967484204	65	lathe.png	\N	t	2026-10-13 18:01:39.682794+05:30	2026-10-13 18:01:39.675698	2026-10-13 18:01:39.682794
42	Test Product 320199400	\N	\N	no_image_available.png	\N	t	2026-10-13 13:25:45.321284+05:30	2026-10-13 13:25:45.32093	2026-10-13 13:25:45.321284
78	Test Product 967484204 Updated	TC967484204	65	lathe.png	\N	t	2026-10-13 18:01:39.684141+05:30	2026-10-13 18:01:39.683687	2026-10-13 18:01:39.684141
80	Test Product 968526307	\N	\N	no_image_available.png	\N	t	2026-10-13 18:01:39.686363+05:30	2026-10-13 18:01:39.685965	2026-10-13 18:01:39.686363
43	Test Product 71778500 Updated	TC71778500	40	lathe.png	\N	t	2026-10-13 14:57:56.080716+05:30	2026-10-13 14:57:56.07259	2026-10-13 14:57:56.080716
46	Test Product 71778500 Updated	TC71778500	40	lathe.png	\N	t	2026-10-13 14:57:56.082301+05:30	2026-10-13 14:57:56.081794	2026-10-13 14:57:56.082301
48	Test Product 85704300	\N	\N	no_image_available.png	\N	t	2026-10-13 14:57:56.087222+05:30	2026-10-13 14:57:56.086616	2026-10-13 14:57:56.087222
81	Test Product 716313705 Updated	TC716313705	70	lathe.png	\N	t	2026-10-13 18:01:47.17125+05:30	2026-10-13 18:01:47.163559	2026-10-13 18:01:47.17125
84	Test Product 716313705 Updated	TC716313705	70	lathe.png	\N	t	2026-10-13 18:01:47.172804+05:30	2026-10-13 18:01:47.172338	2026-10-13 18:01:47.172804
86	Test Product 717401608	\N	\N	no_image_available.png	\N	t	2026-10-13 18:01:47.175048+05:30	2026-10-13 18:01:47.174603	2026-10-13 18:01:47.175048
87	Test Product 830472505 Updated	TC830472505	75	lathe.png	\N	t	2026-10-14 09:37:38.371076+05:30	2026-10-14 09:37:38.306327	2026-10-14 09:37:38.371076
90	Test Product 830472505 Updated	TC830472505	75	lathe.png	\N	t	2026-10-14 09:37:38.469975+05:30	2026-10-14 09:37:38.441006	2026-10-14 09:37:38.469975
92	Test Product 848256608	\N	\N	no_image_available.png	\N	t	2026-10-14 09:37:38.486421+05:30	2026-10-14 09:37:38.483381	2026-10-14 09:37:38.486421
93	Test Product 622839405 Updated	TC622839405	80	lathe.png	\N	t	2026-10-14 09:38:36.238464+05:30	2026-10-14 09:38:36.229222	2026-10-14 09:38:36.238464
96	Test Product 622839405 Updated	TC622839405	80	lathe.png	\N	t	2026-10-14 09:38:36.240119+05:30	2026-10-14 09:38:36.239547	2026-10-14 09:38:36.240119
98	Test Product 624190708	\N	\N	no_image_available.png	\N	t	2026-10-14 09:38:36.242791+05:30	2026-10-14 09:38:36.24243	2026-10-14 09:38:36.242791
99	Test Product 2889005 Updated	TC2889005	85	lathe.png	\N	t	2026-10-14 11:57:30.050063+05:30	2026-10-14 11:57:30.029754	2026-10-14 11:57:30.050063
102	Test Product 2889005 Updated	TC2889005	85	lathe.png	\N	t	2026-10-14 11:57:30.051527+05:30	2026-10-14 11:57:30.051021	2026-10-14 11:57:30.051527
104	Test Product 5291208	\N	\N	no_image_available.png	\N	t	2026-10-14 11:57:30.053887+05:30	2026-10-14 11:57:30.053403	2026-10-14 11:57:30.053887
105	Test Product 58741105 Updated	TC58741105	90	lathe.png	\N	t	2026-10-14 12:07:00.610162+05:30	2026-10-14 12:07:00.588442	2026-10-14 12:07:00.610162
108	Test Product 58741105 Updated	TC58741105	90	lathe.png	\N	t	2026-10-14 12:07:00.613363+05:30	2026-10-14 12:07:00.612591	2026-10-14 12:07:00.613363
110	Test Product 61540908	\N	\N	no_image_available.png	\N	t	2026-10-14 12:07:00.616895+05:30	2026-10-14 12:07:00.616419	2026-10-14 12:07:00.616895
111	Test Product 886809705 Updated	TC886809705	95	lathe.png	\N	t	2026-10-14 12:41:48.887902+05:30	2026-10-14 12:41:48.868843	2026-10-14 12:41:48.887902
114	Test Product 886809705 Updated	TC886809705	95	lathe.png	\N	t	2026-10-14 12:41:48.890674+05:30	2026-10-14 12:41:48.889848	2026-10-14 12:41:48.890674
116	Test Product 889314108	\N	\N	no_image_available.png	\N	t	2026-10-14 12:41:48.895313+05:30	2026-10-14 12:41:48.894464	2026-10-14 12:41:48.895313
117	Test Product 833006005 Updated	TC833006005	100	lathe.png	\N	t	2026-10-14 16:09:28.349827+05:30	2026-10-14 16:09:28.330848	2026-10-14 16:09:28.349827
120	Test Product 833006005 Updated	TC833006005	100	lathe.png	\N	t	2026-10-14 16:09:28.352919+05:30	2026-10-14 16:09:28.35247	2026-10-14 16:09:28.352919
122	Test Product 835559408	\N	\N	no_image_available.png	\N	t	2026-10-14 16:09:28.356814+05:30	2026-10-14 16:09:28.356408	2026-10-14 16:09:28.356814
123	Test Product 756506205 Updated	TC756506205	105	lathe.png	\N	t	2026-10-15 10:20:47.591207+05:30	2026-10-15 10:20:47.566008	2026-10-15 10:20:47.591207
126	Test Product 756506205 Updated	TC756506205	105	lathe.png	\N	t	2026-10-15 10:20:47.5929+05:30	2026-10-15 10:20:47.592376	2026-10-15 10:20:47.5929
128	Test Product 759514408	\N	\N	no_image_available.png	\N	t	2026-10-15 10:20:47.596842+05:30	2026-10-15 10:20:47.596428	2026-10-15 10:20:47.596842
\.


--
-- TOC entry 5646 (class 0 OID 116542)
-- Dependencies: 226
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, description, created_at) FROM stdin;
1	superadmin	Full access to all resources including brands, controllers, products	2026-10-13 17:24:43.756387
2	admin	Can manage controllers and products, cannot manage brands	2026-10-13 17:24:43.756387
3	supervisor	Supervisor with management permissions	2026-10-13 18:13:08.543149
4	operator	Machine operator	2026-10-13 18:13:08.543149
5	customer	Customer with limited access	2026-10-13 18:13:08.543149
6	qc	Quality control staff	2026-10-13 18:13:08.543149
7	shop_head	Shop head/manager	2026-10-13 18:13:08.543149
\.


--
-- TOC entry 5708 (class 0 OID 117890)
-- Dependencies: 288
-- Data for Name: shifts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.shifts (id, user_id, name, shift_type, start_time, end_time, working_days, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
3	85	E-612267903	night	22:00:00	06:00:00	{1,2,3,7}	active	f	\N	2026-10-14 15:51:06.147848	2026-10-14 15:51:06.147848
4	86	O-612267903	night	22:00:00	06:00:00	{1,2,3,7}	active	f	\N	2026-10-14 15:51:06.160686	2026-10-14 15:51:06.160686
1	85	A-612267903	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-14 15:51:06.162537+05:30	2026-10-14 15:51:06.128307	2026-10-14 15:51:06.162537
6	22	SHIFT BBBB	day	18:30:00	19:00:00	{1,2,3,4,5,6,7}	active	f	\N	2026-10-14 16:06:47.851801	2026-10-14 16:06:47.851801
7	23	SHIFT KKKk	day	15:00:00	15:30:00	{1,2,3,4,5,6,7}	active	f	\N	2026-10-14 16:07:04.152334	2026-10-14 16:07:04.152334
8	96	A-862567827	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-14 16:09:28.661243+05:30	2026-10-14 16:09:28.643746	2026-10-14 16:09:28.661243
5	22	SHIFT A	night	00:00:00	23:59:00	{7}	active	f	\N	2026-10-14 16:00:08.946698	2026-10-14 16:25:07.008604
10	49	Olson	night	20:15:00	20:25:00	{1,2}	active	t	2026-10-14 16:25:11.138211+05:30	2026-10-14 16:17:14.092098	2026-10-14 16:25:11.138211
11	111	A-582207604	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-14 16:26:15.864865+05:30	2026-10-14 16:26:15.841396	2026-10-14 16:26:15.864865
13	115	A-98919701	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-14 16:26:31.043887+05:30	2026-10-14 16:26:31.021994	2026-10-14 16:26:31.043887
15	117	A-52540701	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-14 16:26:41.265472+05:30	2026-10-14 16:26:41.243053	2026-10-14 16:26:41.265472
17	132	A-793253827	day	09:00:00	17:30:00	{1,2,3,7}	inactive	t	2026-10-15 10:20:47.982587+05:30	2026-10-15 10:20:47.96171	2026-10-15 10:20:47.982587
\.


--
-- TOC entry 5644 (class 0 OID 116527)
-- Dependencies: 224
-- Data for Name: states; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.states (id, country_id, name, code, created_at) FROM stdin;
1	1	Gujarat	GJ	2026-10-12 10:38:40.418244
2	1	Maharashtra	MH	2026-10-12 10:38:40.418244
3	1	Karnataka	KA	2026-10-12 10:38:40.418244
4	2	California	CA	2026-10-12 10:38:40.418244
\.


--
-- TOC entry 5690 (class 0 OID 117332)
-- Dependencies: 270
-- Data for Name: user_cell; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_cell (id, user_id, user_line_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	49	56	Cell - 1	f	\N	2026-10-15 11:25:35.402747	2026-10-15 11:25:35.402747
2	49	56	Cell - 2	f	\N	2026-10-15 11:25:43.433424	2026-10-15 11:25:43.433424
\.


--
-- TOC entry 5682 (class 0 OID 117182)
-- Dependencies: 262
-- Data for Name: user_city; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_city (id, user_id, user_state_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
4	26	6	City 13336822	t	2026-10-14 11:57:30.215434+05:30	2026-10-14 11:57:30.143354	2026-10-14 11:57:30.215434
5	31	8	City 73394922	t	2026-10-14 12:07:00.816065+05:30	2026-10-14 12:07:00.747493	2026-10-14 12:07:00.816065
7	34	11	City B83246528	t	2026-10-14 12:07:01.143894+05:30	2026-10-14 12:07:00.854015	2026-10-14 12:07:01.143894
6	33	10	City A83246528	t	2026-10-14 12:07:01.148643+05:30	2026-10-14 12:07:00.836085	2026-10-14 12:07:01.148643
8	42	12	City 847312702	t	2026-10-14 12:07:08.559521+05:30	2026-10-14 12:07:08.488777	2026-10-14 12:07:08.559521
10	44	14	City B857087506	t	2026-10-14 12:07:08.861544+05:30	2026-10-14 12:07:08.587935	2026-10-14 12:07:08.861544
9	43	13	City A857087506	t	2026-10-14 12:07:08.865621+05:30	2026-10-14 12:07:08.575001	2026-10-14 12:07:08.865621
11	54	15	City 899315822	t	2026-10-14 12:41:49.092793+05:30	2026-10-14 12:41:49.009211	2026-10-14 12:41:49.092793
13	57	18	City B911479328	t	2026-10-14 12:41:49.445471+05:30	2026-10-14 12:41:49.141964	2026-10-14 12:41:49.445471
12	56	17	City A911479328	t	2026-10-14 12:41:49.451008+05:30	2026-10-14 12:41:49.120174	2026-10-14 12:41:49.451008
15	63	20	City 3279490019	t	2026-10-14 12:55:13.593102+05:30	2026-10-14 12:55:13.495368	2026-10-14 12:55:13.593102
14	62	19	City 327949001	t	2026-10-14 12:55:13.600797+05:30	2026-10-14 12:55:13.377263	2026-10-14 12:55:13.600797
16	64	21	City 769633603	t	2026-10-14 12:55:27.77546+05:30	2026-10-14 12:55:27.714431	2026-10-14 12:55:27.77546
18	66	23	City 7779789049	t	2026-10-14 12:55:27.833756+05:30	2026-10-14 12:55:27.798235	2026-10-14 12:55:27.833756
17	65	22	City 777978904	t	2026-10-14 12:55:27.837582+05:30	2026-10-14 12:55:27.78357	2026-10-14 12:55:27.837582
20	68	25	City 161897019	t	2026-10-14 15:27:40.357406+05:30	2026-10-14 15:27:40.274921	2026-10-14 15:27:40.357406
19	67	24	City 16189701	t	2026-10-14 15:27:40.363905+05:30	2026-10-14 15:27:40.24184	2026-10-14 15:27:40.363905
21	69	26	City 40066304	t	2026-10-14 15:27:40.467454+05:30	2026-10-14 15:27:40.413185	2026-10-14 15:27:40.467454
23	71	28	City 471232059	t	2026-10-14 15:27:40.540228+05:30	2026-10-14 15:27:40.487373	2026-10-14 15:27:40.540228
22	70	27	City 47123205	t	2026-10-14 15:27:40.545707+05:30	2026-10-14 15:27:40.476075	2026-10-14 15:27:40.545707
25	73	30	City 1932926019	t	2026-10-14 15:28:02.062931+05:30	2026-10-14 15:28:01.996603	2026-10-14 15:28:02.062931
24	72	29	City 193292601	t	2026-10-14 15:28:02.068355+05:30	2026-10-14 15:28:01.972099	2026-10-14 15:28:02.068355
26	74	31	City 209017604	t	2026-10-14 15:28:02.13735+05:30	2026-10-14 15:28:02.095358	2026-10-14 15:28:02.13735
28	76	33	City 2140857059	t	2026-10-14 15:28:02.186338+05:30	2026-10-14 15:28:02.156621	2026-10-14 15:28:02.186338
27	75	32	City 214085705	t	2026-10-14 15:28:02.190286+05:30	2026-10-14 15:28:02.145576	2026-10-14 15:28:02.190286
30	78	35	City 1589422019	t	2026-10-14 15:41:11.775796+05:30	2026-10-14 15:41:11.689821	2026-10-14 15:41:11.775796
29	77	34	City 158942201	t	2026-10-14 15:41:11.781852+05:30	2026-10-14 15:41:11.658144	2026-10-14 15:41:11.781852
32	80	37	City 1785631029	t	2026-10-14 15:41:11.836345+05:30	2026-10-14 15:41:11.807187	2026-10-14 15:41:11.836345
31	79	36	City 178563102	t	2026-10-14 15:41:11.840353+05:30	2026-10-14 15:41:11.790933	2026-10-14 15:41:11.840353
34	82	39	City 5730982019	t	2026-10-14 15:51:05.965118+05:30	2026-10-14 15:51:05.837758	2026-10-14 15:51:05.965118
33	81	38	City 573098201	t	2026-10-14 15:51:05.993731+05:30	2026-10-14 15:51:05.807297	2026-10-14 15:51:05.993731
36	84	41	City 6003831029	t	2026-10-14 15:51:06.108951+05:30	2026-10-14 15:51:06.038314	2026-10-14 15:51:06.108951
35	83	40	City 600383102	t	2026-10-14 15:51:06.116687+05:30	2026-10-14 15:51:06.012542	2026-10-14 15:51:06.116687
38	91	43	City 8387817179	t	2026-10-14 16:09:28.479754+05:30	2026-10-14 16:09:28.426596	2026-10-14 16:09:28.479754
37	90	42	City 838781717	t	2026-10-14 16:09:28.484529+05:30	2026-10-14 16:09:28.401992	2026-10-14 16:09:28.484529
39	92	44	City 851444423	t	2026-10-14 16:09:28.558545+05:30	2026-10-14 16:09:28.520287	2026-10-14 16:09:28.558545
41	95	47	City 8567715269	t	2026-10-14 16:09:28.616597+05:30	2026-10-14 16:09:28.587175	2026-10-14 16:09:28.616597
40	94	46	City 856771526	t	2026-10-14 16:09:28.620532+05:30	2026-10-14 16:09:28.574857	2026-10-14 16:09:28.620532
43	101	49	City B868524232	t	2026-10-14 16:09:29.163469+05:30	2026-10-14 16:09:28.699126	2026-10-14 16:09:29.163469
42	100	48	City A868524232	t	2026-10-14 16:09:29.168579+05:30	2026-10-14 16:09:28.68887	2026-10-14 16:09:29.168579
45	107	51	City 5216632019	t	2026-10-14 16:26:15.709305+05:30	2026-10-14 16:26:15.585793	2026-10-14 16:26:15.709305
44	106	50	City 521663201	t	2026-10-14 16:26:15.716102+05:30	2026-10-14 16:26:15.351664	2026-10-14 16:26:15.716102
46	108	52	City 572189502	t	2026-10-14 16:26:15.771755+05:30	2026-10-14 16:26:15.726594	2026-10-14 16:26:15.771755
48	110	54	City 5775058039	t	2026-10-14 16:26:15.815824+05:30	2026-10-14 16:26:15.790656	2026-10-14 16:26:15.815824
47	109	53	City 577505803	t	2026-10-14 16:26:15.819161+05:30	2026-10-14 16:26:15.779997	2026-10-14 16:26:15.819161
49	120	56	Surat	t	2026-10-15 10:15:01.027848+05:30	2026-10-15 10:15:00.924097	2026-10-15 10:15:01.027848
50	122	62	Surat	t	2026-10-15 10:15:10.013943+05:30	2026-10-15 10:15:09.999605	2026-10-15 10:15:10.013943
52	127	68	City 7640435179	t	2026-10-15 10:20:47.781514+05:30	2026-10-15 10:20:47.689258	2026-10-15 10:20:47.781514
51	126	67	City 764043517	t	2026-10-15 10:20:47.788329+05:30	2026-10-15 10:20:47.659181	2026-10-15 10:20:47.788329
53	128	69	City 781488023	t	2026-10-15 10:20:47.85856+05:30	2026-10-15 10:20:47.818917	2026-10-15 10:20:47.85856
55	131	72	City 7870395269	t	2026-10-15 10:20:47.92484+05:30	2026-10-15 10:20:47.885398	2026-10-15 10:20:47.92484
54	130	71	City 787039526	t	2026-10-15 10:20:47.929306+05:30	2026-10-15 10:20:47.874772	2026-10-15 10:20:47.929306
57	137	74	City B801317832	t	2026-10-15 10:20:48.483994+05:30	2026-10-15 10:20:48.184457	2026-10-15 10:20:48.483994
56	136	73	City A801317832	t	2026-10-15 10:20:48.488663+05:30	2026-10-15 10:20:48.016231	2026-10-15 10:20:48.488663
58	142	75	Surat	t	2026-10-15 10:20:48.50288+05:30	2026-10-15 10:20:48.501276	2026-10-15 10:20:48.50288
59	49	84	hi	f	\N	2026-10-15 10:43:14.849158	2026-10-15 10:43:14.849158
60	49	84	hihihih	f	\N	2026-10-15 10:43:19.262234	2026-10-15 10:43:19.262234
61	49	84	98987987	f	\N	2026-10-15 10:43:25.578519	2026-10-15 10:43:25.578519
62	49	84	97656	f	\N	2026-10-15 10:43:38.255573	2026-10-15 10:43:38.255573
63	49	84	47545765	f	\N	2026-10-15 10:44:25.949338	2026-10-15 10:44:25.949338
64	1	91	43ree	f	\N	2026-10-15 10:55:07.864759	2026-10-15 10:55:07.864759
65	1	92	4tret	f	\N	2026-10-15 10:55:13.559792	2026-10-15 10:55:13.559792
66	49	103	08098098	f	\N	2026-10-15 11:24:29.661407	2026-10-15 11:24:29.661407
67	22	104	22332	f	\N	2026-10-15 11:31:40.357318	2026-10-15 11:31:40.357318
\.


--
-- TOC entry 5678 (class 0 OID 117129)
-- Dependencies: 258
-- Data for Name: user_country; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_country (id, user_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
68	107	Country 5216632019	t	2026-10-14 16:26:15.710426+05:30	2026-10-14 16:26:15.584125	2026-10-14 16:26:15.710426
67	106	Country 521663201	t	2026-10-14 16:26:15.716816+05:30	2026-10-14 16:26:15.329609	2026-10-14 16:26:15.716816
69	108	Country 572189502	t	2026-10-14 16:26:15.772391+05:30	2026-10-14 16:26:15.725756	2026-10-14 16:26:15.772391
71	110	Country 5775058039	t	2026-10-14 16:26:15.816386+05:30	2026-10-14 16:26:15.789698	2026-10-14 16:26:15.816386
70	109	Country 577505803	t	2026-10-14 16:26:15.819679+05:30	2026-10-14 16:26:15.779106	2026-10-14 16:26:15.819679
73	119	Owned 60557102	f	\N	2026-10-15 10:15:00.75386	2026-10-15 10:15:00.75386
13	26	Country 13336822	t	2026-10-14 11:57:30.216359+05:30	2026-10-14 11:57:30.137297	2026-10-14 11:57:30.216359
77	120	Germany	t	2026-10-15 10:15:00.865818+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:00.865818
15	27	Owned 22057624	f	\N	2026-10-14 11:57:30.222706	2026-10-14 11:57:30.222706
75	120	India	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:01.040738
76	120	Japan	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:01.040738
17	31	Country 73394922	t	2026-10-14 12:07:00.81704+05:30	2026-10-14 12:07:00.740825	2026-10-14 12:07:00.81704
19	32	Owned 82060724	f	\N	2026-10-14 12:07:00.82242	2026-10-14 12:07:00.82242
79	121	Owned 990636302	f	\N	2026-10-15 10:15:09.915054	2026-10-15 10:15:09.915054
22	34	Country B83246528	t	2026-10-14 12:07:01.144502+05:30	2026-10-14 12:07:00.852959	2026-10-14 12:07:01.144502
21	33	Country A83246528	t	2026-10-14 12:07:01.14924+05:30	2026-10-14 12:07:00.83541	2026-10-14 12:07:01.14924
23	42	Country 847312702	t	2026-10-14 12:07:08.56074+05:30	2026-10-14 12:07:08.483636	2026-10-14 12:07:08.56074
83	122	Germany	t	2026-10-15 10:15:09.98669+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:09.98669
25	44	Country B857087506	t	2026-10-14 12:07:08.862051+05:30	2026-10-14 12:07:08.587273	2026-10-14 12:07:08.862051
24	43	Country A857087506	t	2026-10-14 12:07:08.866507+05:30	2026-10-14 12:07:08.574335	2026-10-14 12:07:08.866507
81	122	India	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:10.026868
26	54	Country 899315822	t	2026-10-14 12:41:49.093848+05:30	2026-10-14 12:41:48.997093	2026-10-14 12:41:49.093848
82	122	Japan	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:10.026868
28	55	Owned 909691024	f	\N	2026-10-14 12:41:49.0992	2026-10-14 12:41:49.0992
85	127	Country 7640435179	t	2026-10-15 10:20:47.782695+05:30	2026-10-15 10:20:47.68838	2026-10-15 10:20:47.782695
31	57	Country B911479328	t	2026-10-14 12:41:49.446137+05:30	2026-10-14 12:41:49.140286	2026-10-14 12:41:49.446137
30	56	Country A911479328	t	2026-10-14 12:41:49.452075+05:30	2026-10-14 12:41:49.118427	2026-10-14 12:41:49.452075
84	126	Country 764043517	t	2026-10-15 10:20:47.788891+05:30	2026-10-15 10:20:47.654325	2026-10-15 10:20:47.788891
33	63	Country 3279490019	t	2026-10-14 12:55:13.594463+05:30	2026-10-14 12:55:13.49408	2026-10-14 12:55:13.594463
32	62	Country 327949001	t	2026-10-14 12:55:13.601444+05:30	2026-10-14 12:55:13.360155	2026-10-14 12:55:13.601444
86	128	Country 781488023	t	2026-10-15 10:20:47.859261+05:30	2026-10-15 10:20:47.818221	2026-10-15 10:20:47.859261
34	64	Country 769633603	t	2026-10-14 12:55:27.776308+05:30	2026-10-14 12:55:27.705501	2026-10-14 12:55:27.776308
88	129	Owned 786226025	f	\N	2026-10-15 10:20:47.864311	2026-10-15 10:20:47.864311
36	66	Country 7779789049	t	2026-10-14 12:55:27.834269+05:30	2026-10-14 12:55:27.797518	2026-10-14 12:55:27.834269
35	65	Country 777978904	t	2026-10-14 12:55:27.838177+05:30	2026-10-14 12:55:27.782875	2026-10-14 12:55:27.838177
38	68	Country 161897019	t	2026-10-14 15:27:40.358848+05:30	2026-10-14 15:27:40.273809	2026-10-14 15:27:40.358848
37	67	Country 16189701	t	2026-10-14 15:27:40.364658+05:30	2026-10-14 15:27:40.22738	2026-10-14 15:27:40.364658
39	69	Country 40066304	t	2026-10-14 15:27:40.468135+05:30	2026-10-14 15:27:40.412223	2026-10-14 15:27:40.468135
91	131	Country 7870395269	t	2026-10-15 10:20:47.925573+05:30	2026-10-15 10:20:47.88457	2026-10-15 10:20:47.925573
90	130	Country 787039526	t	2026-10-15 10:20:47.929954+05:30	2026-10-15 10:20:47.873947	2026-10-15 10:20:47.929954
41	71	Country 471232059	t	2026-10-14 15:27:40.540959+05:30	2026-10-14 15:27:40.486506	2026-10-14 15:27:40.540959
40	70	Country 47123205	t	2026-10-14 15:27:40.546321+05:30	2026-10-14 15:27:40.475257	2026-10-14 15:27:40.546321
43	73	Country 1932926019	t	2026-10-14 15:28:02.063905+05:30	2026-10-14 15:28:01.995818	2026-10-14 15:28:02.063905
42	72	Country 193292601	t	2026-10-14 15:28:02.069032+05:30	2026-10-14 15:28:01.961017	2026-10-14 15:28:02.069032
93	137	Country B801317832	t	2026-10-15 10:20:48.484701+05:30	2026-10-15 10:20:48.183804	2026-10-15 10:20:48.484701
44	74	Country 209017604	t	2026-10-14 15:28:02.138031+05:30	2026-10-14 15:28:02.094607	2026-10-14 15:28:02.138031
92	136	Country A801317832	t	2026-10-15 10:20:48.489215+05:30	2026-10-15 10:20:48.015727	2026-10-15 10:20:48.489215
46	76	Country 2140857059	t	2026-10-14 15:28:02.186951+05:30	2026-10-14 15:28:02.155764	2026-10-14 15:28:02.186951
45	75	Country 214085705	t	2026-10-14 15:28:02.190822+05:30	2026-10-14 15:28:02.14451	2026-10-14 15:28:02.190822
48	78	Country 1589422019	t	2026-10-14 15:41:11.776727+05:30	2026-10-14 15:41:11.688922	2026-10-14 15:41:11.776727
47	77	Country 158942201	t	2026-10-14 15:41:11.782391+05:30	2026-10-14 15:41:11.647378	2026-10-14 15:41:11.782391
96	142	Germany	t	2026-10-15 10:20:48.497521+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.497521
94	142	India	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.503663
50	80	Country 1785631029	t	2026-10-14 15:41:11.836988+05:30	2026-10-14 15:41:11.80643	2026-10-14 15:41:11.836988
49	79	Country 178563102	t	2026-10-14 15:41:11.841544+05:30	2026-10-14 15:41:11.789992	2026-10-14 15:41:11.841544
95	142	Japan	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.503663
51	1	Germany	t	2026-10-15 10:23:05.109394+05:30	2026-10-14 15:45:30.210693	2026-10-15 10:23:05.109394
54	82	Country 5730982019	t	2026-10-14 15:51:05.967252+05:30	2026-10-14 15:51:05.836746	2026-10-14 15:51:05.967252
53	81	Country 573098201	t	2026-10-14 15:51:05.994727+05:30	2026-10-14 15:51:05.797094	2026-10-14 15:51:05.994727
52	1	United Arab Emirates	t	2026-10-15 10:23:05.109394+05:30	2026-10-14 15:45:37.27557	2026-10-15 10:23:05.109394
97	1	British Indian Ocean Territory	f	\N	2026-10-15 10:23:05.109394	2026-10-15 10:23:05.109394
56	84	Country 6003831029	t	2026-10-14 15:51:06.109829+05:30	2026-10-14 15:51:06.036305	2026-10-14 15:51:06.109829
55	83	Country 600383102	t	2026-10-14 15:51:06.118597+05:30	2026-10-14 15:51:06.010492	2026-10-14 15:51:06.118597
98	1	India	f	\N	2026-10-15 10:23:05.109394	2026-10-15 10:23:05.109394
99	1	United Kingdom	f	\N	2026-10-15 10:23:05.109394	2026-10-15 10:23:05.109394
58	91	Country 8387817179	t	2026-10-14 16:09:28.480704+05:30	2026-10-14 16:09:28.425783	2026-10-14 16:09:28.480704
57	90	Country 838781717	t	2026-10-14 16:09:28.485149+05:30	2026-10-14 16:09:28.393285	2026-10-14 16:09:28.485149
100	143	Owned 867567401	f	\N	2026-10-15 10:28:58.687715	2026-10-15 10:28:58.687715
59	92	Country 851444423	t	2026-10-14 16:09:28.559173+05:30	2026-10-14 16:09:28.519546	2026-10-14 16:09:28.559173
61	93	Owned 856238325	f	\N	2026-10-14 16:09:28.564054	2026-10-14 16:09:28.564054
102	1	Afghanistan	f	\N	2026-10-15 10:41:59.80954	2026-10-15 10:41:59.80954
103	1	Åland Islands	f	\N	2026-10-15 10:41:59.860731	2026-10-15 10:41:59.860731
104	1	Albania	f	\N	2026-10-15 10:41:59.869349	2026-10-15 10:41:59.869349
64	95	Country 8567715269	t	2026-10-14 16:09:28.617327+05:30	2026-10-14 16:09:28.586364	2026-10-14 16:09:28.617327
63	94	Country 856771526	t	2026-10-14 16:09:28.621062+05:30	2026-10-14 16:09:28.573963	2026-10-14 16:09:28.621062
105	1	Algeria	f	\N	2026-10-15 10:41:59.876183	2026-10-15 10:41:59.876183
106	49	British Indian Ocean Territory	f	\N	2026-10-15 10:42:15.286589	2026-10-15 10:42:15.286589
66	101	Country B868524232	t	2026-10-14 16:09:29.164212+05:30	2026-10-14 16:09:28.698501	2026-10-14 16:09:29.164212
65	100	Country A868524232	t	2026-10-14 16:09:29.169405+05:30	2026-10-14 16:09:28.688254	2026-10-14 16:09:29.169405
107	49	India	f	\N	2026-10-15 10:42:15.319909	2026-10-15 10:42:15.319909
108	1	American Samoa	f	\N	2026-10-15 10:54:22.733403	2026-10-15 10:54:22.733403
109	1	Andorra	f	\N	2026-10-15 10:54:22.769545	2026-10-15 10:54:22.769545
110	1	Angola	f	\N	2026-10-15 10:54:22.775455	2026-10-15 10:54:22.775455
111	1	Anguilla	f	\N	2026-10-15 10:54:22.780559	2026-10-15 10:54:22.780559
113	144	Owned 182355002	f	\N	2026-10-15 11:08:21.871174	2026-10-15 11:08:21.871174
116	145	Owned 785454202	f	\N	2026-10-15 11:18:37.862156	2026-10-15 11:18:37.862156
118	146	Owned 36042801	f	\N	2026-10-15 11:18:50.367451	2026-10-15 11:18:50.367451
120	49	Åland Islands	f	\N	2026-10-15 11:23:40.679177	2026-10-15 11:23:40.679177
121	49	Albania	f	\N	2026-10-15 11:23:40.847708	2026-10-15 11:23:40.847708
122	49	Algeria	f	\N	2026-10-15 11:23:40.857066	2026-10-15 11:23:40.857066
123	22	Belize	f	\N	2026-10-15 11:27:20.367852	2026-10-15 11:27:20.367852
124	22	Benin	f	\N	2026-10-15 11:27:20.531948	2026-10-15 11:27:20.531948
125	22	Bermuda	f	\N	2026-10-15 11:27:20.540664	2026-10-15 11:27:20.540664
126	22	Bhutan	f	\N	2026-10-15 11:27:20.557373	2026-10-15 11:27:20.557373
127	49	Afghanistan	f	\N	2026-10-15 11:31:13.089964	2026-10-15 11:31:13.089964
128	49	American Samoa	f	\N	2026-10-15 11:31:13.265098	2026-10-15 11:31:13.265098
129	49	Andorra	f	\N	2026-10-15 11:31:13.297312	2026-10-15 11:31:13.297312
130	23	Pakistan	f	\N	2026-10-15 11:43:34.739133	2026-10-15 11:43:34.739133
131	22	Afghanistan	f	\N	2026-10-15 11:57:26.829195	2026-10-15 11:57:26.829195
132	22	Åland Islands	f	\N	2026-10-15 11:57:27.040732	2026-10-15 11:57:27.040732
133	22	Albania	f	\N	2026-10-15 11:57:27.057892	2026-10-15 11:57:27.057892
134	22	Algeria	f	\N	2026-10-15 11:57:27.081223	2026-10-15 11:57:27.081223
\.


--
-- TOC entry 5676 (class 0 OID 117105)
-- Dependencies: 256
-- Data for Name: user_filters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_filters (id, user_id, filter_name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 5688 (class 0 OID 117302)
-- Dependencies: 268
-- Data for Name: user_line; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_line (id, user_id, user_shop_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
4	26	4	Line 13336822	t	2026-10-14 11:57:30.213477+05:30	2026-10-14 11:57:30.16516	2026-10-14 11:57:30.213477
5	31	5	Line 73394922	t	2026-10-14 12:07:00.813734+05:30	2026-10-14 12:07:00.76792	2026-10-14 12:07:00.813734
7	34	7	Line B83246528	t	2026-10-14 12:07:01.142842+05:30	2026-10-14 12:07:00.856794	2026-10-14 12:07:01.142842
6	33	6	Line A83246528	t	2026-10-14 12:07:01.147354+05:30	2026-10-14 12:07:00.837634	2026-10-14 12:07:01.147354
8	42	8	Line 847312702	t	2026-10-14 12:07:08.557272+05:30	2026-10-14 12:07:08.502957	2026-10-14 12:07:08.557272
10	44	10	Line B857087506	t	2026-10-14 12:07:08.860528+05:30	2026-10-14 12:07:08.589495	2026-10-14 12:07:08.860528
9	43	9	Line A857087506	t	2026-10-14 12:07:08.864632+05:30	2026-10-14 12:07:08.576384	2026-10-14 12:07:08.864632
11	54	11	Line 899315822	t	2026-10-14 12:41:49.089514+05:30	2026-10-14 12:41:49.033243	2026-10-14 12:41:49.089514
13	57	13	Line B911479328	t	2026-10-14 12:41:49.444053+05:30	2026-10-14 12:41:49.1437	2026-10-14 12:41:49.444053
12	56	12	Line A911479328	t	2026-10-14 12:41:49.449392+05:30	2026-10-14 12:41:49.123285	2026-10-14 12:41:49.449392
15	63	15	Line 3279490019	t	2026-10-14 12:55:13.591199+05:30	2026-10-14 12:55:13.49745	2026-10-14 12:55:13.591199
14	62	14	Line 327949001	t	2026-10-14 12:55:13.599261+05:30	2026-10-14 12:55:13.413389	2026-10-14 12:55:13.599261
16	64	16	Line 769633603	t	2026-10-14 12:55:27.77361+05:30	2026-10-14 12:55:27.728893	2026-10-14 12:55:27.77361
18	66	18	Line 7779789049	t	2026-10-14 12:55:27.832642+05:30	2026-10-14 12:55:27.799602	2026-10-14 12:55:27.832642
17	65	17	Line 777978904	t	2026-10-14 12:55:27.836375+05:30	2026-10-14 12:55:27.785008	2026-10-14 12:55:27.836375
20	68	20	Line 161897019	t	2026-10-14 15:27:40.35526+05:30	2026-10-14 15:27:40.276836	2026-10-14 15:27:40.35526
19	67	19	Line 16189701	t	2026-10-14 15:27:40.36238+05:30	2026-10-14 15:27:40.264574	2026-10-14 15:27:40.36238
21	69	21	Line 40066304	t	2026-10-14 15:27:40.465876+05:30	2026-10-14 15:27:40.415386	2026-10-14 15:27:40.465876
23	71	23	Line 471232059	t	2026-10-14 15:27:40.538232+05:30	2026-10-14 15:27:40.489589	2026-10-14 15:27:40.538232
22	70	22	Line 47123205	t	2026-10-14 15:27:40.544328+05:30	2026-10-14 15:27:40.477886	2026-10-14 15:27:40.544328
25	73	25	Line 1932926019	t	2026-10-14 15:28:02.060493+05:30	2026-10-14 15:28:01.998139	2026-10-14 15:28:02.060493
24	72	24	Line 193292601	t	2026-10-14 15:28:02.067086+05:30	2026-10-14 15:28:01.988862	2026-10-14 15:28:02.067086
26	74	26	Line 209017604	t	2026-10-14 15:28:02.135668+05:30	2026-10-14 15:28:02.097064	2026-10-14 15:28:02.135668
28	76	28	Line 2140857059	t	2026-10-14 15:28:02.184832+05:30	2026-10-14 15:28:02.158221	2026-10-14 15:28:02.184832
27	75	27	Line 214085705	t	2026-10-14 15:28:02.18921+05:30	2026-10-14 15:28:02.147254	2026-10-14 15:28:02.18921
30	78	30	Line 1589422019	t	2026-10-14 15:41:11.773553+05:30	2026-10-14 15:41:11.691812	2026-10-14 15:41:11.773553
29	77	29	Line 158942201	t	2026-10-14 15:41:11.780646+05:30	2026-10-14 15:41:11.680309	2026-10-14 15:41:11.780646
32	80	32	Line 1785631029	t	2026-10-14 15:41:11.835023+05:30	2026-10-14 15:41:11.808897	2026-10-14 15:41:11.835023
31	79	31	Line 178563102	t	2026-10-14 15:41:11.839214+05:30	2026-10-14 15:41:11.792949	2026-10-14 15:41:11.839214
34	82	34	Line 5730982019	t	2026-10-14 15:51:05.960402+05:30	2026-10-14 15:51:05.841454	2026-10-14 15:51:05.960402
33	81	33	Line 573098201	t	2026-10-14 15:51:05.991341+05:30	2026-10-14 15:51:05.824308	2026-10-14 15:51:05.991341
36	84	36	Line 6003831029	t	2026-10-14 15:51:06.104415+05:30	2026-10-14 15:51:06.041096	2026-10-14 15:51:06.104415
35	83	35	Line 600383102	t	2026-10-14 15:51:06.114052+05:30	2026-10-14 15:51:06.016508	2026-10-14 15:51:06.114052
38	91	38	Line 8387817179	t	2026-10-14 16:09:28.477886+05:30	2026-10-14 16:09:28.428356	2026-10-14 16:09:28.477886
37	90	37	Line 838781717	t	2026-10-14 16:09:28.483466+05:30	2026-10-14 16:09:28.417503	2026-10-14 16:09:28.483466
39	92	39	Line 851444423	t	2026-10-14 16:09:28.557298+05:30	2026-10-14 16:09:28.522048	2026-10-14 16:09:28.557298
41	95	41	Line 8567715269	t	2026-10-14 16:09:28.615095+05:30	2026-10-14 16:09:28.58866	2026-10-14 16:09:28.615095
40	94	40	Line 856771526	t	2026-10-14 16:09:28.619501+05:30	2026-10-14 16:09:28.577946	2026-10-14 16:09:28.619501
43	101	43	Line B868524232	t	2026-10-14 16:09:29.162097+05:30	2026-10-14 16:09:28.866912	2026-10-14 16:09:29.162097
42	100	42	Line A868524232	t	2026-10-14 16:09:29.167411+05:30	2026-10-14 16:09:28.690459	2026-10-14 16:09:29.167411
45	107	45	Line 5216632019	t	2026-10-14 16:26:15.706799+05:30	2026-10-14 16:26:15.588875	2026-10-14 16:26:15.706799
44	106	44	Line 521663201	t	2026-10-14 16:26:15.714721+05:30	2026-10-14 16:26:15.566862	2026-10-14 16:26:15.714721
46	108	46	Line 572189502	t	2026-10-14 16:26:15.770331+05:30	2026-10-14 16:26:15.728254	2026-10-14 16:26:15.770331
48	110	48	Line 5775058039	t	2026-10-14 16:26:15.814546+05:30	2026-10-14 16:26:15.792618	2026-10-14 16:26:15.814546
47	109	47	Line 577505803	t	2026-10-14 16:26:15.818194+05:30	2026-10-14 16:26:15.781745	2026-10-14 16:26:15.818194
50	127	50	Line 7640435179	t	2026-10-15 10:20:47.779502+05:30	2026-10-15 10:20:47.691061	2026-10-15 10:20:47.779502
49	126	49	Line 764043517	t	2026-10-15 10:20:47.786869+05:30	2026-10-15 10:20:47.679431	2026-10-15 10:20:47.786869
51	128	51	Line 781488023	t	2026-10-15 10:20:47.85634+05:30	2026-10-15 10:20:47.820356	2026-10-15 10:20:47.85634
53	131	53	Line 7870395269	t	2026-10-15 10:20:47.92344+05:30	2026-10-15 10:20:47.887264	2026-10-15 10:20:47.92344
52	130	52	Line 787039526	t	2026-10-15 10:20:47.92804+05:30	2026-10-15 10:20:47.876503	2026-10-15 10:20:47.92804
55	137	55	Line B801317832	t	2026-10-15 10:20:48.482627+05:30	2026-10-15 10:20:48.185836	2026-10-15 10:20:48.482627
54	136	54	Line A801317832	t	2026-10-15 10:20:48.487539+05:30	2026-10-15 10:20:48.020434	2026-10-15 10:20:48.487539
56	49	56	Line - 1	f	\N	2026-10-15 11:25:20.104981	2026-10-15 11:25:20.104981
57	49	56	Line - 2	f	\N	2026-10-15 11:25:27.721392	2026-10-15 11:25:27.721392
\.


--
-- TOC entry 5684 (class 0 OID 117242)
-- Dependencies: 264
-- Data for Name: user_location; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_location (id, user_id, name, is_deleted, deleted_at, created_at, updated_at, branch_id) FROM stdin;
4	26	Location 13336822	t	2026-10-14 11:57:30.214557+05:30	2026-10-14 11:57:30.156133	2026-10-14 11:57:30.214557	10
5	31	Location 73394922	t	2026-10-14 12:07:00.81514+05:30	2026-10-14 12:07:00.75856	2026-10-14 12:07:00.81514	11
7	34	Location B83246528	t	2026-10-14 12:07:01.143376+05:30	2026-10-14 12:07:00.855397	2026-10-14 12:07:01.143376	13
6	33	Location A83246528	t	2026-10-14 12:07:01.147898+05:30	2026-10-14 12:07:00.836798	2026-10-14 12:07:01.147898	12
8	42	Location 847312702	t	2026-10-14 12:07:08.558371+05:30	2026-10-14 12:07:08.495993	2026-10-14 12:07:08.558371	14
10	44	Location B857087506	t	2026-10-14 12:07:08.861053+05:30	2026-10-14 12:07:08.588592	2026-10-14 12:07:08.861053	16
9	43	Location A857087506	t	2026-10-14 12:07:08.865141+05:30	2026-10-14 12:07:08.575666	2026-10-14 12:07:08.865141	15
11	54	Location 899315822	t	2026-10-14 12:41:49.09181+05:30	2026-10-14 12:41:49.022742	2026-10-14 12:41:49.09181	17
13	57	Location B911479328	t	2026-10-14 12:41:49.444822+05:30	2026-10-14 12:41:49.142857	2026-10-14 12:41:49.444822	19
12	56	Location A911479328	t	2026-10-14 12:41:49.450129+05:30	2026-10-14 12:41:49.121953	2026-10-14 12:41:49.450129	18
15	63	Location 3279490019	t	2026-10-14 12:55:13.59227+05:30	2026-10-14 12:55:13.496374	2026-10-14 12:55:13.59227	21
14	62	Location 327949001	t	2026-10-14 12:55:13.600104+05:30	2026-10-14 12:55:13.395436	2026-10-14 12:55:13.600104	20
16	64	Location 769633603	t	2026-10-14 12:55:27.774638+05:30	2026-10-14 12:55:27.722281	2026-10-14 12:55:27.774638	22
18	66	Location 7779789049	t	2026-10-14 12:55:27.833257+05:30	2026-10-14 12:55:27.798908	2026-10-14 12:55:27.833257	24
17	65	Location 777978904	t	2026-10-14 12:55:27.837032+05:30	2026-10-14 12:55:27.784294	2026-10-14 12:55:27.837032	23
20	68	Location 161897019	t	2026-10-14 15:27:40.356458+05:30	2026-10-14 15:27:40.275756	2026-10-14 15:27:40.356458	26
19	67	Location 16189701	t	2026-10-14 15:27:40.36297+05:30	2026-10-14 15:27:40.253822	2026-10-14 15:27:40.36297	25
21	69	Location 40066304	t	2026-10-14 15:27:40.466856+05:30	2026-10-14 15:27:40.414451	2026-10-14 15:27:40.466856	27
23	71	Location 471232059	t	2026-10-14 15:27:40.539299+05:30	2026-10-14 15:27:40.488347	2026-10-14 15:27:40.539299	29
22	70	Location 47123205	t	2026-10-14 15:27:40.544979+05:30	2026-10-14 15:27:40.476936	2026-10-14 15:27:40.544979	28
25	73	Location 1932926019	t	2026-10-14 15:28:02.061905+05:30	2026-10-14 15:28:01.997285	2026-10-14 15:28:02.061905	31
24	72	Location 193292601	t	2026-10-14 15:28:02.06776+05:30	2026-10-14 15:28:01.98079	2026-10-14 15:28:02.06776	30
26	74	Location 209017604	t	2026-10-14 15:28:02.136625+05:30	2026-10-14 15:28:02.096193	2026-10-14 15:28:02.136625	32
28	76	Location 2140857059	t	2026-10-14 15:28:02.185595+05:30	2026-10-14 15:28:02.157352	2026-10-14 15:28:02.185595	34
27	75	Location 214085705	t	2026-10-14 15:28:02.189763+05:30	2026-10-14 15:28:02.146357	2026-10-14 15:28:02.189763	33
30	78	Location 1589422019	t	2026-10-14 15:41:11.774819+05:30	2026-10-14 15:41:11.690688	2026-10-14 15:41:11.774819	36
29	77	Location 158942201	t	2026-10-14 15:41:11.781262+05:30	2026-10-14 15:41:11.669447	2026-10-14 15:41:11.781262	35
32	80	Location 1785631029	t	2026-10-14 15:41:11.835681+05:30	2026-10-14 15:41:11.808035	2026-10-14 15:41:11.835681	38
31	79	Location 178563102	t	2026-10-14 15:41:11.839829+05:30	2026-10-14 15:41:11.792026	2026-10-14 15:41:11.839829	37
34	82	Location 5730982019	t	2026-10-14 15:51:05.963536+05:30	2026-10-14 15:51:05.840109	2026-10-14 15:51:05.963536	40
33	81	Location 573098201	t	2026-10-14 15:51:05.992351+05:30	2026-10-14 15:51:05.816405	2026-10-14 15:51:05.992351	39
36	84	Location 6003831029	t	2026-10-14 15:51:06.107929+05:30	2026-10-14 15:51:06.040064	2026-10-14 15:51:06.107929	42
35	83	Location 600383102	t	2026-10-14 15:51:06.114757+05:30	2026-10-14 15:51:06.015091	2026-10-14 15:51:06.114757	41
38	91	Location 8387817179	t	2026-10-14 16:09:28.47887+05:30	2026-10-14 16:09:28.427495	2026-10-14 16:09:28.47887	44
37	90	Location 838781717	t	2026-10-14 16:09:28.484003+05:30	2026-10-14 16:09:28.409286	2026-10-14 16:09:28.484003	43
39	92	Location 851444423	t	2026-10-14 16:09:28.557963+05:30	2026-10-14 16:09:28.521109	2026-10-14 16:09:28.557963	45
41	95	Location 8567715269	t	2026-10-14 16:09:28.615862+05:30	2026-10-14 16:09:28.587908	2026-10-14 16:09:28.615862	47
40	94	Location 856771526	t	2026-10-14 16:09:28.620018+05:30	2026-10-14 16:09:28.577034	2026-10-14 16:09:28.620018	46
43	101	Location B868524232	t	2026-10-14 16:09:29.162733+05:30	2026-10-14 16:09:28.862589	2026-10-14 16:09:29.162733	49
42	100	Location A868524232	t	2026-10-14 16:09:29.168028+05:30	2026-10-14 16:09:28.689636	2026-10-14 16:09:29.168028	48
45	107	Location 5216632019	t	2026-10-14 16:26:15.708078+05:30	2026-10-14 16:26:15.58744	2026-10-14 16:26:15.708078	51
44	106	Location 521663201	t	2026-10-14 16:26:15.71548+05:30	2026-10-14 16:26:15.545253	2026-10-14 16:26:15.71548	50
46	108	Location 572189502	t	2026-10-14 16:26:15.771088+05:30	2026-10-14 16:26:15.727402	2026-10-14 16:26:15.771088	52
48	110	Location 5775058039	t	2026-10-14 16:26:15.815247+05:30	2026-10-14 16:26:15.791592	2026-10-14 16:26:15.815247	54
47	109	Location 577505803	t	2026-10-14 16:26:15.818685+05:30	2026-10-14 16:26:15.780867	2026-10-14 16:26:15.818685	53
50	127	Location 7640435179	t	2026-10-15 10:20:47.780578+05:30	2026-10-15 10:20:47.690182	2026-10-15 10:20:47.780578	56
49	126	Location 764043517	t	2026-10-15 10:20:47.787671+05:30	2026-10-15 10:20:47.670032	2026-10-15 10:20:47.787671	55
51	128	Location 781488023	t	2026-10-15 10:20:47.857918+05:30	2026-10-15 10:20:47.819623	2026-10-15 10:20:47.857918	57
53	131	Location 7870395269	t	2026-10-15 10:20:47.924157+05:30	2026-10-15 10:20:47.886297	2026-10-15 10:20:47.924157	59
52	130	Location 787039526	t	2026-10-15 10:20:47.928687+05:30	2026-10-15 10:20:47.875624	2026-10-15 10:20:47.928687	58
55	137	Location B801317832	t	2026-10-15 10:20:48.483354+05:30	2026-10-15 10:20:48.185137	2026-10-15 10:20:48.483354	61
54	136	Location A801317832	t	2026-10-15 10:20:48.488084+05:30	2026-10-15 10:20:48.019311	2026-10-15 10:20:48.488084	60
56	49	768687687678	f	\N	2026-10-15 11:24:40.032192	2026-10-15 11:24:40.032192	64
57	49	3232342	f	\N	2026-10-15 11:24:49.42587	2026-10-15 11:24:49.42587	66
\.


--
-- TOC entry 5696 (class 0 OID 117680)
-- Dependencies: 276
-- Data for Name: user_machine_access; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_machine_access (id, user_id, machine_id, created_at, updated_at) FROM stdin;
2	37	21	2026-10-14 12:07:01.062939	2026-10-14 12:07:01.062939
4	38	22	2026-10-14 12:07:01.136321	2026-10-14 12:07:01.136321
6	47	26	2026-10-14 12:07:08.783353	2026-10-14 12:07:08.783353
8	48	27	2026-10-14 12:07:08.853665	2026-10-14 12:07:08.853665
10	60	32	2026-10-14 12:41:49.363436	2026-10-14 12:41:49.363436
12	61	33	2026-10-14 12:41:49.435932	2026-10-14 12:41:49.435932
14	104	71	2026-10-14 16:09:29.078007	2026-10-14 16:09:29.078007
16	105	72	2026-10-14 16:09:29.152617	2026-10-14 16:09:29.152617
18	140	87	2026-10-15 10:20:48.399613	2026-10-15 10:20:48.399613
20	141	88	2026-10-15 10:20:48.475144	2026-10-15 10:20:48.475144
\.


--
-- TOC entry 5656 (class 0 OID 116710)
-- Dependencies: 236
-- Data for Name: user_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_machines (id, user_id, machine_id, assigned_at, is_deleted, deleted_at) FROM stdin;
\.


--
-- TOC entry 5652 (class 0 OID 116663)
-- Dependencies: 232
-- Data for Name: user_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_profiles (id, user_id, designation, grade, birth_date, joining_date, blood_group, profile_picture, contact_person, contact_phone, theme, timezone, language, start_day_time, is_deleted, deleted_at, created_at, updated_at, filters) FROM stdin;
2	23	\N	\N	\N	\N	\N	\N	Kolby Skiles	+816596645193	dark_version	America/Iqaluit	en	17:44:00	f	\N	2026-10-14 10:34:33.413564	2026-10-14 10:34:33.413564	{country,state,branch,city,location,shop,line,supervisor}
45	104	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 16:09:29.078007	2026-10-14 16:09:29.078007	{}
46	105	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 16:09:29.152617	2026-10-14 16:09:29.152617	{}
44	103	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	t	2026-10-14 16:09:29.157727+05:30	2026-10-14 16:09:29.004149	2026-10-14 16:09:29.157727	{}
16	35	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:00.916939	2026-10-14 12:07:00.916939	{}
18	37	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:01.062939	2026-10-14 12:07:01.062939	{}
19	38	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:01.136321	2026-10-14 12:07:01.136321	{}
17	36	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	t	2026-10-14 12:07:01.13979+05:30	2026-10-14 12:07:00.988792	2026-10-14 12:07:01.13979	{}
1	22	\N	\N	\N	\N	\N	\N	Edison Stiedemann	+441554577026	dark_version	Europe/Chisinau	en	18:43:00	f	\N	2026-10-14 10:26:06.7819	2026-10-15 09:56:20.930348	{country,state,branch,city,location,shop,line,cell,supervisor}
57	138	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-15 10:20:48.24788	2026-10-15 10:20:48.24788	{}
22	45	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:08.64863	2026-10-14 12:07:08.64863	{}
24	47	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:08.783353	2026-10-14 12:07:08.783353	{}
25	48	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:07:08.853665	2026-10-14 12:07:08.853665	{}
23	46	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	t	2026-10-14 12:07:08.857401+05:30	2026-10-14 12:07:08.715634	2026-10-14 12:07:08.857401	{}
59	140	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-15 10:20:48.399613	2026-10-15 10:20:48.399613	{}
60	141	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-15 10:20:48.475144	2026-10-15 10:20:48.475144	{}
28	49	\N	\N	\N	\N	\N	\N	Carol Schultz	+17137981184	light_version	Pacific/Wallis	en	09:25:00	f	\N	2026-10-14 12:14:25.190123	2026-10-14 12:14:25.190123	{}
30	50	\N	\N	\N	\N	\N	\N	Bernadette Runolfsdottir	+9716700073539	dark_version	Asia/Kuala_Lumpur	en	15:46:00	f	\N	2026-10-14 12:30:26.127054	2026-10-14 12:30:26.127054	{}
31	58	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:41:49.2054	2026-10-14 12:41:49.2054	{}
33	60	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:41:49.363436	2026-10-14 12:41:49.363436	{}
58	139	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	t	2026-10-15 10:20:48.479162+05:30	2026-10-15 10:20:48.322469	2026-10-15 10:20:48.479162	{}
34	61	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 12:41:49.435932	2026-10-14 12:41:49.435932	{}
32	59	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	t	2026-10-14 12:41:49.440609+05:30	2026-10-14 12:41:49.288914	2026-10-14 12:41:49.440609	{}
6	1	\N	\N	\N	\N	\N	\N	\N	\N	dark_version	Asia/Kolkata	en	\N	f	\N	2026-10-14 11:20:33.902672	2026-10-15 11:55:13.932938	{}
43	102	\N	\N	\N	\N	\N	\N	Member	\N	dark_version	Asia/Kolkata	en	06:00:00	f	\N	2026-10-14 16:09:28.932166	2026-10-14 16:09:28.932166	{}
\.


--
-- TOC entry 5686 (class 0 OID 117272)
-- Dependencies: 266
-- Data for Name: user_shop; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_shop (id, user_id, user_location_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
4	26	4	Shop 13336822	t	2026-10-14 11:57:30.214068+05:30	2026-10-14 11:57:30.16003	2026-10-14 11:57:30.214068
5	31	5	Shop 73394922	t	2026-10-14 12:07:00.814576+05:30	2026-10-14 12:07:00.762931	2026-10-14 12:07:00.814576
7	34	7	Shop B83246528	t	2026-10-14 12:07:01.143116+05:30	2026-10-14 12:07:00.855933	2026-10-14 12:07:01.143116
6	33	6	Shop A83246528	t	2026-10-14 12:07:01.147629+05:30	2026-10-14 12:07:00.837183	2026-10-14 12:07:01.147629
8	42	8	Shop 847312702	t	2026-10-14 12:07:08.557799+05:30	2026-10-14 12:07:08.498943	2026-10-14 12:07:08.557799
10	44	10	Shop B857087506	t	2026-10-14 12:07:08.860792+05:30	2026-10-14 12:07:08.589062	2026-10-14 12:07:08.860792
9	43	9	Shop A857087506	t	2026-10-14 12:07:08.864888+05:30	2026-10-14 12:07:08.576019	2026-10-14 12:07:08.864888
11	54	11	Shop 899315822	t	2026-10-14 12:41:49.090827+05:30	2026-10-14 12:41:49.028167	2026-10-14 12:41:49.090827
13	57	13	Shop B911479328	t	2026-10-14 12:41:49.444334+05:30	2026-10-14 12:41:49.143277	2026-10-14 12:41:49.444334
12	56	12	Shop A911479328	t	2026-10-14 12:41:49.449752+05:30	2026-10-14 12:41:49.12256	2026-10-14 12:41:49.449752
15	63	15	Shop 3279490019	t	2026-10-14 12:55:13.591765+05:30	2026-10-14 12:55:13.496924	2026-10-14 12:55:13.591765
14	62	14	Shop 327949001	t	2026-10-14 12:55:13.599666+05:30	2026-10-14 12:55:13.402625	2026-10-14 12:55:13.599666
16	64	16	Shop 769633603	t	2026-10-14 12:55:27.774112+05:30	2026-10-14 12:55:27.725367	2026-10-14 12:55:27.774112
18	66	18	Shop 7779789049	t	2026-10-14 12:55:27.832979+05:30	2026-10-14 12:55:27.799249	2026-10-14 12:55:27.832979
17	65	17	Shop 777978904	t	2026-10-14 12:55:27.836702+05:30	2026-10-14 12:55:27.784649	2026-10-14 12:55:27.836702
20	68	20	Shop 161897019	t	2026-10-14 15:27:40.35586+05:30	2026-10-14 15:27:40.276337	2026-10-14 15:27:40.35586
19	67	19	Shop 16189701	t	2026-10-14 15:27:40.362684+05:30	2026-10-14 15:27:40.258839	2026-10-14 15:27:40.362684
21	69	21	Shop 40066304	t	2026-10-14 15:27:40.466351+05:30	2026-10-14 15:27:40.414913	2026-10-14 15:27:40.466351
23	71	23	Shop 471232059	t	2026-10-14 15:27:40.538815+05:30	2026-10-14 15:27:40.488986	2026-10-14 15:27:40.538815
22	70	22	Shop 47123205	t	2026-10-14 15:27:40.544629+05:30	2026-10-14 15:27:40.47745	2026-10-14 15:27:40.544629
25	73	25	Shop 1932926019	t	2026-10-14 15:28:02.061262+05:30	2026-10-14 15:28:01.997649	2026-10-14 15:28:02.061262
24	72	24	Shop 193292601	t	2026-10-14 15:28:02.06744+05:30	2026-10-14 15:28:01.984445	2026-10-14 15:28:02.06744
26	74	26	Shop 209017604	t	2026-10-14 15:28:02.136244+05:30	2026-10-14 15:28:02.096571	2026-10-14 15:28:02.136244
28	76	28	Shop 2140857059	t	2026-10-14 15:28:02.185211+05:30	2026-10-14 15:28:02.157752	2026-10-14 15:28:02.185211
27	75	27	Shop 214085705	t	2026-10-14 15:28:02.189496+05:30	2026-10-14 15:28:02.146872	2026-10-14 15:28:02.189496
30	78	30	Shop 1589422019	t	2026-10-14 15:41:11.774149+05:30	2026-10-14 15:41:11.691236	2026-10-14 15:41:11.774149
29	77	29	Shop 158942201	t	2026-10-14 15:41:11.780967+05:30	2026-10-14 15:41:11.674251	2026-10-14 15:41:11.780967
32	80	32	Shop 1785631029	t	2026-10-14 15:41:11.835367+05:30	2026-10-14 15:41:11.808481	2026-10-14 15:41:11.835367
31	79	31	Shop 178563102	t	2026-10-14 15:41:11.839513+05:30	2026-10-14 15:41:11.792511	2026-10-14 15:41:11.839513
34	82	34	Shop 5730982019	t	2026-10-14 15:51:05.962423+05:30	2026-10-14 15:51:05.840698	2026-10-14 15:51:05.962423
33	81	33	Shop 573098201	t	2026-10-14 15:51:05.99188+05:30	2026-10-14 15:51:05.819206	2026-10-14 15:51:05.99188
36	84	36	Shop 6003831029	t	2026-10-14 15:51:06.107217+05:30	2026-10-14 15:51:06.040649	2026-10-14 15:51:06.107217
35	83	35	Shop 600383102	t	2026-10-14 15:51:06.114416+05:30	2026-10-14 15:51:06.015837	2026-10-14 15:51:06.114416
38	91	38	Shop 8387817179	t	2026-10-14 16:09:28.478419+05:30	2026-10-14 16:09:28.427915	2026-10-14 16:09:28.478419
37	90	37	Shop 838781717	t	2026-10-14 16:09:28.483731+05:30	2026-10-14 16:09:28.411991	2026-10-14 16:09:28.483731
39	92	39	Shop 851444423	t	2026-10-14 16:09:28.557642+05:30	2026-10-14 16:09:28.521584	2026-10-14 16:09:28.557642
41	95	41	Shop 8567715269	t	2026-10-14 16:09:28.615473+05:30	2026-10-14 16:09:28.588255	2026-10-14 16:09:28.615473
40	94	40	Shop 856771526	t	2026-10-14 16:09:28.619768+05:30	2026-10-14 16:09:28.577508	2026-10-14 16:09:28.619768
43	101	43	Shop B868524232	t	2026-10-14 16:09:29.162423+05:30	2026-10-14 16:09:28.864376	2026-10-14 16:09:29.162423
42	100	42	Shop A868524232	t	2026-10-14 16:09:29.167693+05:30	2026-10-14 16:09:28.690073	2026-10-14 16:09:29.167693
45	107	45	Shop 5216632019	t	2026-10-14 16:26:15.707435+05:30	2026-10-14 16:26:15.588346	2026-10-14 16:26:15.707435
44	106	44	Shop 521663201	t	2026-10-14 16:26:15.715136+05:30	2026-10-14 16:26:15.556835	2026-10-14 16:26:15.715136
46	108	46	Shop 572189502	t	2026-10-14 16:26:15.770709+05:30	2026-10-14 16:26:15.727813	2026-10-14 16:26:15.770709
48	110	48	Shop 5775058039	t	2026-10-14 16:26:15.814887+05:30	2026-10-14 16:26:15.792067	2026-10-14 16:26:15.814887
47	109	47	Shop 577505803	t	2026-10-14 16:26:15.818437+05:30	2026-10-14 16:26:15.781273	2026-10-14 16:26:15.818437
50	127	50	Shop 7640435179	t	2026-10-15 10:20:47.780076+05:30	2026-10-15 10:20:47.690601	2026-10-15 10:20:47.780076
49	126	49	Shop 764043517	t	2026-10-15 10:20:47.787245+05:30	2026-10-15 10:20:47.675011	2026-10-15 10:20:47.787245
51	128	51	Shop 781488023	t	2026-10-15 10:20:47.857532+05:30	2026-10-15 10:20:47.81997	2026-10-15 10:20:47.857532
53	131	53	Shop 7870395269	t	2026-10-15 10:20:47.923777+05:30	2026-10-15 10:20:47.886807	2026-10-15 10:20:47.923777
52	130	52	Shop 787039526	t	2026-10-15 10:20:47.928391+05:30	2026-10-15 10:20:47.876071	2026-10-15 10:20:47.928391
55	137	55	Shop B801317832	t	2026-10-15 10:20:48.483019+05:30	2026-10-15 10:20:48.18547	2026-10-15 10:20:48.483019
54	136	54	Shop A801317832	t	2026-10-15 10:20:48.487815+05:30	2026-10-15 10:20:48.019845	2026-10-15 10:20:48.487815
56	49	57	2222	f	\N	2026-10-15 11:24:56.807886	2026-10-15 11:24:56.807886
57	49	57	Shop - 1	f	\N	2026-10-15 11:25:10.009312	2026-10-15 11:25:10.009312
\.


--
-- TOC entry 5680 (class 0 OID 117152)
-- Dependencies: 260
-- Data for Name: user_state; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_state (id, user_id, user_country_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
89	1	103	Brändö	f	\N	2026-10-15 10:54:54.071675	2026-10-15 10:54:54.071675
90	1	103	Eckerö	f	\N	2026-10-15 10:54:54.094529	2026-10-15 10:54:54.094529
91	1	105	Annaba	f	\N	2026-10-15 10:54:54.101746	2026-10-15 10:54:54.101746
92	1	105	Batna	f	\N	2026-10-15 10:54:54.108738	2026-10-15 10:54:54.108738
93	1	105	Béjaïa	f	\N	2026-10-15 10:54:54.116179	2026-10-15 10:54:54.116179
94	1	105	Guelma	f	\N	2026-10-15 10:54:54.124429	2026-10-15 10:54:54.124429
95	1	105	Tissemsilt	f	\N	2026-10-15 10:54:54.131033	2026-10-15 10:54:54.131033
96	1	105	Tizi Ouzou	f	\N	2026-10-15 10:54:54.136233	2026-10-15 10:54:54.136233
97	1	105	Tlemcen	f	\N	2026-10-15 10:54:54.143595	2026-10-15 10:54:54.143595
6	26	13	State 13336822	t	2026-10-14 11:57:30.215951+05:30	2026-10-14 11:57:30.139069	2026-10-14 11:57:30.215951
7	27	15	State 22057624	f	\N	2026-10-14 11:57:30.224855	2026-10-14 11:57:30.224855
98	144	113	State 182355002	f	\N	2026-10-15 11:08:21.876425	2026-10-15 11:08:21.876425
8	31	17	State 73394922	t	2026-10-14 12:07:00.816558+05:30	2026-10-14 12:07:00.742633	2026-10-14 12:07:00.816558
9	32	19	State 82060724	f	\N	2026-10-14 12:07:00.824184	2026-10-14 12:07:00.824184
99	145	116	State 785454202	f	\N	2026-10-15 11:18:37.86815	2026-10-15 11:18:37.86815
100	146	118	State 36042801	f	\N	2026-10-15 11:18:50.377919	2026-10-15 11:18:50.377919
11	34	22	State B83246528	t	2026-10-14 12:07:01.144217+05:30	2026-10-14 12:07:00.85345	2026-10-14 12:07:01.144217
10	33	21	State A83246528	t	2026-10-14 12:07:01.148961+05:30	2026-10-14 12:07:00.835725	2026-10-14 12:07:01.148961
101	49	122	Adrar	f	\N	2026-10-15 11:24:15.348031	2026-10-15 11:24:15.348031
12	42	23	State 847312702	t	2026-10-14 12:07:08.560133+05:30	2026-10-14 12:07:08.485196	2026-10-14 12:07:08.560133
102	49	122	Aïn Témouchent	f	\N	2026-10-15 11:24:15.399111	2026-10-15 11:24:15.399111
103	49	122	Algiers	f	\N	2026-10-15 11:24:15.40779	2026-10-15 11:24:15.40779
14	44	25	State B857087506	t	2026-10-14 12:07:08.861797+05:30	2026-10-14 12:07:08.587572	2026-10-14 12:07:08.861797
13	43	24	State A857087506	t	2026-10-14 12:07:08.86611+05:30	2026-10-14 12:07:08.574638	2026-10-14 12:07:08.86611
104	22	123	Belize District	f	\N	2026-10-15 11:31:31.319889	2026-10-15 11:31:31.319889
15	54	26	State 899315822	t	2026-10-14 12:41:49.093299+05:30	2026-10-14 12:41:48.999813	2026-10-14 12:41:49.093299
16	55	28	State 909691024	f	\N	2026-10-14 12:41:49.100948	2026-10-14 12:41:49.100948
18	57	31	State B911479328	t	2026-10-14 12:41:49.445836+05:30	2026-10-14 12:41:49.141263	2026-10-14 12:41:49.445836
17	56	30	State A911479328	t	2026-10-14 12:41:49.451504+05:30	2026-10-14 12:41:49.119015	2026-10-14 12:41:49.451504
20	63	33	State 3279490019	t	2026-10-14 12:55:13.593591+05:30	2026-10-14 12:55:13.494762	2026-10-14 12:55:13.593591
19	62	32	State 327949001	t	2026-10-14 12:55:13.601119+05:30	2026-10-14 12:55:13.370464	2026-10-14 12:55:13.601119
21	64	34	State 769633603	t	2026-10-14 12:55:27.775883+05:30	2026-10-14 12:55:27.710777	2026-10-14 12:55:27.775883
23	66	36	State 7779789049	t	2026-10-14 12:55:27.834008+05:30	2026-10-14 12:55:27.797846	2026-10-14 12:55:27.834008
22	65	35	State 777978904	t	2026-10-14 12:55:27.837884+05:30	2026-10-14 12:55:27.783201	2026-10-14 12:55:27.837884
25	68	38	State 161897019	t	2026-10-14 15:27:40.358334+05:30	2026-10-14 15:27:40.274226	2026-10-14 15:27:40.358334
24	67	37	State 16189701	t	2026-10-14 15:27:40.364305+05:30	2026-10-14 15:27:40.235466	2026-10-14 15:27:40.364305
26	69	39	State 40066304	t	2026-10-14 15:27:40.467788+05:30	2026-10-14 15:27:40.412717	2026-10-14 15:27:40.467788
28	71	41	State 471232059	t	2026-10-14 15:27:40.540568+05:30	2026-10-14 15:27:40.486929	2026-10-14 15:27:40.540568
27	70	40	State 47123205	t	2026-10-14 15:27:40.545987+05:30	2026-10-14 15:27:40.475639	2026-10-14 15:27:40.545987
30	73	43	State 1932926019	t	2026-10-14 15:28:02.063415+05:30	2026-10-14 15:28:01.996213	2026-10-14 15:28:02.063415
29	72	42	State 193292601	t	2026-10-14 15:28:02.068698+05:30	2026-10-14 15:28:01.966632	2026-10-14 15:28:02.068698
31	74	44	State 209017604	t	2026-10-14 15:28:02.137716+05:30	2026-10-14 15:28:02.094924	2026-10-14 15:28:02.137716
33	76	46	State 2140857059	t	2026-10-14 15:28:02.186621+05:30	2026-10-14 15:28:02.156181	2026-10-14 15:28:02.186621
32	75	45	State 214085705	t	2026-10-14 15:28:02.190557+05:30	2026-10-14 15:28:02.145151	2026-10-14 15:28:02.190557
35	78	48	State 1589422019	t	2026-10-14 15:41:11.776255+05:30	2026-10-14 15:41:11.68934	2026-10-14 15:41:11.776255
34	77	47	State 158942201	t	2026-10-14 15:41:11.782133+05:30	2026-10-14 15:41:11.65306	2026-10-14 15:41:11.782133
37	80	50	State 1785631029	t	2026-10-14 15:41:11.836679+05:30	2026-10-14 15:41:11.806777	2026-10-14 15:41:11.836679
36	79	49	State 178563102	t	2026-10-14 15:41:11.840656+05:30	2026-10-14 15:41:11.790393	2026-10-14 15:41:11.840656
39	82	54	State 5730982019	t	2026-10-14 15:51:05.966014+05:30	2026-10-14 15:51:05.837127	2026-10-14 15:51:05.966014
38	81	53	State 573098201	t	2026-10-14 15:51:05.994325+05:30	2026-10-14 15:51:05.803815	2026-10-14 15:51:05.994325
41	84	56	State 6003831029	t	2026-10-14 15:51:06.10937+05:30	2026-10-14 15:51:06.037467	2026-10-14 15:51:06.10937
40	83	55	State 600383102	t	2026-10-14 15:51:06.118028+05:30	2026-10-14 15:51:06.011811	2026-10-14 15:51:06.118028
43	91	58	State 8387817179	t	2026-10-14 16:09:28.480226+05:30	2026-10-14 16:09:28.426121	2026-10-14 16:09:28.480226
42	90	57	State 838781717	t	2026-10-14 16:09:28.484804+05:30	2026-10-14 16:09:28.395846	2026-10-14 16:09:28.484804
44	92	59	State 851444423	t	2026-10-14 16:09:28.55887+05:30	2026-10-14 16:09:28.519896	2026-10-14 16:09:28.55887
45	93	61	State 856238325	f	\N	2026-10-14 16:09:28.565652	2026-10-14 16:09:28.565652
47	95	64	State 8567715269	t	2026-10-14 16:09:28.617003+05:30	2026-10-14 16:09:28.586756	2026-10-14 16:09:28.617003
46	94	63	State 856771526	t	2026-10-14 16:09:28.620801+05:30	2026-10-14 16:09:28.574373	2026-10-14 16:09:28.620801
49	101	66	State B868524232	t	2026-10-14 16:09:29.163789+05:30	2026-10-14 16:09:28.698765	2026-10-14 16:09:29.163789
48	100	65	State A868524232	t	2026-10-14 16:09:29.168835+05:30	2026-10-14 16:09:28.688563	2026-10-14 16:09:29.168835
51	107	68	State 5216632019	t	2026-10-14 16:26:15.70986+05:30	2026-10-14 16:26:15.584908	2026-10-14 16:26:15.70986
50	106	67	State 521663201	t	2026-10-14 16:26:15.71646+05:30	2026-10-14 16:26:15.345115	2026-10-14 16:26:15.71646
52	108	69	State 572189502	t	2026-10-14 16:26:15.772081+05:30	2026-10-14 16:26:15.726143	2026-10-14 16:26:15.772081
54	110	71	State 5775058039	t	2026-10-14 16:26:15.816127+05:30	2026-10-14 16:26:15.790166	2026-10-14 16:26:15.816127
53	109	70	State 577505803	t	2026-10-14 16:26:15.819429+05:30	2026-10-14 16:26:15.779561	2026-10-14 16:26:15.819429
55	119	73	State 60557102	f	\N	2026-10-15 10:15:00.75847	2026-10-15 10:15:00.75847
57	120	75	Goa	t	2026-10-15 10:15:00.865818+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:00.865818
56	120	75	Gujarat	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:01.040738
59	120	75	Maharashtra	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.865818	2026-10-15 10:15:01.040738
60	120	76	Osaka	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.865818	2026-10-15 10:15:01.040738
58	120	76	Tokyo	t	2026-10-15 10:15:01.040738+05:30	2026-10-15 10:15:00.847267	2026-10-15 10:15:01.040738
61	121	79	State 990636302	f	\N	2026-10-15 10:15:09.933569	2026-10-15 10:15:09.933569
63	122	81	Goa	t	2026-10-15 10:15:09.98669+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:09.98669
62	122	81	Gujarat	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:10.026868
65	122	81	Maharashtra	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.98669	2026-10-15 10:15:10.026868
66	122	82	Osaka	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.98669	2026-10-15 10:15:10.026868
64	122	82	Tokyo	t	2026-10-15 10:15:10.026868+05:30	2026-10-15 10:15:09.967741	2026-10-15 10:15:10.026868
68	127	85	State 7640435179	t	2026-10-15 10:20:47.782137+05:30	2026-10-15 10:20:47.688791	2026-10-15 10:20:47.782137
67	126	84	State 764043517	t	2026-10-15 10:20:47.788599+05:30	2026-10-15 10:20:47.655936	2026-10-15 10:20:47.788599
69	128	86	State 781488023	t	2026-10-15 10:20:47.858897+05:30	2026-10-15 10:20:47.818555	2026-10-15 10:20:47.858897
70	129	88	State 786226025	f	\N	2026-10-15 10:20:47.866761	2026-10-15 10:20:47.866761
72	131	91	State 7870395269	t	2026-10-15 10:20:47.925231+05:30	2026-10-15 10:20:47.884867	2026-10-15 10:20:47.925231
71	130	90	State 787039526	t	2026-10-15 10:20:47.929627+05:30	2026-10-15 10:20:47.874295	2026-10-15 10:20:47.929627
74	137	93	State B801317832	t	2026-10-15 10:20:48.484349+05:30	2026-10-15 10:20:48.184111	2026-10-15 10:20:48.484349
73	136	92	State A801317832	t	2026-10-15 10:20:48.488944+05:30	2026-10-15 10:20:48.01596	2026-10-15 10:20:48.488944
76	142	94	Goa	t	2026-10-15 10:20:48.497521+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.497521
75	142	94	Gujarat	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.503663
78	142	94	Maharashtra	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.497521	2026-10-15 10:20:48.503663
79	142	95	Osaka	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.497521	2026-10-15 10:20:48.503663
77	142	95	Tokyo	t	2026-10-15 10:20:48.503663+05:30	2026-10-15 10:20:48.494425	2026-10-15 10:20:48.503663
80	1	99	London	f	\N	2026-10-15 10:23:05.109394	2026-10-15 10:23:05.109394
81	143	100	State 867567401	f	\N	2026-10-15 10:28:58.698855	2026-10-15 10:28:58.698855
83	49	107	Andaman and Nicobar Islands	f	\N	2026-10-15 10:43:07.327004	2026-10-15 10:43:07.327004
84	49	107	Andhra Pradesh	f	\N	2026-10-15 10:43:07.345013	2026-10-15 10:43:07.345013
85	49	107	Arunachal Pradesh	f	\N	2026-10-15 10:43:07.361194	2026-10-15 10:43:07.361194
86	49	107	Assam	f	\N	2026-10-15 10:43:07.371447	2026-10-15 10:43:07.371447
87	49	107	Bihar	f	\N	2026-10-15 10:43:07.3781	2026-10-15 10:43:07.3781
88	49	107	Chandigarh	f	\N	2026-10-15 10:43:07.384258	2026-10-15 10:43:07.384258
105	22	123	Cayo District	f	\N	2026-10-15 11:31:31.355882	2026-10-15 11:31:31.355882
106	22	123	Corozal District	f	\N	2026-10-15 11:31:31.374811	2026-10-15 11:31:31.374811
107	23	130	Azad Kashmir	f	\N	2026-10-15 11:44:45.248887	2026-10-15 11:44:45.248887
108	23	130	Balochistan	f	\N	2026-10-15 11:44:45.428204	2026-10-15 11:44:45.428204
109	23	130	Gilgit-Baltistan	f	\N	2026-10-15 11:44:45.437258	2026-10-15 11:44:45.437258
110	23	130	Islamabad	f	\N	2026-10-15 11:44:45.44457	2026-10-15 11:44:45.44457
111	23	130	Khaibar Pakhtunkhwa	f	\N	2026-10-15 11:44:45.479807	2026-10-15 11:44:45.479807
112	23	130	Punjab	f	\N	2026-10-15 11:44:45.501643	2026-10-15 11:44:45.501643
113	23	130	Sindh	f	\N	2026-10-15 11:44:45.509495	2026-10-15 11:44:45.509495
114	23	130	Federally Administered Tribal Areas	f	\N	2026-10-15 11:44:45.515955	2026-10-15 11:44:45.515955
82	49	106	British Indian Ocean Territory	t	2026-10-15 11:45:07.402459+05:30	2026-10-15 10:42:55.459652	2026-10-15 11:45:07.402459
115	49	127	Badghis	f	\N	2026-10-15 11:57:45.002078	2026-10-15 11:57:45.002078
116	49	127	Baghlan	f	\N	2026-10-15 11:57:45.089402	2026-10-15 11:57:45.089402
117	49	127	Balkh	f	\N	2026-10-15 11:57:45.100112	2026-10-15 11:57:45.100112
118	49	127	Bamyan	f	\N	2026-10-15 11:57:45.110437	2026-10-15 11:57:45.110437
\.


--
-- TOC entry 5692 (class 0 OID 117362)
-- Dependencies: 272
-- Data for Name: user_supervisor; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_supervisor (id, user_id, user_cell_id, name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	49	1	Supervisor - 1	f	\N	2026-10-15 11:25:54.444385	2026-10-15 11:25:54.444385
2	49	1	Supervisor - 2	f	\N	2026-10-15 11:26:04.476523	2026-10-15 11:26:04.476523
\.


--
-- TOC entry 5648 (class 0 OID 116580)
-- Dependencies: 228
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password, first_name, last_name, mobile, role_id, branch_id, company_id, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
22	Lucy_Homenick7	102528@gmail.com	$2a$10$KUIefDoGDTbc6t6GtwBYruBpDbmbjtHNJBeAfUgzO7ecvxAc9xN0q	Edison	Stiedemann	+442896779618	5	4	\N	active	f	\N	2026-10-14 10:26:06.7819	2026-10-15 09:56:20.930348
25	auth_customer_986520801	auth_customer_986520801@auth-test.example.com	$2a$04$XSBNZxJCCtOu2Ne8YSX1m.v/1QKwlrmr4RYVzkYlEKLLjeGU5FTES	\N	\N	\N	5	\N	\N	active	t	2026-10-14 11:57:29.975248+05:30	2026-10-14 11:57:29.873392	2026-10-14 11:57:29.975248
5	auth_admin_713517301	auth_admin_713517301@auth-test.example.com	$2a$04$AC24ysLcmPekEM/lObyaNeoZHHv9cDlQjnq0/yO1lCzcFg0R5HM2K	\N	\N	\N	2	\N	\N	active	t	2026-10-13 18:01:47.14162+05:30	2026-10-13 18:01:47.137081	2026-10-13 18:01:47.14162
9	auth_admin_282463701	auth_admin_282463701@auth-test.example.com	$2a$04$8qBjYvKviY/ysmSo7QaH6OUuShcEE7G06T.IEbnXLlZraxzncYD4a	\N	\N	\N	2	\N	\N	active	t	2026-10-13 18:01:52.832984+05:30	2026-10-13 18:01:52.827419	2026-10-13 18:01:52.832984
13	machine_customer_93421215	machine_customer_93421215@auth-test.example.com	$2a$04$OGwqtDinjmn89dhmrKdo4Ol0.PylMkfboIpyFZeHo28tYECs25Ulu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:37:21.081338+05:30	2026-10-14 09:37:20.935952	2026-10-14 09:37:21.081338
15	auth_customer_750072301	auth_customer_750072301@auth-test.example.com	$2a$04$xPDweXCZoNk469dD9Vz6/uduU8fb/5h4b7URw5P6ve2kCiJpsv4wS	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:37:37.633428+05:30	2026-10-14 09:37:37.510447	2026-10-14 09:37:37.633428
14	auth_admin_750072301	auth_admin_750072301@auth-test.example.com	$2a$04$ruQlTZ2UEj6PvhB4A/5gCuet9CiUV5Zd2irogmGjEBMF26t6T4ihi	\N	\N	\N	2	\N	\N	inactive	t	2026-10-14 09:37:37.637579+05:30	2026-10-14 09:37:37.503355	2026-10-14 09:37:37.637579
16	machine_customer_229469522	machine_customer_229469522@auth-test.example.com	$2a$04$Z8L8zLCsp3.nN4Q6MJdYFOOKCJec1pKpZNwv6jJcEMkObAvNPhAai	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:37:54.021184+05:30	2026-10-14 09:37:52.295822	2026-10-14 09:37:54.021184
17	org_owner_403776224	org_owner_403776224@auth-test.example.com	$2a$04$XQCfI2fjjW7UVk8UU1.aeOD9IJTiheJp.Kwfz67hHnqf13ekLgmAC	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:37:54.124049+05:30	2026-10-14 09:37:54.038966	2026-10-14 09:37:54.124049
19	auth_customer_611452401	auth_customer_611452401@auth-test.example.com	$2a$04$K5FH0CzctIURWD6tb/pwoODyEqkpAFTK.K56OYmTVvBM.uDucJV5e	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:38:36.206436+05:30	2026-10-14 09:38:36.120179	2026-10-14 09:38:36.206436
18	auth_admin_611452401	auth_admin_611452401@auth-test.example.com	$2a$04$P6lh6sUpnK7hvFvH4/oUquuiT9QiqEoXfdUAiPIYKm1NWVNgTnSpe	\N	\N	\N	2	\N	\N	inactive	t	2026-10-14 09:38:36.206954+05:30	2026-10-14 09:38:36.116349	2026-10-14 09:38:36.206954
20	machine_customer_631445522	machine_customer_631445522@auth-test.example.com	$2a$04$BmszFPcA4yZkBWozNH9mQOwIjk54EY8yfCOpZb.fKKlgP1/tSAHVi	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:38:36.387372+05:30	2026-10-14 09:38:36.315861	2026-10-14 09:38:36.387372
21	org_owner_638837924	org_owner_638837924@auth-test.example.com	$2a$04$RhJzR.UdoeWU.kjjV6IQnOete2.ZFBHR9s2zW1qYZbxNbWD0fWuFC	\N	\N	\N	5	\N	\N	active	t	2026-10-14 09:38:36.393139+05:30	2026-10-14 09:38:36.389513	2026-10-14 09:38:36.393139
1	superadmin	superadmin@example.com	$2a$10$1PYn7cI1QnlTWY/.whcX0OIvmcgq7jWfCxq1jRyrzFd2pSV52RwE.	SuperAdmin	\N	\N	1	\N	\N	active	f	\N	2026-10-13 17:47:02.917183	2026-10-14 10:12:30.404709
23	cust_93	103424@gmail.com	$2a$10$C/awttlaqD.PMi7aEyCqierUl8n/vJ1gkZgATqAmXgc/xQdPyzpgi	Kolby	Skiles	+81394556979	5	5	\N	active	f	\N	2026-10-14 10:34:33.413564	2026-10-14 10:34:33.413564
24	auth_admin_986520801	auth_admin_986520801@auth-test.example.com	$2a$04$K6tyQJb6mia4aHvXH4KOsek1hrzxuJr17/C5ZqeNbs06KC0FRJKuW	\N	\N	\N	2	\N	\N	inactive	t	2026-10-14 11:57:29.977182+05:30	2026-10-14 11:57:29.867256	2026-10-14 11:57:29.977182
26	machine_customer_13336822	machine_customer_13336822@auth-test.example.com	$2a$04$bUNVCfMx3QivaUyOWVgn1uqN0mAjh3KlUiZ8ltNz1X1b235CkwE6G	\N	\N	\N	5	\N	\N	active	t	2026-10-14 11:57:30.218762+05:30	2026-10-14 11:57:30.134493	2026-10-14 11:57:30.218762
27	org_owner_22057624	org_owner_22057624@auth-test.example.com	$2a$04$b/4UJX4r3LhXFGzk3UKO7OHIrF/VrHsT2TVAb4VwH2Q.3RQs4C7Sq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 11:57:30.227142+05:30	2026-10-14 11:57:30.221932	2026-10-14 11:57:30.227142
30	auth_operator_43110001	auth_operator_43110001@auth-test.example.com	$2a$04$xxy0qm3Wanzc4S8YxLFwc.FdAG1ow9popS/eC6VYL3G1xs.2Ya9Ga	\N	\N	\N	4	\N	\N	active	t	2026-10-14 12:07:00.535709+05:30	2026-10-14 12:07:00.439972	2026-10-14 12:07:00.535709
29	auth_customer_43110001	auth_customer_43110001@auth-test.example.com	$2a$04$/YlKPwoRRcgMTBVD9zUhVuV1EmO/NLHjvKB3DXmMh16ayN.rS0k.K	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:00.539455+05:30	2026-10-14 12:07:00.438405	2026-10-14 12:07:00.539455
28	auth_admin_43110001	auth_admin_43110001@auth-test.example.com	$2a$04$ARL/LkMvFu8VReWkLxNvO.apcqs0uz66JvrXDPFQ2aQsURfgs7Sq2	\N	\N	\N	2	\N	29	inactive	t	2026-10-14 12:07:00.54016+05:30	2026-10-14 12:07:00.433112	2026-10-14 12:07:00.54016
31	machine_customer_73394922	machine_customer_73394922@auth-test.example.com	$2a$04$4xo20L6n7VxzIco4UR9lN.ik6aoXuqabnhKXj2rwgo6e0zgvOzPPu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:00.819442+05:30	2026-10-14 12:07:00.735149	2026-10-14 12:07:00.819442
32	org_owner_82060724	org_owner_82060724@auth-test.example.com	$2a$04$gt9GYk0zSiLnQoeY3gfDheX9vWKd.BNGXClD3DShcHPwlaETNziLa	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:00.826073+05:30	2026-10-14 12:07:00.82174	2026-10-14 12:07:00.826073
37	ua_sup83246528	ua_sup83246528@user-access.example.com	$2a$10$Rm0xkD/Hn0qQyaNKq8UKDe1c/.Q5UQW9sAjEOmLBvuW9uPGt5uaT2	\N	\N	+919876543210	3	12	33	active	t	2026-10-14 12:07:01.066399+05:30	2026-10-14 12:07:01.062939	2026-10-14 12:07:01.066399
34	machine_customer_B83246528	machine_customer_B83246528@auth-test.example.com	$2a$04$7K84xmv6LC9Iq4SlJleUfuOIbrcPXewOWcGuikmQF8cJz18EunNkG	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:01.146449+05:30	2026-10-14 12:07:00.851154	2026-10-14 12:07:01.146449
38	ua_sa83246528	ua_sa83246528@user-access.example.com	$2a$10$JJWW.lbz0AMqEY/VgMhWMeq8J8qErKrAfb/Hp2tAClx.eT9CrbTre	\N	\N	+919876543210	4	13	34	active	t	2026-10-14 12:07:01.139346+05:30	2026-10-14 12:07:01.136321	2026-10-14 12:07:01.139346
36	ua_op83246528	ua_op83246528@user-access.example.com	$2a$10$qFdnz6Nr/imUxvpTSeLbReifDScME/IYzBSjdeWs9tM9Tf3xldOrO	\N	\N	+919876543210	4	12	33	active	t	2026-10-14 12:07:01.142044+05:30	2026-10-14 12:07:00.988792	2026-10-14 12:07:01.142044
35	ua_adm83246528	ua_adm83246528@user-access.example.com	$2a$10$ShkqP7zX8n0KSulMmawKOOUdBWhA06cBl/99IqpYBNFYrna7hjEjm	\N	\N	+919876543210	2	12	33	active	t	2026-10-14 12:07:01.142269+05:30	2026-10-14 12:07:00.916939	2026-10-14 12:07:01.142269
33	machine_customer_A83246528	machine_customer_A83246528@auth-test.example.com	$2a$04$tBeoeIHxzzxo0X9Nxwjcv.JtmQ50oPELlkKGLDdomjMDLTI1Lzkyi	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:01.152179+05:30	2026-10-14 12:07:00.833475	2026-10-14 12:07:01.152179
41	auth_operator_836972301	auth_operator_836972301@auth-test.example.com	$2a$04$WQW.qd2mpsb/qQFkk2lQX.hkX38vaMJ1tKcE9D1H1rUj7ZMk3eWlG	\N	\N	\N	4	\N	\N	active	t	2026-10-14 12:07:08.471588+05:30	2026-10-14 12:07:08.378539	2026-10-14 12:07:08.471588
40	auth_customer_836972301	auth_customer_836972301@auth-test.example.com	$2a$04$axSAWgL6lCJDztSZS7J8duljIGZVO7EPLd8.8mHlWnQZC.r8rUHiu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:08.471996+05:30	2026-10-14 12:07:08.376826	2026-10-14 12:07:08.471996
39	auth_admin_836972301	auth_admin_836972301@auth-test.example.com	$2a$04$JcKMgeHgAMowfXRWMRz2CuND/0wOlEey/uNJpSrGWFqJxAO/c7i4K	\N	\N	\N	2	\N	40	inactive	t	2026-10-14 12:07:08.472319+05:30	2026-10-14 12:07:08.371859	2026-10-14 12:07:08.472319
42	machine_customer_847312702	machine_customer_847312702@auth-test.example.com	$2a$04$keOEIkP3jpyRKPlZ.FE0H.GUpQnhT02lM9CddKKwiNI6TDxUHfuFW	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:08.563503+05:30	2026-10-14 12:07:08.474225	2026-10-14 12:07:08.563503
45	ua_adm857087506	ua_adm857087506@user-access.example.com	$2a$10$FRQ/EzbUpZS959dgdWlTEuzMz0QT4QLVDLpi03eAtpZodAJjZdkDK	\N	\N	+919876543210	2	15	43	active	t	2026-10-14 12:07:08.859908+05:30	2026-10-14 12:07:08.64863	2026-10-14 12:07:08.859908
44	machine_customer_B857087506	machine_customer_B857087506@auth-test.example.com	$2a$04$pKeGOcSOmbGYmRnPXeDYNOEtvcehMELDnEccm2XKTunoI7Nqnu.XK	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:08.863796+05:30	2026-10-14 12:07:08.585344	2026-10-14 12:07:08.863796
43	machine_customer_A857087506	machine_customer_A857087506@auth-test.example.com	$2a$04$C5lJzUlhevf.TbmTYqrIl.1h/gVJe0/fl.pJQkOYUQRsqy5iTZ81G	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:07:08.868665+05:30	2026-10-14 12:07:08.572164	2026-10-14 12:07:08.868665
47	ua_sup857087506	ua_sup857087506@user-access.example.com	$2a$10$9lVIgtMqAAjegmhiXx4rYea9DkQDBBOJl3tloxcZLMkSZxZWqPnIm	\N	\N	+919876543210	3	15	43	active	t	2026-10-14 12:07:08.787201+05:30	2026-10-14 12:07:08.783353	2026-10-14 12:07:08.787201
63	machine_customer_3279490019	machine_customer_3279490019@auth-test.example.com	$2a$04$E81oJq/vja.sCEXeNeMXEOVMMM4w2w/X59YpnBNl8hzbOJf3QgV9m	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:55:13.598564+05:30	2026-10-14 12:55:13.488241	2026-10-14 12:55:13.598564
62	machine_customer_327949001	machine_customer_327949001@auth-test.example.com	$2a$04$NcdC6rfk3DnojKTm3ipvVOl80ZH55Sj9ipeAOLdznWELTHj3U8g26	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:55:13.60348+05:30	2026-10-14 12:55:13.282519	2026-10-14 12:55:13.60348
48	ua_sa857087506	ua_sa857087506@user-access.example.com	$2a$10$VgRU32hltc896NYzK2iSDOhgwLD5nBC66nfQyex7dbkl6ZfWjGL1e	\N	\N	+919876543210	4	16	44	active	t	2026-10-14 12:07:08.856805+05:30	2026-10-14 12:07:08.853665	2026-10-14 12:07:08.856805
46	ua_op857087506	ua_op857087506@user-access.example.com	$2a$10$gHlmJvD1UsGmy8NfEMJNk..I6kBojjoZe3a9CKDUfMtDq77FYBqKm	\N	\N	+919876543210	4	15	43	active	t	2026-10-14 12:07:08.859676+05:30	2026-10-14 12:07:08.715634	2026-10-14 12:07:08.859676
49	cust_15	121411@gmail.com	$2a$10$Jld9rjRYTrbURFGfPTdzs.vxZWpu.LY227ocIxFuzQ6cFbVUPRo1.	Carol	Schultz	+13512223332	5	6	\N	active	f	\N	2026-10-14 12:14:25.190123	2026-10-14 12:14:25.190123
50	cust_22	123003@gmail.com	$2a$10$9h2tO8I3aTLNok.Vlv4RlOq5gDKrT/Kg6MoiUkpVY3nbcnehH3f6e	Bernadette	Runolfsdottir	+9715522717507	3	6	22	active	f	\N	2026-10-14 12:30:26.127054	2026-10-14 12:30:26.127054
64	machine_customer_769633603	machine_customer_769633603@auth-test.example.com	$2a$04$gfoI/HX77qREed1OU7dF0eChmOrPfVIePPO0CH/YcdOZz/egImamu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:55:27.778623+05:30	2026-10-14 12:55:27.697953	2026-10-14 12:55:27.778623
66	machine_customer_7779789049	machine_customer_7779789049@auth-test.example.com	$2a$04$Nqe5dsSqlnLKpq01UROoneiiMUZmpg5KE7P8wjPCLxNlHWVOg8UU2	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:55:27.836056+05:30	2026-10-14 12:55:27.795728	2026-10-14 12:55:27.836056
53	auth_operator_870733501	auth_operator_870733501@auth-test.example.com	$2a$04$PnQ9HlX46oHLf0gfgazeC.9nRN22lubRaxGtbR00ozqaf4xtSt.5K	\N	\N	\N	4	\N	\N	active	t	2026-10-14 12:41:48.816523+05:30	2026-10-14 12:41:48.717694	2026-10-14 12:41:48.816523
52	auth_customer_870733501	auth_customer_870733501@auth-test.example.com	$2a$04$QrpAUhjwGCp4PTf47RMxt.1S3Xq.k9zcuTVhFU9rRCPhI8he4EsJq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:41:48.81693+05:30	2026-10-14 12:41:48.716231	2026-10-14 12:41:48.81693
51	auth_admin_870733501	auth_admin_870733501@auth-test.example.com	$2a$04$XunNLmBvJqeEF/CUi8.uZO2STyq0I2loFuERIbkfgGNZ9u3VDVjcm	\N	\N	\N	2	\N	52	inactive	t	2026-10-14 12:41:48.817267+05:30	2026-10-14 12:41:48.709719	2026-10-14 12:41:48.817267
65	machine_customer_777978904	machine_customer_777978904@auth-test.example.com	$2a$04$4dpGANRfOWuc.qZ/3HmLUeGeRiahbQ.a8U8fh94iv3BiWqalnjJs2	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:55:27.839979+05:30	2026-10-14 12:55:27.780807	2026-10-14 12:55:27.839979
54	machine_customer_899315822	machine_customer_899315822@auth-test.example.com	$2a$04$Fjz8gWc0n91vP30aL44dWeqxVhyY.IvqJirjX1YC52ikZnI0xKFZq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:41:49.096055+05:30	2026-10-14 12:41:48.994343	2026-10-14 12:41:49.096055
55	org_owner_909691024	org_owner_909691024@auth-test.example.com	$2a$04$13SFckoubn0v6OlzZF.cMOxGQlYejR2g7hXsNmM6a6sVFy1Xxh6i2	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:41:49.103578+05:30	2026-10-14 12:41:49.098438	2026-10-14 12:41:49.103578
68	machine_customer_161897019	machine_customer_161897019@auth-test.example.com	$2a$04$o/zvIQEhkBOkh3njmwcaau8xcm2dy6tzD2tocys9XJjC33MMd8nTu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:27:40.361865+05:30	2026-10-14 15:27:40.271013	2026-10-14 15:27:40.361865
67	machine_customer_16189701	machine_customer_16189701@auth-test.example.com	$2a$04$1SayX7uq5D/DXGbQ3.xpXOQjKggqVI9aE/YK/HNt04OTb1RH0VR86	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:27:40.367637+05:30	2026-10-14 15:27:40.165092	2026-10-14 15:27:40.367637
69	machine_customer_40066304	machine_customer_40066304@auth-test.example.com	$2a$04$ZG.lmGhcaeYZ/IyfujzvouaWBKBGKUX643p7DbxIRNg9RlmUtZCQ.	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:27:40.470315+05:30	2026-10-14 15:27:40.402035	2026-10-14 15:27:40.470315
60	ua_sup911479328	ua_sup911479328@user-access.example.com	$2a$10$qxx7b4kUPSw3dNHQPXEmVukPRenWgL7g2GYelBpgZpONQNqe.roEO	\N	\N	+919876543210	3	18	56	active	t	2026-10-14 12:41:49.366372+05:30	2026-10-14 12:41:49.363436	2026-10-14 12:41:49.366372
61	ua_sa911479328	ua_sa911479328@user-access.example.com	$2a$10$VvYZPt8S8Wn86qtZyBINTeS5D2kew4MZPOdqDcK8/s56KSgKlxnGe	\N	\N	+919876543210	4	19	57	active	t	2026-10-14 12:41:49.439981+05:30	2026-10-14 12:41:49.435932	2026-10-14 12:41:49.439981
59	ua_op911479328	ua_op911479328@user-access.example.com	$2a$10$.u.XVIYQKtuBbD0.UP8bZeFkiW7tyP2K0ThcQG8tPZQvexiu9UsqO	\N	\N	+919876543210	4	18	56	active	t	2026-10-14 12:41:49.44316+05:30	2026-10-14 12:41:49.288914	2026-10-14 12:41:49.44316
58	ua_adm911479328	ua_adm911479328@user-access.example.com	$2a$10$SimmwlBtcjIwBknE8dWU3uNLJUvRau4U5IcivwPYZ5gLbxm2pzrNG	\N	\N	+919876543210	2	18	56	active	t	2026-10-14 12:41:49.443441+05:30	2026-10-14 12:41:49.2054	2026-10-14 12:41:49.443441
57	machine_customer_B911479328	machine_customer_B911479328@auth-test.example.com	$2a$04$ORfqYlN5RGIx1JDf8e7dQ.z4MyK/FdiMcfT4AZ1VzcXuo5lbABKEG	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:41:49.448263+05:30	2026-10-14 12:41:49.136741	2026-10-14 12:41:49.448263
56	machine_customer_A911479328	machine_customer_A911479328@auth-test.example.com	$2a$04$j5MO0yeUuRVhAnqq.E0QCuaIK2sLv3wWnwuWEYeiWmKTuJ4Qi8gyu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 12:41:49.454626+05:30	2026-10-14 12:41:49.115991	2026-10-14 12:41:49.454626
71	machine_customer_471232059	machine_customer_471232059@auth-test.example.com	$2a$04$pYT6EYC8xlkiBPkt3pgqkeqsqvonCB/2dffRlKR7nOlgISMt2942O	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:27:40.543964+05:30	2026-10-14 15:27:40.484493	2026-10-14 15:27:40.543964
70	machine_customer_47123205	machine_customer_47123205@auth-test.example.com	$2a$04$sWfsyWa8hDkGsRoRzKGU.uGqpF8qyhx8NE7GVaM5RwVNbnZI508oC	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:27:40.548435+05:30	2026-10-14 15:27:40.472815	2026-10-14 15:27:40.548435
73	machine_customer_1932926019	machine_customer_1932926019@auth-test.example.com	$2a$04$rgokrVHtoAETVw8Wn6xUW.7R5Iq4nRPvgfd1cy55q2vgOX7Cm7ZUe	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:28:02.066597+05:30	2026-10-14 15:28:01.993563	2026-10-14 15:28:02.066597
72	machine_customer_193292601	machine_customer_193292601@auth-test.example.com	$2a$04$B1xgONf9hOT2kbKowqTbHeqQrfoi1JvLJcnwGa4ZRFlXwXwAXzrBO	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:28:02.071122+05:30	2026-10-14 15:28:01.935924	2026-10-14 15:28:02.071122
74	machine_customer_209017604	machine_customer_209017604@auth-test.example.com	$2a$04$TIQTeWqef.1dWYMejYDpde7K2Ye2SoIxLi2DB5YbbICcxkyFYL1cq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:28:02.140209+05:30	2026-10-14 15:28:02.091193	2026-10-14 15:28:02.140209
76	machine_customer_2140857059	machine_customer_2140857059@auth-test.example.com	$2a$04$cC1ygZtWyAUctwoi8hYcQuQC5aYeupShn/2bCK4F7GqTPtNFZ39oS	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:28:02.188859+05:30	2026-10-14 15:28:02.153401	2026-10-14 15:28:02.188859
75	machine_customer_214085705	machine_customer_214085705@auth-test.example.com	$2a$04$2z6xsSlJpt8.1FfgzMwmCOh5sMnwCH5Ga8MWbQvgYHE0UOVuUoPuW	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:28:02.192325+05:30	2026-10-14 15:28:02.142408	2026-10-14 15:28:02.192325
78	machine_customer_1589422019	machine_customer_1589422019@auth-test.example.com	$2a$04$/XAM3UX4bbga6zj05iRpI..DXy9M0QrL/oYt9kA1my5flwy0.cum.	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:41:11.780176+05:30	2026-10-14 15:41:11.686068	2026-10-14 15:41:11.780176
77	machine_customer_158942201	machine_customer_158942201@auth-test.example.com	$2a$04$0BZ3fNB4uaHcO3uUCa16P.afHcJF68WwXQIBMYaHFHCeU.vMM5EI6	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:41:11.784151+05:30	2026-10-14 15:41:11.592012	2026-10-14 15:41:11.784151
80	machine_customer_1785631029	machine_customer_1785631029@auth-test.example.com	$2a$04$P5PHr7AtEwSGUHd0dYieQeA54W88ufsC0002/WPqwtQfXUILX8aTW	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:41:11.838849+05:30	2026-10-14 15:41:11.804188	2026-10-14 15:41:11.838849
79	machine_customer_178563102	machine_customer_178563102@auth-test.example.com	$2a$04$LLrclpLKKyOJc48G5OxWSOS14HoQBdkKbB2fSug4j4flwuE7dAn3i	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:41:11.843356+05:30	2026-10-14 15:41:11.786824	2026-10-14 15:41:11.843356
81	machine_customer_573098201	machine_customer_573098201@auth-test.example.com	$2a$04$nVX2BCajUACxZs5kliIYYuQB1.1I6NRlXfGH0zENyUicrzK1apHzq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:05.997226+05:30	2026-10-14 15:51:05.73377	2026-10-14 15:51:05.997226
82	machine_customer_5730982019	machine_customer_5730982019@auth-test.example.com	$2a$04$dNRLuI1zZD5BnBlPgJavrOD7xQUMN6GiVZy1XyqEe0v78c1u1c/Gu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:05.988746+05:30	2026-10-14 15:51:05.832793	2026-10-14 15:51:05.988746
84	machine_customer_6003831029	machine_customer_6003831029@auth-test.example.com	$2a$04$UIzJrQAd1JoDayudrSjL..8eZSqaH0r9W6Ug0bxpUq85qam9k9Al6	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:06.113569+05:30	2026-10-14 15:51:06.032717	2026-10-14 15:51:06.113569
83	machine_customer_600383102	machine_customer_600383102@auth-test.example.com	$2a$04$EAPHXMPs4bF2VpefKYpzAOpSd20x/opGIPhQ8kKWfMXN4fBVJlQYu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:06.122211+05:30	2026-10-14 15:51:06.005003	2026-10-14 15:51:06.122211
86	shift_other_612267903	shift_other_612267903@auth-test.example.com	$2a$04$JlQKspUwDebJBR3zzQouGe8y9QWBUEjgfoqXMcK2BMW/JHB8MQh4C	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:06.163359+05:30	2026-10-14 15:51:06.126377	2026-10-14 15:51:06.163359
85	shift_customer_612267903	shift_customer_612267903@auth-test.example.com	$2a$04$.XLKsa/wqKBI81t5WuyiJuddchgjenThoWcHJh.gs.w6zsR39SkZ6	\N	\N	\N	5	\N	\N	active	t	2026-10-14 15:51:06.163795+05:30	2026-10-14 15:51:06.124772	2026-10-14 15:51:06.163795
89	auth_operator_815370801	auth_operator_815370801@auth-test.example.com	$2a$04$xiUHSGjcHxiT.KJEqiILVOrKexYdQS7k9fUOiY1.BLGmPJjGWl3Ca	\N	\N	\N	4	\N	\N	active	t	2026-10-14 16:09:28.29515+05:30	2026-10-14 16:09:28.203597	2026-10-14 16:09:28.29515
88	auth_customer_815370801	auth_customer_815370801@auth-test.example.com	$2a$04$RJie3bImoY7qNBuRUeuQ7einmEONQPAbvK5x1WV1.ARrjhK8R2Yee	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.295586+05:30	2026-10-14 16:09:28.201982	2026-10-14 16:09:28.295586
87	auth_admin_815370801	auth_admin_815370801@auth-test.example.com	$2a$04$mxk2rhDrwtu6dDboKml00.Bbovji1t6AwKmzTEt5wo4aZhbRf0aIS	\N	\N	\N	2	\N	88	inactive	t	2026-10-14 16:09:28.295917+05:30	2026-10-14 16:09:28.155884	2026-10-14 16:09:28.295917
91	machine_customer_8387817179	machine_customer_8387817179@auth-test.example.com	$2a$04$VGSZBK/v9uboH4Rz0AHeV.oT730OFULUVCevaIxHWIrWqpXxHKvw.	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.483157+05:30	2026-10-14 16:09:28.422742	2026-10-14 16:09:28.483157
90	machine_customer_838781717	machine_customer_838781717@auth-test.example.com	$2a$04$LSeZekRBEIQziObJiIgFEeP4Hp10l4AH0NU/0Gwu..1X18KacLJ52	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.488874+05:30	2026-10-14 16:09:28.38895	2026-10-14 16:09:28.488874
92	machine_customer_851444423	machine_customer_851444423@auth-test.example.com	$2a$04$gqHbcLjvcMNw11Qam27/we9buYIiE9hsgN/xdWkOjn.mRk0xvewQm	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.561266+05:30	2026-10-14 16:09:28.517205	2026-10-14 16:09:28.561266
93	org_owner_856238325	org_owner_856238325@auth-test.example.com	$2a$04$tC8gJeELQRbKdmVgYWQ4Lut.2Pn2Wwz86qycqv1q4x57whwbQFxjS	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.567672+05:30	2026-10-14 16:09:28.56352	2026-10-14 16:09:28.567672
95	machine_customer_8567715269	machine_customer_8567715269@auth-test.example.com	$2a$04$ow2LtnznVLDL0C36EP.o6OySbCN33.BvKLPX2NrAhrHS3f9Q4XtD.	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.619281+05:30	2026-10-14 16:09:28.583995	2026-10-14 16:09:28.619281
94	machine_customer_856771526	machine_customer_856771526@auth-test.example.com	$2a$04$MIk1cWJt1fDJV3QZhNACOONJkjaNzLeD2f8vZQYR8.IeNfrjMk6oG	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.625045+05:30	2026-10-14 16:09:28.569759	2026-10-14 16:09:28.625045
97	shift_other_862567827	shift_other_862567827@auth-test.example.com	$2a$04$e0kO0fT20FhfTyZzznA2Wurdw6ey1GWQJPujDXbjmHkCeEHBqPDnu	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.662356+05:30	2026-10-14 16:09:28.628362	2026-10-14 16:09:28.662356
96	shift_customer_862567827	shift_customer_862567827@auth-test.example.com	$2a$04$J5HWQE4YTZ6MlT1IGvmomON.ET1fBDelNXhK0AxlIWU72v3B4ADQ2	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.662643+05:30	2026-10-14 16:09:28.627105	2026-10-14 16:09:28.662643
99	holiday_other_866248528	holiday_other_866248528@auth-test.example.com	$2a$04$E9105kGIpQosmJVP6TTRFeEgSnCsGI8yIlL0/aZmSBH6XYI5ZQTzK	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.678834+05:30	2026-10-14 16:09:28.66517	2026-10-14 16:09:28.678834
98	holiday_customer_866248528	holiday_customer_866248528@auth-test.example.com	$2a$04$W68ooGVuN6t4QbEQZAVklO0ON6c.2viU4/vqgv.3YKyP9SGtvTwHC	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:28.679107+05:30	2026-10-14 16:09:28.663927	2026-10-14 16:09:28.679107
104	ua_sup868524232	ua_sup868524232@user-access.example.com	$2a$10$0d9VolfoNz3cJlAdJ7pxJeWtUrG5IcxtQjpTrZv10zRnjBk1fbNg2	\N	\N	+919876543210	3	48	100	active	t	2026-10-14 16:09:29.081908+05:30	2026-10-14 16:09:29.078007	2026-10-14 16:09:29.081908
105	ua_sa868524232	ua_sa868524232@user-access.example.com	$2a$10$ugK2y7jR9xFK0bLz2CUpAOgsnJd8ISGpqa4r.hYT36RbOIRD5BFsO	\N	\N	+919876543210	4	49	101	active	t	2026-10-14 16:09:29.157067+05:30	2026-10-14 16:09:29.152617	2026-10-14 16:09:29.157067
103	ua_op868524232	ua_op868524232@user-access.example.com	$2a$10$Kx7VWmsrAdEbRj4wffVnSuPJ9KV.yu7Gr0PBeTI9SDL6cnQPJIMFy	\N	\N	+919876543210	4	48	100	active	t	2026-10-14 16:09:29.160521+05:30	2026-10-14 16:09:29.004149	2026-10-14 16:09:29.160521
102	ua_adm868524232	ua_adm868524232@user-access.example.com	$2a$10$wNERusis14/J30CMgCrRa.L39E5dYZhENyDx8ksDaGyTdmLGCYZ.y	\N	\N	+919876543210	2	48	100	active	t	2026-10-14 16:09:29.160731+05:30	2026-10-14 16:09:28.932166	2026-10-14 16:09:29.160731
101	machine_customer_B868524232	machine_customer_B868524232@auth-test.example.com	$2a$04$rC9aSkdPnyKPlSt2hKng1..pBhx3oT9iPUOV/vVxPy7qEgFify9wm	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:29.166438+05:30	2026-10-14 16:09:28.696585	2026-10-14 16:09:29.166438
100	machine_customer_A868524232	machine_customer_A868524232@auth-test.example.com	$2a$04$CzMqjJlRaW91RFONd739Duwx54558JxYmzgcu6MVVTC.3rcUaGypy	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:09:29.171391+05:30	2026-10-14 16:09:28.6863	2026-10-14 16:09:29.171391
107	machine_customer_5216632019	machine_customer_5216632019@auth-test.example.com	$2a$04$jgEdffwtr0DceFt2FA4xqumhan1117pxIf5p.JqSh84cOkcPuGwQW	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.71414+05:30	2026-10-14 16:26:15.579032	2026-10-14 16:26:15.71414
106	machine_customer_521663201	machine_customer_521663201@auth-test.example.com	$2a$04$v57isn7IRh4nTtXelApiG.PRApovPQ2PwVk4mG.F2wOLbWXtVkziq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.720657+05:30	2026-10-14 16:26:15.219256	2026-10-14 16:26:15.720657
108	machine_customer_572189502	machine_customer_572189502@auth-test.example.com	$2a$04$.PbtXAYwVlTlNLbQT.pAjublONpM99yfZyesSb/ZiqkBDBSyty7fa	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.774535+05:30	2026-10-14 16:26:15.723102	2026-10-14 16:26:15.774535
110	machine_customer_5775058039	machine_customer_5775058039@auth-test.example.com	$2a$04$dtSr4vZDFmsrM2VCnjbzxeOUg2yx8kBCne3erBNqr8kBVr2LlVbXe	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.817914+05:30	2026-10-14 16:26:15.787265	2026-10-14 16:26:15.817914
109	machine_customer_577505803	machine_customer_577505803@auth-test.example.com	$2a$04$jF3NmYcG/twVSIO4iVoZp.VjTY7RdAiNs6O958upEyAwnZjEjPm8G	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.821118+05:30	2026-10-14 16:26:15.776588	2026-10-14 16:26:15.821118
112	shift_other_582207604	shift_other_582207604@auth-test.example.com	$2a$04$FBluoWAC9FYvB/NMtqSLPeWNYV6gmLnIsEhv1.E6QRdCpTkRw1Gfm	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.866304+05:30	2026-10-14 16:26:15.826836	2026-10-14 16:26:15.866304
111	shift_customer_582207604	shift_customer_582207604@auth-test.example.com	$2a$04$sTIAGQvDK0MB3XmaGeHHAen6hOb.0YCD8CuRb4H6zr1.OfWh323uy	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.866706+05:30	2026-10-14 16:26:15.823661	2026-10-14 16:26:15.866706
114	holiday_other_586692705	holiday_other_586692705@auth-test.example.com	$2a$04$ImPMt7D/E6h2/5GkgjhMtu74syz3B4B9ZkmAgthnEnX2hHSA2xn1y	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.881297+05:30	2026-10-14 16:26:15.869294	2026-10-14 16:26:15.881297
113	holiday_customer_586692705	holiday_customer_586692705@auth-test.example.com	$2a$04$7pm4NI1Fi.WPbSwXwFYzoON5YUGtc8TgsRR/8BFBqK0lFpF5EqkzG	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:15.881519+05:30	2026-10-14 16:26:15.868009	2026-10-14 16:26:15.881519
116	shift_other_98919701	shift_other_98919701@auth-test.example.com	$2a$04$E1CeanP02nflekgDWnoeCOBuaObJkhzmm8.ZbjziWeGT0xt93MK0u	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:31.045301+05:30	2026-10-14 16:26:30.997021	2026-10-14 16:26:31.045301
115	shift_customer_98919701	shift_customer_98919701@auth-test.example.com	$2a$04$36SLcm8Ls2/RM3Uv7Byq..pbIr.07xgCtE2SuijF225EwQlbDXKyq	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:31.045844+05:30	2026-10-14 16:26:30.992002	2026-10-14 16:26:31.045844
117	shift_customer_52540701	shift_customer_52540701@auth-test.example.com	$2a$04$DjjGKX.Q6TsQ0HMeoaeDiOubpiziNJ02ShaBIdb4/ea4Zh4mtLIKK	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:41.267251+05:30	2026-10-14 16:26:40.528042	2026-10-14 16:26:41.267251
118	shift_other_52540701	shift_other_52540701@auth-test.example.com	$2a$04$oi64ov5FVhETpwY9eUQ/guwpR/IOZF2Yysv57ga4N/FX9zI62hjNe	\N	\N	\N	5	\N	\N	active	t	2026-10-14 16:26:41.266752+05:30	2026-10-14 16:26:40.611302	2026-10-14 16:26:41.266752
119	org_owner_60557102	org_owner_60557102@auth-test.example.com	$2a$04$3U0Z6d40PH06vxR.Hir2kuMDq7A8YqrBLxAi/z2yYEKyL1CbS6aKe	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:15:00.834335+05:30	2026-10-15 10:15:00.607603	2026-10-15 10:15:00.834335
120	geo_owner_84338503	geo_owner_84338503@auth-test.example.com	$2a$04$dyulhiDvBUUGCzQhY8RBduhWn/MBSstxoi.Bw/Zw6eaGL0MQhHo6K	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:15:01.081156+05:30	2026-10-15 10:15:00.846013	2026-10-15 10:15:01.081156
121	org_owner_990636302	org_owner_990636302@auth-test.example.com	$2a$04$YuIG8DaVcAoEAExZJzKHpuGmBRdsy2m3OlTte.fvzTiDMFbDBApGe	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:15:09.946514+05:30	2026-10-15 10:15:09.908388	2026-10-15 10:15:09.946514
122	geo_owner_995658803	geo_owner_995658803@auth-test.example.com	$2a$04$efpbh2P7tkbLCJ9gIPk.yOdBkRKduOhg.lCJ8SRoa82fpf/6vhzCS	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:15:10.039702+05:30	2026-10-15 10:15:09.958632	2026-10-15 10:15:10.039702
125	auth_operator_739253001	auth_operator_739253001@auth-test.example.com	$2a$04$tH5VOW9z44HefvAihqQKLuYdGSfH6GQ7UkLpKlSKyvGfgkmunHQlC	\N	\N	\N	4	\N	\N	active	t	2026-10-15 10:20:47.501556+05:30	2026-10-15 10:20:47.402033	2026-10-15 10:20:47.501556
124	auth_customer_739253001	auth_customer_739253001@auth-test.example.com	$2a$04$JShFEXOJfxnbJu9YMDshI.BabBQtfqFWT7j2q2hpbNUVITKLjnhGu	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.502049+05:30	2026-10-15 10:20:47.400582	2026-10-15 10:20:47.502049
123	auth_admin_739253001	auth_admin_739253001@auth-test.example.com	$2a$04$ldhQ7lI0Th0rpyBcf3nkfuVYrUPqWj/zdNSIjMcl0nxGtwZ2DRkei	\N	\N	\N	2	\N	124	inactive	t	2026-10-15 10:20:47.502435+05:30	2026-10-15 10:20:47.394409	2026-10-15 10:20:47.502435
127	machine_customer_7640435179	machine_customer_7640435179@auth-test.example.com	$2a$04$241Q3.OB9Q/6hSphCxP6p.1Z82Vit72xsa8TfmzZYQTUhsGqfAXVK	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.786509+05:30	2026-10-15 10:20:47.685934	2026-10-15 10:20:47.786509
126	machine_customer_764043517	machine_customer_764043517@auth-test.example.com	$2a$04$fLKzf3Qc.ucZOoSlp7mHl.NlTc40iVz1MjGOFGvk8TwEoQ1xspnBe	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.791164+05:30	2026-10-15 10:20:47.641874	2026-10-15 10:20:47.791164
128	machine_customer_781488023	machine_customer_781488023@auth-test.example.com	$2a$04$qgh0UVxtnXvIg0vqG7A3vulY8ikBOA76QAIpJ9UDhqbXBCjYwgexG	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.861469+05:30	2026-10-15 10:20:47.816238	2026-10-15 10:20:47.861469
129	org_owner_786226025	org_owner_786226025@auth-test.example.com	$2a$04$3jlf3DGqFt.Y6Le91ICIyeWLU83GhNGw9otTumt17o6l0WY1jgqY.	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.869974+05:30	2026-10-15 10:20:47.863783	2026-10-15 10:20:47.869974
131	machine_customer_7870395269	machine_customer_7870395269@auth-test.example.com	$2a$04$7mib7SIoNIWrdv7eHWaHPu5tsJifZH3yuQliZGD0d6B3qnQgwvza6	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.927702+05:30	2026-10-15 10:20:47.882323	2026-10-15 10:20:47.927702
130	machine_customer_787039526	machine_customer_787039526@auth-test.example.com	$2a$04$LdvgSDGhr0VrLTcnMqFu8OxlCEDcidsOxFDmAcz6lkBaLx6lG2toa	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.931767+05:30	2026-10-15 10:20:47.871795	2026-10-15 10:20:47.931767
133	shift_other_793253827	shift_other_793253827@auth-test.example.com	$2a$04$gb1IzhcrMO3kO8fU.rr4iuF5qa8jLxz6vudiP16W0hGTi/HHGtw1q	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.983855+05:30	2026-10-15 10:20:47.935094	2026-10-15 10:20:47.983855
132	shift_customer_793253827	shift_customer_793253827@auth-test.example.com	$2a$04$XdkmLf.s37VUYy8rfHUrNuIB7VbVcBxsCURQj79JyK6ENejWNtiPu	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:47.984107+05:30	2026-10-15 10:20:47.933842	2026-10-15 10:20:47.984107
135	holiday_other_798397328	holiday_other_798397328@auth-test.example.com	$2a$04$I.LWuuCBx.W5mP0NpUlndumhYdY3FnDkghkFqfairmRW3hEURTIyW	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:48.007104+05:30	2026-10-15 10:20:47.986491	2026-10-15 10:20:48.007104
134	holiday_customer_798397328	holiday_customer_798397328@auth-test.example.com	$2a$04$9e57gWIMeYEdnOvj84dUw.26OjSAUSg9jZOQIMNnn5fJPjBLJxpf2	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:48.007341+05:30	2026-10-15 10:20:47.985325	2026-10-15 10:20:48.007341
140	ua_sup801317832	ua_sup801317832@user-access.example.com	$2a$10$vRJ.1OFU.5nBQzJUJXzOEe2i68tSMfbm2qf8/ldIj7RuEqWfr/XSS	\N	\N	+919876543210	3	60	136	active	t	2026-10-15 10:20:48.403488+05:30	2026-10-15 10:20:48.399613	2026-10-15 10:20:48.403488
142	geo_owner_849278833	geo_owner_849278833@auth-test.example.com	$2a$04$WxzRAWPfvTdIp41ms.SlxOe78Khmjxyw417g0pMbHmcnvoAZ4hUFi	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:48.506303+05:30	2026-10-15 10:20:48.494123	2026-10-15 10:20:48.506303
141	ua_sa801317832	ua_sa801317832@user-access.example.com	$2a$10$h7xvPHdJEYZZT5t/fd9Ww.quwU9STiv8ZMfcQz/SeVoszrKp6iTf6	\N	\N	+919876543210	4	61	137	active	t	2026-10-15 10:20:48.478637+05:30	2026-10-15 10:20:48.475144	2026-10-15 10:20:48.478637
139	ua_op801317832	ua_op801317832@user-access.example.com	$2a$10$yUvsC4zNjutLuC7oB871kuVSQH2rYyH4j0iESfjrgQfHpeLAWhofC	\N	\N	+919876543210	4	60	136	active	t	2026-10-15 10:20:48.481642+05:30	2026-10-15 10:20:48.322469	2026-10-15 10:20:48.481642
138	ua_adm801317832	ua_adm801317832@user-access.example.com	$2a$10$Qd2vZROybu2Rg2KkcfybJOF5JcjjCsasvelF9BfcGFsIjvv8sONEu	\N	\N	+919876543210	2	60	136	active	t	2026-10-15 10:20:48.481922+05:30	2026-10-15 10:20:48.24788	2026-10-15 10:20:48.481922
137	machine_customer_B801317832	machine_customer_B801317832@auth-test.example.com	$2a$04$cVynC7VD7GdBfGgMXAW4XeMm8H0b1IPr4UmkfMa0r2L8x7plJTNoC	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:48.486581+05:30	2026-10-15 10:20:48.02452	2026-10-15 10:20:48.486581
136	machine_customer_A801317832	machine_customer_A801317832@auth-test.example.com	$2a$04$EHXnQhiiK/K9U31ffTW0VuEHtvePgyzvuqDVZl5DLnRf8Wu/2q0TO	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:20:48.490736+05:30	2026-10-15 10:20:48.014137	2026-10-15 10:20:48.490736
143	org_owner_867567401	org_owner_867567401@auth-test.example.com	$2a$04$LsvU0WjuKBhNJTSfrAxBeu.jDX7ogTLCIp/mHjxTNr5v1AdMvjMxi	\N	\N	\N	5	\N	\N	active	t	2026-10-15 10:28:58.714481+05:30	2026-10-15 10:28:58.678644	2026-10-15 10:28:58.714481
144	org_owner_182355002	org_owner_182355002@auth-test.example.com	$2a$04$iMKNXfphB5RlCNjzvmqH5OGrta/InUkTJgyh6ZiJo2hYrIS9BH97G	\N	\N	\N	5	\N	\N	active	t	2026-10-15 11:08:21.888477+05:30	2026-10-15 11:08:21.826217	2026-10-15 11:08:21.888477
145	org_owner_785454202	org_owner_785454202@auth-test.example.com	$2a$04$I1kha2NpQ6a8A3/LOIkASO7pC07BqvA.PcNTzXb0wOYJyJ3EXElNW	\N	\N	\N	5	\N	\N	active	t	2026-10-15 11:18:37.872306+05:30	2026-10-15 11:18:37.856856	2026-10-15 11:18:37.872306
146	org_owner_36042801	org_owner_36042801@auth-test.example.com	$2a$04$tRYKblKYKNu8bfYU7JrstuKpQgMeH6NTf5Ly.l.dR9nbMNT8M7gS6	\N	\N	\N	5	\N	\N	active	t	2026-10-15 11:18:50.381647+05:30	2026-10-15 11:18:50.362721	2026-10-15 11:18:50.381647
\.


--
-- TOC entry 5760 (class 0 OID 0)
-- Dependencies: 229
-- Name: addresses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.addresses_id_seq', 24, true);


--
-- TOC entry 5761 (class 0 OID 0)
-- Dependencies: 237
-- Name: auth_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.auth_tokens_id_seq', 41, true);


--
-- TOC entry 5762 (class 0 OID 0)
-- Dependencies: 239
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.branches_id_seq', 66, true);


--
-- TOC entry 5763 (class 0 OID 0)
-- Dependencies: 243
-- Name: brands_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.brands_id_seq', 316, true);


--
-- TOC entry 5764 (class 0 OID 0)
-- Dependencies: 289
-- Name: breaks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.breaks_id_seq', 40, true);


--
-- TOC entry 5765 (class 0 OID 0)
-- Dependencies: 241
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 105, true);


--
-- TOC entry 5766 (class 0 OID 0)
-- Dependencies: 253
-- Name: category_machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.category_machines_id_seq', 1, false);


--
-- TOC entry 5767 (class 0 OID 0)
-- Dependencies: 285
-- Name: clienteles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.clienteles_id_seq', 41, true);


--
-- TOC entry 5768 (class 0 OID 0)
-- Dependencies: 233
-- Name: company_profiles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.company_profiles_id_seq', 24, true);


--
-- TOC entry 5769 (class 0 OID 0)
-- Dependencies: 249
-- Name: controller_products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.controller_products_id_seq', 1, false);


--
-- TOC entry 5770 (class 0 OID 0)
-- Dependencies: 247
-- Name: controllers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.controllers_id_seq', 312, true);


--
-- TOC entry 5771 (class 0 OID 0)
-- Dependencies: 221
-- Name: countries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.countries_id_seq', 246, true);


--
-- TOC entry 5772 (class 0 OID 0)
-- Dependencies: 291
-- Name: holidays_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.holidays_id_seq', 16, true);


--
-- TOC entry 5773 (class 0 OID 0)
-- Dependencies: 251
-- Name: machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.machines_id_seq', 88, true);


--
-- TOC entry 5774 (class 0 OID 0)
-- Dependencies: 273
-- Name: models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.models_id_seq', 175, true);


--
-- TOC entry 5775 (class 0 OID 0)
-- Dependencies: 281
-- Name: part_process_machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.part_process_machines_id_seq', 58, true);


--
-- TOC entry 5776 (class 0 OID 0)
-- Dependencies: 279
-- Name: part_processes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.part_processes_id_seq', 37, true);


--
-- TOC entry 5777 (class 0 OID 0)
-- Dependencies: 283
-- Name: part_sub_programs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.part_sub_programs_id_seq', 19, true);


--
-- TOC entry 5778 (class 0 OID 0)
-- Dependencies: 277
-- Name: parts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.parts_id_seq', 28, true);


--
-- TOC entry 5779 (class 0 OID 0)
-- Dependencies: 245
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 128, true);


--
-- TOC entry 5780 (class 0 OID 0)
-- Dependencies: 225
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 8, true);


--
-- TOC entry 5781 (class 0 OID 0)
-- Dependencies: 287
-- Name: shifts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.shifts_id_seq', 18, true);


--
-- TOC entry 5782 (class 0 OID 0)
-- Dependencies: 223
-- Name: states_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.states_id_seq', 4124, true);


--
-- TOC entry 5783 (class 0 OID 0)
-- Dependencies: 269
-- Name: user_cell_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_cell_id_seq', 2, true);


--
-- TOC entry 5784 (class 0 OID 0)
-- Dependencies: 261
-- Name: user_city_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_city_id_seq', 67, true);


--
-- TOC entry 5785 (class 0 OID 0)
-- Dependencies: 257
-- Name: user_country_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_country_id_seq', 134, true);


--
-- TOC entry 5786 (class 0 OID 0)
-- Dependencies: 255
-- Name: user_filters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_filters_id_seq', 1, false);


--
-- TOC entry 5787 (class 0 OID 0)
-- Dependencies: 267
-- Name: user_line_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_line_id_seq', 57, true);


--
-- TOC entry 5788 (class 0 OID 0)
-- Dependencies: 263
-- Name: user_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_location_id_seq', 57, true);


--
-- TOC entry 5789 (class 0 OID 0)
-- Dependencies: 275
-- Name: user_machine_access_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_machine_access_id_seq', 20, true);


--
-- TOC entry 5790 (class 0 OID 0)
-- Dependencies: 235
-- Name: user_machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_machines_id_seq', 1, false);


--
-- TOC entry 5791 (class 0 OID 0)
-- Dependencies: 231
-- Name: user_profiles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_profiles_id_seq', 66, true);


--
-- TOC entry 5792 (class 0 OID 0)
-- Dependencies: 265
-- Name: user_shop_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_shop_id_seq', 57, true);


--
-- TOC entry 5793 (class 0 OID 0)
-- Dependencies: 259
-- Name: user_state_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_state_id_seq', 118, true);


--
-- TOC entry 5794 (class 0 OID 0)
-- Dependencies: 271
-- Name: user_supervisor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_supervisor_id_seq', 2, true);


--
-- TOC entry 5795 (class 0 OID 0)
-- Dependencies: 227
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 146, true);


--
-- TOC entry 5230 (class 2606 OID 116644)
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- TOC entry 5253 (class 2606 OID 116743)
-- Name: auth_tokens auth_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 5255 (class 2606 OID 116745)
-- Name: auth_tokens auth_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_token_key UNIQUE (token);


--
-- TOC entry 5262 (class 2606 OID 116790)
-- Name: branches branches_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_code_key UNIQUE (code);


--
-- TOC entry 5264 (class 2606 OID 116788)
-- Name: branches branches_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_pkey PRIMARY KEY (id);


--
-- TOC entry 5274 (class 2606 OID 116879)
-- Name: brands brands_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (id);


--
-- TOC entry 5403 (class 2606 OID 117974)
-- Name: breaks breaks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.breaks
    ADD CONSTRAINT breaks_pkey PRIMARY KEY (id);


--
-- TOC entry 5269 (class 2606 OID 116862)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 5306 (class 2606 OID 117059)
-- Name: category_machines category_machines_category_id_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_category_id_machine_id_key UNIQUE (category_id, machine_id);


--
-- TOC entry 5308 (class 2606 OID 117057)
-- Name: category_machines category_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5393 (class 2606 OID 117864)
-- Name: clienteles clienteles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clienteles
    ADD CONSTRAINT clienteles_pkey PRIMARY KEY (id);


--
-- TOC entry 5240 (class 2606 OID 116699)
-- Name: company_profiles company_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_pkey PRIMARY KEY (id);


--
-- TOC entry 5242 (class 2606 OID 116701)
-- Name: company_profiles company_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_user_id_key UNIQUE (user_id);


--
-- TOC entry 5288 (class 2606 OID 116939)
-- Name: controller_products controller_products_controller_id_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_controller_id_product_id_key UNIQUE (controller_id, product_id);


--
-- TOC entry 5290 (class 2606 OID 116937)
-- Name: controller_products controller_products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_pkey PRIMARY KEY (id);


--
-- TOC entry 5284 (class 2606 OID 116919)
-- Name: controllers controllers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers
    ADD CONSTRAINT controllers_pkey PRIMARY KEY (id);


--
-- TOC entry 5211 (class 2606 OID 116525)
-- Name: countries countries_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_code_key UNIQUE (code);


--
-- TOC entry 5213 (class 2606 OID 116523)
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- TOC entry 5409 (class 2606 OID 118027)
-- Name: holidays holidays_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.holidays
    ADD CONSTRAINT holidays_pkey PRIMARY KEY (id);


--
-- TOC entry 5303 (class 2606 OID 117008)
-- Name: machines machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5367 (class 2606 OID 117661)
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- TOC entry 5385 (class 2606 OID 117793)
-- Name: part_process_machines part_process_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_process_machines
    ADD CONSTRAINT part_process_machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5387 (class 2606 OID 117795)
-- Name: part_process_machines part_process_machines_process_machine_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_process_machines
    ADD CONSTRAINT part_process_machines_process_machine_key UNIQUE (process_id, machine_id);


--
-- TOC entry 5382 (class 2606 OID 117773)
-- Name: part_processes part_processes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_processes
    ADD CONSTRAINT part_processes_pkey PRIMARY KEY (id);


--
-- TOC entry 5391 (class 2606 OID 117826)
-- Name: part_sub_programs part_sub_programs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_sub_programs
    ADD CONSTRAINT part_sub_programs_pkey PRIMARY KEY (id);


--
-- TOC entry 5377 (class 2606 OID 117727)
-- Name: parts parts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parts
    ADD CONSTRAINT parts_pkey PRIMARY KEY (id);


--
-- TOC entry 5281 (class 2606 OID 116897)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 5217 (class 2606 OID 116554)
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- TOC entry 5219 (class 2606 OID 116552)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 5400 (class 2606 OID 117913)
-- Name: shifts shifts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shifts
    ADD CONSTRAINT shifts_pkey PRIMARY KEY (id);


--
-- TOC entry 5215 (class 2606 OID 116535)
-- Name: states states_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states
    ADD CONSTRAINT states_pkey PRIMARY KEY (id);


--
-- TOC entry 5357 (class 2606 OID 117345)
-- Name: user_cell user_cell_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cell
    ADD CONSTRAINT user_cell_pkey PRIMARY KEY (id);


--
-- TOC entry 5333 (class 2606 OID 117195)
-- Name: user_city user_city_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_city
    ADD CONSTRAINT user_city_pkey PRIMARY KEY (id);


--
-- TOC entry 5321 (class 2606 OID 117141)
-- Name: user_country user_country_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_country
    ADD CONSTRAINT user_country_pkey PRIMARY KEY (id);


--
-- TOC entry 5315 (class 2606 OID 117117)
-- Name: user_filters user_filters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters
    ADD CONSTRAINT user_filters_pkey PRIMARY KEY (id);


--
-- TOC entry 5351 (class 2606 OID 117315)
-- Name: user_line user_line_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_line
    ADD CONSTRAINT user_line_pkey PRIMARY KEY (id);


--
-- TOC entry 5339 (class 2606 OID 117255)
-- Name: user_location user_location_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_location
    ADD CONSTRAINT user_location_pkey PRIMARY KEY (id);


--
-- TOC entry 5370 (class 2606 OID 117690)
-- Name: user_machine_access user_machine_access_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machine_access
    ADD CONSTRAINT user_machine_access_pkey PRIMARY KEY (id);


--
-- TOC entry 5372 (class 2606 OID 117692)
-- Name: user_machine_access user_machine_access_user_machine_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machine_access
    ADD CONSTRAINT user_machine_access_user_machine_key UNIQUE (user_id, machine_id);


--
-- TOC entry 5249 (class 2606 OID 116720)
-- Name: user_machines user_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5251 (class 2606 OID 116722)
-- Name: user_machines user_machines_user_id_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_user_id_machine_id_key UNIQUE (user_id, machine_id);


--
-- TOC entry 5236 (class 2606 OID 116678)
-- Name: user_profiles user_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_pkey PRIMARY KEY (id);


--
-- TOC entry 5238 (class 2606 OID 116680)
-- Name: user_profiles user_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_user_id_key UNIQUE (user_id);


--
-- TOC entry 5345 (class 2606 OID 117285)
-- Name: user_shop user_shop_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_shop
    ADD CONSTRAINT user_shop_pkey PRIMARY KEY (id);


--
-- TOC entry 5327 (class 2606 OID 117165)
-- Name: user_state user_state_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_state
    ADD CONSTRAINT user_state_pkey PRIMARY KEY (id);


--
-- TOC entry 5363 (class 2606 OID 117375)
-- Name: user_supervisor user_supervisor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_supervisor
    ADD CONSTRAINT user_supervisor_pkey PRIMARY KEY (id);


--
-- TOC entry 5227 (class 2606 OID 116597)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 5271 (class 1259 OID 117075)
-- Name: brands_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX brands_code_key ON public.brands USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5272 (class 1259 OID 116952)
-- Name: brands_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX brands_name_key ON public.brands USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5404 (class 1259 OID 117980)
-- Name: breaks_user_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX breaks_user_name_key ON public.breaks USING btree (user_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5266 (class 1259 OID 117074)
-- Name: categories_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX categories_code_key ON public.categories USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5267 (class 1259 OID 116950)
-- Name: categories_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX categories_name_key ON public.categories USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5394 (class 1259 OID 117870)
-- Name: clienteles_user_client_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX clienteles_user_client_code_key ON public.clienteles USING btree (user_id, lower((client_code)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5282 (class 1259 OID 117672)
-- Name: controllers_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX controllers_name_key ON public.controllers USING btree (brand_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5410 (class 1259 OID 118034)
-- Name: holidays_user_name_date_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX holidays_user_name_date_key ON public.holidays USING btree (user_id, lower((name)::text), start_date) WHERE (((holiday_type)::text = 'holiday'::text) AND (is_deleted = false));


--
-- TOC entry 5411 (class 1259 OID 118033)
-- Name: holidays_user_week_day_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX holidays_user_week_day_key ON public.holidays USING btree (user_id, week_day) WHERE (((holiday_type)::text = 'week_off'::text) AND (is_deleted = false));


--
-- TOC entry 5231 (class 1259 OID 116661)
-- Name: idx_addresses_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_addresses_is_deleted ON public.addresses USING btree (is_deleted);


--
-- TOC entry 5232 (class 1259 OID 116660)
-- Name: idx_addresses_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_addresses_user_id ON public.addresses USING btree (user_id);


--
-- TOC entry 5256 (class 1259 OID 116753)
-- Name: idx_auth_tokens_expires_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_expires_at ON public.auth_tokens USING btree (expires_at);


--
-- TOC entry 5257 (class 1259 OID 117419)
-- Name: idx_auth_tokens_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_is_deleted ON public.auth_tokens USING btree (is_deleted);


--
-- TOC entry 5258 (class 1259 OID 117418)
-- Name: idx_auth_tokens_live_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_live_token ON public.auth_tokens USING btree (token) WHERE (is_deleted = false);


--
-- TOC entry 5259 (class 1259 OID 116752)
-- Name: idx_auth_tokens_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_token ON public.auth_tokens USING btree (token);


--
-- TOC entry 5260 (class 1259 OID 116751)
-- Name: idx_auth_tokens_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_user_id ON public.auth_tokens USING btree (user_id);


--
-- TOC entry 5265 (class 1259 OID 117398)
-- Name: idx_branches_user_city_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_branches_user_city_id ON public.branches USING btree (user_city_id);


--
-- TOC entry 5275 (class 1259 OID 116953)
-- Name: idx_brands_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_brands_is_deleted ON public.brands USING btree (is_deleted);


--
-- TOC entry 5405 (class 1259 OID 117982)
-- Name: idx_breaks_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_breaks_is_deleted ON public.breaks USING btree (is_deleted);


--
-- TOC entry 5406 (class 1259 OID 118049)
-- Name: idx_breaks_shift_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_breaks_shift_id ON public.breaks USING btree (shift_id);


--
-- TOC entry 5407 (class 1259 OID 117981)
-- Name: idx_breaks_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_breaks_user_id ON public.breaks USING btree (user_id);


--
-- TOC entry 5270 (class 1259 OID 116951)
-- Name: idx_categories_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categories_is_deleted ON public.categories USING btree (is_deleted);


--
-- TOC entry 5309 (class 1259 OID 117070)
-- Name: idx_category_machines_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_category_id ON public.category_machines USING btree (category_id);


--
-- TOC entry 5310 (class 1259 OID 117072)
-- Name: idx_category_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_is_deleted ON public.category_machines USING btree (is_deleted);


--
-- TOC entry 5311 (class 1259 OID 117071)
-- Name: idx_category_machines_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_machine_id ON public.category_machines USING btree (machine_id);


--
-- TOC entry 5395 (class 1259 OID 117872)
-- Name: idx_clienteles_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_clienteles_is_deleted ON public.clienteles USING btree (is_deleted);


--
-- TOC entry 5396 (class 1259 OID 117871)
-- Name: idx_clienteles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_clienteles_user_id ON public.clienteles USING btree (user_id);


--
-- TOC entry 5243 (class 1259 OID 116708)
-- Name: idx_company_profiles_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_company_profiles_is_deleted ON public.company_profiles USING btree (is_deleted);


--
-- TOC entry 5244 (class 1259 OID 116707)
-- Name: idx_company_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_company_profiles_user_id ON public.company_profiles USING btree (user_id);


--
-- TOC entry 5291 (class 1259 OID 116960)
-- Name: idx_controller_products_controller_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_controller_id ON public.controller_products USING btree (controller_id);


--
-- TOC entry 5292 (class 1259 OID 116962)
-- Name: idx_controller_products_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_is_deleted ON public.controller_products USING btree (is_deleted);


--
-- TOC entry 5293 (class 1259 OID 116961)
-- Name: idx_controller_products_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_product_id ON public.controller_products USING btree (product_id);


--
-- TOC entry 5285 (class 1259 OID 116958)
-- Name: idx_controllers_brand_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controllers_brand_id ON public.controllers USING btree (brand_id);


--
-- TOC entry 5286 (class 1259 OID 116959)
-- Name: idx_controllers_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controllers_is_deleted ON public.controllers USING btree (is_deleted);


--
-- TOC entry 5412 (class 1259 OID 118036)
-- Name: idx_holidays_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_holidays_is_deleted ON public.holidays USING btree (is_deleted);


--
-- TOC entry 5413 (class 1259 OID 118035)
-- Name: idx_holidays_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_holidays_user_id ON public.holidays USING btree (user_id);


--
-- TOC entry 5294 (class 1259 OID 117037)
-- Name: idx_machines_brand_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_brand_id ON public.machines USING btree (brand_id);


--
-- TOC entry 5295 (class 1259 OID 117038)
-- Name: idx_machines_controller_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_controller_id ON public.machines USING btree (controller_id);


--
-- TOC entry 5296 (class 1259 OID 117041)
-- Name: idx_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_is_deleted ON public.machines USING btree (is_deleted);


--
-- TOC entry 5297 (class 1259 OID 117678)
-- Name: idx_machines_model_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_model_id ON public.machines USING btree (model_id);


--
-- TOC entry 5298 (class 1259 OID 117039)
-- Name: idx_machines_online_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_online_status ON public.machines USING btree (online_status);


--
-- TOC entry 5299 (class 1259 OID 117036)
-- Name: idx_machines_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_product_id ON public.machines USING btree (product_id);


--
-- TOC entry 5300 (class 1259 OID 117040)
-- Name: idx_machines_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_status ON public.machines USING btree (status);


--
-- TOC entry 5301 (class 1259 OID 117035)
-- Name: idx_machines_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_user_id ON public.machines USING btree (user_id);


--
-- TOC entry 5364 (class 1259 OID 117668)
-- Name: idx_models_controller_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_models_controller_id ON public.models USING btree (controller_id);


--
-- TOC entry 5383 (class 1259 OID 117806)
-- Name: idx_part_process_machines_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_part_process_machines_machine_id ON public.part_process_machines USING btree (machine_id);


--
-- TOC entry 5379 (class 1259 OID 117780)
-- Name: idx_part_processes_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_part_processes_is_deleted ON public.part_processes USING btree (is_deleted);


--
-- TOC entry 5380 (class 1259 OID 117779)
-- Name: idx_part_processes_part_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_part_processes_part_id ON public.part_processes USING btree (part_id);


--
-- TOC entry 5388 (class 1259 OID 117833)
-- Name: idx_part_sub_programs_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_part_sub_programs_is_deleted ON public.part_sub_programs USING btree (is_deleted);


--
-- TOC entry 5389 (class 1259 OID 117832)
-- Name: idx_part_sub_programs_process_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_part_sub_programs_process_id ON public.part_sub_programs USING btree (process_id);


--
-- TOC entry 5373 (class 1259 OID 117879)
-- Name: idx_parts_clientele_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_parts_clientele_id ON public.parts USING btree (clientele_id);


--
-- TOC entry 5374 (class 1259 OID 117735)
-- Name: idx_parts_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_parts_is_deleted ON public.parts USING btree (is_deleted);


--
-- TOC entry 5375 (class 1259 OID 117734)
-- Name: idx_parts_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_parts_user_id ON public.parts USING btree (user_id);


--
-- TOC entry 5276 (class 1259 OID 116955)
-- Name: idx_products_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_category_id ON public.products USING btree (category_id);


--
-- TOC entry 5277 (class 1259 OID 116956)
-- Name: idx_products_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_is_deleted ON public.products USING btree (is_deleted);


--
-- TOC entry 5397 (class 1259 OID 117921)
-- Name: idx_shifts_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_shifts_is_deleted ON public.shifts USING btree (is_deleted);


--
-- TOC entry 5398 (class 1259 OID 117920)
-- Name: idx_shifts_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_shifts_user_id ON public.shifts USING btree (user_id);


--
-- TOC entry 5352 (class 1259 OID 117359)
-- Name: idx_user_cell_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_cell_is_deleted ON public.user_cell USING btree (is_deleted);


--
-- TOC entry 5353 (class 1259 OID 117357)
-- Name: idx_user_cell_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_cell_user_id ON public.user_cell USING btree (user_id);


--
-- TOC entry 5354 (class 1259 OID 117358)
-- Name: idx_user_cell_user_line_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_cell_user_line_id ON public.user_cell USING btree (user_line_id);


--
-- TOC entry 5328 (class 1259 OID 117209)
-- Name: idx_user_city_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_city_is_deleted ON public.user_city USING btree (is_deleted);


--
-- TOC entry 5329 (class 1259 OID 117207)
-- Name: idx_user_city_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_city_user_id ON public.user_city USING btree (user_id);


--
-- TOC entry 5330 (class 1259 OID 117208)
-- Name: idx_user_city_user_state_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_city_user_state_id ON public.user_city USING btree (user_state_id);


--
-- TOC entry 5317 (class 1259 OID 117149)
-- Name: idx_user_country_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_country_is_deleted ON public.user_country USING btree (is_deleted);


--
-- TOC entry 5318 (class 1259 OID 117148)
-- Name: idx_user_country_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_country_user_id ON public.user_country USING btree (user_id);


--
-- TOC entry 5312 (class 1259 OID 117125)
-- Name: idx_user_filters_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_filters_is_deleted ON public.user_filters USING btree (is_deleted);


--
-- TOC entry 5313 (class 1259 OID 117124)
-- Name: idx_user_filters_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_filters_user_id ON public.user_filters USING btree (user_id);


--
-- TOC entry 5346 (class 1259 OID 117329)
-- Name: idx_user_line_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_line_is_deleted ON public.user_line USING btree (is_deleted);


--
-- TOC entry 5347 (class 1259 OID 117327)
-- Name: idx_user_line_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_line_user_id ON public.user_line USING btree (user_id);


--
-- TOC entry 5348 (class 1259 OID 117328)
-- Name: idx_user_line_user_shop_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_line_user_shop_id ON public.user_line USING btree (user_shop_id);


--
-- TOC entry 5334 (class 1259 OID 117415)
-- Name: idx_user_location_branch_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_location_branch_id ON public.user_location USING btree (branch_id);


--
-- TOC entry 5335 (class 1259 OID 117269)
-- Name: idx_user_location_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_location_is_deleted ON public.user_location USING btree (is_deleted);


--
-- TOC entry 5336 (class 1259 OID 117267)
-- Name: idx_user_location_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_location_user_id ON public.user_location USING btree (user_id);


--
-- TOC entry 5368 (class 1259 OID 117703)
-- Name: idx_user_machine_access_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machine_access_machine_id ON public.user_machine_access USING btree (machine_id);


--
-- TOC entry 5245 (class 1259 OID 116730)
-- Name: idx_user_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_is_deleted ON public.user_machines USING btree (is_deleted);


--
-- TOC entry 5246 (class 1259 OID 116729)
-- Name: idx_user_machines_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_machine_id ON public.user_machines USING btree (machine_id);


--
-- TOC entry 5247 (class 1259 OID 116728)
-- Name: idx_user_machines_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_user_id ON public.user_machines USING btree (user_id);


--
-- TOC entry 5233 (class 1259 OID 116687)
-- Name: idx_user_profiles_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_profiles_is_deleted ON public.user_profiles USING btree (is_deleted);


--
-- TOC entry 5234 (class 1259 OID 116686)
-- Name: idx_user_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_profiles_user_id ON public.user_profiles USING btree (user_id);


--
-- TOC entry 5340 (class 1259 OID 117299)
-- Name: idx_user_shop_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_shop_is_deleted ON public.user_shop USING btree (is_deleted);


--
-- TOC entry 5341 (class 1259 OID 117297)
-- Name: idx_user_shop_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_shop_user_id ON public.user_shop USING btree (user_id);


--
-- TOC entry 5342 (class 1259 OID 117298)
-- Name: idx_user_shop_user_location_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_shop_user_location_id ON public.user_shop USING btree (user_location_id);


--
-- TOC entry 5322 (class 1259 OID 117179)
-- Name: idx_user_state_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_state_is_deleted ON public.user_state USING btree (is_deleted);


--
-- TOC entry 5323 (class 1259 OID 117178)
-- Name: idx_user_state_user_country_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_state_user_country_id ON public.user_state USING btree (user_country_id);


--
-- TOC entry 5324 (class 1259 OID 117177)
-- Name: idx_user_state_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_state_user_id ON public.user_state USING btree (user_id);


--
-- TOC entry 5358 (class 1259 OID 117389)
-- Name: idx_user_supervisor_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_supervisor_is_deleted ON public.user_supervisor USING btree (is_deleted);


--
-- TOC entry 5359 (class 1259 OID 117388)
-- Name: idx_user_supervisor_user_cell_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_supervisor_user_cell_id ON public.user_supervisor USING btree (user_cell_id);


--
-- TOC entry 5360 (class 1259 OID 117387)
-- Name: idx_user_supervisor_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_supervisor_user_id ON public.user_supervisor USING btree (user_id);


--
-- TOC entry 5220 (class 1259 OID 116617)
-- Name: idx_users_branch_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_branch_id ON public.users USING btree (branch_id);


--
-- TOC entry 5221 (class 1259 OID 116618)
-- Name: idx_users_company_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_company_id ON public.users USING btree (company_id);


--
-- TOC entry 5222 (class 1259 OID 116619)
-- Name: idx_users_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_is_deleted ON public.users USING btree (is_deleted);


--
-- TOC entry 5223 (class 1259 OID 116616)
-- Name: idx_users_role_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_role_id ON public.users USING btree (role_id);


--
-- TOC entry 5224 (class 1259 OID 116615)
-- Name: idx_users_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_status ON public.users USING btree (status);


--
-- TOC entry 5304 (class 1259 OID 117420)
-- Name: machines_user_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX machines_user_name_key ON public.machines USING btree (user_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5365 (class 1259 OID 117667)
-- Name: models_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX models_name_key ON public.models USING btree (controller_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5378 (class 1259 OID 117733)
-- Name: parts_user_part_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX parts_user_part_key ON public.parts USING btree (user_id, lower((part_no)::text), lower((part_name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5278 (class 1259 OID 117076)
-- Name: products_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX products_code_key ON public.products USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5279 (class 1259 OID 116954)
-- Name: products_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX products_name_key ON public.products USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5401 (class 1259 OID 117919)
-- Name: shifts_user_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX shifts_user_name_key ON public.shifts USING btree (user_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5355 (class 1259 OID 117356)
-- Name: user_cell_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_cell_name_key ON public.user_cell USING btree (user_id, user_line_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5331 (class 1259 OID 117206)
-- Name: user_city_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_city_name_key ON public.user_city USING btree (user_id, user_state_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5319 (class 1259 OID 117147)
-- Name: user_country_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_country_name_key ON public.user_country USING btree (user_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5316 (class 1259 OID 117123)
-- Name: user_filters_user_id_filter_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_filters_user_id_filter_name_key ON public.user_filters USING btree (user_id, filter_name) WHERE (is_deleted = false);


--
-- TOC entry 5349 (class 1259 OID 117326)
-- Name: user_line_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_line_name_key ON public.user_line USING btree (user_id, user_shop_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5337 (class 1259 OID 117414)
-- Name: user_location_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_location_name_key ON public.user_location USING btree (user_id, branch_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5343 (class 1259 OID 117296)
-- Name: user_shop_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_shop_name_key ON public.user_shop USING btree (user_id, user_location_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5325 (class 1259 OID 117176)
-- Name: user_state_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_state_name_key ON public.user_state USING btree (user_id, user_country_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5361 (class 1259 OID 117386)
-- Name: user_supervisor_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_supervisor_name_key ON public.user_supervisor USING btree (user_id, user_cell_id, lower((name)::text)) WHERE (is_deleted = false);


--
-- TOC entry 5225 (class 1259 OID 116613)
-- Name: users_email_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email) WHERE (is_deleted = false);


--
-- TOC entry 5228 (class 1259 OID 116614)
-- Name: users_username_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_username_key ON public.users USING btree (username) WHERE (is_deleted = false);


--
-- TOC entry 5467 (class 2620 OID 116756)
-- Name: addresses update_addresses_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON public.addresses FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5471 (class 2620 OID 117043)
-- Name: brands update_brands_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_brands_updated_at BEFORE UPDATE ON public.brands FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5492 (class 2620 OID 117983)
-- Name: breaks update_breaks_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_breaks_updated_at BEFORE UPDATE ON public.breaks FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5470 (class 2620 OID 117042)
-- Name: categories update_categories_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5490 (class 2620 OID 117873)
-- Name: clienteles update_clienteles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_clienteles_updated_at BEFORE UPDATE ON public.clienteles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5469 (class 2620 OID 116759)
-- Name: company_profiles update_company_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_company_profiles_updated_at BEFORE UPDATE ON public.company_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5473 (class 2620 OID 117045)
-- Name: controllers update_controllers_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_controllers_updated_at BEFORE UPDATE ON public.controllers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5493 (class 2620 OID 118037)
-- Name: holidays update_holidays_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_holidays_updated_at BEFORE UPDATE ON public.holidays FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5474 (class 2620 OID 117046)
-- Name: machines update_machines_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_machines_updated_at BEFORE UPDATE ON public.machines FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5484 (class 2620 OID 117671)
-- Name: models update_models_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_models_updated_at BEFORE UPDATE ON public.models FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5488 (class 2620 OID 117807)
-- Name: part_process_machines update_part_process_machines_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_part_process_machines_updated_at BEFORE UPDATE ON public.part_process_machines FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5487 (class 2620 OID 117781)
-- Name: part_processes update_part_processes_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_part_processes_updated_at BEFORE UPDATE ON public.part_processes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5489 (class 2620 OID 117834)
-- Name: part_sub_programs update_part_sub_programs_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_part_sub_programs_updated_at BEFORE UPDATE ON public.part_sub_programs FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5486 (class 2620 OID 117736)
-- Name: parts update_parts_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_parts_updated_at BEFORE UPDATE ON public.parts FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5472 (class 2620 OID 117044)
-- Name: products update_products_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON public.products FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5491 (class 2620 OID 117922)
-- Name: shifts update_shifts_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON public.shifts FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5482 (class 2620 OID 117360)
-- Name: user_cell update_user_cell_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_cell_updated_at BEFORE UPDATE ON public.user_cell FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5478 (class 2620 OID 117210)
-- Name: user_city update_user_city_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_city_updated_at BEFORE UPDATE ON public.user_city FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5476 (class 2620 OID 117150)
-- Name: user_country update_user_country_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_country_updated_at BEFORE UPDATE ON public.user_country FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5475 (class 2620 OID 117126)
-- Name: user_filters update_user_filters_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_filters_updated_at BEFORE UPDATE ON public.user_filters FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5481 (class 2620 OID 117330)
-- Name: user_line update_user_line_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_line_updated_at BEFORE UPDATE ON public.user_line FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5479 (class 2620 OID 117270)
-- Name: user_location update_user_location_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_location_updated_at BEFORE UPDATE ON public.user_location FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5485 (class 2620 OID 117704)
-- Name: user_machine_access update_user_machine_access_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_machine_access_updated_at BEFORE UPDATE ON public.user_machine_access FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5468 (class 2620 OID 116758)
-- Name: user_profiles update_user_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON public.user_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5480 (class 2620 OID 117300)
-- Name: user_shop update_user_shop_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_shop_updated_at BEFORE UPDATE ON public.user_shop FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5477 (class 2620 OID 117180)
-- Name: user_state update_user_state_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_state_updated_at BEFORE UPDATE ON public.user_state FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5483 (class 2620 OID 117390)
-- Name: user_supervisor update_user_supervisor_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_supervisor_updated_at BEFORE UPDATE ON public.user_supervisor FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5466 (class 2620 OID 116755)
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5417 (class 2606 OID 116645)
-- Name: addresses addresses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5422 (class 2606 OID 116746)
-- Name: auth_tokens auth_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5423 (class 2606 OID 116791)
-- Name: branches branches_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5424 (class 2606 OID 117393)
-- Name: branches branches_user_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_user_city_id_fkey FOREIGN KEY (user_city_id) REFERENCES public.user_city(id) ON DELETE SET NULL;


--
-- TOC entry 5463 (class 2606 OID 118044)
-- Name: breaks breaks_shift_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.breaks
    ADD CONSTRAINT breaks_shift_id_fkey FOREIGN KEY (shift_id) REFERENCES public.shifts(id);


--
-- TOC entry 5464 (class 2606 OID 117975)
-- Name: breaks breaks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.breaks
    ADD CONSTRAINT breaks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 5434 (class 2606 OID 117060)
-- Name: category_machines category_machines_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- TOC entry 5435 (class 2606 OID 117065)
-- Name: category_machines category_machines_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5461 (class 2606 OID 117865)
-- Name: clienteles clienteles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clienteles
    ADD CONSTRAINT clienteles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 5419 (class 2606 OID 116702)
-- Name: company_profiles company_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5427 (class 2606 OID 116940)
-- Name: controller_products controller_products_controller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_controller_id_fkey FOREIGN KEY (controller_id) REFERENCES public.controllers(id) ON DELETE CASCADE;


--
-- TOC entry 5428 (class 2606 OID 116945)
-- Name: controller_products controller_products_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- TOC entry 5426 (class 2606 OID 116922)
-- Name: controllers controllers_brand_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers
    ADD CONSTRAINT controllers_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brands(id) ON DELETE SET NULL;


--
-- TOC entry 5465 (class 2606 OID 118028)
-- Name: holidays holidays_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.holidays
    ADD CONSTRAINT holidays_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 5429 (class 2606 OID 117019)
-- Name: machines machines_brand_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brands(id) ON DELETE SET NULL;


--
-- TOC entry 5430 (class 2606 OID 117024)
-- Name: machines machines_controller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_controller_id_fkey FOREIGN KEY (controller_id) REFERENCES public.controllers(id) ON DELETE SET NULL;


--
-- TOC entry 5431 (class 2606 OID 117673)
-- Name: machines machines_model_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id) ON DELETE SET NULL;


--
-- TOC entry 5432 (class 2606 OID 117014)
-- Name: machines machines_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE SET NULL;


--
-- TOC entry 5433 (class 2606 OID 117009)
-- Name: machines machines_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5452 (class 2606 OID 117662)
-- Name: models models_controller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_controller_id_fkey FOREIGN KEY (controller_id) REFERENCES public.controllers(id) ON DELETE RESTRICT;


--
-- TOC entry 5458 (class 2606 OID 117801)
-- Name: part_process_machines part_process_machines_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_process_machines
    ADD CONSTRAINT part_process_machines_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5459 (class 2606 OID 117796)
-- Name: part_process_machines part_process_machines_process_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_process_machines
    ADD CONSTRAINT part_process_machines_process_id_fkey FOREIGN KEY (process_id) REFERENCES public.part_processes(id) ON DELETE CASCADE;


--
-- TOC entry 5457 (class 2606 OID 117774)
-- Name: part_processes part_processes_part_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_processes
    ADD CONSTRAINT part_processes_part_id_fkey FOREIGN KEY (part_id) REFERENCES public.parts(id) ON DELETE CASCADE;


--
-- TOC entry 5460 (class 2606 OID 117827)
-- Name: part_sub_programs part_sub_programs_process_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.part_sub_programs
    ADD CONSTRAINT part_sub_programs_process_id_fkey FOREIGN KEY (process_id) REFERENCES public.part_processes(id) ON DELETE CASCADE;


--
-- TOC entry 5455 (class 2606 OID 117874)
-- Name: parts parts_clientele_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parts
    ADD CONSTRAINT parts_clientele_id_fkey FOREIGN KEY (clientele_id) REFERENCES public.clienteles(id);


--
-- TOC entry 5456 (class 2606 OID 117728)
-- Name: parts parts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parts
    ADD CONSTRAINT parts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 5425 (class 2606 OID 116900)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- TOC entry 5462 (class 2606 OID 117914)
-- Name: shifts shifts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shifts
    ADD CONSTRAINT shifts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 5414 (class 2606 OID 116536)
-- Name: states states_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states
    ADD CONSTRAINT states_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.countries(id) ON DELETE CASCADE;


--
-- TOC entry 5448 (class 2606 OID 117346)
-- Name: user_cell user_cell_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cell
    ADD CONSTRAINT user_cell_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5449 (class 2606 OID 117351)
-- Name: user_cell user_cell_user_line_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cell
    ADD CONSTRAINT user_cell_user_line_id_fkey FOREIGN KEY (user_line_id) REFERENCES public.user_line(id) ON DELETE RESTRICT;


--
-- TOC entry 5440 (class 2606 OID 117196)
-- Name: user_city user_city_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_city
    ADD CONSTRAINT user_city_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5441 (class 2606 OID 117201)
-- Name: user_city user_city_user_state_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_city
    ADD CONSTRAINT user_city_user_state_id_fkey FOREIGN KEY (user_state_id) REFERENCES public.user_state(id) ON DELETE RESTRICT;


--
-- TOC entry 5437 (class 2606 OID 117142)
-- Name: user_country user_country_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_country
    ADD CONSTRAINT user_country_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5436 (class 2606 OID 117118)
-- Name: user_filters user_filters_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters
    ADD CONSTRAINT user_filters_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5446 (class 2606 OID 117316)
-- Name: user_line user_line_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_line
    ADD CONSTRAINT user_line_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5447 (class 2606 OID 117321)
-- Name: user_line user_line_user_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_line
    ADD CONSTRAINT user_line_user_shop_id_fkey FOREIGN KEY (user_shop_id) REFERENCES public.user_shop(id) ON DELETE RESTRICT;


--
-- TOC entry 5442 (class 2606 OID 117409)
-- Name: user_location user_location_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_location
    ADD CONSTRAINT user_location_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE RESTRICT;


--
-- TOC entry 5443 (class 2606 OID 117256)
-- Name: user_location user_location_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_location
    ADD CONSTRAINT user_location_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5453 (class 2606 OID 117698)
-- Name: user_machine_access user_machine_access_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machine_access
    ADD CONSTRAINT user_machine_access_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5454 (class 2606 OID 117693)
-- Name: user_machine_access user_machine_access_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machine_access
    ADD CONSTRAINT user_machine_access_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5420 (class 2606 OID 117029)
-- Name: user_machines user_machines_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5421 (class 2606 OID 116723)
-- Name: user_machines user_machines_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5418 (class 2606 OID 116681)
-- Name: user_profiles user_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5444 (class 2606 OID 117286)
-- Name: user_shop user_shop_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_shop
    ADD CONSTRAINT user_shop_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5445 (class 2606 OID 117291)
-- Name: user_shop user_shop_user_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_shop
    ADD CONSTRAINT user_shop_user_location_id_fkey FOREIGN KEY (user_location_id) REFERENCES public.user_location(id) ON DELETE RESTRICT;


--
-- TOC entry 5438 (class 2606 OID 117171)
-- Name: user_state user_state_user_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_state
    ADD CONSTRAINT user_state_user_country_id_fkey FOREIGN KEY (user_country_id) REFERENCES public.user_country(id) ON DELETE RESTRICT;


--
-- TOC entry 5439 (class 2606 OID 117166)
-- Name: user_state user_state_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_state
    ADD CONSTRAINT user_state_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5450 (class 2606 OID 117381)
-- Name: user_supervisor user_supervisor_user_cell_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_supervisor
    ADD CONSTRAINT user_supervisor_user_cell_id_fkey FOREIGN KEY (user_cell_id) REFERENCES public.user_cell(id) ON DELETE RESTRICT;


--
-- TOC entry 5451 (class 2606 OID 117376)
-- Name: user_supervisor user_supervisor_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_supervisor
    ADD CONSTRAINT user_supervisor_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5415 (class 2606 OID 116608)
-- Name: users users_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5416 (class 2606 OID 116598)
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE SET NULL;


-- Completed on 2026-10-15 12:11:20

--
-- PostgreSQL database dump complete
--

\unrestrict v1DMVH0erDHUDQG7KifHiGyRZ0ExeMfDqOTfZd3w8RIvRVhL6sLVBaalPdv3FZ4

