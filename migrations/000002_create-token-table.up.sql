CREATE TABLE IF NOT EXISTS tokens (
    id bigserial PRIMARY KEY,
    hash text NOT NULL,
    id_user integer NOT NULL,
    id_api integer NOT NULL,
    expiry timestamp(0) with time zone NOT NULL,
    
    FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (id_user, id_api)
);

CREATE OR REPLACE FUNCTION delete_expired_tokens()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM tokens WHERE expiry < NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_tokens
BEFORE UPDATE ON tokens
FOR EACH ROW
EXECUTE FUNCTION delete_expired_tokens();

