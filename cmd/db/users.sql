CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON devdaily.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'codebundle';