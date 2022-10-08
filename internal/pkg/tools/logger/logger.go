package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

func New(out io.Writer, formatter logrus.Formatter, level logrus.Level) *logrus.Logger {
	return &logrus.Logger{
		Out:       out,
		Formatter: formatter,
		Level:     level,
	}
}
