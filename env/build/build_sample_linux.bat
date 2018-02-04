ECHO Linux Compile....

set GOOS=linux
set GOARCH=amd64
if not exist %gopath%\bin\pasque\linux\amd64 (
    mkdir  %gopath%\bin\pasque\linux\amd64
) 

if not exist %gopath%\bin\pasque\linux\config\ (
    mkdir  %gopath%\bin\pasque\linux\config
)

go build -o %gopath%\bin\pasque\linux\router %gopath%\src\pasque\bus\router\main.go %gopath%\src\pasque\bus\router\router.go
go build -o %gopath%\bin\pasque\linux\sgate %gopath%\src\pasque\bus\sgate\main.go %gopath%\src\pasque\bus\sgate\gate.go
go build -o %gopath%\bin\pasque\linux\egate %gopath%\src\pasque\bus\egate\main.go %gopath%\src\pasque\bus\egate\gate.go

go build -o %gopath%\bin\pasque\linux\hellosrv %gopath%\src\pasque\services\sample\hellosrv\main.go %gopath%\src\pasque\services\sample\hellosrv\biz.go
go build -o %gopath%\bin\pasque\linux\hellci %gopath%\src\pasque\services\sample\hellocli\main.go %gopath%\src\pasque\services\sample\hellocli\biz.go


copy %gopath%\src\pasque\env\config\system_sample_linux.json %gopath%\bin\pasque\linux\config
copy %gopath%\src\pasque\env\run\run_sample_linux %gopath%\bin\pasque\linux\run
cd %gopath%\bin\pasque\linux\config
del system.json 
ren system_linux.json system.json
cd %gopath%\src\pasque
