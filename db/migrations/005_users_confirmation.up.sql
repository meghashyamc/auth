ALTER TABLE users
ADD COLUMN verified BOOLEAN,
ADD COLUMN confirmation_token VARCHAR(1000),
ADD COLUMN confirmation_validity TIMESTAMP DEFAULT NOW();