package main

import (
	"awesomeProject1/internal/api"
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"awesomeProject1/ioc"
)

func main() {
	l := ioc.InitZapLogger()
	logger := logger2.NewZapLogger(l)
	repo := repository.NewStorageRepository(logger)
	objSvc := service.NewObjectService(repo)
	api.NewObjectHandler(objSvc, logger)
}
