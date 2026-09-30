package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"tugas07/app/model"
)

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	FindByNIM(ctx context.Context, nim string) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	CreateWithOwner(ctx context.Context, s model.Student, ownerID int) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	UpdateRole(ctx context.Context, id int, role string) (model.Student, error)
	Delete(ctx context.Context, id int) error
	FindAfterCursor(ctx context.Context, q model.CursorQuery) ([]model.Student, error)
}

type studentPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgresRepository{pool: pool}
}

const studentColumns = `id, nim, name, email, grade, password, role, owner_id, is_active, created_at`

func scanStudent(row pgx.Row) (model.Student, error) {
	var s model.Student
	err := row.Scan(
		&s.ID, &s.NIM, &s.Name, &s.Email, &s.Grade,
		&s.Password, &s.Role, &s.OwnerID, &s.IsActive, &s.CreatedAt,
	)
	return s, err
}

var sortColumns = map[string]string{
	"id":         "id",
	"nim":        "nim",
	"name":       "name",
	"grade":      "grade",
	"created_at": "created_at",
}

func buildFilter(q model.ListQuery) (string, []interface{}) {
	where := " WHERE 1 = 1"
	args := []interface{}{}

	if q.Search != "" {
		where += fmt.Sprintf(
			" AND (LOWER(name) ILIKE $%d OR LOWER(nim) ILIKE $%d)",
			len(args)+1, len(args)+1,
		)
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
	}
	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}
	if q.MinGrade != nil {
		where += fmt.Sprintf(" AND grade >= $%d", len(args)+1)
		args = append(args, *q.MinGrade)
	}
	if q.MaxGrade != nil {
		where += fmt.Sprintf(" AND grade <= $%d", len(args)+1)
		args = append(args, *q.MaxGrade)
	}
	return where, args
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *studentPostgresRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	var total int
	countSQL := "SELECT COUNT(*) FROM students" + where
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("menghitung student: %w", err)
	}

	sortCol, ok := sortColumns[q.Sort]
	if !ok {
		sortCol = "id"
	}
	order := "ASC"
	if q.Order == "desc" {
		order = "DESC"
	}

	sql := fmt.Sprintf(
		`SELECT `+studentColumns+` FROM students%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		where, sortCol, order, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil daftar student: %w", err)
	}
	defer rows.Close()

	students := []model.Student{}
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("membaca baris student: %w", err)
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error setelah iterasi rows: %w", err)
	}
	return students, total, nil
}

func (r *studentPostgresRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	s, err := scanStudent(r.pool.QueryRow(ctx,
		`SELECT `+studentColumns+` FROM students WHERE id = $1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student: %w", err)
	}
	return s, nil
}

func (r *studentPostgresRepository) FindByNIM(ctx context.Context, nim string) (model.Student, error) {
	s, err := scanStudent(r.pool.QueryRow(ctx,
		`SELECT `+studentColumns+` FROM students WHERE LOWER(nim) = LOWER($1)`,
		nim))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student berdasarkan NIM: %w", err)
	}
	return s, nil
}

func (r *studentPostgresRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO students (nim, name, email, grade, password, role, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING `+studentColumns,
		s.NIM, s.Name, s.Email, s.Grade, s.Password, s.Role, s.IsActive,
	)
	created, err := scanStudent(row)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan student: %w", err)
	}
	return created, nil
}

func (r *studentPostgresRepository) CreateWithOwner(ctx context.Context, s model.Student, ownerID int) (model.Student, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO students (nim, name, email, grade, password, role, owner_id, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING `+studentColumns,
		s.NIM, s.Name, s.Email, s.Grade, s.Password, s.Role, ownerID, s.IsActive,
	)
	created, err := scanStudent(row)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan student dengan owner: %w", err)
	}
	return created, nil
}

func (r *studentPostgresRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE students
		 SET nim = $1, name = $2, email = $3, grade = $4, is_active = $5
		 WHERE id = $6
		 RETURNING `+studentColumns,
		s.NIM, s.Name, s.Email, s.Grade, s.IsActive, s.ID,
	)
	updated, err := scanStudent(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("memperbarui student: %w", err)
	}
	return updated, nil
}

func (r *studentPostgresRepository) UpdateRole(ctx context.Context, id int, role string) (model.Student, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE students SET role = $1 WHERE id = $2
		 RETURNING `+studentColumns,
		role, id,
	)
	updated, err := scanStudent(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengubah role student: %w", err)
	}
	return updated, nil
}

func (r *studentPostgresRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("menghapus student: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *studentPostgresRepository) FindAfterCursor(ctx context.Context, q model.CursorQuery,) ([]model.Student, error) {
	args := []any{}
	where := " WHERE 1 = 1"
	if q.Search != "" {
		args = append(args, "%"+q.Search+"%")
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}
	if q.IsActive != nil {
		args = append(args, *q.IsActive)
		where += fmt.Sprintf(" AND is_active = $%d", len(args))
	}
	if q.After != nil {
		args = append(args, q.After.CreatedAt, q.After.ID)
		where += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)",
			len(args)-1, len(args))
	}
	args = append(args, q.Limit+1)
	query := fmt.Sprintf(
		"SELECT %s FROM students%s ORDER BY created_at DESC, id DESC LIMIT $%d",
		studentColumns, where, len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar student: %w", err)
	}
	defer rows.Close()

	result := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Email, &s.Grade,
			&s.Password, &s.Role, &s.OwnerID, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("membaca row student: %w", err)
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query: %w", err)
	}
	return result, nil
}