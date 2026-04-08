CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    room_id UUID NOT NULL REFERENCES rooms (id) ON DELETE RESTRICT,
    weekday SMALLINT NOT NULL CHECK (weekday >= 0 AND weekday <= 6),
    start_minutes INT NOT NULL CHECK (start_minutes >= 0 AND start_minutes < 1440),
    end_minutes INT NOT NULL CHECK (end_minutes > start_minutes AND end_minutes <= 1440),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_schedules_group_id ON schedules (group_id);
CREATE INDEX IF NOT EXISTS idx_schedules_teacher_id_weekday ON schedules (teacher_id, weekday);
CREATE INDEX IF NOT EXISTS idx_schedules_room_id_weekday ON schedules (room_id, weekday);
