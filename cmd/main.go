package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/billychappell/regselect"
)

var (
	filename  = flag.String("config", "config.json", "json config file to use for editing the registry")
	overwrite = flag.Bool("saveprev", false, "whether to overwrite config file to include previous values")
)

func confirm(filename *string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(`\n WARNING: Changing registry values can break your computer. \n Are you sure you want to proceed? (yes/no)
		\n Config File: %s`, filename)
	answer, _ := reader.ReadString('\n')
	switch answer {
	case "yes", "y", "Y", "YES":
		return true
	case "n", "no", "N", "NO":
		return false
	default:
		fmt.Println("can't understand answer '%s', please answer 'yes' or 'no'.")
		confirm(filename)
	}
	return false
}

func main() {
	flag.Parse()
	cfg, err := regselect.Unmarshal(*filename)
	if err != nil {
		log.Fatalf("error unmarshaling json config file: \n %s \n", err)
	}

	check := confirm(filename)
	if check == false {
		err = fmt.Errorf("failed to confirm changes with user")
		panic(err)
	}

	err = cfg.Set()
	if err != nil {
		log.Fatal(err)
	}

	if *overwrite == true {
		err = cfg.Write(*filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	update := fmt.Sprintf("%s.%s.save", filename, time.Now().String())
	err = cfg.Write(update)
	if err != nil {
		log.Fatalf("can't write updated file, fatal error: %s", err)
	}

	fmt.Println("Finished! Write successful.")
}
