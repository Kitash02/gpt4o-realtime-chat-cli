package structs

//Represent the model response, used for printing the respons
type Response_Text_Delta struct{
	Type string `json:"type"`
	Delta string `json:"delta"`
}

type Response_Done struct {
	Type     string `json:"type"`
	Response struct {
		Output []struct {
			CallID    string `json:"call_id"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		} `json:"output"`
	} `json:"response"`
}
