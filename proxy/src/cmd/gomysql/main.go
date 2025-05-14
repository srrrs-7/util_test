package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Constants for environment variable names
const (
	EnvProxyAddr = "PROXY_ADDR"
	EnvProxyUser = "PROXY_USER"
	EnvProxyPass = "PROXY_PASS"

	EnvTestDBAddr = "TEST_DB_ADDR"
	EnvTestDBUser = "TEST_DB_USER"
	EnvTestDBPass = "TEST_DB_PASS"
	EnvTestDBName = "TEST_DB_NAME"
)

// Env is a struct that holds database connection information loaded from environment variables
type Env struct {
	proxyAddr  string
	proxyUser  string
	proxyPass  string
	testDBAddr string
	testDBUser string
	testDBPass string
	testDBName string
}

// newEnv loads settings from environment variables and generates an Env struct
func newEnv() Env {
	return Env{
		proxyAddr:  os.Getenv(EnvProxyAddr),
		proxyUser:  os.Getenv(EnvProxyUser),
		proxyPass:  os.Getenv(EnvProxyPass),
		testDBAddr: os.Getenv(EnvTestDBAddr),
		testDBUser: os.Getenv(EnvTestDBUser),
		testDBPass: os.Getenv(EnvTestDBPass),
		testDBName: os.Getenv(EnvTestDBName),
	}
}

func (e Env) validate() error {
	if e.proxyAddr == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvProxyAddr)
	}
	if e.proxyUser == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvProxyUser)
	}
	if e.proxyPass == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvProxyPass)
	}
	if e.testDBAddr == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvTestDBAddr)
	}
	if e.testDBUser == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvTestDBUser)
	}
	if e.testDBPass == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvTestDBPass)
	}
	if e.testDBName == "" {
		return fmt.Errorf("missing required environment variable: %s", EnvTestDBName)
	}

	return nil
}

// main is the entry point of the application
func main() {
	// Load settings from environment variables
	env := newEnv()
	if err := env.validate(); err != nil {
		log.Panic("Failed to load environment variables: ", err)
	}

	proxyConf := DBConfig{
		Addr: env.proxyAddr,
		User: env.proxyUser,
		Pass: env.proxyPass,
	}

	testDBConf := DBConfig{
		Addr:   env.testDBAddr,
		User:   env.testDBUser,
		Pass:   env.testDBPass,
		DBName: env.testDBName,
	}

	l, err := net.Listen("tcp", env.proxyAddr)
	if err != nil {
		log.Panic("Failed to listen on proxy address: ", err)
	}
	defer l.Close()
	log.Printf("MySQL proxy server successfully listening on %s", env.proxyAddr)

	ctx, cancel := context.WithCancel(context.Background())
	// Create and start the proxy server
	proxyServer := NewProxyServer(
		proxyConf,
		testDBConf,
		l,
		&sync.WaitGroup{},
	)
	go func(cancel context.CancelFunc) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		defer cancel()
		log.Println("Received shutdown signal, shutting down...")
		proxyServer.Shutdown()
	}(cancel)

	if err := proxyServer.Start(ctx); err != nil {
		log.Panic("Failed to start proxy server: ", err)
	}

	<-ctx.Done()
}
