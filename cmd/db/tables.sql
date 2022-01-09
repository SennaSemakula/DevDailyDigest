CREATE TABLE USERS (
    email VARCHAR(100) PRIMARY KEY NOT NULL,
    loc TEXT NOT NULL,
    git TEXT NOT NULL,
    language TEXT NOT NULL,
    level TEXT NOT NULL,
    frameworks TEXT NOT NULL,
    resume TEXT NOT NULL,
    recruit BOOLEAN NOT NULL
);

CREATE TABLE LANGUAGES (
    email VARCHAR(100) NOT NULL,
    FOREIGN KEY (email) REFERENCES USERS(email) ON DELETE CASCADE,
    language TEXT NOT NULL,
    framework TEXT NOT NULL,
);

CREATE TABLE INTERESTS (
    email VARCHAR(100) NOT NULL,
    FOREIGN KEY (email) REFERENCES USERS(email) ON DELETE CASCADE,
    jobs BOOLEAN NOT NULL,
    events BOOLEAN NOT NULL,
    articles BOOLEAN NOT NULL
);

-- Add an index on the email column.
CREATE INDEX idx_emails ON USERS(email);

-- Add some dummy data
INSERT INTO USERS (email, loc, language, level, frameworks, git, recruit) VALUES (
    "hello@gmail.com",
    "London",
    "Python",
    "Mid",
    "Django,Flask",
    "Mygitprofile",
    FALSE
);

INSERT INTO USERS (email, loc, git, recruit) VALUES (
    "hello@gmail.com",
    "Cali",
    "Mygitprofile",
    TRUE
);

-- -- 
-- ALTER TABLE USERS
-- ADD date TEXT NOT NULL;
-- ALTER TABLE USERS
-- ADD frameworks TEXT NOT NULL;
-- ALTER TABLE USERS
-- ADD level TEXT NOT NULL;
-- ALTER TABLE USERS ADD COLUMN date TEXT NOT NULL;
-- ALTER TABLE USERS ADD COLUMN jobstatus TEXT NOT NULL;
