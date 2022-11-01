--
-- PostgreSQL database dump
--

-- Dumped from database version 15.0 (Debian 15.0-1.pgdg110+1)
-- Dumped by pg_dump version 15.0 (Debian 15.0-1.pgdg110+1)

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
-- Name: accrual; Type: TABLE; Schema: public; Owner: yandex
--

CREATE TABLE public.accrual (
                                login text NOT NULL,
                                points numeric NOT NULL,
                                withdrawn numeric
);


ALTER TABLE public.accrual OWNER TO yandex;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: yandex
--

CREATE TABLE public.orders (
                               orderid bigint NOT NULL,
                               login text NOT NULL,
                               pointsspent boolean DEFAULT false,
                               orderdate timestamp without time zone,
                               accrual numeric,
                               status text
);


ALTER TABLE public.orders OWNER TO yandex;

--
-- Name: users; Type: TABLE; Schema: public; Owner: yandex
--

CREATE TABLE public.users (
                              login text NOT NULL,
                              passwordhash text NOT NULL,
                              passwordsalt text NOT NULL
);


ALTER TABLE public.users OWNER TO yandex;

--
-- Data for Name: accrual; Type: TABLE DATA; Schema: public; Owner: yandex
--

COPY public.accrual (login, points, withdrawn) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: yandex
--

COPY public.orders (orderid, login, pointsspent, orderdate, accrual, status) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: yandex
--

COPY public.users (login, passwordhash, passwordsalt) FROM stdin;
\.


--
-- Name: accrual accrual_pkey; Type: CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.accrual
    ADD CONSTRAINT accrual_pkey PRIMARY KEY (login);


--
-- Name: orders orders_orderid_key; Type: CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_orderid_key UNIQUE (orderid);


--
-- Name: users users_passwordsalt_key; Type: CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_passwordsalt_key UNIQUE (passwordsalt);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (login);


--
-- Name: accrual accrual_login_fkey; Type: FK CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.accrual
    ADD CONSTRAINT accrual_login_fkey FOREIGN KEY (login) REFERENCES public.users(login);


--
-- Name: orders orders_login_fkey; Type: FK CONSTRAINT; Schema: public; Owner: yandex
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_login_fkey FOREIGN KEY (login) REFERENCES public.users(login);


--
-- PostgreSQL database dump complete
--