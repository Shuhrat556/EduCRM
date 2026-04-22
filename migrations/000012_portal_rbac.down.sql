-- Best-effort rollback (may fail if dependent data exists).

ALTER TABLE grades DROP CONSTRAINT IF EXISTS uq_grade_weekly;
ALTER TABLE grades ADD CONSTRAINT uq_grade_weekly UNIQUE (student_id, teacher_id, group_id, week_start_date, grade_type);
ALTER TABLE grades DROP COLUMN IF EXISTS subject_id;

ALTER TABLE attendances DROP CONSTRAINT IF EXISTS uq_attendance_student_group_subject_day;
ALTER TABLE attendances ADD CONSTRAINT uq_attendance_student_group_day UNIQUE (student_id, group_id, lesson_date);
ALTER TABLE attendances DROP COLUMN IF EXISTS subject_id;

DROP TABLE IF EXISTS teacher_group_subject_assignments;
DROP TABLE IF EXISTS student_profiles;
DROP TABLE IF EXISTS admin_profiles;

ALTER TABLE groups DROP COLUMN IF EXISTS academic_year;

DROP INDEX IF EXISTS ux_users_username_lower;
ALTER TABLE users
    DROP COLUMN IF EXISTS created_by_user_id,
    DROP COLUMN IF EXISTS force_password_change,
    DROP COLUMN IF EXISTS username,
    DROP COLUMN IF EXISTS full_name;
