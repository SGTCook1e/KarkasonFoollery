package network

import (
	b "KarkasonFoollery/internal/board"
	"KarkasonFoollery/internal/game"
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

type GameServer struct {
	state *game.GameState

	players   map[b.PlayerID]net.Conn
	playersMu sync.Mutex

	actionsChan chan PlayerAction
	listener    net.Listener

	logger *log.Logger
}

func NewGameServer(l net.Listener, tiles []*b.Tile) *GameServer {
	return &GameServer{
		listener:    l,
		players:     make(map[b.PlayerID]net.Conn),
		actionsChan: make(chan PlayerAction, 10),
		state:       game.NewState(tiles),
		logger:      log.New(os.Stdout, "[Server] ", log.Ltime),
	}
}

func (s *GameServer) Start() {
	s.logger.Println("Listening to new connections...")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		s.playersMu.Lock()

		if len(s.players) >= 5 {
			s.playersMu.Unlock()
			conn.Close()
			continue
		}

		id := s.state.NewPlayerId

		s.players[id] = conn
		s.playersMu.Unlock()

		s.logger.Printf("Player %d connected successfully!\n", id)

		go s.listenToPlayer(id, conn)

		s.state.NewPlayerId = s.state.AddPlayer()

		if len(s.players) == 1 {
			go s.runGameLogic()
		}
	}
}

func (s *GameServer) listenToPlayer(id b.PlayerID, conn net.Conn) {
	defer func() {
		conn.Close()
		s.handleDisconnect(conn)
	}()

	reader := bufio.NewReader(conn)

	for {
		packet, err := reader.ReadBytes('\n')
		if err != nil {
			s.logger.Printf("Player %d disconnected!\n", id)
			return
		}

		cleanPacket := packet[:len(packet)-1]

		packetCopy := make([]byte, len(cleanPacket))
		copy(packetCopy, cleanPacket)

		s.actionsChan <- PlayerAction{
			PlayerID: id,
			Data:     packetCopy,
		}
	}
}

func (s *GameServer) runGameLogic() {
	s.sendInitialPlayersData()
	s.logger.Println("Game started")

	for action := range s.actionsChan {
		if s.state.CurrPlayer != action.PlayerID {
			continue
		}

		s.logger.Printf("CurrPlayer: %d\n", s.state.CurrPlayer)
		for i, player := range s.state.Players {
			s.logger.Printf("Player [%d]: %+v\n", i, player)
		}

		// result := game.ResolvePlacement(*s.state)
		// s.state.ApplyPlacement(*result, 1)

		// s.state.PassTurn()
		// s.state.CurrTile = s.state.Deck.Draw()

		// s.broadcast(action.Data)
	}
}

func (s *GameServer) broadcast(data []byte) {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()

	for id, conn := range s.players {
		_, err := conn.Write(data)
		if err != nil {
			s.logger.Printf("Error sending data to player %d: %v\n", id, err)
		}
	}
}

func (s *GameServer) handleDisconnect(conn net.Conn) {

}

func (s *GameServer) sendInitialPlayersData() {
	payload := make([]b.PlayerID, 0, len(s.state.Players))
	for _, p := range s.state.Players {
		payload = append(payload, p.Id)
	}

	msgBytes, err := makeMessageBytes(MsgStartPacket, payload)
	if err != nil {
		s.logger.Printf("Message marshal error: %v", err)
		return
	}
	s.broadcast(msgBytes)
}
