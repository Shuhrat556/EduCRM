DROP TABLE IF EXISTS student_group_memberships;

CREATE TABLE IF NOT EXISTS teacher_groups (
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    PRIMARY KEY (teacher_id, group_id)
);

INSERT INTO teacher_groups (teacher_id, group_id)
SELECT teacher_id, id
FROM groups
WHERE teacher_id IS NOT NULL
ON CONFLICT (teacher_id, group_id) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_teacher_groups_group_id ON teacher_groups (group_id);

ALTER TABLE groups DROP COLUMN IF EXISTS status;
ALTER TABLE groups DROP COLUMN IF EXISTS monthly_fee_minor;
ALTER TABLE groups DROP COLUMN IF EXISTS end_date;
ALTER TABLE groups DROP COLUMN IF EXISTS start_date;
ALTER TABLE groups DROP COLUMN IF EXISTS room_id;
ALTER TABLE groups DROP COLUMN IF EXISTS teacher_id;
ALTER TABLE groups DROP COLUMN IF EXISTS subject_id;

DROP TABLE IF EXISTS subjects;
