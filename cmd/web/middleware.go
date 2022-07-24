package main

import (
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' font.googleapis.com; font-src fonts.gstatic.com",
		)

		writer.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "deny")
		writer.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(writer, req)

	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())

		next.ServeHTTP(writer, req)
	})
}
