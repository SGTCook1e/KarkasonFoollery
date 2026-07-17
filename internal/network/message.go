package network

import (
	b "KarkasonFoollery/internal/board"
	"encoding/json"
)

const (
	MsgStartPacket    = "START_PACKET"
	MsgStartGame      = "START_GAME"
	MsgMovePermission = "MOVE_PERMISSION"
	MsgBroadcastState = "BROADCAST_STATE"

	MsgPlayerMove = "PLAYER_MOVE"
)

type NetworkMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PlayerAction struct {
	PlayerID b.PlayerID
	Data     []byte
}

func makeMessageBytes(msgType string, payload any) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	msg := NetworkMessage{
		Type: msgType,
		Data: payloadBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	msgBytes = append(msgBytes, '\n')
	return msgBytes, nil
}
