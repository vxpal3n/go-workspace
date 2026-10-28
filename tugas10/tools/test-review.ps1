# ============================================================
# Test End-to-End Review Endpoints (Modul 10)
# ============================================================
$baseUrl = "http://localhost:3000/api/v1"

function Api {
    param(
        [string]$Method, [string]$Uri,
        [hashtable]$Headers = @{}, [string]$Body = $null
    )
    $p = @{
        Uri = $Uri; Method = $Method; Headers = $Headers
        ContentType = "application/json"; UseBasicParsing = $true
        ErrorAction = "Stop"
    }
    if ($Body) { $p.Body = $Body }
    try {
        $r = Invoke-WebRequest @p
        @{ Code = [int]$r.StatusCode; Body = $r.Content; Headers = $r.Headers }
    } catch {
        $code = 0
        if ($_.Exception.Response) { $code = [int]$_.Exception.Response.StatusCode }
        $errBody = $_.ErrorDetails.Message
        if (-not $errBody -and $_.Exception.Response) {
            try {
                $errBody = (New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())).ReadToEnd()
            } catch { }
        }
        @{ Code = $code; Body = $errBody; Headers = @{} }
    }
}

function Show($label, $expect, $actual) {
    $ok = $expect -eq $actual
    $color = if ($ok) { "Green" } else { "Red" }
    Write-Host ("{0} {1,-55} harap={2,-4} dapat={3}" -f
        ($(if ($ok) {"[OK]  "} else {"[FAIL]"})), $label, $expect, $actual) -ForegroundColor $color
}

# ============================================================
# Persiapan: register dua mahasiswa baru + login
# ============================================================
Write-Host "`n=== Persiapan akun ===" -ForegroundColor Cyan

function RegisterUser($nim, $name, $email) {
    $body = @{
        nim = $nim; name = $name; email = $email
        grade = 80; password = "Rahasia123"
    } | ConvertTo-Json -Compress
    $r = Api POST "$baseUrl/auth/register" @{} $body
    if ($r.Code -eq 201) {
        Write-Host "  [OK] Register $nim" -ForegroundColor Green
    } elseif ($r.Code -eq 409) {
        Write-Host "  [INFO] $nim sudah terdaftar" -ForegroundColor Yellow
    } else {
        # PERBAIKAN: Menggunakan ${nim} bukan $nim:
        Write-Host "  [FAIL] Register ${nim}: $($r.Code) - $($r.Body)" -ForegroundColor Red
    }
}

RegisterUser "T200" "ReviewTester1" "t200@test.com"
RegisterUser "T201" "ReviewTester2" "t201@test.com"

function Login($nim, $pw) {
    $body = @{ nim = $nim; password = $pw } | ConvertTo-Json -Compress
    $r = Api POST "$baseUrl/auth/login" @{} $body
    if ($r.Code -ne 200) {
        # PERBAIKAN: Menggunakan ${nim} bukan $nim:
        Write-Host "  [FAIL] Login ${nim}: $($r.Code) - $($r.Body)" -ForegroundColor Red
        return $null
    }
    return ($r.Body | ConvertFrom-Json).data.access_token
}

Start-Sleep -Seconds 1
$token1 = Login "T200" "Rahasia123"

Start-Sleep -Seconds 1
$token2 = Login "T201" "Rahasia123"

if (-not $token1 -or -not $token2) {
    Write-Host "`n❌ Login gagal! Jika error 429 (Too Many Requests), TUNGGU 60 DETIK lalu jalankan ulang." -ForegroundColor Red
    return
}
Write-Host "✅ Login berhasil" -ForegroundColor Green

$h1 = @{ Authorization = "Bearer $token1" }
$h2 = @{ Authorization = "Bearer $token2" }

# ============================================================
# Uji 1: Tidak enrolled → 403 NOT_ENROLLED
# ============================================================
Write-Host "`n=== Uji 1: Not enrolled ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/courses/1/reviews" $h1 (@{
    rating = 5; comment = "Materi sangat jelas dan menarik"
    tags = @("mudah")
} | ConvertTo-Json -Compress)
Show "T200 review tanpa enroll" 403 $r.Code

# ============================================================
# Uji 2: Enroll T200 ke IF101 (course_id=1)
# ============================================================
Write-Host "`n=== Uji 2: Enroll T200 ke IF101 ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/courses/1/enroll" $h1 "{}"
Show "T200 enroll IF101" 201 $r.Code

# ============================================================
# Uji 3: Create review sukses → 201
# ============================================================
Write-Host "`n=== Uji 3: Create review ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/courses/1/reviews" $h1 (@{
    rating = 5
    comment = "Materi Go-nya jelas, tugasnya cukup banyak"
    tags = @("Tugas-Banyak", " dosen-jelas ", "tugas-banyak")
} | ConvertTo-Json -Compress)

# Fallback jika review sudah ada dari run sebelumnya
if ($r.Code -eq 409) {
    Write-Host "  [INFO] Review T200 sudah ada, mengambil ID yang ada..." -ForegroundColor Yellow
    $listR = Api GET "$baseUrl/courses/1/reviews" $h1
    $listObj = $listR.Body | ConvertFrom-Json
    $existing = $listObj.data | Where-Object { $_.nim -eq "T200" } | Select-Object -First 1
    $reviewID = $existing.id
    Write-Host "  Review ID: $reviewID" -ForegroundColor Yellow
} else {
    Show "T200 create review" 201 $r.Code
    $created = $r.Body | ConvertFrom-Json
    $reviewID = $created.data.id
    Write-Host "  Review ID: $reviewID" -ForegroundColor Yellow
    Write-Host "  Tags tersimpan: $($created.data.tags -join ', ')" -ForegroundColor Yellow
}

# ============================================================
# Uji 4: Duplikasi review → 409 ALREADY_REVIEWED
# ============================================================
Write-Host "`n=== Uji 4: Duplikasi review ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/courses/1/reviews" $h1 (@{
    rating = 4; comment = "Review kedua (harus gagal)"
} | ConvertTo-Json -Compress)
Show "T200 review kedua di kelas sama" 409 $r.Code

# ============================================================
# Uji 5: GET daftar review → 200
# ============================================================
Write-Host "`n=== Uji 5: List reviews ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/courses/1/reviews" $h1
Show "GET /courses/1/reviews" 200 $r.Code

# ============================================================
# Uji 6: Filter by tag → 200
# ============================================================
Write-Host "`n=== Uji 6: Filter by tag ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/courses/1/reviews?tag=tugas-banyak" $h1
Show "Filter tag=tugas-banyak" 200 $r.Code

# ============================================================
# Uji 7: Summary → 200
# ============================================================
Write-Host "`n=== Uji 7: Summary ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/courses/1/reviews/summary" $h1
Show "GET summary" 200 $r.Code

# ============================================================
# Uji 8: T201 helpful review T200 → 200
# ============================================================
Write-Host "`n=== Uji 8: Mark helpful ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/reviews/$reviewID/helpful" $h2 "{}"
Show "T201 helpful review T200" 200 $r.Code

# ============================================================
# Uji 9: T201 helpful kedua kali → 409 ALREADY_VOTED
# ============================================================
Write-Host "`n=== Uji 9: Helpful kedua kali ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/reviews/$reviewID/helpful" $h2 "{}"
Show "T201 helpful kedua kali" 409 $r.Code

# ============================================================
# Uji 10: T200 helpful dirinya sendiri → 409 OWN_REVIEW
# ============================================================
Write-Host "`n=== Uji 10: Self helpful ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/reviews/$reviewID/helpful" $h1 "{}"
Show "T200 self helpful" 409 $r.Code

# ============================================================
# Uji 11: T201 update review T200 → 403 NotOwner
# ============================================================
Write-Host "`n=== Uji 11: Update orang lain ===" -ForegroundColor Cyan
$r = Api PATCH "$baseUrl/reviews/$reviewID" $h2 (@{ rating = 1 } | ConvertTo-Json -Compress)
Show "T201 update review T200" 403 $r.Code

# ============================================================
# Uji 12: T200 update sendiri → 200
# ============================================================
Write-Host "`n=== Uji 12: Update sendiri ===" -ForegroundColor Cyan
$r = Api PATCH "$baseUrl/reviews/$reviewID" $h1 (@{
    rating = 4
    comment = "Update: materi oke, tugas lebih banyak dari ekspektasi"
} | ConvertTo-Json -Compress)
Show "T200 update sendiri" 200 $r.Code

# ============================================================
# Uji 13: Reply → 201
# ============================================================
Write-Host "`n=== Uji 13: Reply ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/reviews/$reviewID/replies" $h2 (@{
    body = "Setuju, tugasnya memang banyak"
} | ConvertTo-Json -Compress)
Show "T201 reply review T200" 201 $r.Code

# ============================================================
# Uji 14: Invalid ObjectID → 400
# ============================================================
Write-Host "`n=== Uji 14: Invalid ID ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/reviews/bukan-objectid" $h1
Show "GET /reviews/bukan-objectid" 400 $r.Code

# ============================================================
# Uji 15: NoSQL injection dalam tag → 200
# ============================================================
Write-Host "`n=== Uji 15: NoSQL injection ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/courses/1/reviews?tag=%7B%22%24ne%22%3A%22%22%7D" $h1
Show "Filter tag={`$ne:''}" 200 $r.Code

# ============================================================
# Uji 16: NoSQL injection di body → 400
# ============================================================
Write-Host "`n=== Uji 16: NoSQL injection di body ===" -ForegroundColor Cyan
$r = Api POST "$baseUrl/courses/1/reviews" $h1 '{"rating":5,"comment":{"$ne":""}}'
Show "Body dengan operator di comment" 400 $r.Code

# ============================================================
# Ringkasan
# ============================================================
Write-Host "`n============= SELESAI =============" -ForegroundColor Cyan
if ($reviewID) {
    Write-Host "Review ID untuk EXPLAIN: $reviewID" -ForegroundColor Green
}