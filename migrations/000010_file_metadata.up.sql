CREATE TABLE IF NOT EXISTS file_metadata (
    id UUID PRIMARY KEY,
    owner_type VARCHAR(32) NOT NULL,
    owner_id UUID NOT NULL,
    file_name VARCHAR(512) NOT NULL,
    storage_key VARCHAR(1024) NOT NULL,
    file_url TEXT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_file_metadata_storage_key UNIQUE (storage_key)
);

CREATE INDEX IF NOT EXISTS idx_file_metadata_owner ON file_metadata (owner_type, owner_id);
CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata (created_at DESC);
