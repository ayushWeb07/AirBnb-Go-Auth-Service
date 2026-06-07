package controllers

import "net/http"

func CheckHealthStatus(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Working fine as always"))
}
