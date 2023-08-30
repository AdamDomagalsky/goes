package gapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// another way: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/interceptors/logging/examples/zap/example_test.go

// 0log logger
func Setup0logGrpcLogger(env string) []grpc.ServerOption {
	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return []grpc.ServerOption{grpc.UnaryInterceptor(grpcLogger)}
}

func grpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode) // call original WriteHeader
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, r)
		duration := time.Since(startTime)
		logger := log.Info()

		if rec.StatusCode >= 400 {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Str("url", r.URL.String()).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP (gRPC GW) request")

	})
}

// ZAP logger
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func SetupZapGrpcLogger(env string) (serverOption []grpc.ServerOption, err error) {

	var logger *zap.Logger
	var opts []logging.Option
	if env == "development" {
		conf := zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = conf.Build()
		if err != nil {
			return nil, err
		}
		opts = []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			// Add any other option (check functions starting with logging.With).
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
		defer logger.Sync() // flushes buffer, if any
		// sugar := logger.Sugar()
		// sugar.Infow("failed to fetch URL",
		// 	// Structured context as loosely typed key-value pairs.
		// 	"url", url,
		// 	"attempt", 3,
		// 	"backoff", time.Second,
		// )
		// sugar.Infof("Failed to fetch URL: %s", url)
	}

	serverOption = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
			// Add any other interceptor you want.
		),
	}

	return
}
