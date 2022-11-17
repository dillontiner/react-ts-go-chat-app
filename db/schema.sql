CREATE TABLE users(
	uuid VARCHAR(36) PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
	password VARCHAR(64) NOT NULL
);

CREATE TABLE messages(
	uuid VARCHAR(36) PRIMARY KEY,
	sent_at TIMESTAMP,
	sender_uuid VARCHAR(36) REFERENCES users(uuid),
	body text,
	upvote_user_uuids VARCHAR(36)[],
	downvote_user_uuids VARCHAR(36)[]
);

CREATE TABLE sessions(
	uuid VARCHAR(36) PRIMARY KEY, 
	user_uuid VARCHAR(36) REFERENCES users(uuid), 
	created_at TIMESTAMP, 
	ends_at TIMESTAMP
);