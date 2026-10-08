-- ============================================================================
-- SIAKAD Mini — Seeder
-- ============================================================================
-- Membutuhkan ekstensi pgcrypto untuk menghasilkan bcrypt hash.
-- Password yang dihasilkan kompatibel dengan verifikasi bcrypt di Go.
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Bersihkan data agar seeder dapat dijalankan berulang.
TRUNCATE TABLE enrollments RESTART IDENTITY CASCADE;
TRUNCATE TABLE students    RESTART IDENTITY CASCADE;
TRUNCATE TABLE users       RESTART IDENTITY CASCADE;
TRUNCATE TABLE courses     RESTART IDENTITY CASCADE;

-- ---------------------------------------------------------------------------
-- 1 Admin
-- ---------------------------------------------------------------------------
INSERT INTO users (email, password, role)
VALUES (
    'admin@siakad.test',
    crypt('admin12345', gen_salt('bf', 12)),
    'admin'
);

-- ---------------------------------------------------------------------------
-- 20 Mahasiswa
-- Password awal = NIM masing-masing (di-hash dengan bcrypt cost 12).
-- ---------------------------------------------------------------------------
WITH seed(nim, nama, prodi, angkatan, ipk) AS (
    VALUES
        ('187221000001', 'Rina Putri',         'Sistem Informasi', 2022, 3.45),
        ('187221000002', 'Budi Santoso',       'Sistem Informasi', 2022, 3.12),
        ('187221000003', 'Citra Dewi',         'Sistem Informasi', 2022, 2.87),
        ('187221000004', 'Doni Pratama',       'Sistem Informasi', 2022, 2.34),
        ('187221000005', 'Eka Wijaya',         'Sistem Informasi', 2022, 3.78),
        ('187221000006', 'Fajar Hidayat',      'Sistem Informasi', 2023, 3.56),
        ('187221000007', 'Gita Lestari',       'Sistem Informasi', 2023, 3.21),
        ('187221000008', 'Hendra Kusuma',      'Sistem Informasi', 2023, 2.65),
        ('187221000009', 'Indah Permata',      'Sistem Informasi', 2023, 3.90),
        ('187221000010', 'Joko Susilo',        'Sistem Informasi', 2023, 2.45),
        ('187222000001', 'Kartika Sari',       'Teknik Informatika', 2022, 3.67),
        ('187222000002', 'Lukman Hakim',       'Teknik Informatika', 2022, 3.05),
        ('187222000003', 'Maya Anggraini',     'Teknik Informatika', 2022, 2.78),
        ('187222000004', 'Nanda Pratama',      'Teknik Informatika', 2023, 3.44),
        ('187222000005', 'Oki Ramadhan',       'Teknik Informatika', 2023, 3.89),
        ('187222000006', 'Putri Amelia',       'Teknik Informatika', 2023, 3.11),
        ('187222000007', 'Qori Anwar',         'Teknik Informatika', 2023, 2.99),
        ('187222000008', 'Rizki Firmansyah',   'Teknik Informatika', 2023, 3.55),
        ('187222000009', 'Sinta Rahayu',       'Teknik Informatika', 2024, 3.72),
        ('187222000010', 'Taufik Hidayat',     'Teknik Informatika', 2024, 2.88)
),
inserted_users AS (
    INSERT INTO users (email, password, role)
    SELECT
        nim || '@student.siakad.test',
        crypt(nim, gen_salt('bf', 12)),
        'mahasiswa'
    FROM seed
    RETURNING id, email
)
INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
SELECT
    iu.id, s.nim, s.nama, s.prodi, s.angkatan, s.ipk
FROM inserted_users iu
JOIN seed s ON iu.email = s.nim || '@student.siakad.test';

-- ---------------------------------------------------------------------------
-- 10 Mata Kuliah
-- ---------------------------------------------------------------------------
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota) VALUES
    ('IF1234', 'Pemrograman Dasar',              3, 1, 40),
    ('IF1235', 'Logika Informatika',             3, 1, 40),
    ('IF2101', 'Struktur Data',                  3, 3, 35),
    ('IF2102', 'Basis Data',                     3, 3, 35),
    ('IF2103', 'Pemrograman Web',                3, 4, 30),
    ('IF2104', 'Algoritma Lanjut',               3, 4, 30),
    ('IF3101', 'Pemrograman Backend',            4, 5, 25),
    ('IF3102', 'Keamanan Informasi',             3, 5, 25),
    ('IF3103', 'Machine Learning',               3, 5, 20),
    ('IF3201', 'Cloud Computing',                3, 6, 20);

-- ---------------------------------------------------------------------------
-- Verifikasi
-- ---------------------------------------------------------------------------
DO $$
DECLARE
    v_admin       INTEGER;
    v_mahasiswa   INTEGER;
    v_courses     INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_admin     FROM users WHERE role = 'admin';
    SELECT COUNT(*) INTO v_mahasiswa FROM users WHERE role = 'mahasiswa';
    SELECT COUNT(*) INTO v_courses   FROM courses;

    RAISE NOTICE 'Seeder selesai: admin=%, mahasiswa=%, courses=%',
        v_admin, v_mahasiswa, v_courses;
END $$;