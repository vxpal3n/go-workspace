INSERT INTO students (nim, name, email, password, role, is_active, grade)
SELECT
    'M' || LPAD(g::TEXT, 3, '0'),
    'Load Test ' || LPAD(g::TEXT, 3, '0'),
    'mhs' || LPAD(g::TEXT, 3, '0') || '@load.test',
    '!',
    'user',
    true,
    0
FROM generate_series(1, 60) AS g
ON CONFLICT DO NOTHING;