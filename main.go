package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

var sugarLogger *zap.SugaredLogger

func bruteDefender(tokenID string, requestID string, password string) error {
	ctx := context.Background()

	var myPass string
	myPass = "mypass"
	isBlocked := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	userKey := tokenID + requestID

	_, err := rdb.Get(ctx, userKey).Result()

	if err != nil {
		var index int
		index = 0
		rdb.Set(ctx, userKey, index, 0)
	}

	attempt, err := rdb.Get(ctx, userKey).Result()
	if err != nil {
		panic(err)
	}

	intAttempt, err := strconv.Atoi(attempt)
	if err != nil {
		panic(err)
	}

	if intAttempt != 3 {
		intAttempt++
		rdb.Set(ctx, userKey, intAttempt, 10*time.Minute)
	} else {
		isBlocked = 1
	}

	infoMsg := "[TokenID: " + fmt.Sprint(tokenID) + " & RequestID: " + fmt.Sprint(requestID) + "]  |  "

	if isBlocked == 1 {
		sugarLogger.Error(infoMsg + "Error Blocked!")
	} else if myPass == password {
		sugarLogger.Info(infoMsg + "Info Success.. Entered to the system!")
	} else {
		sugarLogger.Error(infoMsg + "Error Wrong Password!")
	}
	return err
}

func intToString(value int) string {
	return strconv.Itoa(value)
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// Print function lines
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// The format time can be customized
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	bruteDefender("token056", "request272", "mypass")
}
