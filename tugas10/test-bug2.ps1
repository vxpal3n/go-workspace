# ============================================================
# Bug #2 — Missing defer tx.Rollback(ctx)
# Evidence: pg_stat_activity + health check, before vs after
# ============================================================
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("before", "after")]
    [string]$Phase
)

$baseUrl  = "http://localhost:3000/api/v1"
$healthUrl = "http://localhost:3000/health"
$db       = "db_students"
$pgUser   = "postgres"

if (-not $env:PGPASSWORD) {
    Write-Host "❌ PGPASSWORD belum di-set." -ForegroundColor Red
    Write-Host '   Jalankan: $env:PGPASSWORD = "password_postgres_anda"' -ForegroundColor Yellow
    exit 1
}

function Invoke-Psql {
    param([string]$Sql)
    & psql -U $pgUser -d $db -c $Sql 2>&1
}

function Show-ActiveConns {
    Write-Host "`n--- pg_stat_activity (state <> 'idle', exclude self) ---" -ForegroundColor Cyan
    $sql = @"
SELECT pid,
       state,
       COALESCE(now() - xact_start, interval '0') AS age,
       left(query, 60) AS query
FROM pg_stat_activity
WHERE datname = '$db'
  AND state <> 'idle'
  AND pid <> pg_backend_pid()
ORDER BY pid;
"@
    $raw = Invoke-Psql $sql
    $raw | Write-Host
    $idleTx = @($raw | Select-String -SimpleMatch 'idle in transaction').Count
    Write-Host "  → idle-in-transaction rows: $idleTx" `
        -ForegroundColor $(if ($idleTx -ge 3) { "Red" } elseif ($idleTx -ge 1) { "Yellow" } else { "Green" })
    return $idleTx
}

function Show-PoolSummary {
    Write-Host "`n--- Ringkasan state koneksi ---" -ForegroundColor Cyan
    Invoke-Psql @"
SELECT state, count(*) AS n
FROM pg_stat_activity
WHERE datname = '$db' AND pid <> pg_backend_pid()
GROUP BY state ORDER BY n DESC;
"@ | Write-Host
}

function Test-Health {
    param([string]$Label)
    Write-Host "`n--- Health check ($Label) ---" -ForegroundColor Cyan
    $sw = [System.Diagnostics.Stopwatch]::StartNew()
    try {
        $r = Invoke-WebRequest -Uri $healthUrl -UseBasicParsing -TimeoutSec 15 -ErrorAction Stop
        $sw.Stop()
        $ok = ($r.StatusCode -eq 200) -and ($sw.Elapsed.TotalSeconds -lt 1)
        Write-Host ("HTTP {0}  {1:N3}s  body={2}" -f $r.StatusCode, $sw.Elapsed.TotalSeconds, $r.Content.Trim()) `
            -ForegroundColor $(if ($ok) { "Green" } else { "Yellow" })
    } catch {
        $sw.Stop()
        $code = if ($_.Exception.Response) { [int]$_.Exception.Response.StatusCode } else { 0 }
        Write-Host ("HTTP {0}  {1:N3}s  error: {2}" -f $code, $sw.Elapsed.TotalSeconds, $_.Exception.Message) `
            -ForegroundColor Red
    }
}

# ============================================================
Write-Host "`n============================================" -ForegroundColor Magenta
Write-Host " Bug #2 Reproduction — Phase: $($Phase.ToUpper())" -ForegroundColor Magenta
Write-Host "============================================`n" -ForegroundColor Magenta

# ---------- STEP 1: Baseline ----------
Write-Host "STEP 1: Baseline pg_stat_activity (sebelum uji)" -ForegroundColor Cyan
$baseline = Show-ActiveConns
Show-PoolSummary
Write-Host "`nBaseline idle-in-tx: $baseline" -ForegroundColor Yellow

# ---------- STEP 2: Health pra ----------
Write-Host "`nSTEP 2: Health pra-uji" -ForegroundColor Cyan
Test-Health -Label "pra-uji"

# ---------- STEP 3: Reset ----------
Write-Host "`nSTEP 3: Reset DB" -ForegroundColor Cyan
Invoke-Psql "DELETE FROM enrollments; UPDATE courses SET enrolled_count = 0;" | Write-Host

# ---------- STEP 4: 3x credits mode ----------
Write-Host "`nSTEP 4: Jalankan 3× 'credits' mode (TANPA reset di antara)" -ForegroundColor Cyan
Push-Location $PSScriptRoot
try {
    for ($i = 1; $i -le 3; $i++) {
        Write-Host "`n--- Iterasi $i / 3 ---" -ForegroundColor Yellow
        $out = & go run ./tools/loadtest -mode credits 2>&1
        $out | Write-Host
        # Cek cepat apakah ada CREDIT_LIMIT di output
        $creditLimitHits = @($out | Select-String -SimpleMatch 'CREDIT_LIMIT').Count
        Write-Host "  (CREDIT_LIMIT hits di output: $creditLimitHits)" -ForegroundColor DarkGray
    }
} finally {
    Pop-Location
}

# ---------- STEP 5: pg_stat_activity pasca ----------
Write-Host "`nSTEP 5: pg_stat_activity SETELAH 3× credits" -ForegroundColor Cyan
$post = Show-ActiveConns
Show-PoolSummary

# ---------- STEP 6: Health pasca ----------
Write-Host "`nSTEP 6: Health pasca-uji" -ForegroundColor Cyan
Test-Health -Label "pasca-uji"

# ---------- STEP 7: Verdict ----------
Write-Host "`n============================================" -ForegroundColor Magenta
Write-Host " VERDICT — phase $Phase" -ForegroundColor Magenta
Write-Host "============================================" -ForegroundColor Magenta
Write-Host "Baseline idle-in-tx : $baseline"
Write-Host "Setelah uji         : $post"

if ($Phase -eq "before") {
    if ($post -ge 3) {
        Write-Host "`n✅ BUG #2 CONFIRMED" -ForegroundColor Red
        Write-Host "   $post koneksi tertinggal di state 'idle in transaction'." -ForegroundColor Red
        Write-Host "   → pool bisa habis; health check melambat/gagal setelah beberapa iterasi." -ForegroundColor Red
    } else {
        Write-Host "`n⚠️  Bug #2 TIDAK terlihat di fase 'before'." -ForegroundColor Yellow
        Write-Host "   Cek: (a) apakah 'defer tx.Rollback(ctx)' benar-benar di-comment?" -ForegroundColor Yellow
        Write-Host "        (b) apakah server sudah di-restart setelah comment?" -ForegroundColor Yellow
        Write-Host "        (c) apakah 'credits' mode menghasilkan error (CREDIT_LIMIT)?" -ForegroundColor Yellow
    }
} else {
    if ($post -le 1) {
        Write-Host "`n✅ BUG #2 FIXED" -ForegroundColor Green
        Write-Host "   Tidak ada koneksi tertinggal; rollback berjalan." -ForegroundColor Green
    } else {
        Write-Host "`n⚠️  Masih ada $post koneksi idle-in-tx." -ForegroundColor Yellow
        Write-Host "   Cek: apakah 'defer tx.Rollback(ctx)' sudah di-uncomment & server di-restart?" -ForegroundColor Yellow
    }
}