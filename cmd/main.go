package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/effective_mobile_test/docs"
	"github.com/rlapenok/effective_mobile_test/internal/api/handlers"
	"github.com/rlapenok/effective_mobile_test/internal/config"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(level *string) {
	var zapLevel zapcore.Level
	zapLevel, err := zapcore.ParseLevel(*level)
	if err != nil {
		zapLevel = zapcore.DebugLevel
	}
	encoderConfig := zapcore.EncoderConfig{

		TimeKey:      "time",
		LevelKey:     "level",
		NameKey:      "logger",
		CallerKey:    "caller",
		FunctionKey:  zapcore.OmitKey,
		MessageKey:   "msg",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderConfig
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.Development = true
	cfg.Encoding = "console"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}

	logger, err := cfg.Build()
	if err != nil {
		slog.Error("Error while build logger", "err", err)

	}
	zap.ReplaceGlobals(logger)
}

// @title		Effective mobile
// @BasePath	/
func main() {
	//parse flag
	configPath := getFlag()
	//load config from file and map into struct
	config := config.Load(configPath)
	//get log level
	level := config.GetLevel()
	//init logger
	initLogger(level)
	defer zap.L().Sync()
	//create state
	config.ToState()
	gin.SetMode(gin.ReleaseMode)
	//create router
	router := gin.Default()
	//create endpoints
	router.Use(cors.Default())
	router.POST("/add_song", handlers.AddSong)
	router.DELETE("/delete_song/:id", handlers.DeleteSong)
	router.PATCH("/change_song/:id", handlers.ChangeSong)
	router.GET("/get_info", handlers.GetInfo)
	router.GET("/lyrics/:id", handlers.GetLyrics)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", 8080)
	adress := fmt.Sprintf(":%d", config.Server.Port)
	zap.L().Info("Server start on", zap.String("address", adress))
	//start server
	router.Run(adress)

}

func getFlag() string {
	configPath := flag.String("path_to_config", ".env", "Path to file with config for app")
	flag.Parse()
	return *configPath
}
