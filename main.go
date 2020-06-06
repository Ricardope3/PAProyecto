package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

type person struct {
	id     int
	speed  float32
	exited bool
}

/*
Matrix building representation:
	0 represents empty
	1 represents wall/obstacle
	2 represents a person -- static matrix will place them in initial position
	3 represents exit
*/

var (
	building = [12][12]int{
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
	numberOfPeople int
	minSpeed       float32 = 0.5
	maxSpeed       float32 = 1.5
)

func printBuilding() {
	for _, row := range building {
		for _, col := range row {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func getNumOfPeople() {
	for _, row := range building {
		for _, col := range row {
			if col == 2 {
				numberOfPeople++
			}
		}
	}
	fmt.Println("Number of people: ", numberOfPeople)
}

func createWindow() *pixelgl.Window {

	// Specify configuration window
	cfg := pixelgl.WindowConfig{
		Title:  "Eathquake Evacuation Simulator",
		Bounds: pixel.R(0, 0, 840, 840),
		VSync:  true,
	}

	// Create a new window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)
	return win
}

func drawFloor(win *pixelgl.Window) *imdraw.IMDraw {

	floor := imdraw.New(nil)

	floor.Color = colornames.Lightgray
	floor.Push(pixel.V(60, 60))
	floor.Push(pixel.V(780, 780))
	floor.Rectangle(0)

	var x = 60.0
	var y = 60.0

	for i := len(building) - 1; i >= 0; i-- {
		for _, col := range building[i] {
			if col == 1 {
				floor.Color = colornames.Gray
				floor.Push(pixel.V(x, y))
				floor.Push(pixel.V(x+60.0, y+60.0))
				floor.Rectangle(0)
			} else if col == 3 {
				floor.Color = colornames.Red
				floor.Push(pixel.V(x, y))
				floor.Push(pixel.V(x+60.0, y+60.0))
				floor.Rectangle(0)
			}
			x += 60.0
		}
		x = 60.0
		y += 60.0
	}

	floor.Draw(win)
	win.Update()

	return floor
}

func drawPeople(win *pixelgl.Window) *imdraw.IMDraw {

	people := imdraw.New(nil)
	people.Color = colornames.Limegreen

	var x = 90.0
	var y = 90.0

	for i := len(building) - 1; i >= 0; i-- {
		for _, col := range building[i] {
			if col == 2 {
				people.Push(pixel.V(x, y))
				people.Circle(20, 0)
			}
			x += 60.0
		}
		x = 90.0
		y += 60.0
	}

	people.Draw(win)
	win.Update()

	return people
}

func run() {

	win := createWindow()

	drawFloor(win)
	drawPeople(win)

	for !win.Closed() {

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
	getNumOfPeople()
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

	pixelgl.Run(run)

}
