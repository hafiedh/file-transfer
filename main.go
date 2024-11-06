package main

import (
	"fmt"
	"time"

	"hafiedh.com/downloader/cmd"
	"hafiedh.com/downloader/internal/config"
)

const banner = `
Ⱨ₳₣łɆĐⱧ
`

func main() {
	if tz := config.GetString("tz"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			fmt.Printf("error loading location '%s': %v\n", tz, err)
		} else {
			fmt.Printf("location loaded '%s'\n", tz)
		}
	}

	fmt.Print(banner)
	cmd.Run()
}
