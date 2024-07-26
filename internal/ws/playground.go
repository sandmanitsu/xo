package ws

import (
	"encoding/json"
	"fmt"
)

type Playground struct {
	Winner   string `json:"winner"`
	Turn     string `json:"turn"`
	Username string `json:"username"`
	UserTurn string `json:"userTurn"`
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
	var prevPlayground Playground

	err := json.Unmarshal(p, &playground)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(hub.Rooms[roomId].Playground), &prevPlayground)
	if err != nil {
		fmt.Println(err)
	}

	hasWinner, winner := checkWin(playground)
	if !hasWinner {
		playground = changeUserTurn(playground, prevPlayground, hub.Rooms[roomId].Clients)
	} else {
		playground.Winner = winner
	}

	result, err := json.Marshal(playground)
	if err != nil {
		fmt.Println(err)
	}

	hub.Rooms[roomId].Playground = string(result)
}

func changeUserTurn(playground Playground, prevPlayground Playground, clients map[string]*Client) Playground {
	if playground.UserTurn == "" {
		return playground
	}

	if playground.Turn != prevPlayground.Turn {
		for _, client := range clients {
			if playground.UserTurn != client.User {
				playground.UserTurn = client.User
			}
		}
	}

	return playground
}

func checkWin(p Playground) (bool, string) {
	fmt.Printf("\ncheck if win user: %s, with : %s", p.UserTurn, p.Turn)

	if p.Row1.A == p.Turn && p.Row1.B == p.Turn && p.Row1.C == p.Turn {
		return true, p.UserTurn
	}
	if p.Row2.A == p.Turn && p.Row2.B == p.Turn && p.Row2.C == p.Turn {
		return true, p.UserTurn
	}
	if p.Row3.A == p.Turn && p.Row3.B == p.Turn && p.Row3.C == p.Turn {
		return true, p.UserTurn
	}

	if p.Row1.A == p.Turn && p.Row2.A == p.Turn && p.Row3.A == p.Turn {
		return true, p.UserTurn
	}
	if p.Row1.B == p.Turn && p.Row2.B == p.Turn && p.Row3.B == p.Turn {
		return true, p.UserTurn
	}
	if p.Row1.C == p.Turn && p.Row2.C == p.Turn && p.Row3.C == p.Turn {
		return true, p.UserTurn
	}

	if p.Row1.A == p.Turn && p.Row2.B == p.Turn && p.Row3.C == p.Turn {
		return true, p.UserTurn
	}
	if p.Row1.C == p.Turn && p.Row2.B == p.Turn && p.Row3.A == p.Turn {
		return true, p.UserTurn
	}

	return false, ""
}
