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

//мидлвэйр нам нужны чтобы поймать запрос до того как он вошел в апи, пока не взломал,
//может мошенник, в нашем случе

// Структура для логирования ответа
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Переопределение метода WriteHeader для захвата статуса
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Переопределение метода Write для захвата размера контента
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		// Если статус-код еще не установлен, устанавливаем 200 OK
		lrw.WriteHeader(http.StatusOK)
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// Метод Flush для логирования и завершения вывода
func (lrw *loggingResponseWriter) Flush() {
	if fl, ok := lrw.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

// Middleware для логирования запросов
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

// for web socket,
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()

}
