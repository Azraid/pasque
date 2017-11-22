ECHO Test Compile....

set GOOS=windows
set GOARCH=amd64

if not exist %gopath%\bin\pasque\windows\ (
    mkdir  %gopath%\bin\pasque\windows
) 

if not exist %gopath%\bin\pasque\env\config\ (
    mkdir  %gopath%\bin\pasque\env\config
)


go build -o %gopath%\bin\pasque\windows\router.exe %gopath%\src\pasque\server\router\main.go %gopath%\src\pasque\server\router\router.go
go build -o %gopath%\bin\pasque\windows\svcgate.exe %gopath%\src\pasque\server\svcgate\main.go %gopath%\src\pasque\server\svcgate\gate.go
go build -o %gopath%\bin\pasque\windows\apigate.exe %gopath%\src\pasque\server\apigate\main.go %gopath%\src\pasque\server\apigate\gate.go

copy %gopath%\src\pasque\env\config\system_sample.json %gopath%\bin\pasque\windows\config\system.json
copy %gopath%\src\pasque\env\run\runw_sample.bat %gopath%\bin\pasque\windows

