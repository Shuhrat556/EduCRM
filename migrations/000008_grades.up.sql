CREATE TABLE IF NOT EXISTS grades (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    week_start_date DATE NOT NULL,
    grade_type VARCHAR(32) NOT NULL,
    grade_value DOUBLE PRECISION NOT NULL,
    comment TEXT,
    graded_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_grade_weekly UNIQUE (student_id, teacher_id, group_id, week_start_date, grade_type)
);

CREATE INDEX IF NOT EXISTS idx_grades_student_id ON grades (student_id);
CREATE INDEX IF NOT EXISTS idx_grades_group_id ON grades (group_id);
CREATE INDEX IF NOT EXISTS idx_grades_week_start ON grades (week_start_date);
