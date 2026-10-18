-- 001_initial_schema.sql
-- CustardFin Cloud API Initial Schema

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Edge Nodes
CREATE TABLE IF NOT EXISTS nodes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    node_token_hash VARCHAR(255) NOT NULL,
    frp_auth_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nodes_user_id ON nodes(user_id);
CREATE INDEX IF NOT EXISTS idx_nodes_token_hash ON nodes(node_token_hash);

-- Catalog Items (Metadata & Manifest only - NO video binaries!)
CREATE TABLE IF NOT EXISTS catalog_items (
    media_id VARCHAR(36) NOT NULL,
    node_id VARCHAR(36) NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    duration_s INT NOT NULL,
    content_hash VARCHAR(64) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    manifest_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (node_id, media_id)
);

CREATE INDEX IF NOT EXISTS idx_catalog_items_title ON catalog_items(title);
CREATE INDEX IF NOT EXISTS idx_catalog_items_is_public ON catalog_items(is_public);
CREATE INDEX IF NOT EXISTS idx_catalog_items_media_id ON catalog_items(media_id);

-- Catalog Renditions
CREATE TABLE IF NOT EXISTS catalog_renditions (
    id VARCHAR(36) PRIMARY KEY,
    media_id VARCHAR(36) NOT NULL,
    node_id VARCHAR(36) NOT NULL,
    resolution VARCHAR(16) NOT NULL,
    bitrate_bps INT NOT NULL,
    playlist_path VARCHAR(255) NOT NULL,
    chunk_duration_s INT NOT NULL DEFAULT 4,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (node_id, media_id) REFERENCES catalog_items(node_id, media_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_catalog_renditions_media ON catalog_renditions(node_id, media_id);
