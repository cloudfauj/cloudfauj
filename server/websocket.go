package server

import "github.com/gorilla/websocket"

// wsManager adds some convenience methods on top of a Gorilla websocket connection
type wsManager struct {
	*websocket.Conn
}

func (m *wsManager) sendTextMsg(msg string) error {
	return m.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (m *wsManager) sendClosureMsg(code int) error {
	return m.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, ""),
	)
}

func (m *wsManager) sendFailure(msg string, code int) error {
	m.sendTextMsg(msg)
	return m.sendClosureMsg(code)
}

func (m *wsManager) sendFailureISE() error {
	return m.sendFailure("An internal server error occured", websocket.CloseInternalServerErr)
}

func (m *wsManager) sendSuccess(msg string) error {
	m.sendTextMsg(msg)
	return m.sendClosureMsg(websocket.CloseNormalClosure)
}
