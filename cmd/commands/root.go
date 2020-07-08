package commands

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/binacsgo/log"

	"github.com/OpenKikCoc/raftkv/config"
)

var (
	configFile string
	zaplogger  *zap.Logger
	logger     log.Logger
	cfg        *config.Config
)

func init() {
	rootCmdFlags(RootCmd)
}

func rootCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&configFile, "configFile", "config.toml", "config file (default is ./config.toml)")
}

var (
	// RootCmd the root command
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Root Command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Println("Load config from file path:", configFile)
			cfg, err = config.LoadFromFile(configFile)
			if err != nil {
				fmt.Println("LoadFromFile err: ", err)
			}
			zaplogger = initLogger(cfg.WorkSpace, cfg.LogConfig)
			logger = log.NewZapLoggerWapper(zaplogger.Sugar())
			logger.Info("Root init finished")
			return nil
		},
	}
)

func initLogger(rootPath string, logConfig config.LogConfig) *zap.Logger {
	logpath := logConfig.File
	if !path.IsAbs(logpath) {
		logpath = path.Join(rootPath, logConfig.File)
	}
	fmt.Printf("Log path : %s\n", logpath)
	hook := lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    logConfig.Maxsize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.Maxage,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "_time",
		LevelKey:       "_level",
		NameKey:        "_logger",
		CallerKey:      "_caller",
		MessageKey:     "_message",
		StacktraceKey:  "_stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(stringToXZapLoggerLevel(logConfig.Level))
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(&hook), atomicLevel)
	logger := zap.New(core)
	return logger
}

func stringToXZapLoggerLevel(level string) zapcore.Level {
	lower := strings.ToLower(level)
	switch lower {
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
