package main

import "github.com/lobshunter/dorisctl/pkg/log"

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
