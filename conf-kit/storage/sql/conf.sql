CREATE DATABASE IF NOT EXISTS test;

USE test;

# DROP TABLE IF EXISTS tb_config;
CREATE TABLE IF NOT EXISTS tb_config (
	id BIGINT NOT NULL AUTO_INCREMENT,
	pid BIGINT NOT NULL,
	name VARCHAR(128) NOT NULL,
	value TEXT NOT NULL,
	PRIMARY KEY(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
