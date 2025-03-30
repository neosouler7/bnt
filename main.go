package main

import (
	"bnt/bithumb"
	"bnt/commons"
	"bnt/config"
	"bnt/tgmanager"
	"bnt/upbit"
	"flag"
	"fmt"
	"os"
	"sync"
)

func usage() {
	fmt.Print("Welcome to bnt\n\n")
	fmt.Print("Please use the following commands\n\n")
	fmt.Print("-e : Set exchange code to run (or all)\n")
	os.Exit(0)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		usage()
	}

	exchange := flag.String("e", "", "Set exchange code to run")
	flag.Parse()

	tgConfig := config.GetTg()
	tgmanager.InitBot(
		tgConfig.Token,
		tgConfig.Chat_ids,
		commons.SetTimeZone("Tg"),
	)

	tgMsg := fmt.Sprintf("## START bnt %s %s", config.GetName(), *exchange)

	switch *exchange {
	case "upb":
		tgmanager.SendMsg(tgMsg)
		upbit.Run(*exchange)
	case "bmb":
		tgmanager.SendMsg(tgMsg)
		bithumb.Run(*exchange)
	case "all":
		tgmanager.SendMsg("## START bnt all")
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			upbit.Run("upb")
		}()

		go func() {
			defer wg.Done()
			bithumb.Run("bmb")
		}()

		wg.Wait()
	default:
		usage()
	}
}
