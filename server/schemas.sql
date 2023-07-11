--
-- PostgreSQL database dump
--

-- Dumped from database version 14.8
-- Dumped by pg_dump version 14.8

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: repositories; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.repositories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    full_name text NOT NULL,
    url text NOT NULL,
    user_id bigint
);


ALTER TABLE public.repositories OWNER TO username;

--
-- Name: repositories_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.repositories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.repositories_id_seq OWNER TO username;

--
-- Name: repositories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.repositories_id_seq OWNED BY public.repositories.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at text,
    vendor text NOT NULL,
    user_name text
);


ALTER TABLE public.users OWNER TO username;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO username;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: variables; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.variables (
    id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    nonce text NOT NULL,
    repository_id bigint
);


ALTER TABLE public.variables OWNER TO username;

--
-- Name: variables_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.variables_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variables_id_seq OWNER TO username;

--
-- Name: variables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.variables_id_seq OWNED BY public.variables.id;


--
-- Name: repositories id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories ALTER COLUMN id SET DEFAULT nextval('public.repositories_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: variables id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables ALTER COLUMN id SET DEFAULT nextval('public.variables_id_seq'::regclass);


--
-- Data for Name: repositories; Type: TABLE DATA; Schema: public; Owner: username
--

COPY public.repositories (id, created_at, full_name, url, user_id) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: username
--

COPY public.users (id, created_at, vendor, user_name) FROM stdin;
\.


--
-- Data for Name: variables; Type: TABLE DATA; Schema: public; Owner: username
--

COPY public.variables (id, created_at, updated_at, key, value, nonce, repository_id) FROM stdin;
\.


--
-- Name: repositories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: username
--

SELECT pg_catalog.setval('public.repositories_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: username
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: variables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: username
--

SELECT pg_catalog.setval('public.variables_id_seq', 1, false);


--
-- Name: repositories repositories_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT repositories_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: variables variables_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables
    ADD CONSTRAINT variables_pkey PRIMARY KEY (id);


--
-- Name: variables fk_repositories_variables; Type: FK CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables
    ADD CONSTRAINT fk_repositories_variables FOREIGN KEY (repository_id) REFERENCES public.repositories(id);


--
-- Name: repositories fk_users_repositories; Type: FK CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT fk_users_repositories FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 14.8
-- Dumped by pg_dump version 14.8

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: repositories; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.repositories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    full_name text NOT NULL,
    url text NOT NULL,
    user_id bigint
);


ALTER TABLE public.repositories OWNER TO username;

--
-- Name: repositories_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.repositories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.repositories_id_seq OWNER TO username;

--
-- Name: repositories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.repositories_id_seq OWNED BY public.repositories.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at text,
    vendor text NOT NULL,
    user_name text
);


ALTER TABLE public.users OWNER TO username;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO username;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: variables; Type: TABLE; Schema: public; Owner: username
--

CREATE TABLE public.variables (
    id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    nonce text NOT NULL,
    repository_id bigint
);


ALTER TABLE public.variables OWNER TO username;

--
-- Name: variables_id_seq; Type: SEQUENCE; Schema: public; Owner: username
--

CREATE SEQUENCE public.variables_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variables_id_seq OWNER TO username;

--
-- Name: variables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: username
--

ALTER SEQUENCE public.variables_id_seq OWNED BY public.variables.id;


--
-- Name: repositories id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories ALTER COLUMN id SET DEFAULT nextval('public.repositories_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: variables id; Type: DEFAULT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables ALTER COLUMN id SET DEFAULT nextval('public.variables_id_seq'::regclass);


--
-- Name: repositories repositories_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT repositories_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: variables variables_pkey; Type: CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables
    ADD CONSTRAINT variables_pkey PRIMARY KEY (id);


--
-- Name: variables fk_repositories_variables; Type: FK CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.variables
    ADD CONSTRAINT fk_repositories_variables FOREIGN KEY (repository_id) REFERENCES public.repositories(id);


--
-- Name: repositories fk_users_repositories; Type: FK CONSTRAINT; Schema: public; Owner: username
--

ALTER TABLE ONLY public.repositories
    ADD CONSTRAINT fk_users_repositories FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

