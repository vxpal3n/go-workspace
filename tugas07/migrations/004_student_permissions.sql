ALTER TABLE students
    ADD COLUMN IF NOT EXISTS owner_id INTEGER;

ALTER TABLE students
    DROP CONSTRAINT IF EXISTS students_owner_id_fkey;
ALTER TABLE students
    ADD CONSTRAINT students_owner_id_fkey
    FOREIGN KEY (owner_id) REFERENCES students(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS students_owner_id_idx ON students (owner_id);