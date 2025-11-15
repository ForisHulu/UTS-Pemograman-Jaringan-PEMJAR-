package main

import (
	"log"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// --- Konfigurasi ---
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Izinkan semua
	},
}

const (
	floorRR         = 100 // Rank terendah adalah Bronze (100 RR)
	defaultRR       = 100 // Pemain baru mulai di 100 RR
	maxRRDifference = 100  // Server akan mencoba mencari lawan dalam 100 RR
	baseRRChange    = 20  // RR dasar yang didapat/hilang
	rrBonusFactor   = 10  // Seberapa besar pengaruh perbedaan rank (dibagi 10)
)

// --- Structs (Blueprint) ---

// Message: Blueprint untuk semua pesan JSON
type Message struct {
	Type          string   `json:"type"` // "join", "move", "rematch"
	Position      int      `json:"position"`
	Symbol        string   `json:"symbol,omitempty"`
	Board         []string `json:"board,omitempty"`
	Turn          string   `json:"turn,omitempty"`
	Winner        string   `json:"winner,omitempty"`
	RematchStatus []bool   `json:"rematch_status,omitempty"`

	RankName string `json:"rankName,omitempty"`
	//
	// <-- PERBAIKAN BUG "undefined" ADA DI SINI
	//
	RR     int `json:"rr"`     // "omitempty" DIHAPUS
	Change int `json:"change"` // "omitempty" DIHAPUS
}

// Client: Mewakili satu pemain
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
	rr   int // Ini adalah memory rank Anda
	game *Game
}

// Game: Mewakili satu sesi game
type Game struct {
	players  [2]*Client
	board    [9]string
	turn     string
	gameOver bool
	rematch  [2]bool
	lock     sync.Mutex
}

// Hub: Lobi utama! Mengelola semua klien dan matchmaking
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	queue      []*Client // Satu antrian untuk semua
	lock       sync.Mutex
}

// --- Logika Hub (Lobi) ---

// newHub: Membuat Hub baru
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		queue:      make([]*Client, 0),
	}
}

// run: Menjalankan Hub (loop utama)
func (h *Hub) run() {
	go h.runMatchmaker()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Klien terhubung (Total: %d)", len(h.clients))
			// Kirim info rank awal
			client.send <- Message{
				Type:     "rank_info",
				RankName: getRankName(client.rr),
				RR:       client.rr,
				Change:   0, // Kirim change 0
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("Klien terputus.")
				delete(h.clients, client)
				close(client.send)
				h.removeFromQueue(client)
				if client.game != nil {
					client.game.handlePlayerDisconnect(client)
				}
			}
		}
	}
}

// addToQueue: Menambahkan klien ke pool matchmaking
func (h *Hub) addToQueue(client *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, c := range h.queue {
		if c == client {
			log.Println("Klien sudah ada di antrian.")
			return
		}
	}

	h.queue = append(h.queue, client)
	log.Printf("Klien (RR: %d) ditambahkan ke antrian. (total antri %d)", client.rr, len(h.queue))
	client.send <- Message{Type: "wait"}
}

// removeFromQueue: Menghapus klien dari pool
func (h *Hub) removeFromQueue(client *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for i, c := range h.queue {
		if c == client {
			h.queue = append(h.queue[:i], h.queue[i+1:]...)
			log.Printf("Klien (RR: %d) dihapus dari antrian.", client.rr)
			return
		}
	}
}

// runMatchmaker: Goroutine yang terus menerus mencari jodoh
func (h *Hub) runMatchmaker() {
	for {
		time.Sleep(2 * time.Second)
		h.lock.Lock()

		if len(h.queue) < 2 {
			h.lock.Unlock()
			continue
		}

		sort.Slice(h.queue, func(i, j int) bool {
			return h.queue[i].rr < h.queue[j].rr
		})

		var p1, p2 *Client
		var matchedIndex int = -1

		for i := 0; i < len(h.queue)-1; i++ {
			p1 = h.queue[i]
			p2 = h.queue[i+1]

			diff := int(math.Abs(float64(p1.rr - p2.rr)))

			if diff <= maxRRDifference {
				matchedIndex = i
				break
			}
		}

		if matchedIndex != -1 {
			p1 = h.queue[matchedIndex]
			p2 = h.queue[matchedIndex+1]

			h.queue = append(h.queue[:matchedIndex+1], h.queue[matchedIndex+2:]...)
			h.queue = append(h.queue[:matchedIndex], h.queue[matchedIndex+1:]...)

			log.Printf("Game ditemukan! RR %d vs RR %d", p1.rr, p2.rr)

			game := newGame(p1, p2)
			p1.game = game
			p2.game = game
		}

		h.lock.Unlock()
	}
}

// --- Logika Game (Satu Sesi) ---

// BARU: Fungsi untuk mapping 8-Rank Valorant
func getRankName(rr int) string {
	if rr < 200 {
		return "Bronze"
	} // 100-199
	if rr < 300 {
		return "Silver"
	} // 200-299
	if rr < 400 {
		return "Gold"
	} // 300-399
	if rr < 500 {
		return "Platinum"
	} // 400-499
	if rr < 600 {
		return "Diamond"
	} // 500-599
	if rr < 700 {
		return "Ascendant"
	} // 600-699
	if rr < 800 {
		return "Immortal"
	} // 700-799
	return "Radiant" // 800+
}

// calculateRRChange: Menghitung perubahan RR
func calculateRRChange(winnerRR, loserRR int) int {
	rrDiff := loserRR - winnerRR
	bonus := rrDiff / rrBonusFactor
	change := baseRRChange + bonus
	if change < 5 {
		return 5
	} // Perubahan minimal adalah 5
	return change
}

// newGame: Membuat game baru untuk 2 pemain
func newGame(p1, p2 *Client) *Game {
	g := &Game{
		players:  [2]*Client{p1, p2},
		board:    [9]string{"", "", "", "", "", "", "", "", ""},
		turn:     "X",
		gameOver: false,
		rematch:  [2]bool{false, false},
	}

	p1.send <- Message{Type: "start", Symbol: "X", Board: g.board[:], Turn: g.turn}
	p2.send <- Message{Type: "start", Symbol: "O", Board: g.board[:], Turn: g.turn}
	return g
}

// broadcast: Mengirim pesan ke KEDUA pemain di game ini
func (g *Game) broadcast(msg Message) {
	for _, client := range g.players {
		if client != nil {
			select {
			case client.send <- msg:
			default:
				log.Println("Gagal kirim pesan ke klien, channel penuh.")
			}
		}
	}
}

// broadcastState: Mengirim update papan
func (g *Game) broadcastState() {
	msg := Message{
		Type:  "update",
		Board: g.board[:],
		Turn:  g.turn,
	}
	g.broadcast(msg)
}

// checkWinner, checkDraw
func (g *Game) checkWinner() string {
	lines := [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
	for _, line := range lines {
		if g.board[line[0]] != "" && g.board[line[0]] == g.board[line[1]] && g.board[line[1]] == g.board[line[2]] {
			return g.board[line[0]]
		}
	}
	return ""
}
func (g *Game) checkDraw() bool {
	for _, cell := range g.board {
		if cell == "" {
			return false
		}
	}
	return g.checkWinner() == ""
}

// handleMove: DIPERBARUI dengan logika RANK FLOOR
func (g *Game) handleMove(client *Client, position int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.gameOver {
		return
	}

	var clientSymbol string
	var opponent *Client

	if client == g.players[0] {
		clientSymbol = "X"
		opponent = g.players[1]
	} else if client == g.players[1] {
		clientSymbol = "O"
		opponent = g.players[0]
	} else {
		return
	}

	if clientSymbol != g.turn {
		return
	}
	if position < 0 || position > 8 || g.board[position] != "" {
		return
	}

	g.board[position] = clientSymbol
	log.Printf("Gerakan %s ke %d", clientSymbol, position)

	winner := g.checkWinner()
	if winner != "" {
		g.gameOver = true
		g.broadcast(Message{Type: "winner", Winner: winner, Board: g.board[:]})

		// --- LOGIKA RR BARU DENGAN RANK FLOOR ---
		var winnerClient, loserClient *Client
		if winner == clientSymbol {
			winnerClient = client
			loserClient = opponent
		} else {
			winnerClient = opponent
			loserClient = client
		}

		change := calculateRRChange(winnerClient.rr, loserClient.rr)

		// Pemenang selalu dapat RR
		winnerClient.rr += change

		var loserChange = -change // Simpan perubahan negatif

		// Pecundang hanya kehilangan RR jika di atas "lantai"
		if loserClient.rr > floorRR {
			loserClient.rr -= change
			if loserClient.rr < floorRR { // Jangan biarkan jatuh di bawah lantai
				loserClient.rr = floorRR
			}
		} else {
			// Jika RR mereka sudah di lantai, perubahan RR mereka 0
			loserChange = 0 // <-- PERBAIKAN EKSPLISIT DI SINI
		}
		// --- AKHIR LOGIKA RR BARU ---

		log.Printf("Pemenang (%s) RR: %d (+%d) | Kalah (%s) RR: %d (%d)",
			winner, winnerClient.rr, change, loserClient.rr, loserChange)

		winnerClient.send <- Message{Type: "rank_update", RankName: getRankName(winnerClient.rr), RR: winnerClient.rr, Change: change}
		loserClient.send <- Message{Type: "rank_update", RankName: getRankName(loserClient.rr), RR: loserClient.rr, Change: loserChange}
		return
	}

	if g.checkDraw() {
		g.gameOver = true
		g.broadcast(Message{Type: "draw", Board: g.board[:]})
		return
	}

	if g.turn == "X" {
		g.turn = "O"
	} else {
		g.turn = "X"
	}
	g.broadcastState()
}

// handleRematch
func (g *Game) handleRematch(client *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if !g.gameOver {
		return
	}

	playerIndex := -1
	if g.players[0] == client {
		playerIndex = 0
	} else {
		playerIndex = 1
	}
	if playerIndex == -1 || g.rematch[playerIndex] {
		return
	}

	g.rematch[playerIndex] = true

	var clientSymbol string
	if playerIndex == 0 {
		clientSymbol = "X"
	} else {
		clientSymbol = "O"
	}
	log.Printf("Pemain %s ingin rematch.", clientSymbol)

	g.broadcast(Message{Type: "rematch_status", RematchStatus: g.rematch[:]})

	if g.rematch[0] && g.rematch[1] {
		log.Println("Rematch dimulai!")
		g.board = [9]string{"", "", "", "", "", "", "", "", ""}
		g.turn = "X"
		g.gameOver = false
		g.rematch = [2]bool{false, false}
		g.broadcastState()
	}
}

// handlePlayerDisconnect
func (g *Game) handlePlayerDisconnect(client *Client) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.gameOver {
		if client == g.players[0] {
			g.players[0] = nil
		}
		if client == g.players[1] {
			g.players[1] = nil
		}
		return
	}

	var remainingPlayer *Client
	if g.players[0] == client {
		remainingPlayer = g.players[1]
	} else {
		remainingPlayer = g.players[0]
	}

	g.gameOver = true
	g.players[0] = nil
	g.players[1] = nil

	if remainingPlayer != nil {
		log.Println("Lawan terputus. Mengirim pemain kembali ke lobi.")
		remainingPlayer.game = nil
		select {
		case remainingPlayer.send <- Message{Type: "opponent_left"}:
		default:
			log.Println("Gagal kirim 'opponent_left', channel penuh.")
		}
		remainingPlayer.hub.addToQueue(remainingPlayer)
	}
}

// --- Logika Klien (Satu Koneksi) ---

// readPump: Goroutine untuk membaca pesan DARI klien
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Println("Koneksi terputus:", err)
			break
		}

		switch msg.Type {
case "set_rr":
    if msg.RR >= floorRR {
        c.rr = msg.RR
        log.Printf("Pemain mengatur RR: %d → %s", c.rr, getRankName(c.rr))
        c.send <- Message{
            Type:     "rank_info",
            RankName: getRankName(c.rr),
            RR:       c.rr,
            Change:   0,
        }
    }

case "join":
    c.hub.addToQueue(c)

case "move":
    if c.game != nil {
        c.game.handleMove(c, msg.Position)
    }

case "rematch":
    if c.game != nil {
        c.game.handleRematch(c)
    }
}
	}
}

// writePump: Goroutine untuk menulis pesan KE klien
func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Println("Gagal menulis ke klien:", err)
			break
		}
	}
}

// handleConnection: Fungsi utama saat klien terhubung
func handleConnection(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan Message, 256),
		rr:   defaultRR, // Klien baru mulai dengan RR default (100)
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// --- Main ---
func main() {
	log.Println("Server Matchmaking (v4.3 - FINAL FIX) dimulai di http://localhost:8080")
	hub := newHub()
	go hub.run()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(hub, w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
