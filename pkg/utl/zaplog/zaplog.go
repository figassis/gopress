package zaplog

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	echo "github.com/labstack/echo/v4"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log represents zerolog logger
type (
	Log struct {
		logger  *zap.Logger
		version string
		// sugar   *zap.SugaredLogger
	}
)

var log *Log

// New instantiates new zero logger
func New(version string) *Log {

	cfg := zap.NewProductionConfig()
	if os.Getenv("ENVIRONMENT") != "production" {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		return nil
	}

	log = &Log{logger: logger, version: version}
	return log
}

// Log logs using zerolog
func (z *Log) Log(ctx echo.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if id, ok := ctx.Get("id").(int); ok {
		params["id"] = id
		params["user"] = ctx.Get("username").(string)
	}

	var fields []zap.Field
	fields = append(fields, zap.String("version", z.version))
	for name, param := range params {
		fields = append(fields, zap.Any(name, param))
	}

	if err != nil {
		params["error"] = err
		z.logger.Error(msg, fields...)
		return
	}

	z.logger.Info(msg, fields...)
}

// Log logs using zerolog
func (z *Log) SimpleLog(msg interface{}) {
	switch msgType := msg.(type) {
	case error:
		if msgType == nil {
			return
		}
		z.logger.Error(msgType.Error())
	case string:
		z.logger.Info(msgType)
	default:
		z.logger.Info(fmt.Sprint(msg))
	}
}

func ZLog(msg interface{}) (err error) {
	switch msgType := msg.(type) {
	case nil:
		return
	case error:
		log.logger.Error(msgType.Error())
		return msgType
	case string:
		log.logger.Info(msgType)
	default:
		if JSON(msg) != "" {
			log.logger.Info(JSON(msg))
		}
	}
	return nil
}

func JSON(value interface{}) string {
	bytes, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (z Log) GetLogger() *zap.Logger {
	return z.logger
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// MyCaller returns the caller of the function that called it :)
func MyCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2).Function
}
