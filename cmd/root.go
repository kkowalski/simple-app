package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"simple/api"
	mongo "simple/db"
)

var rootCmd = &cobra.Command{
	Use:   "simple",
	Short: "simply simple",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() {
	ctx := context.Background()

	db, cleanup := mongo.Setup(ctx)
	defer cleanup()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := api.SetupServer(db).Run()
		if err != nil {
			panic(fmt.Errorf("error setuping web server: %v", err))
		}
		wg.Done()
	}()
	wg.Wait()
}
