package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

var logLevelNames = map[LogLevel]string{
	DebugLevel:   "DEBUG",
	InfoLevel:    "INFO",
	WarningLevel: "WARNING",
	ErrorLevel:   "ERROR",
	FatalLevel:   "FATAL",
}

type Logger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	level         LogLevel
	file          *os.File
}

type Config struct {
	EnableConsole    bool
	EnableFile       bool
	FilePath         string
	Level            LogLevel
	IncludeTimestamp bool
	IncludeFileLine  bool
}

func DefaultConfig() Config {
	return Config{
		EnableConsole:    true,
		EnableFile:       true,
		FilePath:         "application.log",
		Level:            InfoLevel,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
	}
}

func New(config Config) (*Logger, error) {
	logger := &Logger{
		level: config.Level,
	}

	if config.EnableConsole {
		logger.consoleLogger = log.New(os.Stdout, "", 0)
	}

	if config.EnableFile {
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = file
		logger.fileLogger = log.New(file, "", 0)
	}

	return logger, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) formatMessage(level LogLevel, message string, includeCallerInfo bool) string {
	var builder strings.Builder

	builder.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	builder.WriteString(" ")

	builder.WriteString("[")
	builder.WriteString(logLevelNames[level])
	builder.WriteString("] ")

	if includeCallerInfo {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			shortFile := filepath.Base(file)
			builder.WriteString(fmt.Sprintf("%s:%d ", shortFile, line))
		}
	}

	builder.WriteString(message)

	return builder.String()
}

func (l *Logger) log(level LogLevel, message string, includeCallerInfo bool) {
	if level < l.level {
		return
	}

	formattedMessage := l.formatMessage(level, message, includeCallerInfo)

	if l.consoleLogger != nil {
		l.consoleLogger.Println(formattedMessage)
	}

	if l.fileLogger != nil {
		l.fileLogger.Println(formattedMessage)
	}
}

func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message, true)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...), true)
}

func (l *Logger) Info(message string) {
	l.log(InfoLevel, message, true)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...), true)
}

func (l *Logger) Warning(message string) {
	l.log(WarningLevel, message, true)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.log(WarningLevel, fmt.Sprintf(format, args...), true)
}

func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message, true)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), true)
}

func (l *Logger) Fatal(message string) {
	l.log(FatalLevel, message, true)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...), true)
	os.Exit(1)
}

func (l *Logger) SetOutput(w io.Writer) {
	if l.consoleLogger != nil {
		l.consoleLogger.SetOutput(io.MultiWriter(os.Stdout, w))
	}

	if l.fileLogger != nil && l.file != nil {
		l.fileLogger.SetOutput(io.MultiWriter(l.file, w))
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) GetLevel() LogLevel {
	return l.level
}

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = New(DefaultConfig())
	if err != nil {
		defaultLogger = &Logger{
			consoleLogger: log.New(os.Stdout, "", 0),
			level:         InfoLevel,
		}
		defaultLogger.Error(fmt.Sprintf("Failed to initialize default logger: %v", err))
	}
}

func LogDebug(message string) {
	defaultLogger.Debug(message)
}

func LogDebugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func LogInfo(message string) {
	defaultLogger.Info(message)
}

func LogInfof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func LogWarning(message string) {
	defaultLogger.Warning(message)
}

func LogWarningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

func LogError(message string) {
	defaultLogger.Error(message)
}

func LogErrorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func LogFatal(message string) {
	defaultLogger.Fatal(message)
}

func LogFatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func SetGlobalLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func CloseLogger() error {
	return defaultLogger.Close()
}
