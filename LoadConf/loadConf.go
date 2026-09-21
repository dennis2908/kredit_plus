package loadconf

import (
	"log"

	"github.com/joho/godotenv"
)

func Connects() { // init instead of int

	err := godotenv.Load("conf/app.conf")
	if err != nil {
		err1 := godotenv.Load("../conf/app.conf")
		if err1 != nil {
			err2 := godotenv.Load("../../conf/app.conf")
			if err2 != nil {
				log.Fatal("Error loading .conf file")
			}
		}
	}

}
