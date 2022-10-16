package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go jsonUpdateScheduler()

	r := gin.Default()
	r.LoadHTMLGlob("templates/**")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/status", func(ctx *gin.Context) {

		jsonFile, err := os.Open("status.json")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		read := Data{}
		if err := json.Unmarshal(byteValue, &read); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, read)
	})

	if err := r.Run(":5000"); err != nil {
		panic(err)
	}
}

func jsonUpdateScheduler() error {
	for {
		time.Sleep(time.Second * 15)
		if err := writeStatusToJSON(); err != nil {
			return err
		}
	}
}

func writeStatusToJSON() error {
	min := 1
	max := 100
	input := Data{}
	input.Status.Wind = randomRange(min, max)
	input.Status.Water = randomRange(min, max)

	file, err := json.MarshalIndent(input, "", "")
	if err != nil {
		return err
	}
	if err := os.WriteFile("status.json", file, 0644); err != nil {
		return err
	}
	return nil
}

func randomRange(min, max int) int {
	return rand.Intn(max-min) + min
}
