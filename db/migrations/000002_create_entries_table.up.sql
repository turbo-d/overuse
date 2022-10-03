CREATE TABLE IF NOT EXISTS entries(
    entry_id INT GENERATED ALWAYS AS IDENTITY,
    activity_id INT NOT NULL,
    occured_on DATE NOT NULL CHECK (occured_on <= CURRENT_DATE),
    volume INT NOT NULL CHECK (volume > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (entry_id),
    FOREIGN KEY (activity_id)
        REFERENCES activities (activity_id)
        ON DELETE CASCADE
);