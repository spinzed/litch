package main

import "fmt"

func main() {
	data, err := readyAllData()

	if err != nil {
		fmt.Println(err)
		return
	}
	app := newApp(data)

	app.Run()
}
