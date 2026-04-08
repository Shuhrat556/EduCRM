-- Links a users row (role teacher) to teachers.id so JWT subjects can be authorized against groups.teacher_id.
CREATE TABLE IF NOT EXISTS user_teacher_links (
    user_id UUID NOT NULL PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE CASCADE,
    CONSTRAINT uq_user_teacher_links_teacher UNIQUE (teacher_id)
);

CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    lesson_date DATE NOT NULL,
    status VARCHAR(16) NOT NULL,
    comment TEXT,
    marked_by_teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_attendance_student_group_day UNIQUE (student_id, group_id, lesson_date)
);

CREATE INDEX IF NOT EXISTS idx_attendances_student_id ON attendances (student_id);
CREATE INDEX IF NOT EXISTS idx_attendances_group_id ON attendances (group_id);
CREATE INDEX IF NOT EXISTS idx_attendances_lesson_date ON attendances (lesson_date);
CREATE INDEX IF NOT EXISTS idx_attendances_marked_by_teacher_id ON attendances (marked_by_teacher_id);
