CREATE TABLE IF NOT EXISTS entries(
    id INT GENERATED ALWAYS AS IDENTITY,
    activity_id INT NOT NULL,
    occured_on DATE NOT NULL CHECK (occured_on <= CURRENT_DATE),
    volume INT NOT NULL CHECK (volume > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    PRIMARY KEY (id)
    FOREIGN KEY (activity_id)
        REFERENCES activities (id)
        ON DELETE CASCADE
);