@echo off
cd /d "%~dp0"
if exist "tmp\main.exe" (
    start "" "tmp\main.exe"
) else (
    echo Binary not found
)
