create database wallet_db;

create table wallets (
	id SERIAL primary key,
	created_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	updated_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL,
	number integer unique,
	balance integer
);

create table transactions (
	id SERIAL primary KEY,
	created_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	updated_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL,
	amount INTEGER,
	description text,
	type text,
	datetime TIMESTAMP not null default CURRENT_TIMESTAMP,
	source_id integer,
	from_number integer,
	to_number integer,
	constraint fk_transactions_to_wallet
	foreign key (to_number)
	references wallets(number) on update CASCADE
);

create table users (
	id SERIAL primary key,
	created_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	updated_at TIMESTAMP not null default CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL,
	name VARCHAR(256),
	email VARCHAR(50) unique,
	password VARCHAR(256),
	wallet_number integer,
	constraint fk_users_wallet foreign key (wallet_number)
	references wallets(number) on update cascade
);

INSERT INTO wallets (id, created_at, updated_at, deleted_at, number, balance) VALUES (4, '2022-09-09 20:31:56.965255+07', '2022-09-09 20:31:56.976004+07', NULL, 100004, 0);
INSERT INTO wallets (id, created_at, updated_at, deleted_at, number, balance) VALUES (3, '2022-09-09 19:01:07.102157+07', '2022-09-09 19:29:59.805945+07', NULL, 100003, 123000);
INSERT INTO wallets (id, created_at, updated_at, deleted_at, number, balance) VALUES (1, '2022-09-09 19:00:35.378232+07', '2022-09-09 19:31:26.174968+07', NULL, 100001, 1396000);
INSERT INTO wallets (id, created_at, updated_at, deleted_at, number, balance) VALUES (5, '2022-09-09 19:01:32.814824+07', '2022-09-09 19:31:47.12707+07', NULL, 100005, 617000);
INSERT INTO wallets (id, created_at, updated_at, deleted_at, number, balance) VALUES (2, '2022-09-09 19:00:51.142768+07', '2022-09-09 19:32:54.14518+07', NULL, 100002, 76000);

INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (1, '2022-09-09 19:06:17.939292+07', '2022-09-09 19:06:17.939292+07', NULL, 70000, 'Top Up from Bank Transfer', 'TOP_UP', '2022-09-09 19:06:17.939064+07', 1, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (2, '2022-09-09 19:06:58.268476+07', '2022-09-09 19:06:58.268476+07', NULL, 100000, 'Bayar Earphone', 'TRANSFER', '2022-09-09 19:06:58.255882+07', NULL, 100005, 100003);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (3, '2022-09-09 19:08:47.959677+07', '2022-09-09 19:08:47.959677+07', NULL, 412000, 'Top Up from Cash', 'TOP_UP', '2022-09-09 19:08:47.959526+07', 3, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (4, '2022-09-09 19:11:06.2444+07', '2022-09-09 19:11:06.2444+07', NULL, 12000, 'Beli mi ayam', 'TRANSFER', '2022-09-09 19:11:06.23216+07', NULL, 100001, 100002);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (5, '2022-09-09 19:27:31.684149+07', '2022-09-09 19:27:31.684149+07', NULL, 30000, 'Traktir makanan', 'TRANSFER', '2022-09-09 19:27:31.672177+07', NULL, 100004, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (6, '2022-09-09 19:29:59.806637+07', '2022-09-09 19:29:59.806637+07', NULL, 10000, 'Beli Baso', 'TRANSFER', '2022-09-09 19:29:59.793631+07', NULL, 100004, 100003);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (7, '2022-09-09 19:31:26.175796+07', '2022-09-09 19:31:26.175796+07', NULL, 50000, 'Belanja ke pasar', 'TRANSFER', '2022-09-09 19:31:26.1635+07', NULL, 100002, 100001);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (8, '2022-09-09 19:31:47.127514+07', '2022-09-09 19:31:47.127514+07', NULL, 50000, 'Ngedate', 'TRANSFER', '2022-09-09 19:31:47.125637+07', NULL, 100002, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (9, '2022-09-09 19:32:38.816666+07', '2022-09-09 19:32:38.816666+07', NULL, 20000, 'Denda parkir', 'TRANSFER', '2022-09-09 19:32:38.80453+07', NULL, 100002, 100004);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (10, '2022-09-09 19:32:54.156388+07', '2022-09-09 19:32:54.156388+07', NULL, 75000, 'Denda tilang', 'TRANSFER', '2022-09-09 19:32:54.144617+07', NULL, 100002, 100004);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (11, '2022-09-09 19:03:02.832213+07', '2022-09-09 19:03:02.832213+07', NULL, 500000, 'Top Up from Cash', 'TOP_UP', '2021-08-08 19:03:02.832044+07', 3, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (12, '2022-09-09 19:05:48.85285+07', '2022-09-09 19:05:48.85285+07', NULL, 500000, 'Top Up from Cash', 'TOP_UP', '2022-05-05 19:05:48.852673+07', 3, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (13, '2022-09-09 19:07:09.677093+07', '2022-09-09 19:07:09.677093+07', NULL, 13000, 'Nalangin makan', 'TRANSFER', '2022-07-03 19:07:09.665119+07', NULL, 100005, 100003);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (14, '2022-09-09 19:07:29.145566+07', '2022-09-09 19:07:29.145566+07', NULL, 27000, 'Makan ramen', 'TRANSFER', '2022-08-04 19:07:29.133674+07', NULL, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (15, '2022-09-09 19:08:04.348018+07', '2022-09-09 19:08:04.348018+07', NULL, 80000, 'Top Up from Credit Card', 'TOP_UP', '2022-05-01 19:08:04.347645+07', 2, 100005, 100005);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (16, '2022-09-09 19:08:36.855118+07', '2022-09-09 19:08:36.855118+07', NULL, 412000, 'Beli switch keyboard', 'TRANSFER', '2021-10-10 19:08:36.84233+07', NULL, 100005, 100001);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (20, '2022-09-09 19:11:16.872236+07', '2022-09-09 19:11:16.872236+07', NULL, 15000, 'Beli baso', 'TRANSFER', '2021-05-09 19:11:16.86034+07', NULL, 100001, 100002);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (19, '2022-09-09 19:11:29.805512+07', '2022-09-09 19:11:29.805512+07', NULL, 100000, 'Beli stand laptop', 'TRANSFER', '2022-07-09 19:11:29.793611+07', NULL, 100001, 100002);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (18, '2022-09-09 19:11:40.162878+07', '2022-09-09 19:11:40.162878+07', NULL, 50000, 'Top Up from Cash', 'TOP_UP', '2021-07-07 19:11:40.162676+07', 3, 100001, 100001);
INSERT INTO transactions (id, created_at, updated_at, deleted_at, amount, description, type, datetime, source_id, from_number, to_number) VALUES (17, '2022-09-09 19:11:49.354784+07', '2022-09-09 19:11:49.354784+07', NULL, 1001000, 'Top Up from Bank Transfer', 'TOP_UP', '2020-12-12 19:11:49.35465+07', 1, 100001, 100001);

INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, wallet_number) VALUES (1, '2022-09-09 19:00:35.39051+07', '2022-09-09 19:00:35.39051+07', NULL, 'Jonathan', 'jonathan@gmail.com', '$2a$04$kiZ.6xYFGoZ8wRHSNZllbOyHWiu6bcHB7GMtVGTxHa3FTHq64qSSS', 100001);
INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, wallet_number) VALUES (2, '2022-09-09 19:00:51.154668+07', '2022-09-09 19:00:51.154668+07', NULL, 'Thoriq', 'thoriq@gmail.com', '$2a$04$piqyhed0ifSSzwthT8auyOW8GJyOOsKaDpCVxcS1Rn7y0V/Rtd0WS', 100002);
INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, wallet_number) VALUES (3, '2022-09-09 19:01:07.113586+07', '2022-09-09 19:01:07.113586+07', NULL, 'Andri Winata', 'andri.winata@gmail.com', '$2a$04$.T9L6cT.0hc9Pxs5HTrExObVpwjY5GAZVVtwcXFjGZM3bUYtbQxXy', 100003);
INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, wallet_number) VALUES (4, '2022-09-09 19:01:21.874727+07', '2022-09-09 19:01:21.874727+07', NULL, 'Tafia', 'tafia@gmail.com', '$2a$04$6zSKoCy0bQgsAaRHvvDXy.wNfzL9nE4jDlIrl7vTE6sXfkyENZjVG', 100004);
INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, wallet_number) VALUES (5, '2022-09-09 19:01:32.826546+07', '2022-09-09 19:01:32.826546+07', NULL, 'Example Name', 'example@gmail.com', '$2a$04$fCdi9GZjwKQH8IhtHyRA3ODLQRj5Fru/doksRn8OXkdQA3VnU5KBm', 100005);

alter sequence wallets_id_seq restart with 6;

alter sequence users_id_seq restart with 6;

alter sequence transactions_id_seq restart with 21;