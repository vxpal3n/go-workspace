-- ============================================================================
-- SIAKAD Mini — Schema Initialization
-- ============================================================================
-- Empat tabel utama:
--   users        — akun login (admin & mahasiswa)
--   students     — data mahasiswa (1-1 dengan users)
--   courses      — mata kuliah
--   enrollments  — KRS (mahasiswa × mata kuliah × tahun akademik)
-- ============================================================================

-- ---------------------------------------------------------------------------
-- Tabel: users
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL       PRIMARY KEY,
    email       VARCHAR(255) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'mahasiswa')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Email unik tanpa membedakan huruf besar/kecil.
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_key
    ON users (LOWER(email));

-- ---------------------------------------------------------------------------
-- Tabel: students
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS students (
    id             SERIAL       PRIMARY KEY,
    user_id        INTEGER      NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    nim            VARCHAR(20)  NOT NULL,
    nama           VARCHAR(100) NOT NULL,
    prodi          VARCHAR(100) NOT NULL,
    angkatan       INTEGER      NOT NULL,
    ipk_terakhir   NUMERIC(3,2) NOT NULL DEFAULT 0.00
                                CHECK (ipk_terakhir >= 0 AND ipk_terakhir <= 4),
    deleted_at     TIMESTAMPTZ,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- NIM unik tanpa membedakan huruf besar/kecil.
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key
    ON students (LOWER(nim));

-- Filter daftar: hanya tampilkan yang belum dihapus.
CREATE INDEX IF NOT EXISTS students_not_deleted_idx
    ON students (deleted_at)
    WHERE deleted_at IS NULL;

-- Filter by prodi dan angkatan untuk endpoint GET /students.
CREATE INDEX IF NOT EXISTS students_prodi_angkatan_idx
    ON students (prodi, angkatan);

-- ---------------------------------------------------------------------------
-- Tabel: courses
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS courses (
    id          SERIAL       PRIMARY KEY,
    kode_mk     VARCHAR(20)  NOT NULL,
    nama_mk     VARCHAR(150) NOT NULL,
    sks         INTEGER      NOT NULL CHECK (sks >= 1 AND sks <= 6),
    semester    INTEGER      NOT NULL CHECK (semester >= 1 AND semester <= 8),
    kuota       INTEGER      NOT NULL CHECK (kuota > 0),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS courses_kode_mk_lower_key
    ON courses (LOWER(kode_mk));

CREATE INDEX IF NOT EXISTS courses_semester_idx
    ON courses (semester);

-- ---------------------------------------------------------------------------
-- Tabel: enrollments (KRS)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS enrollments (
    id              SERIAL       PRIMARY KEY,
    student_id      INTEGER      NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    course_id       INTEGER      NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    tahun_akademik  VARCHAR(20)  NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- Business rule #2: satu mahasiswa tidak boleh mengambil mata kuliah
    -- yang sama dua kali pada tahun akademik yang sama.
    CONSTRAINT enrollments_unique_per_year
        UNIQUE (student_id, course_id, tahun_akademik)
);

-- Hitung kuota terpakai per mata kuliah — dipanggil di setiap POST /enrollments.
CREATE INDEX IF NOT EXISTS enrollments_course_id_idx
    ON enrollments (course_id);

-- Daftar KRS seorang mahasiswa — untuk GET /students/{id}.
CREATE INDEX IF NOT EXISTS enrollments_student_id_idx
    ON enrollments (student_id);    