package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"wonderful/infra"
	"wonderful/structs"
	"wonderful/core"
	"github.com/gorilla/websocket"
)

func main(){
	
	//Pull the key from the env variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
    	fmt.Println("OPENAI KEY environment variable is required")
		return
	}

	syncChannel := make(chan bool) //Prevents the user from sending request before the previous is ended

	//Conect to openai server through websocket
	conn ,err := infra.SessionConnect(apiKey)
	if err != nil{
		fmt.Println("Connection attempt failed:", err)
		return
	}

	//Add tools to the session
	err = infra.AddToolFunctions(conn, []structs.Tool{core.MultiplierTool()})
	
	if err != nil{
		fmt.Println("Add functions failed:", err)
	}
	
	fmt.Println("Welcome to your CLI interactive chatbot!\n" +
				"After we get started you will be able to ask any question.\n" +
				"You can exit in any time using the input `Exit` or press ctrl+C.")
	prompt := "Enter your question:"

	//Manage the user side of the websocket
	go websocketListner(conn, syncChannel)

	//User CLI loop, get request from the user and send it to the server
	for{
		Userinput := userPrompt(prompt)

		if Userinput == "" {
			fmt.Println("No input. please provide input")
			continue

		} else if strings.ToLower(Userinput) == "exit"{ //enable to EXIT for both upper and lower case
			break
			
		} else{
			message := structs.User_Conversation_Item_Create{
				Type: "conversation.item.create",
				Item: structs.Item_User{
					Type: "message",
					Role: "user",
					Content: []structs.Content{
						{
						Type: "input_text",
						Text: Userinput,
						},	
					},
				},
			}			
		
			messageRequest, err := json.Marshal(message)

			if err != nil {
				fmt.Println("Marshal error:", err)
				continue
			}

			//Sends the user request to the server
			conn.WriteMessage(websocket.TextMessage, messageRequest)
		}

		message := structs.Response_Create{
		Type: "response.create",
		Response: structs.Response{
			Modalities: []string{"text"},
				},
		}
		
		responseRequest, err := json.Marshal(message)

		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}
		
		//Sends reminder to the server to answer
		conn.WriteMessage(websocket.TextMessage, responseRequest)

		<-syncChannel //Barrier, blocks the user requests until responses the previous one

	}
}

//This function manage the interact with the user through the terminal
func userPrompt(prompt string) string{
	reader := bufio.NewReader(os.Stdin)
		fmt.Fprint(os.Stderr,"\n" + prompt + "\n")
		s, err := reader.ReadString('\n')

		if err != nil {
				fmt.Println("Error reading your input: ",err)
				s = "exit"
				return s
		}

		s = strings.TrimSpace(s)
		return s
}

//manage the whole user request cycle
func websocketListner(conn *websocket.Conn, syncChannel chan bool){
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Input reading error: ", err)
			syncChannel <- false
			return
		}

		// Handler the received message
		core.ResponseEvent(conn, message, syncChannel)
	}
}