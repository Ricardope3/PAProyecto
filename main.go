package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type person struct {
	id     int
	speed  float32
	exited bool
	// path   []coordinate
}

type coordinate struct {
	row int
	col int
}

/*
Matrix building representation:
	0 represents empty
	1 represents wall/obstacle
	2 represents a person -- static matrix will place them in initial position
	3 represents exit
*/

var (
	building = [][]int{
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
	numberOfExits  int
	exits          [5]coordinate
	past           [12][12][2]bool //0:whether it has been visited 1:whether it formas part of the path
	path           []coordinate
	timeout        float64 = 5
)

func initializePast() {
	for i, row := range past {
		for j, _ := range row {
			past[i][j][0] = false
			past[i][j][1] = true
		}
	}
}

func printBuilding() {
	for _, row := range building {
		for _, col := range row {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func printPast() {
	for _, row := range past {
		for _, col := range row {
			fmt.Print("[", col[0], ",", col[1], "]")
		}
		fmt.Println()
	}
}

func printPathMatrix() {
	for i, row := range building {
		for j, col := range row {
			if past[i][j][1] == true && past[i][j][0] == true {
				fmt.Print("Y")
			} else {
				fmt.Print(col)
			}

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

func insertExit(floor [][]int, sideExits []int, side, indexExit, lenF int) bool {
	switch side {
	case 0:
		if floor[0][indexExit] == 3 || floor[1][indexExit] == 1 {
			return false
		}
		floor[0][indexExit] = 3
		exits[numberOfExits] = coordinate{row: 0, col: indexExit}
		numberOfExits++
	case 1:
		if floor[indexExit][lenF-1] == 3 || floor[indexExit][lenF-2] == 1 {
			return false
		}
		floor[indexExit][lenF-1] = 3
		exits[numberOfExits] = coordinate{row: indexExit, col: lenF - 1}
		numberOfExits++
	case 2:
		if floor[lenF-1][indexExit] == 3 || floor[lenF-2][indexExit] == 1 {
			return false
		}
		floor[lenF-1][indexExit] = 3
		exits[numberOfExits] = coordinate{row: lenF - 1, col: indexExit}
		numberOfExits++
	case 3:
		if floor[indexExit][0] == 3 || floor[indexExit][1] == 1 {
			return false
		}
		floor[indexExit][0] = 3
		exits[numberOfExits] = coordinate{row: indexExit, col: 0}
		numberOfExits++
	default:
		sideExits[side]++
	}
	return true
}

func generateExits(floor [][]int) {
	sideExits := make([]int, 4)
	lenF := len(floor)
	rand.Seed(time.Now().UnixNano())
	nExits := rand.Intn(5) + 1
	for i := 0; i < nExits; {
		side := rand.Intn(4)
		indexExit := rand.Intn(lenF)
		if sideExits[side] > 1 {
			continue
		}
		valid := insertExit(floor, sideExits, side, indexExit, lenF)
		if valid {
			i++
		}
	}
	fmt.Println("Number of exits: ", numberOfExits)
}

func distance(a coordinate, b coordinate) float64 {
	res := math.Pow(float64(a.col)-float64(b.col), 2) + math.Pow(float64(a.row)-float64(b.row), 2)
	res = math.Sqrt(res)
	return res
}

func findClosestExit(position coordinate) int {
	targetExit := 0
	dist := distance(position, exits[0])
	for i := 1; i < numberOfExits; i++ {
		d := distance(position, exits[i])
		if d < dist {
			targetExit = i
			dist = d
		}
	}
	return targetExit
}

func validate(row int, col int) bool {
	if row < 0 || row >= len(building) {
		return false
	} //not a valid row
	if col < 0 || col >= len(building[row]) {
		return false
	} //not a valid column
	if building[row][col] == 1 {
		return false
	} //wall/obstacle
	if past[row][col][0] == true {
		return false
	} //visited
	return true
}

func searchPath(row int, col int) bool {
	position := coordinate{row, col}
	e := findClosestExit(position) //index of target exit
	initializePast()
	return searchPathRec(row, col, e)
}

func searchPathRec(row int, col int, e int) bool {
	past[row][col][0] = true                        //mark visited
	if row == exits[e].row && col == exits[e].col { //reached end
		path = append(path, coordinate{row, col})
		return true
	} else {
		if validate(row, col+1) {
			if searchPathRec(row, col+1, e) {
				path = append(path, coordinate{row, col})
				return true
			} //;
		} //available path at right
		if validate(row, col-1) {
			if searchPathRec(row, col-1, e) {
				path = append(path, coordinate{row, col})
				return true
			} //;
		} //available path at left
		if validate(row+1, col) {
			if searchPathRec(row+1, col, e) {
				path = append(path, coordinate{row, col})
				return true
			} //;
		} //available path down
		if validate(row-1, col) {
			if searchPathRec(row-1, col, e) {
				path = append(path, coordinate{row, col})
				return true
			} //;
		} //available path up
		past[row][col][1] = false
		return false //not the way
	}
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
	initializePast()
	generateExits(building)
	getNumOfPeople()
	trapped := make([]person, numberOfPeople)
	safe := make([]person, 0)

	onMove := make(chan person)
	onExit := make(chan person)
	start := time.Now()

	for i := 0; i < numberOfPeople; i++ {
		trapped[i] = person{i, float32(i + 2), false}
		go initiatePerson(trapped[i], onMove, onExit)
	}
	go func() {
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
			default:
				elapsed := time.Since(start)
				seconds := elapsed.Seconds()
				if seconds > timeout {
					win.Clear(colornames.White)
				}
			}
		}
	}()
	for !win.Closed() {
		win.Update()
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

	pixelgl.Run(run)
}
