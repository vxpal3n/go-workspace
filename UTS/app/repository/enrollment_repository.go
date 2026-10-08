package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"UTS/app/model"
)

type EnrollmentRepository interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	LockCourseForUpdate(ctx context.Context, tx pgx.Tx, courseID int) (model.Course, error)
	CountByCourseAndYear(ctx context.Context, tx pgx.Tx, courseID int, tahun string) (int, error)
	Exists(ctx context.Context, tx pgx.Tx, studentID, courseID int, tahun string) (bool, error)
	SumSKSForStudentYear(ctx context.Context, tx pgx.Tx, studentID int, tahun string) (int, error)
	Insert(ctx context.Context, tx pgx.Tx, studentID, courseID int, tahun string) (model.Enrollment, error)

	FindByID(ctx context.Context, id int) (model.Enrollment, error)
	Delete(ctx context.Context, id int) error
	ListEnrichedByStudent(ctx context.Context, studentID int) ([]model.EnrolledCourse, error)
}

type enrollmentPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewEnrollmentRepository(pool *pgxpool.Pool) EnrollmentRepository {
	return &enrollmentPostgresRepository{pool: pool}
}

func (r *enrollmentPostgresRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *enrollmentPostgresRepository) LockCourseForUpdate(
	ctx context.Context, tx pgx.Tx, courseID int,
) (model.Course, error) {
	var c model.Course
	err := tx.QueryRow(ctx,
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota, created_at, updated_at
		 FROM courses WHERE id = $1 FOR UPDATE`,
		courseID,
	).Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Course{}, ErrNotFound
		}
		return model.Course{}, fmt.Errorf("mengunci course: %w", err)
	}
	return c, nil
}

func (r *enrollmentPostgresRepository) CountByCourseAndYear(
	ctx context.Context, tx pgx.Tx, courseID int, tahun string,
) (int, error) {
	var count int
	err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM enrollments
		 WHERE course_id = $1 AND tahun_akademik = $2`,
		courseID, tahun,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("menghitung enrollment: %w", err)
	}
	return count, nil
}

func (r *enrollmentPostgresRepository) Exists(
	ctx context.Context, tx pgx.Tx, studentID, courseID int, tahun string,
) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM enrollments
			WHERE student_id = $1 AND course_id = $2 AND tahun_akademik = $3
		 )`,
		studentID, courseID, tahun,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("mengecek duplikasi enrollment: %w", err)
	}
	return exists, nil
}

func (r *enrollmentPostgresRepository) SumSKSForStudentYear(
	ctx context.Context, tx pgx.Tx, studentID int, tahun string,
) (int, error) {
	var total *int
	err := tx.QueryRow(ctx,
		`SELECT SUM(c.sks)
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1 AND e.tahun_akademik = $2`,
		studentID, tahun,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("menghitung total SKS: %w", err)
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

func (r *enrollmentPostgresRepository) Insert(
	ctx context.Context, tx pgx.Tx, studentID, courseID int, tahun string,
) (model.Enrollment, error) {
	var e model.Enrollment
	err := tx.QueryRow(ctx,
		`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, student_id, course_id, tahun_akademik, created_at`,
		studentID, courseID, tahun,
	).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik, &e.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Enrollment{}, ErrDuplicate
		}
		return model.Enrollment{}, fmt.Errorf("menyimpan enrollment: %w", err)
	}
	return e, nil
}

func (r *enrollmentPostgresRepository) FindByID(ctx context.Context, id int) (model.Enrollment, error) {
	var e model.Enrollment
	err := r.pool.QueryRow(ctx,
		`SELECT id, student_id, course_id, tahun_akademik, created_at
		 FROM enrollments WHERE id = $1`,
		id,
	).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik, &e.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Enrollment{}, ErrNotFound
		}
		return model.Enrollment{}, fmt.Errorf("mencari enrollment: %w", err)
	}
	return e, nil
}

func (r *enrollmentPostgresRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM enrollments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("menghapus enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *enrollmentPostgresRepository) ListEnrichedByStudent(
	ctx context.Context, studentID int,
) ([]model.EnrolledCourse, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT e.id, e.course_id, c.kode_mk, c.nama_mk, c.sks, c.semester, e.tahun_akademik
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1
		 ORDER BY e.created_at DESC`,
		studentID,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar enrollment: %w", err)
	}
	defer rows.Close()

	result := []model.EnrolledCourse{}
	for rows.Next() {
		var ec model.EnrolledCourse
		if err := rows.Scan(&ec.EnrollmentID, &ec.CourseID, &ec.KodeMK,
			&ec.NamaMK, &ec.SKS, &ec.Semester, &ec.TahunAkademik); err != nil {
			return nil, fmt.Errorf("membaca baris enrollment: %w", err)
		}
		result = append(result, ec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query: %w", err)
	}
	return result, nil
}