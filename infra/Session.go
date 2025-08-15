package infra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"wonderful/structs"
)

//Connect to openAI through websocket
func SessionConnect(openAIKey string) (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+openAIKey)
	headers.Set("OpenAI-Beta", "realtime=v1")

	url := "wss://api.openai.com/v1/realtime?model=gpt-4o-realtime-preview-2025-06-03"
	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	return conn, err
}

//Sends functions call tools to the server after session is on
func AddToolFunctions(conn *websocket.Conn, toolfuncs []structs.Tool) error{
	session := structs.Session_Update{
		Type: "session.update",
		Session: structs.Session_Tools{Tools: toolfuncs,},
	}

	updateRequest, err := json.Marshal(session)

	if err != nil {
		fmt.Println("Marshal error:", err)
		return err
	}
	//Sends function call tools to the server 
	return conn.WriteMessage(websocket.TextMessage, updateRequest)
}
