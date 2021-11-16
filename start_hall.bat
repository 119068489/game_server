@echo off
title hall
go run -race .\execute\hall\main.go 2>>logs\hall_std_err.log
pause