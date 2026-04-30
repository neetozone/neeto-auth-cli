@echo off
setlocal

set "BASE_URL=https://neeto-downloads.s3.amazonaws.com/cli/NeetoAuth/latest"
set "INSTALL_DIR=%LOCALAPPDATA%\Programs\neetoauth"

set "ARCH=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"
set "ARCHIVE=neetoauth_windows_%ARCH%.zip"

echo Downloading NeetoAuth CLI...
set "TMPDIR=%TEMP%\neetoauth-install"
if exist "%TMPDIR%" rmdir /s /q "%TMPDIR%"
mkdir "%TMPDIR%"

curl -fsSL "%BASE_URL%/%ARCHIVE%" -o "%TMPDIR%\%ARCHIVE%"
if %errorlevel% neq 0 (
    echo Failed to download NeetoAuth CLI.
    exit /b 1
)

echo Extracting...
powershell -Command "Expand-Archive -Path '%TMPDIR%\%ARCHIVE%' -DestinationPath '%TMPDIR%' -Force"

echo Installing to %INSTALL_DIR%...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /y "%TMPDIR%\neetoauth.exe" "%INSTALL_DIR%\neetoauth.exe" >nul

:: Add to user PATH if not already present
echo %PATH% | findstr /i /c:"%INSTALL_DIR%" >nul
if %errorlevel% neq 0 (
    for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%B"
    setx PATH "%USER_PATH%;%INSTALL_DIR%" >nul
    echo Added %INSTALL_DIR% to user PATH.
)

rmdir /s /q "%TMPDIR%"

echo NeetoAuth CLI installed successfully. Restart your terminal and run 'neetoauth --help' to get started.
