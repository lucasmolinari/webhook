package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PongMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]string)

		data, _ := io.ReadAll(r.Body)
		json.Unmarshal(data, &resp)
		r.Body = io.NopCloser(bytes.NewBuffer(data))

		ping, ok := resp["ping"]
		if !ok {
			next(w, r)
			return

		}
		if ping == "ping" {
			fmt.Println("ping recebido")
			resp := map[string]string{"result": "pong"}
			json.NewEncoder(w).Encode(resp)
			return

		}
		next(w, r)
		return
	}

}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(string(data))
	fmt.Println()
	fmt.Println()

	w.WriteHeader(200)
	w.Write([]byte("recebido do webhook"))

}

func main() {
	m := http.NewServeMux()
	m.Handle("POST /webhook", PongMiddleware(http.HandlerFunc(WebhookHandler)))
	http.ListenAndServe(":5000", m)
}
