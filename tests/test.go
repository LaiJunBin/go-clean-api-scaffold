package main

import (
	"fmt"
	"os"
	"os/exec"

	"go-clean-api-scaffold/internal/config"
	"go-clean-api-scaffold/tests/setup"
)

func main() {
	args := os.Args[1:]
	goTestArgs := []string{"test", "-failfast", "-p", "1"}

	config := config.NewConfig()
	setup := setup.NewSetup(config)

	defer setup.TerminateInfrastructure()
	if err := setup.SetupInfrastructure(); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", append(goTestArgs, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
