package main

import (
	"fmt"
	"math/rand"
	"math"
	"time"
)

type person struct {
	speed float32
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
				{1,1,1,1,1,1,1,1,1,1,1,1},//0
				{1,0,0,0,1,0,0,1,0,0,0,1},//1
				{1,2,1,1,1,1,0,1,1,1,2,1},//2
				{1,0,0,0,0,0,0,0,0,0,0,1},//3
				{1,1,1,1,1,0,1,1,1,1,1,1},//4
				{1,0,0,0,0,0,0,0,0,0,0,1},//5
				{1,0,1,1,0,0,0,0,2,1,0,1},//6
				{1,1,1,1,1,0,0,1,1,1,1,1},//7
				{1,0,0,0,0,0,0,0,0,0,0,1},//8
				{1,2,1,1,1,0,0,1,1,1,2,1},//9
				{1,0,0,0,1,0,0,1,0,0,0,1},//10
				{1,1,1,1,1,1,1,1,1,1,1,1},//11
				}
	numberOfPeople int
	numberOfExits int
	exits [5]coordinate
	past [12][12][2]bool //0:whether it has been visited 1:whether it formas part of the path
)

func initializePast(){
	for i,row := range past{
		for j,_ := range row{
			past[i][j][0] = false
			past[i][j][1] = true
		}
	}
}

func printBuilding(){
	for _,row := range building {
		for _,col := range row{
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func printPast(){
	for _,row := range past {
		for _,col := range row{
			fmt.Print("[",col[0],",",col[1],"]")
		}
		fmt.Println()
	}
}

func printPath(){
	for i,row := range building {
		for j,col := range row{
			if past[i][j][1] == true && past[i][j][0] == true{
				fmt.Print("Y")
			}else{
				fmt.Print(col)
			}
			
		}
		fmt.Println()
	}
}

func getNumOfPeople(){
	for _,row := range building {
		for _,col := range row{
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
		if (floor[0][indexExit] == 3 || floor[1][indexExit] == 1){
			return false
		}
		floor[0][indexExit] = 3
		exits[numberOfExits] = coordinate{row:0 , col: indexExit}
		numberOfExits++
	case 1:
		if (floor[indexExit][lenF-1]==3 || floor[indexExit][lenF-2] == 1){
			return false
		}
		floor[indexExit][lenF-1] = 3
		exits[numberOfExits] = coordinate{row:indexExit , col: lenF-1}
		numberOfExits++
	case 2:
		if (floor[lenF-1][indexExit]==3 || floor[lenF-2][indexExit] == 1){
			return false
		}
		floor[lenF-1][indexExit] = 3
		exits[numberOfExits] = coordinate{row: lenF-1, col: indexExit}
		numberOfExits++
	case 3:
		if (floor[indexExit][0]==3 || floor[indexExit][1] == 1){
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

func generateExits(floor [][]int){
	sideExits := make([]int, 4)
	lenF := len(floor)
	rand.Seed(time.Now().UnixNano())
	nExits := rand.Intn(5) + 1
	for i:=0;i<nExits;{
		side := rand.Intn(4)
		indexExit := rand.Intn(lenF)
		if (sideExits[side]>1){
			continue
		}
		valid := insertExit(floor, sideExits, side, indexExit, lenF)
		if(valid){
			i++
		}
	}
	fmt.Println("Number of exits: ", numberOfExits)
}

func distance(a coordinate, b coordinate) float64{
	res := math.Pow(float64(a.col) - float64(b.col), 2) + math.Pow(float64(a.row) - float64(b.row), 2)
	res = math.Sqrt(res)
	return res
}

func findClosestExit(position coordinate) int{
	targetExit := 0
	dist := distance(position,exits[0]) 
	for i := 1; i < numberOfExits; i++ {
		d :=distance(position,exits[i])
		if d < dist {
			targetExit = i
			dist = d 
		}
	}
	return targetExit
}

func validate(row int, col int) bool{
	if row < 0 || row >= len(building) {
		return false;
	}//not a valid row
	if col < 0 || col >= len(building[row]) {
		return false;
	}//not a valid column
	if building[row][col] == 1 {
		return false;
	}//wall/obstacle
	if past[row][col][0] == true {
		return false;
	}//visited
	return true;
}

func searchPath(row int, col int) bool{
	position := coordinate{row,col}
	e := findClosestExit(position) //index of target exit
	initializePast()
	return searchPathRec(row,col,e)
}

func searchPathRec(row int, col int,e int) bool{
	past[row][col][0] = true;//mark visited
	if row == exits[e].row && col == exits[e].col { //reached end
		return true;
	}else{
		if validate(row,col+1) {
			if searchPathRec(row,col+1,e){
				return true
			}//;
		}//available path at right
		if validate(row,col-1) {
			if searchPathRec(row,col-1,e){
				return true
			}//;
		}//available path at left
		if validate(row+1,col) {
			if searchPathRec(row+1,col,e) {
				return true
			}//;
		}//available path down
		if validate(row-1,col) {
			if searchPathRec(row-1,col,e) {
				return true
			}//;
		}//available path up
		past[row][col][1] = false
		return false//not the way
	}
}

func main() {
	initializePast()
	printPast()
	generateExits(building)
	printBuilding()
	getNumOfPeople()
	searchPath(2,1)
	printPath()
	/*
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
	}*/
}
