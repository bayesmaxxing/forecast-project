CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS scores (
    id BIGSERIAL PRIMARY KEY,
    brier_score DOUBLE PRECISION NOT NULL,
    log2_score DOUBLE PRECISION NOT NULL,
    logn_score DOUBLE PRECISION NOT NULL,
    user_id BIGINT NOT NULL,
    forecast_id BIGINT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (forecast_id) REFERENCES forecasts(id)
);

CREATE TABLE IF NOT EXISTS points (
    update_id BIGSERIAL PRIMARY KEY,
    forecast_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    point_forecast DOUBLE PRECISION NOT NULL,
    reason TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (forecast_id) REFERENCES forecasts(id)
);

CREATE TABLE IF NOT EXISTS forecasts (
    id BIGSERIAL PRIMARY KEY,
    question TEXT,
    category TEXT,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL,
    resolution_criteria TEXT,
    resolution TEXT,
    resolved TIMESTAMP,
    resolution TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
