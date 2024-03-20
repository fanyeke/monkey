package input

import (
	"encoding/json"
	"os"
)

type ReadInput struct {
	Input []string `json:"input"`
}

var Ri ReadInput

func init() {
	file, err := os.ReadFile("input/input.json")
	if err != nil {
		panic(err)
	}
	Ri = ReadInput{}
	err = json.Unmarshal(file, &Ri)
	if err != nil {
		panic(err)
	}
}
