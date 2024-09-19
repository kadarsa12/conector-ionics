package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

type program struct {
	DisplayName       string
	db                *database
	baseUrl           string
	authEndpoint      string
	writerEndpoint    string
	username          string
	password          string
	executionInterval int
	qtySendData       int
	exit              chan struct{}
}

func (p *program) Start(s service.Service) error {
	go run_intervaled_job(p.executionInterval, RUN(p.db, p.baseUrl, p.authEndpoint, p.writerEndpoint, p.username, p.password, p.qtySendData))
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping...")
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func run_program(prg *program) {
	svcConfig := &service.Config{
		Name:        "InovaKPIService",
		DisplayName: "InovaKPIService",
		Description: "InovaKPIService",
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
