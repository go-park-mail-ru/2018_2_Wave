@echo off
go build -o ./build/auth-serv.exe ./cmd/auth/
go build -o ./build/game-serv.exe ./cmd/game/
go build -o ./build/api-serv.exe ./cmd/api/