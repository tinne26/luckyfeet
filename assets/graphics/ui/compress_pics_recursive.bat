:: Exit if optipng not available
@echo off
echo Starting compression...
optipng -version >nul
if %errorlevel% neq 0 timeout 10 & exit

:: Gets each file in directory and compresses it with OptiPNG
for /R %%a in (*.png) do (
	echo ...compressing %%a
	optipng "%%a" -strip all -o8 -quiet
)

echo[
echo Done!
timeout 10
