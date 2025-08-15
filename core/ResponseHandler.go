package core

import (
	"encoding/json"
	"fmt"
	"wonderful/structs"
	"github.com/gorilla/websocket"
)

// Counting the number of functions calls that waits
var FunctionCallRequests int = 0

// MessageReceiver processes incoming messages
func ResponseEvent(conn *websocket.Conn, message []byte, syncChannel chan bool) {
	var modelResponse struct{ Type string }

	if err := json.Unmarshal(message, &modelResponse); err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	switch modelResponse.Type {

	// Server send answer,
	case "response.text.delta":
		var responseText structs.Response_Text_Delta

		if err := json.Unmarshal(message, &responseText); err == nil {
			fmt.Printf("%s", responseText.Delta)
		}

	case "response.function_call_arguments.done":
		FunctionCallRequests++ // Server sends function call request and the user detect it

	case "response.done":
		if FunctionCallRequests > 0 {
			var response structs.Response_Done

			if err := json.Unmarshal(message, &response); err != nil {
				fmt.Println("Unmarshal error:", err)
				return
			}

			//Allows handling more then one function call request
			for _, output := range response.Response.Output {
				call := structs.Pending_Call{
					CallID: output.CallID,
					Name:   output.Name,
					Args:   output.Arguments,
				}
				FunctionCallEvent(conn, call)
			}
			FunctionCallRequests-- // Finished to handle the request
		} else {
			syncChannel <- true // Allows new client request
		}

	case "error":
		fmt.Printf("Error: %s", string(message))
		syncChannel <- true
	}
}

// Handle function call responses
func FunctionCallEvent(conn *websocket.Conn, call structs.Pending_Call) {

	var result string

	//Calling the correct function by name
	switch call.Name {
	case "Multiplier":
		result = Multiplier(call.Args)

	default:
		fmt.Println("Function not found")
	}

	functionResult := structs.Funcation_Conversation_Item_Create{
		Type: "conversation.item.create",
		Item: structs.Item_Function{Type: "function_call_output",
			CallID: call.CallID,
			Output: fmt.Sprintf(`{"result": %s}`, result),
		},
	}

	message, err := json.Marshal(functionResult)
	if err != nil {
		fmt.Println("Marshal error:", err)
	}
	//Send the result of function to the server
	conn.WriteMessage(websocket.TextMessage, message)

	response := structs.Response_Create{
		Type: "response.create",
		Response: structs.Response{
			Modalities: []string{"text"},
		},
	}

	message, err = json.Marshal(response)
	if err != nil {
		fmt.Println("Marshal error:", err)
	}

	conn.WriteMessage(websocket.TextMessage, message)
}
