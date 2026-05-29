package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	matrix := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}

	var requestedPattern1 string = "ABCCED"
	var requestedPattern2 string = "ABCD"
	var requestedPattern3 string = "SEE"
	var requestedPattern4 string = "SABA"

	println(
		requestedPattern1,
		requestedPattern2,
		requestedPattern3,
		requestedPattern4,
	)

	ws, err := NewWordSearch(matrix, 3, 4)
	if err != nil {
		fmt.Println(err)
	}
	results := ws.Search(requestedPattern1)
	println(results)
}
