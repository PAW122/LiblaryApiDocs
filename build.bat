@echo off
REM Budowanie projektu Go na Windows
echo Building Go project for Windows...
set GOOS=windows
set GOARCH=amd64

go build -o docs.exe

REM Sprawdź, czy build zakończył się sukcesem
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

REM Budowanie projektu Go na Linuxa
echo Building Go project for Linux...
set GOOS=linux
set GOARCH=amd64

go build -o docs-linux

REM Sprawdź, czy build zakończył się sukcesem
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Linux build completed successfully!

pause
