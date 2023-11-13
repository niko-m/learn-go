package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}

func main() {
	data := []byte(`
		{
			"name": "do not call me that",
			"email": "hide@my.email",
			"age": 9123
		}
	`)

	var user User

	err := json.Unmarshal(data, &user)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(user)
}
