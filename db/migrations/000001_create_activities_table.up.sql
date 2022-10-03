CREATE TABLE IF NOT EXISTS activities(
    activity_id INT GENERATED ALWAYS AS IDENTITY,
    activity_name VARCHAR(100) NOT NULL,
    description TEXT,
    units VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (activity_id)
);