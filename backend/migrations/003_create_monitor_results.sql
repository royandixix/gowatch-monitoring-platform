CREATE TABLE IF NOT EXISTS monitor_results (
    id BIGSERIAL PRIMARY KEY,
    monitor_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    http_status INTEGER,
    response_time_ms BIGINT,
    error TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_monitor_results_monitor
        FOREIGN KEY (monitor_id)
        REFERENCES monitors(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_monitor_results_monitor_id
ON monitor_results(monitor_id);

CREATE INDEX IF NOT EXISTS idx_monitor_results_checked_at
ON monitor_results(checked_at);