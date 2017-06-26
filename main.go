package main

import (
	"flag"
	"fmt"

	"github.com/lunny/log"
	"github.com/lunny/tango"
	"github.com/tango-contrib/debug"
)

var (
	addr    = flag.String("addr", "127.0.0.1:9191", "Listen IP and Port of web service")
	path    = flag.String("path", "./data", "directory to store index")
	isDebug = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()

	var logger = log.Std
	if *isDebug {
		logger.SetOutputLevel(log.Ldebug)
	}

	if err := initDB(*path); err != nil {
		logger.Fatal(err)
	}

	t := tango.NewWithLog(logger)
	if *isDebug {
		fmt.Println("debug mode")
		t.Use(debug.Debug())
	}
	t.Use(tango.Logging())
	t.Use(tango.Recovery(*isDebug))

	t.Get("/status/:prefix", func(ctx *tango.Context) {
		total := indexStatus(ctx.Param("prefix"))
		ctx.ServeJSON(map[string]interface{}{
			"total": total,
		})
	})

	t.Put("/:prefix/:unit_id/:word/:id", func(ctx *tango.Context) {
		prefix := ctx.Param("prefix")
		word := ctx.Param("word")
		id := ctx.ParamInt64("id")
		unitID := ctx.ParamInt64("unit_id")
		err := addIndex(prefix, word, id, unitID)
		if err != nil {
			ctx.WriteHeader(500)
			ctx.ServeJSON(map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		ctx.ServeJSON(map[string]interface{}{
			"status": "ok",
		})
	})

	t.Delete("/:prefix/:unit_id/:word/:id", func(ctx *tango.Context) {
		prefix := ctx.Param("prefix")
		word := ctx.Param("word")
		id := ctx.ParamInt64("id")
		unitID := ctx.ParamInt64("unit_id")
		err := delIndex(prefix, word, id, unitID)
		if err != nil {
			ctx.WriteHeader(500)
			ctx.ServeJSON(map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		ctx.ServeJSON(map[string]interface{}{
			"status": "ok",
		})
	})

	t.Get("/:prefix/:unit_id/:word", func(ctx *tango.Context) {
		prefix := ctx.Param("prefix")
		kw := ctx.Param("word")
		unitID := ctx.ParamInt64("unit_id")
		limit := ctx.FormInt("limit", 20)
		results, err := search(prefix, kw, unitID, limit)
		if err != nil {
			ctx.WriteHeader(500)
			ctx.ServeJSON(map[string]interface{}{
				"err": err.Error(),
			})
			return
		}

		ctx.ServeJSON(map[string]interface{}{
			"results": results,
		})
	})

	t.Run(*addr)
}
