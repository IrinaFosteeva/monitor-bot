CREATE TABLE regions
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO regions (name)
VALUES ('default');

ALTER TABLE targets
DROP
COLUMN IF EXISTS region_restriction;

ALTER TABLE targets
    ADD COLUMN region_id INT NOT NULL DEFAULT 1 REFERENCES regions (id);
