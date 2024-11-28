package code

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
	"shared/interfaces"
	"strings"
)

func GetPlayerConstructor(code string, constructor string) (func(int) interfaces.Player, error) {
	if err := os.MkdirAll("./players", 0755); err != nil {
		return nil, fmt.Errorf("failed to create players directory: %w", err)
	}

	moduleDir := "./players/plugin"
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	modFile := fmt.Sprintf(`module player_plugin
		go 1.22.5
		require shared v0.0.0
		replace shared => ../../../shared
	`)
	if err := os.WriteFile(moduleDir+"/go.mod", []byte(modFile), 0644); err != nil {
		return nil, fmt.Errorf("failed to write go.mod: %w", err)
	}

	filename := moduleDir + "/main.go"
	code = strings.ReplaceAll(code, "package players", "package main")
	if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write plugin code: %w", err)
	}
	defer os.RemoveAll("./players") // Clean up the entire directory

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", moduleDir+"/impl.so", filename)
	// cmd.Dir = moduleDir // Set working directory to the module directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return nil, err
	}
	plug, err := plugin.Open(moduleDir + "/impl.so")
	if err != nil {
		fmt.Printf("error plugin: %s\n", err)
		return nil, err
	}
	defer os.Remove(moduleDir + "/impl.so")
	playerImpl, err := plug.Lookup(constructor)
	if err != nil {
		fmt.Printf("error in lookup: %s\n", err)
		return nil, err
	}

	cmd.Dir = ""

	return playerImpl.(func(int) interfaces.Player), nil
}

func GetGameConstructor(code string, constructor string) (func([]interfaces.Player) interfaces.Game, error) {
	if err := os.MkdirAll("./games", 0755); err != nil {
		return nil, fmt.Errorf("failed to create players directory: %w", err)
	}

	moduleDir := "./games/plugin"
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	modFile := fmt.Sprintf(`module games_plugin
		go 1.22.5
		require shared v0.0.0
		replace shared => ../../../shared
	`)
	if err := os.WriteFile(moduleDir+"/go.mod", []byte(modFile), 0644); err != nil {
		return nil, fmt.Errorf("failed to write go.mod: %w", err)
	}

	filename := moduleDir + "/main.go"
	code = strings.ReplaceAll(code, "package games", "package main")
	if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write plugin code: %w", err)
	}
	defer os.RemoveAll("./games")

	// TODO: add unique identifier to players uploaded
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", moduleDir+"/impl.so", filename)
	// cmd.Dir = moduleDir // Set working directory to the module directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return nil, err
	}
	plug, err := plugin.Open(moduleDir + "/impl.so")
	if err != nil {
		fmt.Printf("error plugin: %s\n", err)
		return nil, err
	}
	defer os.Remove(moduleDir + "/impl.so")
	gameImpl, err := plug.Lookup(constructor)
	if err != nil {
		fmt.Printf("error in lookup: %s\n", err)
		return nil, err
	}

	cmd.Dir = ""

	return gameImpl.(func([]interfaces.Player) interfaces.Game), nil
}
