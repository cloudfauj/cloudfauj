package wsmanager

import "github.com/gorilla/websocket"

// WSManager adds some convenience methods on top of a Gorilla websocket connection
type WSManager struct {
	*websocket.Conn
}

func (m *WSManager) Write(p []byte) (int, error) {
	return len(p), m.WriteMessage(websocket.TextMessage, p)
}

func (m *WSManager) SendTextMsg(msg string) error {
	return m.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (m *WSManager) SendClosureMsg(code int) error {
	return m.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, ""),
	)
}

func (m *WSManager) SendFailure(msg string, code int) error {
	m.SendTextMsg(msg)
	return m.SendClosureMsg(code)
}

func (m *WSManager) SendFailureISE() error {
	return m.SendFailure("An internal server error occured", websocket.CloseInternalServerErr)
}

func (m *WSManager) SendSuccess(msg string) error {
	m.SendTextMsg(msg)
	return m.SendClosureMsg(websocket.CloseNormalClosure)
}
