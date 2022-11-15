package hardlogger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func (f Fields) Extend(f2 Fields) {
	for k, v := range f2 {
		f[k] = v
	}
}

func (f Fields) Copy() Fields {
	ret := make(Fields, len(f))
	for k, v := range f {
		ret[k] = v
	}
	return ret
}

type Level string

const (
	LevelError Level = "error"
	LevelWarn  Level = "warn"
	LevelInfo  Level = "info"
	LevelDebug Level = "debug"
)

var configLevels = map[Level]logrus.Level{
	LevelError: logrus.ErrorLevel,
	LevelWarn:  logrus.WarnLevel,
	LevelInfo:  logrus.InfoLevel,
	LevelDebug: logrus.DebugLevel,
}

type Logger interface {
	WithCtx(ctx context.Context) Logger
	WithField(key string, val interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	//Debug(args ...interface{})
	//Debugf(format string, args ...interface{})
	Info(args ...interface{})
	//Infof(format string, args ...interface{})
	//Warn(args ...interface{})
	//Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	//Fatal(args ...interface{})
	//Fatalf(format string, args ...interface{})
	//Print(args ...interface{})
	//Printf(format string, args ...interface{})
	//Log(level Level, args ...interface{})
	//Logf(level Level, format string, args ...interface{})
	WithCaller(ctx context.Context) context.Context
}

type fieldsCtxKey struct{}

func SetCtxFields(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, fieldsCtxKey{}, fields)
}

func GetCtxFields(ctx context.Context) Fields {
	fields := getCtxFields(ctx)
	if fields == nil {
		return Fields{}
	}

	return fields.Copy()
}

func getCtxFields(ctx context.Context) Fields {
	fieldsVal := ctx.Value(fieldsCtxKey{})
	if fieldsVal == nil {
		return nil
	}

	return fieldsVal.(Fields)
}

func AddCtxFields(ctx context.Context, fields Fields) context.Context {
	ctxFields := getCtxFields(ctx)
	if ctxFields == nil {
		return SetCtxFields(ctx, fields)
	}

	for k, v := range fields {
		ctxFields[k] = v
	}

	return ctx
}
