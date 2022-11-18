CREATE TABLE users(
	uuid VARCHAR(36) PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
	password VARCHAR(64) NOT NULL
);

CREATE TABLE messages(
	uuid VARCHAR(36) PRIMARY KEY,
	sent_at TIMESTAMP,
	sender_uuid VARCHAR(36),
	body text
);

CREATE TABLE votes(
	uuid VARCHAR(36) PRIMARY KEY,
	updated_at TIMESTAMP,
	message_uuid VARCHAR(36),
	voter_uuid VARCHAR(36),
	vote BOOLEAN
);