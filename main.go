package main

import (
	"os"

	"github.com/lifeym/she/mail"
)

func main() {
	smtp := mail.New("87363255@qq.com", "password", "host", 587)
	msg := mail.NewMessage()
	smtp.Send(msg)
	os.Exit(0)
}
