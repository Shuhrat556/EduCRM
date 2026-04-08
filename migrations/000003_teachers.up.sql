-- Teaching groups (classes/cohorts). Insert rows before assigning teachers.
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(32) UNIQUE,
    email VARCHAR(255) UNIQUE,
    specialization VARCHAR(255),
    photo_url VARCHAR(2048),
    photo_storage_key VARCHAR(512),
    photo_content_type VARCHAR(128),
    photo_original_name VARCHAR(255),
    status VARCHAR(16) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_teachers_status ON teachers (status);

-- Many-to-many: one teacher can belong to multiple groups.
CREATE TABLE IF NOT EXISTS teacher_groups (
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    PRIMARY KEY (teacher_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_teacher_groups_group_id ON teacher_groups (group_id);
