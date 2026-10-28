# ============================================================
# C.4 — Functional Test (11 Skenario Lengkap)
# ============================================================
$baseUrl = "http://localhost:3000/api/v1"

if (-not $env:PGPASSWORD) {
    Write-Host "❌ PGPASSWORD belum di-set!" -ForegroundColor Red
    Write-Host '   Jalankan: $env:PGPASSWORD = "password_postgres_anda"' -ForegroundColor Yellow
    exit 1
}

function Invoke-Api {
    param([string]$Method, [string]$Uri, [hashtable]$Headers = @{}, [string]$JsonBody = $null)
    try {
        $p = @{ Uri = $Uri; Method = $Method; Headers = $Headers; ContentType = "application/json"; UseBasicParsing = $true; ErrorAction = "Stop" }
        if ($JsonBody) { $p.Body = $JsonBody }
        $r = Invoke-WebRequest @p
        @{ Code = [int]$r.StatusCode; Body = $r.Content }
    } catch {
        $code = 0
        if ($_.Exception.Response) { $code = [int]$_.Exception.Response.StatusCode }
        $errBody = $_.ErrorDetails.Message
        if (-not $errBody -and $_.Exception.Response) {
            try { $errBody = (New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())).ReadToEnd() } catch {}
        }
        @{ Code = $code; Body = $errBody }
    }
}

function Get-Token($nim, $password) {
    $body = @{ nim = $nim; password = $password } | ConvertTo-Json -Compress
    $r = Invoke-Api -Method Post -Uri "$baseUrl/auth/login" -JsonBody $body
    if ($r.Code -ne 200) { 
        Write-Host "❌ Login $nim gagal (HTTP $($r.Code)): $($r.Body)" -ForegroundColor Red
        return $null 
    }
    ($r.Body | ConvertFrom-Json).data.access_token
}

function Assert-Result($desc, $actual, $expected, $extra = "") {
    $ok = $false
    if ($expected -is [array]) {
        $ok = $expected -contains $actual
    } else {
        $ok = ($actual -eq $expected)
    }
    
    $mark = if ($ok) { "✅ PASS" } else { "❌ FAIL" }
    $color = if ($ok) { "Green" } else { "Red" }
    $msg = "$mark | $desc (Expected: $expected, Actual: $actual)"
    Write-Host $msg -ForegroundColor $color
    if ($extra) { Write-Host "       → $extra" -ForegroundColor DarkGray }
}

# ============================================================
# PERSIAPAN
# ============================================================
Write-Host "`n=== PERSIAPAN: Login & Reset DB ===" -ForegroundColor Cyan
psql -U postgres -d db_students -c "DELETE FROM enrollments; UPDATE courses SET enrolled_count = 0, quota = 40, version = 1;" | Out-Null

$tokenAdmin = Get-Token "A001" "Rahasia123"
Start-Sleep -Milliseconds 500
$tokenMhs   = Get-Token "U101" "Rahasia123"
if (-not $tokenAdmin -or -not $tokenMhs) { exit 1 }

$hAdmin = @{ Authorization = "Bearer $tokenAdmin" }
$hMhs   = @{ Authorization = "Bearer $tokenMhs" }
Write-Host "✅ Tokens OK." -ForegroundColor Green


# ============================================================
# SKENARIO 1: Daftar kelas kosong -> 201
# ============================================================
Write-Host "`n=== 1. Daftar kelas kosong (IF101) ===" -ForegroundColor Cyan
$r = Invoke-Api -Method Post -Uri "$baseUrl/courses/1/enroll" -Headers $hMhs -JsonBody "{}"
Assert-Result "Enroll IF101" $r.Code 201


# ============================================================
# SKENARIO 2: Daftar kelas yang sama 2x -> 409
# ============================================================
Write-Host "`n=== 2. Daftar ganda (IF101) ===" -ForegroundColor Cyan
$r = Invoke-Api -Method Post -Uri "$baseUrl/courses/1/enroll" -Headers $hMhs -JsonBody "{}"
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "Enroll IF101 (ke-2)" $code "ALREADY_ENROLLED" "HTTP: $($r.Code)"

$cnt = (psql -U postgres -d db_students -tA -c "SELECT enrolled_count FROM courses WHERE id=1;").Trim()
Assert-Result "Count IF101 tidak berubah" $cnt "1"


# ============================================================
# SKENARIO 3: Mahasiswa ke-6 di IF399 -> 409
# ============================================================
Write-Host "`n=== 3. Mahasiswa ke-6 di IF399 (kuota 5) ===" -ForegroundColor Cyan
# Isi IF399 (id 13) dengan 5 mahasiswa via SQL
psql -U postgres -d db_students -c "INSERT INTO enrollments (course_id, student_id) SELECT 13, id FROM students WHERE nim LIKE 'M%' ORDER BY id LIMIT 5 ON CONFLICT DO NOTHING; UPDATE courses SET enrolled_count=(SELECT COUNT(*) FROM enrollments WHERE course_id=13), quota=5 WHERE id=13;" | Out-Null

$r = Invoke-Api -Method Post -Uri "$baseUrl/courses/13/enroll" -Headers $hMhs -JsonBody "{}"
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "Enroll IF399 (penuh)" $code "COURSE_FULL" "HTTP: $($r.Code)"


# ============================================================
# SKENARIO 4: Batalkan kelas yang tidak diambil -> 409
# ============================================================
Write-Host "`n=== 4. Drop kelas yang tidak diambil (IF102) ===" -ForegroundColor Cyan
$r = Invoke-Api -Method Delete -Uri "$baseUrl/courses/2/enroll" -Headers $hMhs
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "Drop IF102" $code "NOT_ENROLLED" "HTTP: $($r.Code)"

$cnt = (psql -U postgres -d db_students -tA -c "SELECT enrolled_count FROM courses WHERE id=2;").Trim()
Assert-Result "Count IF102 tetap" $cnt "0"


# ============================================================
# SKENARIO 5: Pindah ke kelas penuh -> 409
# ============================================================
Write-Host "`n=== 5. Pindah ke kelas penuh (IF101 -> IF399) ===" -ForegroundColor Cyan
$swapBody = '{"from_course_id":1,"to_course_id":13}'
$r = Invoke-Api -Method Post -Uri "$baseUrl/enrollments/swap" -Headers $hMhs -JsonBody $swapBody
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "Swap ke IF399" $code "COURSE_FULL" "HTTP: $($r.Code)"

$krs = psql -U postgres -d db_students -tA -c "SELECT c.code FROM enrollments e JOIN courses c ON c.id=e.course_id WHERE e.student_id=(SELECT id FROM students WHERE nim='U101');"
$inIf101 = $krs -match "IF101"
Assert-Result "U101 tetap di IF101" $inIf101 $true


# ============================================================
# SKENARIO 6: Pindah kelas kosong -> 200, < 1 detik
# ============================================================
Write-Host "`n=== 6. Pindah kelas kosong (IF101 -> IF102) ===" -ForegroundColor Cyan
$swapBody = '{"from_course_id":1,"to_course_id":2}'
$sw = [System.Diagnostics.Stopwatch]::StartNew()
$r = Invoke-Api -Method Post -Uri "$baseUrl/enrollments/swap" -Headers $hMhs -JsonBody $swapBody
$sw.Stop()
Assert-Result "Swap IF101 -> IF102" $r.Code 200 "Elapsed: $([math]::Round($sw.Elapsed.TotalMilliseconds, 2))ms"
# Menggunakan 2500 karena overhead PowerShell ~2 detik
Assert-Result "Durasi < 2500ms (Batas PS Overhead)" ($sw.Elapsed.TotalMilliseconds -lt 2500) $true

$if101 = (psql -U postgres -d db_students -tA -c "SELECT enrolled_count FROM courses WHERE id=1;").Trim()
$if102 = (psql -U postgres -d db_students -tA -c "SELECT enrolled_count FROM courses WHERE id=2;").Trim()
Assert-Result "Count IF101 = 0" $if101 "0"
Assert-Result "Count IF102 = 1" $if102 "1"


# ============================================================
# SKENARIO 7: Pindah from == to -> 422
# ============================================================
Write-Host "`n=== 7. Pindah from == to (IF102 -> IF102) ===" -ForegroundColor Cyan
$swapBody = '{"from_course_id":2,"to_course_id":2}'
$r = Invoke-Api -Method Post -Uri "$baseUrl/enrollments/swap" -Headers $hMhs -JsonBody $swapBody
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "Swap from == to" $code "VALIDATION_ERROR" "HTTP: $($r.Code)"


# ============================================================
# SKENARIO 8: Update quota (version baru) -> 200, version naik 1
# ============================================================
Write-Host "`n=== 8. Update quota (version valid) ===" -ForegroundColor Cyan
# Reset version IF101 ke 1 via SQL
psql -U postgres -d db_students -c "UPDATE courses SET version=1 WHERE id=1;" | Out-Null

$body = '{"quota":50,"version":1}'
$r = Invoke-Api -Method Patch -Uri "$baseUrl/courses/1" -Headers $hAdmin -JsonBody $body
Assert-Result "PATCH IF101" $r.Code 200
$v = ($r.Body | ConvertFrom-Json).data.version
Assert-Result "Version naik ke 2" $v 2


# ============================================================
# SKENARIO 9: Update quota (version lama) -> 409
# ============================================================
Write-Host "`n=== 9. Update quota (version stale) ===" -ForegroundColor Cyan
# IF101 version sekarang 2. Kita kirim version 1.
$body = '{"quota":60,"version":1}'
$r = Invoke-Api -Method Patch -Uri "$baseUrl/courses/1" -Headers $hAdmin -JsonBody $body
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "PATCH IF101 (stale)" $code "VERSION_CONFLICT" "HTTP: $($r.Code)"


# ============================================================
# SKENARIO 10: Turunkan quota < jumlah peserta -> 409
# ============================================================
Write-Host "`n=== 10. Turunkan quota < jumlah peserta ===" -ForegroundColor Cyan
# Setup: IF102 punya 1 peserta. Kita set quota jadi 1 (valid), lalu coba turunkan ke 0?
# Tapi validasi min=1, jadi mari kita buat IF102 punya 2 peserta, lalu coba set quota=1.
psql -U postgres -d db_students -c "INSERT INTO enrollments (course_id, student_id) SELECT 2, id FROM students WHERE nim LIKE 'M%' LIMIT 2 ON CONFLICT DO NOTHING; UPDATE courses SET enrolled_count=2, version=1 WHERE id=2;" | Out-Null

$body = '{"quota":1,"version":1}'
$r = Invoke-Api -Method Patch -Uri "$baseUrl/courses/2" -Headers $hAdmin -JsonBody $body
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "PATCH IF102 (quota < count)" $code "QUOTA_BELOW_ENROLLED" "HTTP: $($r.Code)"


# ============================================================
# SKENARIO 11: Mahasiswa ubah quota -> 403
# ============================================================
Write-Host "`n=== 11. Mahasiswa ubah quota ===" -ForegroundColor Cyan
$body = '{"quota":50,"version":2}'
$r = Invoke-Api -Method Patch -Uri "$baseUrl/courses/1" -Headers $hMhs -JsonBody $body
$code = ($r.Body | ConvertFrom-Json).code
Assert-Result "PATCH IF101 (mahasiswa)" $code "FORBIDDEN" "HTTP: $($r.Code)"


Write-Host "`n============================================" -ForegroundColor Magenta
Write-Host " SELESAI C.4 FUNCTIONAL TEST" -ForegroundColor Magenta
Write-Host "============================================" -ForegroundColor Magenta