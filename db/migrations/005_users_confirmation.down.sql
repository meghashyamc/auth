ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS verified,
DROP COLUMN IF EXISTS confirmation_token,
DROP COLUMN IF EXISTS confirmation_validity;