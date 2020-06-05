package main

import (
	"fmt"
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

func printBuilding(){
	for _,row := range building {
		for _,col := range row{
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
}
