package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Gagal upgrade ke WebSocket: ", err)
		return
	}
	defer conn.Close()

	log.Println("Client terhubung!")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Gagal membaca pesan: ", err)
			break
		}

		log.Printf("Menerima pesan: %s", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Gagal menerima pesan:", err)
			break
		}
	}

	log.Println("Client terputus!")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Could not open requested file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})

	http.HandleFunc("/ws", websocketHandler)

	fmt.Println("Server berjalan di http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
