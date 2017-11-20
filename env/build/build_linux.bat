ECHO Linux Compile....

set GOOS=linux
set GOARCH=amd64
if not exist %gopath%\bin\pasque\linux\amd64 (
    mkdir  %gopath%\bin\pasque\linux\amd64
) 

if not exist %gopath%\bin\pasque\linux\config\ (
    mkdir  %gopath%\bin\pasque\linux\config
)

go build -o %gopath%\bin\pasque\linux\router %gopath%\src\pasque\server\router\main.go %gopath%\src\pasque\server\router\router.go
go build -o %gopath%\bin\pasque\linux\svcgate %gopath%\src\pasque\server\svcgate\main.go %gopath%\src\pasque\server\svcgate\gate.go
go build -o %gopath%\bin\pasque\linux\apigate %gopath%\src\pasque\server\apigate\main.go %gopath%\src\pasque\server\apigate\gate.go

copy %gopath%\src\pasque\env\config\system_linux.json %gopath%\bin\pasque\linux\config
copy %gopath%\src\pasque\env\run\run_linux %gopath%\bin\pasque\linux\run
cd %gopath%\bin\pasque\linux\config
del system.json 
ren system_linux.json system.json
cd %gopath%\src\pasque
