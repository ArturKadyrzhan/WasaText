package api

import (
	"WasaText/internal/consts"
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

//I used middleware for catch request until it enter api,for security maybe if someone hacks

// Structure of logging response
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader Override  method to catch status
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Overriding the Write method to capture the content size
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		// Если статус-код еще не установлен, устанавливаем 200 OK
		lrw.WriteHeader(http.StatusOK)
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// Flush method for logging and finalizing output
func (lrw *loggingResponseWriter) Flush() {
	if fl, ok := lrw.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

// Middleware for logging requests
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter
		lrw := &loggingResponseWriter{ResponseWriter: w}

		// Передаем управление следующему обработчику
		next.ServeHTTP(lrw, r)

		// Логирование запроса
		methodColor := consts.Green // Устанавливаем цвет метода
		if lrw.statusCode >= 400 && lrw.statusCode < 500 {
			methodColor = consts.Yellow // Ошибка клиента
		} else if lrw.statusCode >= 500 {
			methodColor = consts.Red // Ошибка сервера
		}

		log.Printf("%s[%s%s%s%s] \"%s\" - status - %s%d%s, size %d bytes in %v second%s",
			consts.Cyan,                         // Цвет метки времени
			methodColor, r.Method, consts.Reset, // Цвет метода
			consts.Cyan, r.URL.Path, // Путь
			methodColor, lrw.statusCode, consts.Reset, // Цветной статус
			lrw.size, time.Since(start).Seconds(), consts.Reset) // Время выполнения
	})
}

// for web
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()

}
