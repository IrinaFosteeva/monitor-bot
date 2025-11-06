ALTER TABLE targets
DROP
COLUMN IF EXISTS region_id;

ALTER TABLE targets
    ADD COLUMN region_restriction TEXT;

DROP TABLE IF EXISTS regions;
