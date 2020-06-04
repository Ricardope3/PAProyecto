package main

var (
	numberOfPeople int = 10
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
