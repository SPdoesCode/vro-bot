@echo off

echo Building...

go build -o vro-bot.exe  bot.go

if %ERRORLEVEL% neq 0 (
    echo Build failed!
    pause
    exit /b %ERRORLEVEL%
)

echo Done! To use run vro-bot.exe

pause
