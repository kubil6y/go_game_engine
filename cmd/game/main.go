package main

import "os"

func main() {
	game := NewGame()
	if err := game.Initialize(); err != nil {
        game.logger.Fatal(err, "something is terribly wrong", nil)
        os.Exit(1)
	}
	defer game.Destroy()
	game.Run()
}
