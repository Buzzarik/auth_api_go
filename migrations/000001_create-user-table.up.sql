CREATE TABLE IF NOT EXISTS users (
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	name text NOT NULL,
	phone_number text UNIQUE NOT NULL,
	hash_password text NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users (phone_number);