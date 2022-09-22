package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"kwd/kernel/app"
	"os"
)

func InitSnowflake() {

	var err error = nil

	app.Snowflake, err = snowflake.NewNode(app.Cfg.Server.Node)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return
}
