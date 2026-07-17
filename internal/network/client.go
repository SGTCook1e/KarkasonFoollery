package network

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type ConnStatus string

const (
	Connected    ConnStatus = "connected"
	Connecting   ConnStatus = "connecting"
	Disconnected ConnStatus = "disconnected"
)

type GameClient struct {
	conn        net.Conn
	InboundChan chan []byte
	isClosed    bool

	mu     sync.RWMutex
	status ConnStatus

	logger *log.Logger
}

func NewGameClient() *GameClient {
	return &GameClient{
		InboundChan: make(chan []byte, 100),
		logger:      log.New(os.Stdout, "[Client] ", log.Ltime),
	}
}

func (c *GameClient) Connect(address string) {
	c.setStatus(Connecting)

	go func() {
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			c.setStatus(Disconnected)
			c.logger.Printf("Can't connect to server: %v", err)
			return
		}

		c.conn = conn
		c.setStatus(Connected)
		c.logger.Println("Connected to server successfully!")

		c.readLoop()
	}()
}

func (c *GameClient) GetStatus() ConnStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *GameClient) setStatus(s ConnStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = s
}

func (c *GameClient) readLoop() {
	defer c.Close()
	reader := bufio.NewReader(c.conn)

	for {
		packet, err := reader.ReadBytes('\n')
		if err != nil {
			if !c.isClosed {
				c.logger.Printf("Connection to server broken: %v", err)
			}
			return
		}

		cleanPacket := packet[:len(packet)-1]

		packetCopy := make([]byte, len(cleanPacket))
		copy(packetCopy, cleanPacket)

		c.InboundChan <- packetCopy
	}
}

func (c *GameClient) Send(data []byte) error {
	if c.conn == nil {
		return net.ErrClosed
	}
	_, err := c.conn.Write(data)
	return err
}

func (c *GameClient) Close() {
	c.isClosed = true
	if c.conn != nil {
		c.conn.Close()
	}
}
