package web

import (
	"context"
	"log"
	"net/http"
	"time"
	"xstation/router"
)

var (
	hs *http.Server
)

func Start(port uint16) {
	hs = router.NewServer(port)
}

func Stop() {
	if hs == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hs.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
}
