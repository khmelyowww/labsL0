--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.5

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
-- Name: delivery; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.delivery (
    order_uid character varying(50) NOT NULL,
    name character varying(100),
    phone character varying(20),
    zip character varying(20),
    city character varying(50),
    address text,
    region character varying(50),
    email character varying(100)
);


ALTER TABLE public.delivery OWNER TO postgres;

--
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    rid character varying(50) NOT NULL,
    order_uid character varying(50),
    chrt_id integer,
    track_number character varying(50),
    price numeric,
    name character varying(100),
    sale integer,
    size character varying(10),
    total_price numeric,
    nm_id integer,
    brand character varying(50),
    status integer
);


ALTER TABLE public.items OWNER TO postgres;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    order_uid character varying(50) NOT NULL,
    track_number character varying(50),
    entry character varying(50),
    locale character varying(10),
    internal_signature text,
    customer_id character varying(50),
    delivery_service character varying(50),
    shardkey character varying(10),
    sm_id integer,
    date_created timestamp without time zone,
    oof_shard character varying(10)
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
    transaction character varying(50) NOT NULL,
    order_uid character varying(50),
    request_id character varying(50),
    currency character varying(10),
    provider character varying(50),
    amount numeric,
    payment_dt bigint,
    bank character varying(50),
    delivery_cost numeric,
    goods_total numeric,
    custom_fee numeric
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- Data for Name: delivery; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.delivery (order_uid, name, phone, zip, city, address, region, email) FROM stdin;
b563feb7b2b84b6test	Test Testov	+9720000000	2639809	Kiryat Mozkin	Ploshad Mira 15	Kraiot	test@gmail.com
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (rid, order_uid, chrt_id, track_number, price, name, sale, size, total_price, nm_id, brand, status) FROM stdin;
ab4219087a764ae0btest	b563feb7b2b84b6test	9934930	WBILMTESTTRACK	453	Mascaras	30	0	317	2389212	Vivienne Sabo	202
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) FROM stdin;
b563feb7b2b84b6test	WBILMTESTTRACK	WBIL	en		test	meest	9	99	2021-11-26 06:22:19	1
\.


--
-- Data for Name: payment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment (transaction, order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) FROM stdin;
b563feb7b2b84b6test	b563feb7b2b84b6test		USD	wbpay	1817	1637907727	alpha	1500	317	0
\.


--
-- Name: delivery delivery_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery
    ADD CONSTRAINT delivery_pkey PRIMARY KEY (order_uid);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (rid);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (order_uid);


--
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (transaction);


--
-- Name: delivery delivery_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery
    ADD CONSTRAINT delivery_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


--
-- Name: items items_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


--
-- Name: payment payment_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid);


--
-- PostgreSQL database dump complete
--
