package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Matrix building representation:
	0 represents empty
	1 represents wall/obstacle
	2 represents a person -- static matrix will place them in initial position
	3 represents exit
*/
var building = [12][12]int{
	//0 1 2 3 4 5 6 7 8 9 10 11
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, //0
	{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1}, //1
	{1, 2, 1, 1, 1, 1, 0, 1, 1, 1, 2, 1}, //2
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, //3
	{1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1}, //4
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, //5
	{1, 0, 1, 1, 0, 0, 0, 0, 2, 1, 0, 1}, //6
	{1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1}, //7
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, //8
	{1, 2, 1, 1, 1, 0, 0, 1, 1, 1, 2, 1}, //9
	{1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1}, //10
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, //11
}

var (
	numberOfPeople int     = 5
	minSpeed       float32 = 0.5
	maxSpeed       float32 = 1.5
)

type person struct {
	id     int
	speed  float32
	exited bool
}

func printBuilding() {
	for _, row := range building {
		for _, col := range row {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func generateRandomSpeed() float32 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return minSpeed + r1.Float32()*(maxSpeed-minSpeed)
}

func initiatePerson(p person, onMove, onExit chan person) {
	go func() {
		for {
			time.Sleep(time.Duration(p.speed) * time.Second)
			//MOVERTE
			//VALIDAR SI LLEGASTE A LA SALIDA
			if generateRandomSpeed() < 1 {
				onExit <- p
				return
			}
			onMove <- p
		}
	}()
}

func main() {

	trapped := make([]person, numberOfPeople)
	safe := make([]person, 0)

	onMove := make(chan person)
	onExit := make(chan person)

	for i := 0; i < numberOfPeople; i++ {
		trapped[i] = person{i, float32(i + 2), false}
		go initiatePerson(trapped[i], onMove, onExit)
	}

	for {
		select {
		case person := <-onMove:
			//REPINTAR CANVAS
			fmt.Println(person.id, "Me movi")
		case person := <-onExit:
			//REPINTAR CANVAS
			fmt.Println(person.id, "Me sali")
			safe = append(safe, person)
			if len(safe) >= numberOfPeople {
				close(onMove)
				return
			}
		}
	}

}
