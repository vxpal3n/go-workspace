# ============================================================
# Test End-to-End Portfolio Endpoints (Modul 10 - D.2)
# ============================================================
$baseUrl = "http://localhost:3000/api/v1"

function Api {
    param([string]$Method, [string]$Uri, [hashtable]$Headers = @{}, [string]$Body = $null)
    $p = @{
        Uri = $Uri; Method = $Method; Headers = $Headers
        ContentType = "application/json"; UseBasicParsing = $true; ErrorAction = "Stop"
    }
    if ($Body) { $p.Body = $Body }
    try {
        $r = Invoke-WebRequest @p
        @{ Code = [int]$r.StatusCode; Body = $r.Content }
    } catch {
        $code = 0
        if ($_.Exception.Response) { $code = [int]$_.Exception.Response.StatusCode }
        $errBody = $_.ErrorDetails.Message
        if (-not $errBody -and $_.Exception.Response) {
            try { $errBody = (New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())).ReadToEnd() } catch { }
        }
        @{ Code = $code; Body = $errBody }
    }
}

function Show($label, $expect, $actual, $extra = "") {
    $ok = $expect -eq $actual
    $color = if ($ok) { "Green" } else { "Red" }
    Write-Host ("{0} {1,-50} harap={2,-4} dapat={3}" -f
        ($(if ($ok) {"[OK]  "} else {"[FAIL]"})), $label, $expect, $actual) -ForegroundColor $color
    if ($extra) { Write-Host "       → $extra" -ForegroundColor DarkGray }
}

# ============================================================
# 1. Persiapan: Login & Ambil ID
# ============================================================
Write-Host "`n=== Persiapan ===" -ForegroundColor Cyan
$loginBody = @{ nim = "T200"; password = "Rahasia123" } | ConvertTo-Json -Compress
$loginResp = Api POST "$baseUrl/auth/login" @{} $loginBody

if ($loginResp.Code -ne 200) {
    Write-Host "❌ Login gagal: $($loginResp.Body)" -ForegroundColor Red
    return
}

$token = ($loginResp.Body | ConvertFrom-Json).data.access_token
$h = @{ Authorization = "Bearer $token" }

# Ambil ID T200 dari endpoint /auth/me
$meResp = Api GET "$baseUrl/auth/me" $h
$studentId = ($meResp.Body | ConvertFrom-Json).data.user.id
Write-Host "✅ Login OK. Student ID T200 = $studentId" -ForegroundColor Green

# ============================================================
# Uji 1: GET portfolio (belum ada) → 404
# ============================================================
Write-Host "`n=== Uji 1: GET Portfolio (Kosong) ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/students/$studentId/portfolio" $h
Show "GET portfolio kosong" 404 $r.Code

# ============================================================
# Uji 2: PUT portfolio (buat baru) → 200
# ============================================================
Write-Host "`n=== Uji 2: PUT Portfolio ===" -ForegroundColor Cyan
$putBody = @{
    headline = "Backend Enthusiast"
    bio = "Suka Go dan PostgreSQL"
    projects = @(
        @{ title = "Siakad Mini"; technologies = @("Go", "PostgreSQL"); year = 2026 }
    )
    certificates = @()
} | ConvertTo-Json -Depth 5 -Compress

$r = Api PUT "$baseUrl/students/$studentId/portfolio" $h $putBody
Show "PUT portfolio" 200 $r.Code
$portfolio = $r.Body | ConvertFrom-Json
if ($portfolio.success) {
    Write-Host "       → ID Portfolio: $($portfolio.data.id)" -ForegroundColor DarkGray
}

# ============================================================
# Uji 3: GET portfolio (sudah ada) → 200
# ============================================================
Write-Host "`n=== Uji 3: GET Portfolio (Ada) ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/students/$studentId/portfolio" $h
Show "GET portfolio terisi" 200 $r.Code

# ============================================================
# Uji 4: POST project baru → 201
# ============================================================
Write-Host "`n=== Uji 4: POST Project ===" -ForegroundColor Cyan
$postBody = @{
    title = "SIAKAD Mini v2"
    technologies = @("Go", "MongoDB")
    year = 2026
} | ConvertTo-Json -Compress

$r = Api POST "$baseUrl/students/$studentId/portfolio/projects" $h $postBody
Show "POST project baru" 201 $r.Code
$proj = $r.Body | ConvertFrom-Json
if ($proj.success) {
    Write-Host "       → Total proyek sekarang: $($proj.data.projects.Count)" -ForegroundColor DarkGray
}

# ============================================================
# Uji 5: GET stats technologies → 200
# ============================================================
Write-Host "`n=== Uji 5: GET Stats ===" -ForegroundColor Cyan
$r = Api GET "$baseUrl/portfolios/stats/technologies" $h
Show "GET stats tech" 200 $r.Code
$stats = $r.Body | ConvertFrom-Json
if ($stats.success) {
    $top = $stats.data | Select-Object -First 3 | ForEach-Object { "$($_.technology)=$($_.total)" }
    Write-Host "       → Top tech: $($top -join ', ')" -ForegroundColor DarkGray
}

Write-Host "`n============= SELESAI =============" -ForegroundColor Cyan