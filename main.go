package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

func main() {
	baseUrl := flag.String("base_url", "http://127.0.0.1:3333", "API URL")
	authEndpoint := flag.String("auth_endpoint", "/auth/login", "Auth Endpoint")
	writerEndpoint := flag.String("writer_endpoint", "/v1/writer", "Send Data Endpoint")
	username := flag.String("username", "marcel", "API Username Auth")
	password := flag.String("password", "marcel", "API Password Auth")
	dbName := flag.String("db_name", "pgsql", "Database Name (oracle, pgsql, sqlsrv)")
	dbHost := flag.String("db_host", "127.0.0.1", "Database Host")
	dbSID := flag.String("db_sid", "ionics", "Database SID")
	dbPort := flag.Int("db_port", 5432, "Database Port")
	dbServiceName := flag.String("db_service_name", "ionics", "Database Service Name")
	dbUsername := flag.String("db_username", "dev", "Database Username")
	dbPassword := flag.String("db_password", "dev", "Database Password")
	executionInterval := flag.Int("execution_interval", 0, "Execution Interval in Hours")
	qtySendData := flag.Int("qty_sent_data", 100, "Quantity of data to send")
	logPath := flag.String("log_path", `./output.log`, "Log Path") // C:\inovakpi\output.log
	flag.Parse()

	var db *database
	var err error

	file, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger = slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), nil))

	logger.Info("Starting...")

	if *dbName == "oracle" {
		db, err = oracle_connection(dbHost, dbPort, dbServiceName, dbUsername, dbPassword, dbSID)
	} else if *dbName == "pgsql" {
		db, err = pgsql_connection(dbHost, dbPort, dbServiceName, dbUsername, dbPassword, dbSID)
	} else if *dbName == "sqlsrv" {
		db, err = sqlsrv_connection(dbHost, dbPort, dbServiceName, dbUsername, dbPassword, dbSID)
	}

	if err != nil {
		panic(err)
	}

	defer db.Close()

	logger.Info("Connected to DB...")

	prg := &program{
		DisplayName:       "InovaKPIService",
		db:                db,
		baseUrl:           *baseUrl,
		authEndpoint:      *authEndpoint,
		writerEndpoint:    *writerEndpoint,
		username:          *username,
		password:          *password,
		executionInterval: *executionInterval,
		qtySendData:       *qtySendData,
		exit:              make(chan struct{}),
	}

	run_program(prg)

	<-prg.exit
}

func run_intervaled_job(interval int, f func()) {
	logger.Info("Waiting for jobs...")

	// Run the job immediately
	if interval == 0 {
		f()
		return
	}

	// Run the job every interval hours
	ticker := time.NewTicker(time.Duration(interval) * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		f()
	}
}
