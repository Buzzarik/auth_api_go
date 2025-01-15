DROP TRIGGER IF EXISTS before_update_tokens ON tokens;
DROP FUNCTION IF EXISTS delete_expired_tokens();
DROP TABLE IF EXISTS tokens;
