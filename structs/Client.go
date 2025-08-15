package structs

type Pending_Call struct {
	CallID    string
	Name      string
	Args 	 string
}

type Session_Update struct {
	Type    string `json:"type"`
	Session Session_Tools `json:"session"`
}

type Session_Tools struct {
	Tools []Tool `json:"tools"`
}

type Funcation_Conversation_Item_Create struct{
	Type string `json:"type"`
	Item Item_Function `json:"item"` 
}

type Item_Function struct{
	Type string `json:"type"`
	CallID string `json:"call_id"`
	Output string `json:"output"`
}

type User_Conversation_Item_Create struct {
	Type string `json:"type"`
	Item Item_User `json:"item"`
}

type Item_User struct{
	Type string `json:"type"`
	Role string `json:"role"`
	Content []Content `json:"content"`
}

type Content struct{
	Type 	string `json:"type"`
	Text    string `json:"text"`
}

type Tool struct{
	Type        string    	`json:"type"`
	Name        string     	`json:"name"`
	Description string     	`json:"description"`
	Parameters  Tool_Parameters `json:"parameters"`
}

type Tool_Parameters struct{
	Type       string               `json:"type"`
	Properties map[string]Property  `json:"properties"`
	Required   []string             `json:"required"`
}

type Property struct{
	Type string `json:"type"`
}

type Response_Create struct {
	Type string 	`json:"type"`
	Response Response   `json:"response"`
}

type Response struct {
	Modalities   []string `json:"modalities"`
}
