package ws

import (
	"encoding/json"
	"fmt"
)

type Playground struct {
	Turn     string `json:"turn"`
	Username string `json:"username"`
	Row1     Row1   `json:"row1"`
	Row2     Row2   `json:"row2"`
	Row3     Row3   `json:"row3"`
}

type Row1 struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

type Row2 struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

type Row3 struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

func UpdatePlayground(p []byte, roomId string, hub *Hub) {
	var playground Playground

	err := json.Unmarshal(p, &playground)
	if err != nil {
		fmt.Println(err)
	}

	result, err := json.Marshal(playground)
	if err != nil {
		fmt.Println(err)
	}

	hub.Rooms[roomId].Playground = string(result)
}
