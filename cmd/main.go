package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mdm/logger"
)

const (
	Sep = "----------"
	Ext = "toml"
)

var WorkDir string

func init() {
	// Set Profile Path
	WorkDir, _ = os.Getwd()
	viper.AddConfigPath(WorkDir + "/conf")

	viper.SetConfigType(Ext)
}

func main() {
	fmt.Println("Thanks for using the mdm")

	/* ---------- Main ---------- */

	// Setting
	viper.SetConfigName("main")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper.ReadInConfig(): %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Logger
	logger.Init(fmt.Sprintf(`%s\%s`, WorkDir, viper.GetString("log.dir")))
	go logger.Serve(ctx)

	logger.Loading(zapcore.InfoLevel, "loading setting")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	req.DevMode()
	client := req.C()
	resp, err := client.R().Get("https://httpbin.org/uuid")
	if err != nil {
		logger.Loading(zapcore.ErrorLevel, "fetch uuid")
	}
	logger.Loading(
		zapcore.InfoLevel,
		"fetch uuid",
		zap.String("resp", resp.String()),
	)

	select {
	case s := <-sig:
		cancel()
		log.Printf("%v<-sig", s)
	}
}
