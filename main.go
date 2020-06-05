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

func main() {
	printBuilding()
}
