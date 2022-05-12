package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
)

func newRecoveryMiddleware(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				err := recover()

				if err != nil {
					logger.Println(err, string(debug.Stack()))
					response := createDomainErrorResponse(domain.ErrInternalError)
					bytes, err := json.Marshal(&response)

					if err != nil {
						writer.WriteHeader(http.StatusInternalServerError)
						logger.Println(err)
						return
					}

					writer.Header().Set("Content-Type", contentTypeJSON)
					writer.WriteHeader(http.StatusOK)
					_, err = writer.Write(bytes)

					if err != nil {
						writer.WriteHeader(http.StatusInternalServerError)
						logger.Println(err)
					}
				}
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
