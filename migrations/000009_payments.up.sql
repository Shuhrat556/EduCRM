CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    amount_minor BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL,
    payment_date DATE,
    month_for DATE NOT NULL,
    payment_type VARCHAR(32) NOT NULL,
    comment TEXT,
    is_free BOOLEAN NOT NULL DEFAULT false,
    discount_amount_minor BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payments_student_id ON payments (student_id);
CREATE INDEX IF NOT EXISTS idx_payments_group_id ON payments (group_id);
CREATE INDEX IF NOT EXISTS idx_payments_month_for ON payments (month_for);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_payment_type ON payments (payment_type);
CREATE INDEX IF NOT EXISTS idx_payments_is_free ON payments (is_free);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
