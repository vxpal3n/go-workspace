-- ============================================================================
-- SIAKAD Mini — Schema Initialization
-- ============================================================================
-- Empat tabel: users, students, courses, enrollments.
-- Soft delete dipakai pada students (kolom deleted_at) agar riwayat KRS
-- tidak hilang ketika mahasiswa dihapus dari daftar aktif.
-- ============================================================================

-- 1. users — akun login untuk admin dan mahasiswa
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL       PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'mahasiswa'
                            CHECK (role IN ('admin', 'mahasiswa')),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_key
    ON users (LOWER(email));

-- 2. students — data akademik mahasiswa; relasi 1-1 dengan users
CREATE TABLE IF NOT EXISTS students (
    id            SERIAL       PRIMARY KEY,
    user_id       INTEGER      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nim           VARCHAR(12)  NOT NULL,
    nama          VARCHAR(100) NOT NULL,
    prodi         VARCHAR(100) NOT NULL,
    angkatan      INTEGER      NOT NULL,
    ipk_terakhir  DECIMAL(3,2),
    deleted_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS students_nim_key
    ON students (nim);
CREATE UNIQUE INDEX IF NOT EXISTS students_user_id_key
    ON students (user_id);
CREATE INDEX IF NOT EXISTS students_deleted_at_idx
    ON students (deleted_at);

-- 3. courses — daftar mata kuliah
CREATE TABLE IF NOT EXISTS courses (
    id         SERIAL       PRIMARY KEY,
    kode_mk    VARCHAR(20)  NOT NULL,
    nama_mk    VARCHAR(150) NOT NULL,
    sks        INTEGER      NOT NULL CHECK (sks > 0 AND sks <= 6),
    semester   INTEGER      NOT NULL CHECK (semester >= 1 AND semester <= 8),
    kuota      INTEGER      NOT NULL CHECK (kuota >= 0),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS courses_kode_mk_key
    ON courses (kode_mk);

-- 4. enrollments — KRS; relasi 1-N ke students dan courses
CREATE TABLE IF NOT EXISTS enrollments (
    id             SERIAL       PRIMARY KEY,
    student_id     INTEGER      NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    course_id      INTEGER      NOT NULL REFERENCES courses(id)  ON DELETE CASCADE,
    tahun_akademik VARCHAR(20)  NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT enrollments_unique_per_year
        UNIQUE (student_id, course_id, tahun_akademik)
);

CREATE INDEX IF NOT EXISTS enrollments_student_id_idx
    ON enrollments (student_id);
CREATE INDEX IF NOT EXISTS enrollments_course_id_idx
    ON enrollments (course_id);