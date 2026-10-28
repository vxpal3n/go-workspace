CREATE TABLE IF NOT EXISTS supervisors (
    id                SERIAL       PRIMARY KEY,
    name              VARCHAR(120) NOT NULL,
    quota             INTEGER      NOT NULL CHECK (quota > 0),
    supervised_count  INTEGER      NOT NULL DEFAULT 0,
    version           INTEGER      NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT supervisors_supervised_within_quota
        CHECK (supervised_count >= 0 AND supervised_count <= quota)
);

CREATE TABLE IF NOT EXISTS supervisions (
    id            BIGSERIAL   PRIMARY KEY,
    student_id    INTEGER     NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    supervisor_id INTEGER     NOT NULL REFERENCES supervisors(id) ON DELETE RESTRICT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS supervisions_supervisor_idx
    ON supervisions (supervisor_id);

INSERT INTO supervisors (name, quota) VALUES
    ('Dr. Ahmad Fauzi',   3),
    ('Dr. Budi Santoso',  5),
    ('Dr. Citra Dewi',    4),
    ('Dr. Doni Pratama',  3),
    ('Dr. Eka Wijaya',   10)
ON CONFLICT DO NOTHING;

INSERT INTO permissions (name, description) VALUES
    ('supervision:manage', 'Menetapkan dan mengganti pembimbing TA')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'supervision:manage')
ON CONFLICT DO NOTHING;