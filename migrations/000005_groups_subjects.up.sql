-- Subjects catalog (one subject per group).
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    code VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(16) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subjects_status ON subjects (status);

-- Seed default subject for migrating legacy groups (code GEN).
INSERT INTO subjects (id, name, code, status, created_at, updated_at)
SELECT '00000000-0000-4000-8000-000000000001'::uuid,
       'General',
       'GEN',
       'active',
       now(),
       now()
WHERE NOT EXISTS (SELECT 1 FROM subjects WHERE code = 'GEN');

-- Extend groups: one teacher, one subject, optional room, schedule, fee, status.
ALTER TABLE groups ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects (id);
ALTER TABLE groups ADD COLUMN IF NOT EXISTS teacher_id UUID REFERENCES teachers (id);
ALTER TABLE groups ADD COLUMN IF NOT EXISTS room_id UUID REFERENCES rooms (id) ON DELETE SET NULL;
ALTER TABLE groups ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE groups ADD COLUMN IF NOT EXISTS end_date DATE;
ALTER TABLE groups ADD COLUMN IF NOT EXISTS monthly_fee_minor BIGINT;
ALTER TABLE groups ADD COLUMN IF NOT EXISTS status VARCHAR(16);

-- Backfill from teacher_groups (pick one teacher per group).
UPDATE groups g
SET teacher_id = tg.teacher_id
FROM (
    SELECT DISTINCT ON (group_id) group_id, teacher_id
    FROM teacher_groups
    ORDER BY group_id, teacher_id
) tg
WHERE g.id = tg.group_id
  AND g.teacher_id IS NULL;

UPDATE groups
SET subject_id = '00000000-0000-4000-8000-000000000001'::uuid
WHERE subject_id IS NULL;

UPDATE groups
SET teacher_id = (SELECT id FROM teachers ORDER BY created_at ASC LIMIT 1)
WHERE teacher_id IS NULL;

UPDATE groups
SET start_date = (created_at::date)
WHERE start_date IS NULL;

UPDATE groups
SET end_date = ((created_at::date) + interval '1 year')::date
WHERE end_date IS NULL;

UPDATE groups
SET monthly_fee_minor = 0
WHERE monthly_fee_minor IS NULL;

UPDATE groups
SET status = 'active'
WHERE status IS NULL OR trim(status) = '';

ALTER TABLE groups ALTER COLUMN subject_id SET NOT NULL;
ALTER TABLE groups ALTER COLUMN teacher_id SET NOT NULL;
ALTER TABLE groups ALTER COLUMN start_date SET NOT NULL;
ALTER TABLE groups ALTER COLUMN end_date SET NOT NULL;
ALTER TABLE groups ALTER COLUMN monthly_fee_minor SET NOT NULL;
ALTER TABLE groups ALTER COLUMN status SET NOT NULL;

DROP TABLE IF EXISTS teacher_groups;

CREATE INDEX IF NOT EXISTS idx_groups_subject_id ON groups (subject_id);
CREATE INDEX IF NOT EXISTS idx_groups_teacher_id ON groups (teacher_id);
CREATE INDEX IF NOT EXISTS idx_groups_room_id ON groups (room_id);
CREATE INDEX IF NOT EXISTS idx_groups_status ON groups (status);

-- At most one group per student user (enrollment API can be added later).
CREATE TABLE IF NOT EXISTS student_group_memberships (
    user_id UUID NOT NULL PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_student_group_memberships_group_id ON student_group_memberships (group_id);
