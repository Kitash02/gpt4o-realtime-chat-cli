package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"wonderful/structs"
)

// multiplies two integers and return the result as string 
func Multiplier (args string) string{
	var numbers structs.Multiplier
	err := json.Unmarshal([]byte(args), &numbers)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return ""
	}
	result := numbers.First * numbers.Second

	return strconv.Itoa(result)
}

// Multiplier function tool, sended to the server after init session
func MultiplierTool() structs.Tool{
	return structs.Tool{
		Type: "function",
		Name: "Multiplier",
		Description: "Multiply two integers",
		Parameters: structs.Tool_Parameters{
			Type: "object",
			Properties: map[string]structs.Property{
				"first": {Type: "integer"},
				"second": {Type: "integer"},
			},
			Required: []string{"first", "second"},
			},
	}
}