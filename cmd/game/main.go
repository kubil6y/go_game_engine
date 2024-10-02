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

	// type a struct{}
	// type b struct{}
	// registry := ecs.NewRegistry()
	// registry.AddComponent(a{})
	// registry.AddComponent(a{})
	// registry.AddComponent(a{})
	// registry.AddComponent(a{})
	// registry.AddComponent(b{})
}
