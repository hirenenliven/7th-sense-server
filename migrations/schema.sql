--
-- PostgreSQL database dump
--

\restrict t8PlDayqVn2eJzr0sMpgUkS9Dm1DfeuUBJxeii1ICX88BShFGN5pwJ2gzAvGeAK

-- Dumped from database version 18.6
-- Dumped by pg_dump version 18.6

-- Started on 2026-10-13 13:01:24

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
-- TOC entry 266 (class 1255 OID 116754)
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
-- TOC entry 228 (class 1259 OID 116630)
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
-- TOC entry 227 (class 1259 OID 116629)
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
-- TOC entry 5343 (class 0 OID 0)
-- Dependencies: 227
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- TOC entry 236 (class 1259 OID 116732)
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
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.auth_tokens OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 116731)
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
-- TOC entry 5344 (class 0 OID 0)
-- Dependencies: 235
-- Name: auth_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.auth_tokens_id_seq OWNED BY public.auth_tokens.id;


--
-- TOC entry 238 (class 1259 OID 116773)
-- Name: branches; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.branches (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100),
    address text,
    customer_id integer,
    city character varying(100),
    state_id integer,
    gst_number character varying(30),
    pan_number character varying(30),
    color_code character varying(10),
    status character varying(20) DEFAULT 'active'::character varying,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.branches OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 116772)
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
-- TOC entry 5345 (class 0 OID 0)
-- Dependencies: 237
-- Name: branches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.branches_id_seq OWNED BY public.branches.id;


--
-- TOC entry 242 (class 1259 OID 116866)
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
-- TOC entry 5346 (class 0 OID 0)
-- Dependencies: 242
-- Name: TABLE brands; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.brands IS 'Machine brands - Added in migration 005';


--
-- TOC entry 241 (class 1259 OID 116865)
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
-- TOC entry 5347 (class 0 OID 0)
-- Dependencies: 241
-- Name: brands_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.brands_id_seq OWNED BY public.brands.id;


--
-- TOC entry 240 (class 1259 OID 116849)
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
-- TOC entry 5348 (class 0 OID 0)
-- Dependencies: 240
-- Name: TABLE categories; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.categories IS 'Product categories - Added in migration 005';


--
-- TOC entry 239 (class 1259 OID 116848)
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
-- TOC entry 5349 (class 0 OID 0)
-- Dependencies: 239
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 252 (class 1259 OID 117048)
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
-- TOC entry 251 (class 1259 OID 117047)
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
-- TOC entry 5350 (class 0 OID 0)
-- Dependencies: 251
-- Name: category_machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.category_machines_id_seq OWNED BY public.category_machines.id;


--
-- TOC entry 232 (class 1259 OID 116689)
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
-- TOC entry 231 (class 1259 OID 116688)
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
-- TOC entry 5351 (class 0 OID 0)
-- Dependencies: 231
-- Name: company_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.company_profiles_id_seq OWNED BY public.company_profiles.id;


--
-- TOC entry 248 (class 1259 OID 116928)
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
-- TOC entry 5352 (class 0 OID 0)
-- Dependencies: 248
-- Name: TABLE controller_products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.controller_products IS 'Controller-Product relationships - Added in migration 005';


--
-- TOC entry 247 (class 1259 OID 116927)
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
-- TOC entry 5353 (class 0 OID 0)
-- Dependencies: 247
-- Name: controller_products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.controller_products_id_seq OWNED BY public.controller_products.id;


--
-- TOC entry 246 (class 1259 OID 116906)
-- Name: controllers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.controllers (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50),
    variant character varying(100),
    brand_id integer,
    description text,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.controllers OWNER TO postgres;

--
-- TOC entry 5354 (class 0 OID 0)
-- Dependencies: 246
-- Name: TABLE controllers; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.controllers IS 'Machine controllers - Added in migration 005';


--
-- TOC entry 245 (class 1259 OID 116905)
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
-- TOC entry 5355 (class 0 OID 0)
-- Dependencies: 245
-- Name: controllers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.controllers_id_seq OWNED BY public.controllers.id;


--
-- TOC entry 220 (class 1259 OID 116513)
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
-- TOC entry 219 (class 1259 OID 116512)
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
-- TOC entry 5356 (class 0 OID 0)
-- Dependencies: 219
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.countries_id_seq OWNED BY public.countries.id;


--
-- TOC entry 250 (class 1259 OID 116968)
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
    bill_date date NOT NULL,
    warranty_date date NOT NULL,
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
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.machines OWNER TO postgres;

--
-- TOC entry 5357 (class 0 OID 0)
-- Dependencies: 250
-- Name: TABLE machines; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.machines IS 'Machine inventory and details - Added in migration 002';


--
-- TOC entry 249 (class 1259 OID 116967)
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
-- TOC entry 5358 (class 0 OID 0)
-- Dependencies: 249
-- Name: machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.machines_id_seq OWNED BY public.machines.id;


--
-- TOC entry 244 (class 1259 OID 116883)
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
-- TOC entry 5359 (class 0 OID 0)
-- Dependencies: 244
-- Name: TABLE products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.products IS 'Machine products - Added in migration 005';


--
-- TOC entry 243 (class 1259 OID 116882)
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
-- TOC entry 5360 (class 0 OID 0)
-- Dependencies: 243
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 224 (class 1259 OID 116542)
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
-- TOC entry 223 (class 1259 OID 116541)
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
-- TOC entry 5361 (class 0 OID 0)
-- Dependencies: 223
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 222 (class 1259 OID 116527)
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
-- TOC entry 221 (class 1259 OID 116526)
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
-- TOC entry 5362 (class 0 OID 0)
-- Dependencies: 221
-- Name: states_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.states_id_seq OWNED BY public.states.id;


--
-- TOC entry 254 (class 1259 OID 117105)
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
-- TOC entry 253 (class 1259 OID 117104)
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
-- TOC entry 5363 (class 0 OID 0)
-- Dependencies: 253
-- Name: user_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_filters_id_seq OWNED BY public.user_filters.id;


--
-- TOC entry 234 (class 1259 OID 116710)
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
-- TOC entry 233 (class 1259 OID 116709)
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
-- TOC entry 5364 (class 0 OID 0)
-- Dependencies: 233
-- Name: user_machines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_machines_id_seq OWNED BY public.user_machines.id;


--
-- TOC entry 230 (class 1259 OID 116663)
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
    theme character varying(20) DEFAULT 'light_version'::character varying,
    timezone character varying(50) DEFAULT 'Asia/Kolkata'::character varying,
    language character varying(10) DEFAULT 'en'::character varying,
    start_day_time time without time zone,
    is_deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_profiles OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 116662)
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
-- TOC entry 5365 (class 0 OID 0)
-- Dependencies: 229
-- Name: user_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_profiles_id_seq OWNED BY public.user_profiles.id;


--
-- TOC entry 226 (class 1259 OID 116580)
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
-- TOC entry 225 (class 1259 OID 116579)
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
-- TOC entry 5366 (class 0 OID 0)
-- Dependencies: 225
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4953 (class 2604 OID 116633)
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- TOC entry 4973 (class 2604 OID 116735)
-- Name: auth_tokens id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens ALTER COLUMN id SET DEFAULT nextval('public.auth_tokens_id_seq'::regclass);


--
-- TOC entry 4975 (class 2604 OID 116776)
-- Name: branches id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches ALTER COLUMN id SET DEFAULT nextval('public.branches_id_seq'::regclass);


--
-- TOC entry 4984 (class 2604 OID 116869)
-- Name: brands id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.brands ALTER COLUMN id SET DEFAULT nextval('public.brands_id_seq'::regclass);


--
-- TOC entry 4980 (class 2604 OID 116852)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 5012 (class 2604 OID 117051)
-- Name: category_machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines ALTER COLUMN id SET DEFAULT nextval('public.category_machines_id_seq'::regclass);


--
-- TOC entry 4966 (class 2604 OID 116692)
-- Name: company_profiles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles ALTER COLUMN id SET DEFAULT nextval('public.company_profiles_id_seq'::regclass);


--
-- TOC entry 4997 (class 2604 OID 116931)
-- Name: controller_products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products ALTER COLUMN id SET DEFAULT nextval('public.controller_products_id_seq'::regclass);


--
-- TOC entry 4993 (class 2604 OID 116909)
-- Name: controllers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers ALTER COLUMN id SET DEFAULT nextval('public.controllers_id_seq'::regclass);


--
-- TOC entry 4942 (class 2604 OID 116516)
-- Name: countries id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries ALTER COLUMN id SET DEFAULT nextval('public.countries_id_seq'::regclass);


--
-- TOC entry 5000 (class 2604 OID 116971)
-- Name: machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines ALTER COLUMN id SET DEFAULT nextval('public.machines_id_seq'::regclass);


--
-- TOC entry 4988 (class 2604 OID 116886)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 4946 (class 2604 OID 116545)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 4944 (class 2604 OID 116530)
-- Name: states id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states ALTER COLUMN id SET DEFAULT nextval('public.states_id_seq'::regclass);


--
-- TOC entry 5015 (class 2604 OID 117108)
-- Name: user_filters id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters ALTER COLUMN id SET DEFAULT nextval('public.user_filters_id_seq'::regclass);


--
-- TOC entry 4970 (class 2604 OID 116713)
-- Name: user_machines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines ALTER COLUMN id SET DEFAULT nextval('public.user_machines_id_seq'::regclass);


--
-- TOC entry 4959 (class 2604 OID 116666)
-- Name: user_profiles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles ALTER COLUMN id SET DEFAULT nextval('public.user_profiles_id_seq'::regclass);


--
-- TOC entry 4948 (class 2604 OID 116583)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 5311 (class 0 OID 116630)
-- Dependencies: 228
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.addresses (id, user_id, address_type, address_line1, address_line2, city, pin_code, contact_person, contact_phone, is_default, is_deleted, deleted_at, created_at, updated_at, country_name, state_name) FROM stdin;
9	6	work	456 Business Park	\N	Ahmedabad	380001	Jane Smith	+919876543211	f	t	2026-10-13 11:20:45.582439+05:30	2026-10-13 11:20:45.575655	2026-10-13 12:08:04.49362	India	Gujarat
10	10	work	44915 Jackson Avenue	h0dudouifjoi sdoifud iou	Chino Hills	978987	Owen	+496880155856	t	t	2026-10-13 12:20:10.070292+05:30	2026-10-13 12:09:56.313649	2026-10-13 12:20:10.070292	India	Chhattisgarh
\.


--
-- TOC entry 5319 (class 0 OID 116732)
-- Dependencies: 236
-- Data for Name: auth_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.auth_tokens (id, user_id, token, token_type, device_type, device_info, expires_at, last_used_at, created_at) FROM stdin;
\.


--
-- TOC entry 5321 (class 0 OID 116773)
-- Dependencies: 238
-- Data for Name: branches; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.branches (id, code, name, email, address, customer_id, city, state_id, gst_number, pan_number, color_code, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
2	HO_0200	Metoda	metoda.service@jyoti.co.in	Plot No.,P-4, Lodhika GIDC, Vill.-Metoda, Dist.-Rajkot - 360 021.	6	Metoda	1	24AABCJ1947R1Z5	AABCJ1947R	#f39c12	active	f	\N	2026-10-12 12:45:04.069535	2026-10-12 12:45:04.069535
4	BR003223	Updated Branch Name	updated@example.com	456 New Address	6	Delhi	1	27AAPFU0939G1Z1	AAPFU0939G	#33FF57	active	f	2026-10-12 14:18:26.260462+05:30	2026-10-12 14:17:39.581263	2026-10-12 14:17:39.581263
1	METO123	RAJKOT-METODA	hello@gmail.com	Samrat Industrial Area Main Road, B/h S.T. Work Shop,Gondal Road, Rajkot - 360 004.	6	Rajkot	1	24AABCJ1947R1Z5	AABCJ1947R	#381E1A	active	f	\N	2026-10-12 12:45:04.069535	2026-10-12 12:45:04.069535
5	RK123	BR-234234	br23432@gmail.com	IN/110001, Ma Anandmayee Marg,, Okhla Phase III, Kalkaji	6	Jamnagar	\N	\N	CKLJF2394C	#D42525	active	f	\N	2026-10-13 12:21:44.84865	2026-10-13 12:21:44.84865
\.


--
-- TOC entry 5325 (class 0 OID 116866)
-- Dependencies: 242
-- Data for Name: brands; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.brands (id, name, code, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Brand 138251500 Updated	TC138251500	\N	t	2026-10-12 17:48:44.146918+05:30	2026-10-12 17:48:44.139025	2026-10-12 17:48:44.146918
4	Test Brand 138251500 Updated	TC138251500	\N	t	2026-10-12 17:48:44.148185+05:30	2026-10-12 17:48:44.147868	2026-10-12 17:48:44.148185
5	Test Brand Ref 162588400	\N	\N	t	2026-10-12 17:48:44.17353+05:30	2026-10-12 17:48:44.162812	2026-10-12 17:48:44.17353
6	Test Brand Ref 174094300	\N	\N	t	2026-10-12 17:48:44.175785+05:30	2026-10-12 17:48:44.174538	2026-10-12 17:48:44.175785
7	Test Brand 672113800 Updated	TC672113800	\N	t	2026-10-12 17:48:50.679916+05:30	2026-10-12 17:48:50.672868	2026-10-12 17:48:50.679916
10	Test Brand 672113800 Updated	TC672113800	\N	t	2026-10-12 17:48:50.681956+05:30	2026-10-12 17:48:50.681304	2026-10-12 17:48:50.681956
11	Test Brand Ref 697801600	\N	\N	t	2026-10-12 17:48:50.711368+05:30	2026-10-12 17:48:50.698191	2026-10-12 17:48:50.711368
12	Test Brand Ref 712340000	\N	\N	t	2026-10-12 17:48:50.714369+05:30	2026-10-12 17:48:50.712745	2026-10-12 17:48:50.714369
13	Test Brand 261232500 Updated	TC261232500	\N	t	2026-10-12 17:48:53.272033+05:30	2026-10-12 17:48:53.261771	2026-10-12 17:48:53.272033
16	Test Brand 261232500 Updated	TC261232500	\N	t	2026-10-12 17:48:53.273392+05:30	2026-10-12 17:48:53.273057	2026-10-12 17:48:53.273392
17	Test Brand Ref 285218700	\N	\N	t	2026-10-12 17:48:53.296425+05:30	2026-10-12 17:48:53.285774	2026-10-12 17:48:53.296425
18	Test Brand Ref 297068900	\N	\N	t	2026-10-12 17:48:53.298916+05:30	2026-10-12 17:48:53.297581	2026-10-12 17:48:53.298916
19	Test Brand 906849900 Updated	TC906849900	\N	t	2026-10-12 17:48:55.914131+05:30	2026-10-12 17:48:55.907655	2026-10-12 17:48:55.914131
22	Test Brand 906849900 Updated	TC906849900	\N	t	2026-10-12 17:48:55.915584+05:30	2026-10-12 17:48:55.915222	2026-10-12 17:48:55.915584
23	Test Brand Ref 928165900	\N	\N	t	2026-10-12 17:48:55.938048+05:30	2026-10-12 17:48:55.928399	2026-10-12 17:48:55.938048
24	Test Brand Ref 939142900	\N	\N	t	2026-10-12 17:48:55.942397+05:30	2026-10-12 17:48:55.939378	2026-10-12 17:48:55.942397
25	Test Brand 601793700 Updated	TC601793700	\N	t	2026-10-13 11:20:45.610841+05:30	2026-10-13 11:20:45.602418	2026-10-13 11:20:45.610841
28	Test Brand 601793700 Updated	TC601793700	\N	t	2026-10-13 11:20:45.613197+05:30	2026-10-13 11:20:45.612789	2026-10-13 11:20:45.613197
29	Test Brand Ref 626772900	\N	\N	t	2026-10-13 11:20:45.637674+05:30	2026-10-13 11:20:45.627099	2026-10-13 11:20:45.637674
30	Test Brand Ref 638294600	\N	\N	t	2026-10-13 11:20:45.640009+05:30	2026-10-13 11:20:45.638794	2026-10-13 11:20:45.640009
\.


--
-- TOC entry 5323 (class 0 OID 116849)
-- Dependencies: 240
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, code, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Category 121986100 Updated	TC121986100	\N	t	2026-10-12 17:48:44.136062+05:30	2026-10-12 17:48:44.124319	2026-10-12 17:48:44.136062
4	Test Category 121986100 Updated	TC121986100	\N	t	2026-10-12 17:48:44.138245+05:30	2026-10-12 17:48:44.137075	2026-10-12 17:48:44.138245
5	Test Category Ref 148465800	\N	\N	t	2026-10-12 17:48:44.160127+05:30	2026-10-12 17:48:44.148591	2026-10-12 17:48:44.160127
6	Test Category 662000300 Updated	TC662000300	\N	t	2026-10-12 17:48:50.67064+05:30	2026-10-12 17:48:50.663036	2026-10-12 17:48:50.67064
9	Test Category 662000300 Updated	TC662000300	\N	t	2026-10-12 17:48:50.672009+05:30	2026-10-12 17:48:50.671633	2026-10-12 17:48:50.672009
10	Test Category Ref 682313000	\N	\N	t	2026-10-12 17:48:50.695098+05:30	2026-10-12 17:48:50.682786	2026-10-12 17:48:50.695098
11	Test Category 250434800 Updated	TC250434800	\N	t	2026-10-12 17:48:53.259031+05:30	2026-10-12 17:48:53.251486	2026-10-12 17:48:53.259031
14	Test Category 250434800 Updated	TC250434800	\N	t	2026-10-12 17:48:53.260799+05:30	2026-10-12 17:48:53.260328	2026-10-12 17:48:53.260799
15	Test Category Ref 273755100	\N	\N	t	2026-10-12 17:48:53.2836+05:30	2026-10-12 17:48:53.27392	2026-10-12 17:48:53.2836
16	Test Category 895651600 Updated	TC895651600	\N	t	2026-10-12 17:48:55.90516+05:30	2026-10-12 17:48:55.896477	2026-10-12 17:48:55.90516
19	Test Category 895651600 Updated	TC895651600	\N	t	2026-10-12 17:48:55.906783+05:30	2026-10-12 17:48:55.906387	2026-10-12 17:48:55.906783
20	Test Category Ref 915794800	\N	\N	t	2026-10-12 17:48:55.926216+05:30	2026-10-12 17:48:55.91606	2026-10-12 17:48:55.926216
21	Test Category 582972500 Updated	TC582972500	\N	t	2026-10-13 11:20:45.600236+05:30	2026-10-13 11:20:45.5841	2026-10-13 11:20:45.600236
24	Test Category 582972500 Updated	TC582972500	\N	t	2026-10-13 11:20:45.601671+05:30	2026-10-13 11:20:45.60127	2026-10-13 11:20:45.601671
25	Test Category Ref 613390200	\N	\N	t	2026-10-13 11:20:45.625141+05:30	2026-10-13 11:20:45.613693	2026-10-13 11:20:45.625141
\.


--
-- TOC entry 5335 (class 0 OID 117048)
-- Dependencies: 252
-- Data for Name: category_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.category_machines (id, category_id, machine_id, is_deleted, deleted_at, created_at) FROM stdin;
\.


--
-- TOC entry 5315 (class 0 OID 116689)
-- Dependencies: 232
-- Data for Name: company_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.company_profiles (id, user_id, gst_number, pan_number, website, finix_id, is_deleted, deleted_at, created_at, updated_at, company_name) FROM stdin;
12	6	27AAACC1206D1ZN	\N	https://example.org	\N	t	2026-10-13 11:20:45.646281+05:30	2026-10-13 11:20:45.644422	2026-10-13 11:20:45.646281	\N
13	10	07ABCDE1234F1Z5	CHDOD9599C	www.goole.com	34343	t	2026-10-13 12:20:10.070292+05:30	2026-10-13 12:09:56.313649	2026-10-13 12:20:10.070292	Corkery - Update it
\.


--
-- TOC entry 5331 (class 0 OID 116928)
-- Dependencies: 248
-- Data for Name: controller_products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.controller_products (id, controller_id, product_id, is_deleted, deleted_at, created_at) FROM stdin;
\.


--
-- TOC entry 5329 (class 0 OID 116906)
-- Dependencies: 246
-- Data for Name: controllers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.controllers (id, name, code, variant, brand_id, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Controller 163130400 Updated	TC163130400	828D	5	\N	t	2026-10-12 17:48:44.171761+05:30	2026-10-12 17:48:44.163534	2026-10-12 17:48:44.171761
4	Test Controller 163130400 Updated	TC163130400	828D	5	\N	t	2026-10-12 17:48:44.173176+05:30	2026-10-12 17:48:44.172729	2026-10-12 17:48:44.173176
6	Test Controller 174606100	\N	840D	6	\N	t	2026-10-12 17:48:44.175375+05:30	2026-10-12 17:48:44.174976	2026-10-12 17:48:44.175375
7	Test Controller 698618000 Updated	TC698618000	828D	11	\N	t	2026-10-12 17:48:50.708379+05:30	2026-10-12 17:48:50.700372	2026-10-12 17:48:50.708379
10	Test Controller 698618000 Updated	TC698618000	828D	11	\N	t	2026-10-12 17:48:50.710898+05:30	2026-10-12 17:48:50.71021	2026-10-12 17:48:50.710898
12	Test Controller 712896300	\N	840D	12	\N	t	2026-10-12 17:48:50.713993+05:30	2026-10-12 17:48:50.713416	2026-10-12 17:48:50.713993
13	Test Controller 285756000 Updated	TC285756000	828D	17	\N	t	2026-10-12 17:48:53.29444+05:30	2026-10-12 17:48:53.286487	2026-10-12 17:48:53.29444
16	Test Controller 285756000 Updated	TC285756000	828D	17	\N	t	2026-10-12 17:48:53.296038+05:30	2026-10-12 17:48:53.295567	2026-10-12 17:48:53.296038
18	Test Controller 297585500	\N	840D	18	\N	t	2026-10-12 17:48:53.298551+05:30	2026-10-12 17:48:53.298112	2026-10-12 17:48:53.298551
19	Test Controller 928699000 Updated	TC928699000	828D	23	\N	t	2026-10-12 17:48:55.936011+05:30	2026-10-12 17:48:55.929205	2026-10-12 17:48:55.936011
22	Test Controller 928699000 Updated	TC928699000	828D	23	\N	t	2026-10-12 17:48:55.937541+05:30	2026-10-12 17:48:55.936972	2026-10-12 17:48:55.937541
24	Test Controller 940726800	\N	840D	24	\N	t	2026-10-12 17:48:55.941992+05:30	2026-10-12 17:48:55.941473	2026-10-12 17:48:55.941992
25	Test Controller 627318800 Updated	TC627318800	828D	29	\N	t	2026-10-13 11:20:45.635866+05:30	2026-10-13 11:20:45.627802	2026-10-13 11:20:45.635866
28	Test Controller 627318800 Updated	TC627318800	828D	29	\N	t	2026-10-13 11:20:45.63736+05:30	2026-10-13 11:20:45.636947	2026-10-13 11:20:45.63736
30	Test Controller 638819500	\N	840D	30	\N	t	2026-10-13 11:20:45.639696+05:30	2026-10-13 11:20:45.639241	2026-10-13 11:20:45.639696
\.


--
-- TOC entry 5303 (class 0 OID 116513)
-- Dependencies: 220
-- Data for Name: countries; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.countries (id, name, code, phone_code, created_at) FROM stdin;
1	India	IND	+91	2026-10-12 10:38:40.418244
2	United States	USA	+1	2026-10-12 10:38:40.418244
\.


--
-- TOC entry 5333 (class 0 OID 116968)
-- Dependencies: 250
-- Data for Name: machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.machines (id, user_id, machine_type_id, product_id, name, serial_number, live_ip, port, username, password, local_ip, brand_id, controller_id, controller_version, accuracy_id, line_id, installation_date, bill_date, warranty_date, amc, is_insured, online_status, maintenance_required, hourly_cost, status, is_monitor, online_date, offline_date, mode_id, controller_version_type, health_ids, longitude, latitude, cycle_analysis, machine_system_type, security_mode, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
4	6	1	1	CNC Machine 01	SN-2024-001	192.168.1.100	8080	admin_user	securePass123	192.168.1.101	1	1	v2.1.0	1	1	2024-01-15	2024-01-10	2026-01-15	Annual	t	t	f	50.00	active	t	2024-01-20	\N	1	2	1,2,3	-73.9857	40.7484	t	NewPlc	Basic256Sha256	f	\N	2026-10-13 10:49:49.205738	2026-10-13 10:49:49.205738
6	6	1	1	CNC Machine 01 Updated	SN-2024-001	192.168.1.100	8080	admin_user	newSecurePass456	192.168.1.101	1	1	v2.2.0	1	1	2024-01-15	2024-01-10	2026-01-15	Annual	t	t	t	55.00	inactive	t	2024-01-20	\N	1	2	1,2,3,4	-73.9857	40.7484	t	NewPlc	Basic256Sha256	f	2026-10-13 10:56:28.283508+05:30	2026-10-13 10:55:22.027533	2026-10-13 10:56:53.157488
7	\N	\N	\N	Machine-647281800-Updated	SN-647281800	192.168.1.102	8081	admin2	newpassword	192.168.1.103	\N	\N	2.0.0	\N	\N	2024-02-01	2024-02-01	2025-02-01	\N	f	f	f	150.75	inactive	f	\N	\N	\N	\N	\N	\N	\N	f	OldPlc	None	t	2026-10-13 11:20:45.660061+05:30	2026-10-13 11:20:45.648327	2026-10-13 11:20:45.660061
9	\N	\N	\N	Machine-647281800	SN-647281800-REUSE	192.168.1.100	8080	admin	password123	192.168.1.101	\N	\N	1.0.0	\N	\N	2024-01-01	2024-01-01	2025-01-01		f	f	f	100.50	active	f	\N	\N	\N	\N				f	NewPlc	Basic256Sha256	t	2026-10-13 11:20:45.662094+05:30	2026-10-13 11:20:45.66155	2026-10-13 11:20:45.662094
\.


--
-- TOC entry 5327 (class 0 OID 116883)
-- Dependencies: 244
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, code, category_id, image, description, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
1	Test Product 148465800 Updated	TC148465800	5	lathe.png	\N	t	2026-10-12 17:48:44.157782+05:30	2026-10-12 17:48:44.149268	2026-10-12 17:48:44.157782
4	Test Product 148465800 Updated	TC148465800	5	lathe.png	\N	t	2026-10-12 17:48:44.15966+05:30	2026-10-12 17:48:44.159117	2026-10-12 17:48:44.15966
6	Test Product 161014200	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:44.162273+05:30	2026-10-12 17:48:44.161869	2026-10-12 17:48:44.162273
7	Test Product 682313000 Updated	TC682313000	10	lathe.png	\N	t	2026-10-12 17:48:50.692503+05:30	2026-10-12 17:48:50.683849	2026-10-12 17:48:50.692503
10	Test Product 682313000 Updated	TC682313000	10	lathe.png	\N	t	2026-10-12 17:48:50.694535+05:30	2026-10-12 17:48:50.693926	2026-10-12 17:48:50.694535
12	Test Product 696127300	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:50.697607+05:30	2026-10-12 17:48:50.696914	2026-10-12 17:48:50.697607
13	Test Product 274269500 Updated	TC274269500	15	lathe.png	\N	t	2026-10-12 17:48:53.281889+05:30	2026-10-12 17:48:53.274804	2026-10-12 17:48:53.281889
16	Test Product 274269500 Updated	TC274269500	15	lathe.png	\N	t	2026-10-12 17:48:53.2832+05:30	2026-10-12 17:48:53.282777	2026-10-12 17:48:53.2832
18	Test Product 284189700	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:53.285268+05:30	2026-10-12 17:48:53.284917	2026-10-12 17:48:53.285268
19	Test Product 916338700 Updated	TC916338700	20	lathe.png	\N	t	2026-10-12 17:48:55.924207+05:30	2026-10-12 17:48:55.916769	2026-10-12 17:48:55.924207
22	Test Product 916338700 Updated	TC916338700	20	lathe.png	\N	t	2026-10-12 17:48:55.925796+05:30	2026-10-12 17:48:55.925295	2026-10-12 17:48:55.925796
24	Test Product 927046000	\N	\N	no_image_available.png	\N	t	2026-10-12 17:48:55.927792+05:30	2026-10-12 17:48:55.927456	2026-10-12 17:48:55.927792
25	Test Product 613950100 Updated	TC613950100	25	lathe.png	\N	t	2026-10-13 11:20:45.623118+05:30	2026-10-13 11:20:45.614442	2026-10-13 11:20:45.623118
28	Test Product 613950100 Updated	TC613950100	25	lathe.png	\N	t	2026-10-13 11:20:45.624731+05:30	2026-10-13 11:20:45.624358	2026-10-13 11:20:45.624731
30	Test Product 625742900	\N	\N	no_image_available.png	\N	t	2026-10-13 11:20:45.62664+05:30	2026-10-13 11:20:45.626309	2026-10-13 11:20:45.62664
\.


--
-- TOC entry 5307 (class 0 OID 116542)
-- Dependencies: 224
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, description, created_at) FROM stdin;
1	admin	System administrator with full access	2026-10-12 10:38:40.418244
2	supervisor	Supervisor with management permissions	2026-10-12 10:38:40.418244
3	operator	Machine operator	2026-10-12 10:38:40.418244
4	customer	Customer with limited access	2026-10-12 10:38:40.418244
5	qc	Quality control staff	2026-10-12 10:38:40.418244
6	shop_head	Shop head/manager	2026-10-12 10:38:40.418244
\.


--
-- TOC entry 5305 (class 0 OID 116527)
-- Dependencies: 222
-- Data for Name: states; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.states (id, country_id, name, code, created_at) FROM stdin;
1	1	Gujarat	GJ	2026-10-12 10:38:40.418244
2	1	Maharashtra	MH	2026-10-12 10:38:40.418244
3	1	Karnataka	KA	2026-10-12 10:38:40.418244
4	2	California	CA	2026-10-12 10:38:40.418244
\.


--
-- TOC entry 5337 (class 0 OID 117105)
-- Dependencies: 254
-- Data for Name: user_filters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_filters (id, user_id, filter_name, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 5317 (class 0 OID 116710)
-- Dependencies: 234
-- Data for Name: user_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_machines (id, user_id, machine_id, assigned_at, is_deleted, deleted_at) FROM stdin;
\.


--
-- TOC entry 5313 (class 0 OID 116663)
-- Dependencies: 230
-- Data for Name: user_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_profiles (id, user_id, designation, grade, birth_date, joining_date, blood_group, profile_picture, contact_person, contact_phone, theme, timezone, language, start_day_time, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
11	6	Senior Software Engineer	G3	1995-03-14	\N	O+	\N	\N	\N	light_version	Asia/Kolkata	hi	\N	t	2026-10-13 11:20:45.668883+05:30	2026-10-13 11:20:45.66667	2026-10-13 11:20:45.668883
12	10	\N	\N	\N	\N	\N	\N	Owen	+49408218573	light_version	Asia/Kabul	en	11:15:00	t	2026-10-13 12:20:10.070292+05:30	2026-10-13 12:09:56.313649	2026-10-13 12:20:10.070292
\.


--
-- TOC entry 5309 (class 0 OID 116580)
-- Dependencies: 226
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password, first_name, last_name, mobile, role_id, branch_id, company_id, status, is_deleted, deleted_at, created_at, updated_at) FROM stdin;
6	hiren	hiren11@example.com	$2a$10$kqgtdBjBgFusN8rOaBU9cuL7T4PvsZWvX1tlBdvaOjUFRdZAryWvG	hiren	chavda	+919876543210	1	\N	\N	active	f	2026-10-12 12:26:33.775246+05:30	2026-10-12 12:23:43.959479	2026-10-12 12:26:53.827569
10	Hassie_Price88	update@gmail.com	$2a$10$M83hZnn2jvTqYolzVaIixOGgnAcZVvtH5aDsHNxhoY8LkMo9KmgV6	Owen	\N	+496880155856	4	1	\N	active	f	2026-10-13 12:20:10.070292+05:30	2026-10-13 12:09:56.313649	2026-10-13 12:25:11.310404
\.


--
-- TOC entry 5367 (class 0 OID 0)
-- Dependencies: 227
-- Name: addresses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.addresses_id_seq', 11, true);


--
-- TOC entry 5368 (class 0 OID 0)
-- Dependencies: 235
-- Name: auth_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.auth_tokens_id_seq', 1, false);


--
-- TOC entry 5369 (class 0 OID 0)
-- Dependencies: 237
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.branches_id_seq', 5, true);


--
-- TOC entry 5370 (class 0 OID 0)
-- Dependencies: 241
-- Name: brands_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.brands_id_seq', 30, true);


--
-- TOC entry 5371 (class 0 OID 0)
-- Dependencies: 239
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 25, true);


--
-- TOC entry 5372 (class 0 OID 0)
-- Dependencies: 251
-- Name: category_machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.category_machines_id_seq', 1, false);


--
-- TOC entry 5373 (class 0 OID 0)
-- Dependencies: 231
-- Name: company_profiles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.company_profiles_id_seq', 14, true);


--
-- TOC entry 5374 (class 0 OID 0)
-- Dependencies: 247
-- Name: controller_products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.controller_products_id_seq', 1, false);


--
-- TOC entry 5375 (class 0 OID 0)
-- Dependencies: 245
-- Name: controllers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.controllers_id_seq', 30, true);


--
-- TOC entry 5376 (class 0 OID 0)
-- Dependencies: 219
-- Name: countries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.countries_id_seq', 246, true);


--
-- TOC entry 5377 (class 0 OID 0)
-- Dependencies: 249
-- Name: machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.machines_id_seq', 9, true);


--
-- TOC entry 5378 (class 0 OID 0)
-- Dependencies: 243
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 30, true);


--
-- TOC entry 5379 (class 0 OID 0)
-- Dependencies: 223
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 6, true);


--
-- TOC entry 5380 (class 0 OID 0)
-- Dependencies: 221
-- Name: states_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.states_id_seq', 4124, true);


--
-- TOC entry 5381 (class 0 OID 0)
-- Dependencies: 253
-- Name: user_filters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_filters_id_seq', 6, true);


--
-- TOC entry 5382 (class 0 OID 0)
-- Dependencies: 233
-- Name: user_machines_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_machines_id_seq', 1, false);


--
-- TOC entry 5383 (class 0 OID 0)
-- Dependencies: 229
-- Name: user_profiles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_profiles_id_seq', 13, true);


--
-- TOC entry 5384 (class 0 OID 0)
-- Dependencies: 225
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- TOC entry 5039 (class 2606 OID 116644)
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- TOC entry 5062 (class 2606 OID 116743)
-- Name: auth_tokens auth_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 5064 (class 2606 OID 116745)
-- Name: auth_tokens auth_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_token_key UNIQUE (token);


--
-- TOC entry 5069 (class 2606 OID 116790)
-- Name: branches branches_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_code_key UNIQUE (code);


--
-- TOC entry 5071 (class 2606 OID 116788)
-- Name: branches branches_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_pkey PRIMARY KEY (id);


--
-- TOC entry 5080 (class 2606 OID 116879)
-- Name: brands brands_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (id);


--
-- TOC entry 5075 (class 2606 OID 116862)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 5112 (class 2606 OID 117059)
-- Name: category_machines category_machines_category_id_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_category_id_machine_id_key UNIQUE (category_id, machine_id);


--
-- TOC entry 5114 (class 2606 OID 117057)
-- Name: category_machines category_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5049 (class 2606 OID 116699)
-- Name: company_profiles company_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_pkey PRIMARY KEY (id);


--
-- TOC entry 5051 (class 2606 OID 116701)
-- Name: company_profiles company_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_user_id_key UNIQUE (user_id);


--
-- TOC entry 5095 (class 2606 OID 116939)
-- Name: controller_products controller_products_controller_id_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_controller_id_product_id_key UNIQUE (controller_id, product_id);


--
-- TOC entry 5097 (class 2606 OID 116937)
-- Name: controller_products controller_products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_pkey PRIMARY KEY (id);


--
-- TOC entry 5091 (class 2606 OID 116919)
-- Name: controllers controllers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers
    ADD CONSTRAINT controllers_pkey PRIMARY KEY (id);


--
-- TOC entry 5020 (class 2606 OID 116525)
-- Name: countries countries_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_code_key UNIQUE (code);


--
-- TOC entry 5022 (class 2606 OID 116523)
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- TOC entry 5110 (class 2606 OID 117008)
-- Name: machines machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5087 (class 2606 OID 116897)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 5026 (class 2606 OID 116554)
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- TOC entry 5028 (class 2606 OID 116552)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 5024 (class 2606 OID 116535)
-- Name: states states_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states
    ADD CONSTRAINT states_pkey PRIMARY KEY (id);


--
-- TOC entry 5121 (class 2606 OID 117117)
-- Name: user_filters user_filters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters
    ADD CONSTRAINT user_filters_pkey PRIMARY KEY (id);


--
-- TOC entry 5058 (class 2606 OID 116720)
-- Name: user_machines user_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_pkey PRIMARY KEY (id);


--
-- TOC entry 5060 (class 2606 OID 116722)
-- Name: user_machines user_machines_user_id_machine_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_user_id_machine_id_key UNIQUE (user_id, machine_id);


--
-- TOC entry 5045 (class 2606 OID 116678)
-- Name: user_profiles user_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_pkey PRIMARY KEY (id);


--
-- TOC entry 5047 (class 2606 OID 116680)
-- Name: user_profiles user_profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_user_id_key UNIQUE (user_id);


--
-- TOC entry 5036 (class 2606 OID 116597)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 5077 (class 1259 OID 117075)
-- Name: brands_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX brands_code_key ON public.brands USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5078 (class 1259 OID 116952)
-- Name: brands_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX brands_name_key ON public.brands USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5072 (class 1259 OID 117074)
-- Name: categories_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX categories_code_key ON public.categories USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5073 (class 1259 OID 116950)
-- Name: categories_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX categories_name_key ON public.categories USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5088 (class 1259 OID 117077)
-- Name: controllers_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX controllers_code_key ON public.controllers USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5089 (class 1259 OID 116957)
-- Name: controllers_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX controllers_name_key ON public.controllers USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5040 (class 1259 OID 116661)
-- Name: idx_addresses_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_addresses_is_deleted ON public.addresses USING btree (is_deleted);


--
-- TOC entry 5041 (class 1259 OID 116660)
-- Name: idx_addresses_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_addresses_user_id ON public.addresses USING btree (user_id);


--
-- TOC entry 5065 (class 1259 OID 116753)
-- Name: idx_auth_tokens_expires_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_expires_at ON public.auth_tokens USING btree (expires_at);


--
-- TOC entry 5066 (class 1259 OID 116752)
-- Name: idx_auth_tokens_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_token ON public.auth_tokens USING btree (token);


--
-- TOC entry 5067 (class 1259 OID 116751)
-- Name: idx_auth_tokens_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_auth_tokens_user_id ON public.auth_tokens USING btree (user_id);


--
-- TOC entry 5081 (class 1259 OID 116953)
-- Name: idx_brands_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_brands_is_deleted ON public.brands USING btree (is_deleted);


--
-- TOC entry 5076 (class 1259 OID 116951)
-- Name: idx_categories_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categories_is_deleted ON public.categories USING btree (is_deleted);


--
-- TOC entry 5115 (class 1259 OID 117070)
-- Name: idx_category_machines_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_category_id ON public.category_machines USING btree (category_id);


--
-- TOC entry 5116 (class 1259 OID 117072)
-- Name: idx_category_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_is_deleted ON public.category_machines USING btree (is_deleted);


--
-- TOC entry 5117 (class 1259 OID 117071)
-- Name: idx_category_machines_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_category_machines_machine_id ON public.category_machines USING btree (machine_id);


--
-- TOC entry 5052 (class 1259 OID 116708)
-- Name: idx_company_profiles_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_company_profiles_is_deleted ON public.company_profiles USING btree (is_deleted);


--
-- TOC entry 5053 (class 1259 OID 116707)
-- Name: idx_company_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_company_profiles_user_id ON public.company_profiles USING btree (user_id);


--
-- TOC entry 5098 (class 1259 OID 116960)
-- Name: idx_controller_products_controller_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_controller_id ON public.controller_products USING btree (controller_id);


--
-- TOC entry 5099 (class 1259 OID 116962)
-- Name: idx_controller_products_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_is_deleted ON public.controller_products USING btree (is_deleted);


--
-- TOC entry 5100 (class 1259 OID 116961)
-- Name: idx_controller_products_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controller_products_product_id ON public.controller_products USING btree (product_id);


--
-- TOC entry 5092 (class 1259 OID 116958)
-- Name: idx_controllers_brand_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controllers_brand_id ON public.controllers USING btree (brand_id);


--
-- TOC entry 5093 (class 1259 OID 116959)
-- Name: idx_controllers_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_controllers_is_deleted ON public.controllers USING btree (is_deleted);


--
-- TOC entry 5101 (class 1259 OID 117037)
-- Name: idx_machines_brand_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_brand_id ON public.machines USING btree (brand_id);


--
-- TOC entry 5102 (class 1259 OID 117038)
-- Name: idx_machines_controller_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_controller_id ON public.machines USING btree (controller_id);


--
-- TOC entry 5103 (class 1259 OID 117041)
-- Name: idx_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_is_deleted ON public.machines USING btree (is_deleted);


--
-- TOC entry 5104 (class 1259 OID 117039)
-- Name: idx_machines_online_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_online_status ON public.machines USING btree (online_status);


--
-- TOC entry 5105 (class 1259 OID 117036)
-- Name: idx_machines_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_product_id ON public.machines USING btree (product_id);


--
-- TOC entry 5106 (class 1259 OID 117040)
-- Name: idx_machines_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_status ON public.machines USING btree (status);


--
-- TOC entry 5107 (class 1259 OID 117035)
-- Name: idx_machines_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_machines_user_id ON public.machines USING btree (user_id);


--
-- TOC entry 5082 (class 1259 OID 116955)
-- Name: idx_products_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_category_id ON public.products USING btree (category_id);


--
-- TOC entry 5083 (class 1259 OID 116956)
-- Name: idx_products_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_is_deleted ON public.products USING btree (is_deleted);


--
-- TOC entry 5118 (class 1259 OID 117125)
-- Name: idx_user_filters_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_filters_is_deleted ON public.user_filters USING btree (is_deleted);


--
-- TOC entry 5119 (class 1259 OID 117124)
-- Name: idx_user_filters_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_filters_user_id ON public.user_filters USING btree (user_id);


--
-- TOC entry 5054 (class 1259 OID 116730)
-- Name: idx_user_machines_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_is_deleted ON public.user_machines USING btree (is_deleted);


--
-- TOC entry 5055 (class 1259 OID 116729)
-- Name: idx_user_machines_machine_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_machine_id ON public.user_machines USING btree (machine_id);


--
-- TOC entry 5056 (class 1259 OID 116728)
-- Name: idx_user_machines_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_machines_user_id ON public.user_machines USING btree (user_id);


--
-- TOC entry 5042 (class 1259 OID 116687)
-- Name: idx_user_profiles_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_profiles_is_deleted ON public.user_profiles USING btree (is_deleted);


--
-- TOC entry 5043 (class 1259 OID 116686)
-- Name: idx_user_profiles_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_profiles_user_id ON public.user_profiles USING btree (user_id);


--
-- TOC entry 5029 (class 1259 OID 116617)
-- Name: idx_users_branch_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_branch_id ON public.users USING btree (branch_id);


--
-- TOC entry 5030 (class 1259 OID 116618)
-- Name: idx_users_company_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_company_id ON public.users USING btree (company_id);


--
-- TOC entry 5031 (class 1259 OID 116619)
-- Name: idx_users_is_deleted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_is_deleted ON public.users USING btree (is_deleted);


--
-- TOC entry 5032 (class 1259 OID 116616)
-- Name: idx_users_role_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_role_id ON public.users USING btree (role_id);


--
-- TOC entry 5033 (class 1259 OID 116615)
-- Name: idx_users_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_status ON public.users USING btree (status);


--
-- TOC entry 5108 (class 1259 OID 117034)
-- Name: machines_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX machines_name_key ON public.machines USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5084 (class 1259 OID 117076)
-- Name: products_code_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX products_code_key ON public.products USING btree (code) WHERE (is_deleted = false);


--
-- TOC entry 5085 (class 1259 OID 116954)
-- Name: products_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX products_name_key ON public.products USING btree (name) WHERE (is_deleted = false);


--
-- TOC entry 5122 (class 1259 OID 117123)
-- Name: user_filters_user_id_filter_name_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_filters_user_id_filter_name_key ON public.user_filters USING btree (user_id, filter_name) WHERE (is_deleted = false);


--
-- TOC entry 5034 (class 1259 OID 116613)
-- Name: users_email_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email) WHERE (is_deleted = false);


--
-- TOC entry 5037 (class 1259 OID 116614)
-- Name: users_username_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_username_key ON public.users USING btree (username) WHERE (is_deleted = false);


--
-- TOC entry 5146 (class 2620 OID 116756)
-- Name: addresses update_addresses_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON public.addresses FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5150 (class 2620 OID 117043)
-- Name: brands update_brands_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_brands_updated_at BEFORE UPDATE ON public.brands FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5149 (class 2620 OID 117042)
-- Name: categories update_categories_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5148 (class 2620 OID 116759)
-- Name: company_profiles update_company_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_company_profiles_updated_at BEFORE UPDATE ON public.company_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5152 (class 2620 OID 117045)
-- Name: controllers update_controllers_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_controllers_updated_at BEFORE UPDATE ON public.controllers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5153 (class 2620 OID 117046)
-- Name: machines update_machines_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_machines_updated_at BEFORE UPDATE ON public.machines FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5151 (class 2620 OID 117044)
-- Name: products update_products_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON public.products FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5154 (class 2620 OID 117126)
-- Name: user_filters update_user_filters_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_filters_updated_at BEFORE UPDATE ON public.user_filters FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5147 (class 2620 OID 116758)
-- Name: user_profiles update_user_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON public.user_profiles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5145 (class 2620 OID 116755)
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 5126 (class 2606 OID 116645)
-- Name: addresses addresses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5131 (class 2606 OID 116746)
-- Name: auth_tokens auth_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auth_tokens
    ADD CONSTRAINT auth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5132 (class 2606 OID 116791)
-- Name: branches branches_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5133 (class 2606 OID 116796)
-- Name: branches branches_state_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_state_id_fkey FOREIGN KEY (state_id) REFERENCES public.states(id) ON DELETE SET NULL;


--
-- TOC entry 5142 (class 2606 OID 117060)
-- Name: category_machines category_machines_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- TOC entry 5143 (class 2606 OID 117065)
-- Name: category_machines category_machines_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.category_machines
    ADD CONSTRAINT category_machines_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5128 (class 2606 OID 116702)
-- Name: company_profiles company_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5136 (class 2606 OID 116940)
-- Name: controller_products controller_products_controller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_controller_id_fkey FOREIGN KEY (controller_id) REFERENCES public.controllers(id) ON DELETE CASCADE;


--
-- TOC entry 5137 (class 2606 OID 116945)
-- Name: controller_products controller_products_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controller_products
    ADD CONSTRAINT controller_products_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- TOC entry 5135 (class 2606 OID 116922)
-- Name: controllers controllers_brand_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.controllers
    ADD CONSTRAINT controllers_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brands(id) ON DELETE SET NULL;


--
-- TOC entry 5138 (class 2606 OID 117019)
-- Name: machines machines_brand_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brands(id) ON DELETE SET NULL;


--
-- TOC entry 5139 (class 2606 OID 117024)
-- Name: machines machines_controller_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_controller_id_fkey FOREIGN KEY (controller_id) REFERENCES public.controllers(id) ON DELETE SET NULL;


--
-- TOC entry 5140 (class 2606 OID 117014)
-- Name: machines machines_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE SET NULL;


--
-- TOC entry 5141 (class 2606 OID 117009)
-- Name: machines machines_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.machines
    ADD CONSTRAINT machines_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5134 (class 2606 OID 116900)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- TOC entry 5123 (class 2606 OID 116536)
-- Name: states states_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states
    ADD CONSTRAINT states_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.countries(id) ON DELETE CASCADE;


--
-- TOC entry 5144 (class 2606 OID 117118)
-- Name: user_filters user_filters_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_filters
    ADD CONSTRAINT user_filters_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5129 (class 2606 OID 117029)
-- Name: user_machines user_machines_machine_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_machine_id_fkey FOREIGN KEY (machine_id) REFERENCES public.machines(id) ON DELETE CASCADE;


--
-- TOC entry 5130 (class 2606 OID 116723)
-- Name: user_machines user_machines_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_machines
    ADD CONSTRAINT user_machines_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5127 (class 2606 OID 116681)
-- Name: user_profiles user_profiles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 5124 (class 2606 OID 116608)
-- Name: users users_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 5125 (class 2606 OID 116598)
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE SET NULL;


-- Completed on 2026-10-13 13:01:24

--
-- PostgreSQL database dump complete
--

\unrestrict t8PlDayqVn2eJzr0sMpgUkS9Dm1DfeuUBJxeii1ICX88BShFGN5pwJ2gzAvGeAK

