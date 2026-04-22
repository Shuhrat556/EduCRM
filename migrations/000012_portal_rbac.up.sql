-- Portal auth fields, profiles, teacher–group–subject assignments, attendance/grade subject scoping.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS username VARCHAR(64),
    ADD COLUMN IF NOT EXISTS force_password_change BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS created_by_user_id UUID REFERENCES users (id) ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_username_lower ON users (lower(username)) WHERE username IS NOT NULL AND trim(username) <> '';

ALTER TABLE groups ADD COLUMN IF NOT EXISTS academic_year VARCHAR(32);

CREATE TABLE IF NOT EXISTS admin_profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS student_profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    student_code VARCHAR(64),
    group_id UUID REFERENCES groups (id) ON DELETE SET NULL,
    parent_phone VARCHAR(32),
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_student_profiles_group_id ON student_profiles (group_id);

CREATE TABLE IF NOT EXISTS teacher_group_subject_assignments (
    id UUID PRIMARY KEY,
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_tgsa_teacher_group_subject UNIQUE (teacher_id, group_id, subject_id)
);

CREATE INDEX IF NOT EXISTS idx_tgsa_teacher ON teacher_group_subject_assignments (teacher_id);
CREATE INDEX IF NOT EXISTS idx_tgsa_group ON teacher_group_subject_assignments (group_id);

-- Backfill assignments from existing groups (one subject per group).
INSERT INTO teacher_group_subject_assignments (id, teacher_id, group_id, subject_id)
SELECT gen_random_uuid(), g.teacher_id, g.id, g.subject_id
FROM groups g
WHERE NOT EXISTS (
    SELECT 1 FROM teacher_group_subject_assignments t
    WHERE t.teacher_id = g.teacher_id AND t.group_id = g.id AND t.subject_id = g.subject_id
);

-- Attendance: add subject_id, replace uniqueness.
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects (id);

UPDATE attendances a
SET subject_id = g.subject_id
FROM groups g
WHERE a.group_id = g.id AND a.subject_id IS NULL;

ALTER TABLE attendances ALTER COLUMN subject_id SET NOT NULL;

ALTER TABLE attendances DROP CONSTRAINT IF EXISTS uq_attendance_student_group_day;

ALTER TABLE attendances ADD CONSTRAINT uq_attendance_student_group_subject_day UNIQUE (student_id, group_id, subject_id, lesson_date);

-- Grades: add subject_id, replace uniqueness.
ALTER TABLE grades ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects (id);

UPDATE grades gr
SET subject_id = g.subject_id
FROM groups g
WHERE gr.group_id = g.id AND gr.subject_id IS NULL;

ALTER TABLE grades ALTER COLUMN subject_id SET NOT NULL;

ALTER TABLE grades DROP CONSTRAINT IF EXISTS uq_grade_weekly;

ALTER TABLE grades ADD CONSTRAINT uq_grade_weekly UNIQUE (student_id, teacher_id, group_id, subject_id, week_start_date, grade_type);

CREATE INDEX IF NOT EXISTS idx_grades_subject_id ON grades (subject_id);
