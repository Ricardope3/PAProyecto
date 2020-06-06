package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
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

func printBuilding() {
	for _, row := range building {
		for _, col := range row {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

var (
	numberOfPeople int = 5
)

type person struct {
	speed float32
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

func main() {

	//Slice containing all the trapped people inside the building
	//TODO: POPULATE THIS ARRAY
	trapped := make([]person, numberOfPeople)
	//Slice containing all the people that have left the building
	safe := make([]person, numberOfPeople)
	//Channel to handle when a person wants to move
	onMove := make(chan person)
	onExit := make(chan person)

	for {
		select {
		//Case when a person tryes to move from one place to another
		case person := <-onMove:
			//here one should modify the matrix to reflect the movement
			//the canvas should be redrawn
		//Case when a person left the building
		case person := <-onExit:
			//here the person that left the building should be erased from trapped slice
			//here the person that left the building should be added to the safe slice
			//the canvas should be redrawn

		}
	}

	pixelgl.Run(run)
}
