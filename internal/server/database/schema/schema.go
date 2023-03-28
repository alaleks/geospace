// Package schema is the structure of a database described.
package schema

var City = `
	CREATE TABLE cities (
		cid INT auto_increment NULL,
		name varchar(100) NULL,
		name_ascii varchar(100) NULL,
		alternative_names TEXT NULL,
		country_code varchar(2) NULL,
		country varchar(100) NULL,
		timezone varchar(100) NULL,
		latitude FLOAT NULL,
		longitude FLOAT NULL,
		created_at INT NULL,
		CONSTRAINT cities_PK PRIMARY KEY (cid),
		FULLTEXT KEY (name,alternative_names),
		INDEX latitude_idx (latitude),
		INDEX longitude_idx (longitude)
	)

		ENGINE=InnoDB
		DEFAULT CHARSET=utf8mb4
		COLLATE=utf8mb4_general_ci;
`
var User = `
	CREATE TABLE users (
		uid INT auto_increment NULL,
		name varchar(100) NULL,
		email varchar(100) NULL,
		password varchar(256) NULL,
		created_at INT NULL,
		CONSTRAINT users_PK PRIMARY KEY (uid),
		FULLTEXT KEY (name,email)
	)

		ENGINE=InnoDB
		DEFAULT CHARSET=utf8mb4
		COLLATE=utf8mb4_general_ci;
`
