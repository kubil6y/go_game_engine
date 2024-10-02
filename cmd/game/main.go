package main

import "log"

func main() {
    game := NewGame()
    if err := game.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer game.Destroy()
    game.Run()
}
