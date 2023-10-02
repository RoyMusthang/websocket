package main

import (
	"encoding/json"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type User struct {
	Nick string
	con  *websocket.Conn
}

var clients map[*User]bool = make(map[*User]bool)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	nick := r.URL.Query().Get("nick")
	con, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Nick: nick,
		con:  con,
	}

	clients[&user] = true

	for {
		_, data, err := user.con.Read(r.Context())
		if err != nil {
			log.Println("conexao do client encerrada")
			delete(clients, &user)
			break
		}

		for client := range clients {
			client.con.Write(r.Context(), websocket.MessageText, []byte(data))
		}
	}
}

//firebird

type Pedido struct {
	CNPJ           string    `json:"CNPJ"`
	NUMPED         int       `json:"NUMPED"`
	CODUSUR        int       `json:"CODUSUR"`
	CODCLI         int       `json:"CODCLI"`
	VENDA          string    `json:"VENDA"`
	DEVOLUCAO      string    `json:"DEVOLUCAO"`
	DT_PEDIDO      time.Time `json:"DT_PEDIDO"`
	DT_FAT         time.Time `json:"DT_FAT"`
	DT_IMPLANTACAO time.Time `json:"DT_IMPLANTACAO"`
	DT_REFERENCIA  time.Time `json:"DT_REFERENCIA"`
	PERIODO        string    `json:"PERIODO"`
	CODPROD        string    `json:"CODPROD"`
	CODFORNEC      int       `json:"CODFORNEC"`
	CODDIVISAO     int       `json:"CODDIVISAO"`
	UND            int       `json:"UND"`
	CX             int       `json:"CX"`
	VL_UNIT        float64   `json:"VL_UNIT"`
	ACRES          float64   `json:"ACRES"`
	DESCONTO       float64   `json:"DESCONTO"`
	DESCONTOPED    float64   `json:"DESCONTOPED"`
	CANALVENDA     string    `json:"CANALVENDA"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var res []*User
		for c := range clients {
			res = append(res, c)
		}
		json.NewEncoder(w).Encode(res)
	})

	http.ListenAndServe(":3001", nil)
}
