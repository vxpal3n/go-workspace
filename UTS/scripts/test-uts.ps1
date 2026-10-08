# ============================================================================
# test-uts.ps1 — Uji end-to-end SIAKAD Mini (v2)
# ============================================================================

$ErrorActionPreference = "Stop"
$B = "http://localhost:3000/api/v1"
$J = @{ "Content-Type" = "application/json" }

function Hr($title) {
    Write-Host ""
    Write-Host "=======================================================" -ForegroundColor Cyan
    Write-Host $title -ForegroundColor Cyan
    Write-Host "=======================================================" -ForegroundColor Cyan
}

function ShowResult($label, $expect, $actual) {
    $ok = $expect -eq $actual
    $color = if ($ok) { "Green" } else { "Red" }
    $mark = if ($ok) { "[OK]  " } else { "[FAIL]" }
    Write-Host ("{0} {1,-60} harap={2,-6} dapat={3}" -f $mark, $label, $expect, $actual) -ForegroundColor $color
}

function StatusCode($method, $url, $headers, $body = $null) {
    $params = @{
        Method = $method
        Uri = $url
        Headers = $headers
        ErrorAction = "Stop"
    }
    if ($null -ne $body) {
        $params["Body"] = ($body | ConvertTo-Json -Compress)
        $params["ContentType"] = "application/json"
    }
    try {
        return [int](Invoke-WebRequest @params).StatusCode
    } catch {
        if ($_.Exception.Response) {
            return [int]$_.Exception.Response.StatusCode
        }
        return -1
    }
}

function GetJSON($method, $url, $headers, $body = $null) {
    $params = @{
        Method = $method
        Uri = $url
        Headers = $headers
        ErrorAction = "Stop"
    }
    if ($null -ne $body) {
        $params["Body"] = ($body | ConvertTo-Json -Compress)
        $params["ContentType"] = "application/json"
    }
    return Invoke-RestMethod @params
}

# ============================================================================
# 0. Reset (bersihkan sisa data run sebelumnya)
# ============================================================================
Hr "0. RESET (hapus enrollment + mahasiswa uji)"

$envFile = Join-Path (Split-Path $PSScriptRoot -Parent) ".env"
if ((Test-Path $envFile) -and -not $env:PGPASSWORD) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*DB_PASSWORD\s*=\s*(.+?)\s*$') {
            $env:PGPASSWORD = $Matches[1]
        }
    }
}

# 0.1 Bersihkan enrollment mahasiswa uji tetap
psql -U postgres -d siakad_mini -c `
  "DELETE FROM enrollments WHERE student_id IN (1, 2);" | Out-Null

# 0.2 Bersihkan sisa mahasiswa uji (termasuk yang soft-deleted dari run sebelumnya).
#     Urutan penting: enrollments -> students -> users (FK).
psql -U postgres -d siakad_mini -c `
  "DELETE FROM enrollments WHERE student_id IN (SELECT id FROM students WHERE nim = '187221000099');" | Out-Null

psql -U postgres -d siakad_mini -c `
  "DELETE FROM students WHERE nim = '187221000099';" | Out-Null

psql -U postgres -d siakad_mini -c `
  "DELETE FROM users WHERE email = '187221000099@student.siakad.test';" | Out-Null

Write-Host "  [OK]  Data mahasiswa uji dibersihkan" -ForegroundColor Green

# ============================================================================
# 1. Login
# ============================================================================
Hr "1. LOGIN"

$adminLogin = GetJSON "POST" "$B/auth/login" $J @{
    email = "admin@siakad.test"; password = "admin12345"
}
$TOKEN_ADMIN = $adminLogin.data.access_token

$mhsLogin = GetJSON "POST" "$B/auth/login" $J @{
    email = "187221000001@student.siakad.test"; password = "187221000001"
}
$TOKEN_MHS = $mhsLogin.data.access_token
$MHS_USER_ID = $mhsLogin.data.user.id  # = 2

# Ambil student_id milik mahasiswa ini dari /auth/me (kita butuh untuk ownership test)
$meMhs = GetJSON "GET" "$B/auth/me" @{ "Authorization" = "Bearer $TOKEN_MHS" }
$MHS_STUDENT_ID = $null

# Kita dapat id dengan menanyakan langsung via psql
$MHS_STUDENT_ID = (psql -U postgres -d siakad_mini -tA -c `
  "SELECT id FROM students WHERE user_id = $MHS_USER_ID;").Trim()
Write-Host "  Mahasiswa login: user_id=$MHS_USER_ID, student_id=$MHS_STUDENT_ID" -ForegroundColor Yellow

# Mahasiswa lain (bukan dirinya) — kita pakai student_id = 2
$OTHER_STUDENT_ID = 2
if ($MHS_STUDENT_ID -eq "2") { $OTHER_STUDENT_ID = 3 }

Write-Host "  Admin login -> 200, role = $($adminLogin.data.user.role)" -ForegroundColor Green
Write-Host "  Mahasiswa login -> 200, role = $($mhsLogin.data.user.role)" -ForegroundColor Green

$st = StatusCode "POST" "$B/auth/login" $J @{
    email = "admin@siakad.test"; password = "salahbanget"
}
ShowResult "Login password salah" 401 $st

$st = StatusCode "POST" "$B/auth/login" $J @{
    email = "bukan-email"; password = "admin12345"
}
ShowResult "Login email format salah" 422 $st

# ============================================================================
# 2. Authorization Matrix
# ============================================================================
Hr "2. AUTHORIZATION MATRIX"

$H_AUTH_ADMIN = @{ "Authorization" = "Bearer $TOKEN_ADMIN"; "Content-Type" = "application/json" }
$H_AUTH_MHS   = @{ "Authorization" = "Bearer $TOKEN_MHS";   "Content-Type" = "application/json" }

$st = StatusCode "GET" "$B/auth/me" $H_AUTH_ADMIN
ShowResult "GET /auth/me (admin)" 200 $st
$st = StatusCode "GET" "$B/auth/me" $H_AUTH_MHS
ShowResult "GET /auth/me (mahasiswa)" 200 $st
$st = StatusCode "GET" "$B/auth/me" @{}
ShowResult "GET /auth/me (no token)" 401 $st

$st = StatusCode "GET" "$B/students" $H_AUTH_ADMIN
ShowResult "GET /students (admin)" 200 $st
$st = StatusCode "GET" "$B/students" $H_AUTH_MHS
ShowResult "GET /students (mahasiswa -> 403)" 403 $st

$st = StatusCode "GET" "$B/students/$MHS_STUDENT_ID" $H_AUTH_MHS
ShowResult "GET /students/$MHS_STUDENT_ID (mahasiswa self)" 200 $st

$st = StatusCode "GET" "$B/students/$OTHER_STUDENT_ID" $H_AUTH_MHS
ShowResult "GET /students/$OTHER_STUDENT_ID (mahasiswa other -> 403)" 403 $st

$st = StatusCode "GET" "$B/students/$MHS_STUDENT_ID" $H_AUTH_ADMIN
ShowResult "GET /students/$MHS_STUDENT_ID (admin)" 200 $st

$st = StatusCode "GET" "$B/students/99999" $H_AUTH_ADMIN
ShowResult "GET /students/99999 (not found -> 404)" 404 $st

$st = StatusCode "PUT" "$B/students/$MHS_STUDENT_ID" $H_AUTH_ADMIN @{
    nama = "Rina Putri Updated"; prodi = "Sistem Informasi"; angkatan = 2022; ipk_terakhir = 3.55
}
ShowResult "PUT /students/$MHS_STUDENT_ID (admin)" 200 $st

$st = StatusCode "PUT" "$B/students/$MHS_STUDENT_ID" $H_AUTH_MHS @{
    nama = "Hacked"; prodi = "X"; angkatan = 2022
}
ShowResult "PUT /students/$MHS_STUDENT_ID (mahasiswa -> 403)" 403 $st

# ============================================================================
# 3. GET /courses
# ============================================================================
Hr "3. GET /courses"

$st = StatusCode "GET" "$B/courses" $H_AUTH_ADMIN
ShowResult "GET /courses (admin)" 200 $st
$st = StatusCode "GET" "$B/courses" $H_AUTH_MHS
ShowResult "GET /courses (mahasiswa)" 200 $st
$st = StatusCode "GET" "$B/courses" @{}
ShowResult "GET /courses (no token -> 401)" 401 $st

$courses = GetJSON "GET" "$B/courses" $H_AUTH_MHS
Write-Host "`n  Daftar 3 mata kuliah pertama:" -ForegroundColor Yellow
$courses.data | Select-Object -First 3 | Format-Table id, kode_mk, nama_mk, sks, kuota, terisi, sisa_kuota

$coursesAvailable = GetJSON "GET" "$B/courses?available=true" $H_AUTH_MHS
Write-Host "`n  Dengan available=true -> $($coursesAvailable.data.Count) mata kuliah" -ForegroundColor Yellow

# ============================================================================
# 4. POST /enrollments — Business Rules
# ============================================================================
Hr "4. BUSINESS RULES — ENROLLMENT"

$courseId1 = $courses.data[0].id
$courseId2 = $courses.data[1].id

# 4.1 Berhasil mengambil (setelah reset, pasti sukses)
$st = StatusCode "POST" "$B/enrollments" $H_AUTH_MHS @{
    course_id = $courseId1; tahun_akademik = "2026/2027-Ganjil"
}
ShowResult "POST /enrollments (valid pertama kali)" 201 $st

# 4.2 Duplikasi
$st = StatusCode "POST" "$B/enrollments" $H_AUTH_MHS @{
    course_id = $courseId1; tahun_akademik = "2026/2027-Ganjil"
}
ShowResult "POST /enrollments (duplikat -> 409)" 409 $st

# 4.3 Tahun akademik salah format
$st = StatusCode "POST" "$B/enrollments" $H_AUTH_MHS @{
    course_id = $courseId2; tahun_akademik = "2026-2027"
}
ShowResult "POST /enrollments (tahun akademik salah -> 422)" 422 $st

# 4.4 Course tidak ada
$st = StatusCode "POST" "$B/enrollments" $H_AUTH_MHS @{
    course_id = 99999; tahun_akademik = "2026/2027-Ganjil"
}
ShowResult "POST /enrollments (course_id tidak ada -> 422)" 422 $st

# 4.5 Admin tidak boleh akses
$st = StatusCode "POST" "$B/enrollments" $H_AUTH_ADMIN @{
    course_id = $courseId2; tahun_akademik = "2026/2027-Ganjil"
}
ShowResult "POST /enrollments (admin -> 403)" 403 $st

# 4.6 Cek detail mahasiswa
$detail = GetJSON "GET" "$B/students/$MHS_STUDENT_ID" $H_AUTH_MHS
Write-Host "`n  Detail mahasiswa setelah 1 enrollment:" -ForegroundColor Yellow
Write-Host "    NIM        = $($detail.data.nim)"
Write-Host "    Nama       = $($detail.data.nama)"
Write-Host "    IPK        = $($detail.data.ipk_terakhir)"
Write-Host "    Total SKS  = $($detail.data.total_sks)"
Write-Host "    Batas SKS  = $($detail.data.batas_sks)"
Write-Host "    Mata kuliah = $($detail.data.mata_kuliah.Count)"

$ENROLLMENT_ID = $detail.data.mata_kuliah[0].enrollment_id

# ============================================================================
# 5. DELETE /enrollments/{id}
# ============================================================================
Hr "5. DELETE /enrollments"

$st = StatusCode "DELETE" "$B/enrollments/$ENROLLMENT_ID" $H_AUTH_ADMIN
ShowResult "DELETE /enrollments (admin -> 403)" 403 $st

$st = StatusCode "DELETE" "$B/enrollments/$ENROLLMENT_ID" $H_AUTH_MHS
ShowResult "DELETE /enrollments (mahasiswa milik sendiri -> 204)" 204 $st

$st = StatusCode "DELETE" "$B/enrollments/$ENROLLMENT_ID" $H_AUTH_MHS
ShowResult "DELETE /enrollments (id sama lagi -> 404)" 404 $st

# ============================================================================
# 6. Soft Delete Student
# ============================================================================
Hr "6. SOFT DELETE STUDENT"

$newEmail = "187221000099@student.siakad.test"
$st = StatusCode "POST" "$B/students" $H_AUTH_ADMIN @{
    nim = "187221000099"; nama = "Test Student"; email = $newEmail
    prodi = "Sistem Informasi"; angkatan = 2024; ipk_terakhir = 3.20
}
ShowResult "POST /students (admin buat baru -> 201)" 201 $st

$list = GetJSON "GET" "$B/students?search=187221000099" $H_AUTH_ADMIN
$NEW_STUDENT_ID = $list.data[0].id
Write-Host "  Mahasiswa baru ID = $NEW_STUDENT_ID" -ForegroundColor Yellow

$st = StatusCode "DELETE" "$B/students/$NEW_STUDENT_ID" $H_AUTH_ADMIN
ShowResult "DELETE /students/$NEW_STUDENT_ID (soft delete -> 204)" 204 $st

$st = StatusCode "POST" "$B/auth/login" $J @{
    email = $newEmail; password = "187221000099"
}
ShowResult "Login mahasiswa terhapus -> 401" 401 $st

$listAfter = GetJSON "GET" "$B/students?search=187221000099" $H_AUTH_ADMIN
$count = if ($listAfter.data) { $listAfter.data.Count } else { 0 }
ShowResult "Mahasiswa terhapus tidak muncul di list" 0 $count

$st = StatusCode "GET" "$B/students/$NEW_STUDENT_ID" $H_AUTH_ADMIN
ShowResult "GET /students/$NEW_STUDENT_ID (soft deleted -> 404)" 404 $st

# ============================================================================
# 7. Rate Limiting
# ============================================================================
Hr "7. RATE LIMITING (6x login gagal)"

$codes = @()
for ($i = 0; $i -lt 6; $i++) {
    $st = StatusCode "POST" "$B/auth/login" $J @{
        email = "admin@siakad.test"; password = "salahterus"
    }
    $codes += $st
}
Write-Host "  Hasil: $($codes -join ', ')" -ForegroundColor Yellow
$count429 = ($codes | Where-Object { $_ -eq 429 }).Count
if ($count429 -ge 1) {
    Write-Host "  [OK]  Rate limiter bekerja ($count429 request di-429)" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Rate limiter tidak bekerja" -ForegroundColor Red
}

# ============================================================================
# Ringkasan
# ============================================================================
Hr "SELESAI"
Write-Host "  Semua pengujian selesai." -ForegroundColor Green