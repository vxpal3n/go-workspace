package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"tugas09/app/model"
	"tugas09/helper"
)

const (
	baseURL    = "http://localhost:3000/api/v1"
	loadPrefix = "M"
)

func main() {
	mode := flag.String("mode", "crowd", "crowd | credits")
	flag.Parse()

	_ = godotenv.Load()
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fatal("connect db", err)
	}
	defer pool.Close()

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		fatal("JWT_SECRET", fmt.Errorf("kurang dari 32 karakter"))
	}
	jwt := helper.NewJWTManager(secret, os.Getenv("JWT_ISSUER"), 60*time.Minute)

	students, err := loadStudents(ctx, pool)
	if err != nil {
		fatal("load students", err)
	}
	if len(students) < 60 {
		fatal("load students", fmt.Errorf("hanya %d dari 60 ditemukan", len(students)))
	}
	fmt.Printf("Dimuat %d students, membuat %d token...\n", len(students), len(students))

	tokens := make([]string, len(students))
	for i, s := range students {
		tok, err := jwt.GenerateAccess(s)
		if err != nil {
			fatal("generate token", err)
		}
		tokens[i] = tok
	}

	switch *mode {
	case "crowd":
		runCrowd(ctx, pool, tokens)
	case "credits":
		runCredits(ctx, pool, students, tokens)
	default:
		fatal("mode", fmt.Errorf("tidak dikenal: %s", *mode))
	}
}

func loadStudents(ctx context.Context, pool *pgxpool.Pool) ([]model.Student, error) {
	rows, err := pool.Query(ctx,
		"SELECT id, nim, name, role FROM students WHERE nim LIKE $1 ORDER BY nim LIMIT 60",
		loadPrefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Role); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func findCourseID(ctx context.Context, pool *pgxpool.Pool, code string) (int, error) {
	var id int
	err := pool.QueryRow(ctx, "SELECT id FROM courses WHERE code = $1", code).Scan(&id)
	return id, err
}

func runCrowd(ctx context.Context, pool *pgxpool.Pool, tokens []string) {
	courseID, err := findCourseID(ctx, pool, "IF399")
	if err != nil {
		fatal("find IF399", err)
	}

	fmt.Printf("60 students berebut IF399 (course_id=%d) secara serentak...\n", courseID)

	start := time.Now()
	tally := fire(len(tokens), func(i int) string {
		return doEnroll(tokens[i], courseID)
	})
	elapsed := time.Since(start)

	fmt.Println("\nHasil request:")
	printTally(tally)
	fmt.Printf("Waktu total: %v\n", elapsed)

	reportConsistency(ctx, pool, courseID, "IF399")
}

func runCredits(ctx context.Context, pool *pgxpool.Pool, students []model.Student, tokens []string) {
	rows, err := pool.Query(ctx,
		"SELECT id FROM courses WHERE code <> 'IF399' ORDER BY code LIMIT 12")
	if err != nil {
		fatal("load courses", err)
	}
	var courseIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		courseIDs = append(courseIDs, id)
	}
	rows.Close()

	type job struct {
		token    string
		courseID int
	}
	var jobs []job
	for i := 0; i < 10; i++ {
		for _, cid := range courseIDs {
			jobs = append(jobs, job{tokens[i], cid})
		}
	}

	fmt.Printf("10 students × %d courses serentak (batas 24 SKS = 8 courses)...\n",
		len(courseIDs))

	start := time.Now()
	tally := fire(len(jobs), func(i int) string {
		return doEnroll(jobs[i].token, jobs[i].courseID)
	})
	elapsed := time.Since(start)

	fmt.Println("\nHasil request:")
	printTally(tally)
	fmt.Printf("Waktu total: %v\n", elapsed)

	fmt.Println("\nCek total SKS per student (batas 24):")
	over := 0
	for i := 0; i < 10; i++ {
		var total int
		_ = pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(c.credits), 0)
			 FROM enrollments e JOIN courses c ON c.id = e.course_id
			 WHERE e.student_id = $1`, students[i].ID).Scan(&total)
		mark := "OK"
		if total > 24 {
			mark = "MELEBIHI"
			over++
		}
		fmt.Printf("  %s total SKS = %2d [%s]\n", students[i].NIM, total, mark)
	}
	if over == 0 {
		fmt.Printf("\nseluruh 10 mahasiswa <= 24 SKS -> KONSISTEN\n")
	} else {
		fmt.Printf("\n%d mahasiswa MELEBIHI 24 SKS -> TIDAK KONSISTEN\n", over)
	}
}

func doEnroll(token string, courseID int) string {
	url := fmt.Sprintf("%s/courses/%d/enroll", baseURL, courseID)
	req, _ := http.NewRequest("POST", url, strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "ERR_NETWORK"
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	switch resp.StatusCode {
	case 201:
		return "201 OK"
	case 409:
		return "409 CONFLICT"
	case 500:
		return "500 INTERNAL"
	default:
		return fmt.Sprintf("HTTP_%d", resp.StatusCode)
	}
}

func fire(n int, job func(i int) string) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			results <- job(i)
		}(i)
	}

	close(start)
	wg.Wait()
	close(results)

	tally := map[string]int{}
	for r := range results {
		tally[r]++
	}
	return tally
}

func reportConsistency(ctx context.Context, pool *pgxpool.Pool, courseID int, code string) {
	var quota, enrolled, rows int
	_ = pool.QueryRow(ctx,
		"SELECT quota, enrolled_count FROM courses WHERE id = $1", courseID).
		Scan(&quota, &enrolled)
	_ = pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM enrollments WHERE course_id = $1", courseID).
		Scan(&rows)

	verdict := "TIDAK KONSISTEN"
	if enrolled == rows && enrolled <= quota {
		verdict = "KONSISTEN"
	}
	fmt.Printf("%s quota=%d enrolled_count=%d baris enrollments=%d -> %s\n",
		code, quota, enrolled, rows, verdict)
}

func printTally(tally map[string]int) {
	for status, count := range tally {
		fmt.Printf("  %-24s %d\n", status, count)
	}
}

func fatal(where string, err error) {
	fmt.Fprintf(os.Stderr, "[%s] %v\n", where, err)
	os.Exit(1)
}