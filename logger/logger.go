package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Init()(err error) {

	writesyncer := getLogWriter(viper.GetString("log.filename"),
		viper.GetInt("log.max_size"), viper.GetInt("log.max_backups"),
		viper.GetInt("log.max_age"))
	encoder := getEncoder()

	var l=new(zapcore.Level)
	err=l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}
	corn := zapcore.NewCore(encoder, writesyncer, l)
	lg := zap.New(corn, zap.AddCaller())

	zap.ReplaceGlobals(lg)
	return
}

//Filename: 日志文件的位置
//MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
//MaxBackups：保留旧文件的最大个数
//MaxAges：保留旧文件的最大天数
//Compress：是否压缩/归档旧文件
func getLogWriter(fileName string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	encoderConfig.EncodeCaller=zapcore.FullCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
