@echo off
title SILAYANG - Web Layanan Kelurahan Online
echo ========================================================
echo         MENJALANKAN SILAYANG ^& CLOUDFLARE TUNNEL
echo ========================================================
echo.

:: 1. Cek apakah server SILAYANG sudah jalan
tasklist /FI "IMAGENAME eq main.exe" 2>NUL | find /I /N "main.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [OK] Web server SILAYANG sudah aktif di port 8080.
) else (
    echo [INFO] Menjalankan web server SILAYANG...
    start /B "" "tmp\main.exe"
    timeout /t 2 /nobreak >nul
)

:: 2. Jalankan Cloudflare Tunnel
echo.
echo [INFO] Menghubungkan ke Cloudflare Tunnel (Akses Internet Publik)...
echo [INFO] URL publik akan muncul di bawah ini (akhiran .trycloudflare.com):
echo ========================================================
echo.
tmp\cloudflared.exe tunnel --url http://127.0.0.1:8080
pause
