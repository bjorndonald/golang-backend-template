CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
			id UUID NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255),
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			photo VARCHAR(255) DEFAULT NULL,
			ip VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login VARCHAR(255) NULL,
		 role VARCHAR(255) DEFAULT 'user',
		 email_verified BOOLEAN DEFAULT FALSE,
		 country VARCHAR(255) DEFAULT NULL,
		 phone_number VARCHAR(255) DEFAULT NULL,
		 status VARCHAR(255) DEFAULT 'Inactive'
			);

INSERT INTO users (id, email, password, first_name, last_name, photo, ip, role, email_verfied, status)
VALUES ('38ee1579-4454-4b4d-8163-32370713789e', 'bjorndonaldb@gmail.com', '$2a$04$e26kbguUzcZRdbmfC7mBmORPQ7wkzIb1vlsus32k3.Xg5kNLRPD66', 'admin', 'user', '', "192.168.65.1", "admin", "true", "active");

CREATE TABLE user_agents (
	id UUID NOT NULL,
	user_id UUID NOT NULL,
	agent VARCHAR(255) NOT NULL
)

INSERT INTO user_agents (id, user_id, agent) VALUES ("019350c0-f209-70d0-9368-6798e5b86b0e", "38ee1579-4454-4b4d-8163-32370713789e", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")