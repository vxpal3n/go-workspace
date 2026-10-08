package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"UTS/app/model"
)

type CourseRepository interface {
	List(ctx context.Context, q model.ListCoursesQuery) ([]model.CourseWithQuota, error)
	FindByID(ctx context.Context, id int) (model.Course, error)
}

type coursePostgresRepository struct {
	pool *pgxpool.Pool
}

func NewCourseRepository(pool *pgxpool.Pool) CourseRepository {
	return &coursePostgresRepository{pool: pool}
}

func (r *coursePostgresRepository) List(
	ctx context.Context, q model.ListCoursesQuery,
) ([]model.CourseWithQuota, error) {
	where := " WHERE 1 = 1"
	args := []any{}

	if q.Semester > 0 {
		args = append(args, q.Semester)
		where += fmt.Sprintf(" AND c.semester = $%d", len(args))
	}
	if q.Search != "" {
		args = append(args, "%"+q.Search+"%")
		where += fmt.Sprintf(" AND (c.kode_mk ILIKE $%d OR c.nama_mk ILIKE $%d)", len(args), len(args))
	}
	if q.Available {
		where += ` AND (SELECT COUNT(*) FROM enrollments e WHERE e.course_id = c.id) < c.kuota`
	}

	query := `
		SELECT
			c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota,
			COALESCE(sub.terisi, 0) AS terisi,
			c.kuota - COALESCE(sub.terisi, 0) AS sisa_kuota
		FROM courses c
		LEFT JOIN (
			SELECT course_id, COUNT(*) AS terisi
			FROM enrollments
			GROUP BY course_id
		) sub ON sub.course_id = c.id` + where + `
		ORDER BY c.semester ASC, c.kode_mk ASC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar course: %w", err)
	}
	defer rows.Close()

	result := []model.CourseWithQuota{}
	for rows.Next() {
		var c model.CourseWithQuota
		if err := rows.Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester,
			&c.Kuota, &c.Terisi, &c.SisaKuota); err != nil {
			return nil, fmt.Errorf("membaca baris course: %w", err)
		}
		result = append(result, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query: %w", err)
	}
	return result, nil
}

func (r *coursePostgresRepository) FindByID(ctx context.Context, id int) (model.Course, error) {
	var c model.Course
	err := r.pool.QueryRow(ctx,
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota, created_at, updated_at
		 FROM courses WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Course{}, ErrNotFound
		}
		return model.Course{}, fmt.Errorf("mencari course: %w", err)
	}
	return c, nil
}