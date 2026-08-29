package http_middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	http_response "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type Middleware func(http.Handler) http.Handler

const (
	requestIDHeader            = "X-Request-ID"
	requestAuthorizationHeader = "Authorization"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			alowedOrigins := map[string]struct{}{
				"http://localhost:5051": {},
			}

			origin := r.Header.Get("Origin")
			if _, ok := alowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				core_logger.String("request_id", requestID),
				core_logger.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				core_logger.String("http_method", r.Method),
				core_logger.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				core_logger.Int("status_code", rw.GetStatusCode()),
				core_logger.Duration("latency", time.Since(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Auth(tokenGenerator *core_jwt.TokenGenerator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewResponseHandler(log, w)

			authHeader := r.Header.Get(requestAuthorizationHeader)
			if authHeader == "" {
				err := fmt.Errorf("missing authorization token: %w", core_errors.ErrUnauthorized)
				responseHandler.ErrorResponse(err, "authorization token is missing")
				return
			}

			fields := strings.Fields(authHeader)
			if len(fields) != 2 || fields[0] != "Bearer" {
				err := fmt.Errorf("invalid authorization header: %w", core_errors.ErrUnauthorized)
				responseHandler.ErrorResponse(err, "invalid authorization header")
				return
			}

			token := fields[1]
			claims, err := tokenGenerator.VerifyToken(token)
			if err != nil {
				err = fmt.Errorf("invalid jwt token: %v: %w", err, core_errors.ErrUnauthorized)
				responseHandler.ErrorResponse(err, "invalid jwt token")
				return
			}

			ctx = core_jwt.ToContext(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
