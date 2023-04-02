package serve

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"xsolla/libs/logkey"
)

// HTTP starts HTTP server on addr.
// It runs until failed or ctx.Done.
func HTTP(host string, port uint16, handler http.Handler) func(context.Context) error {
	return func(ctx context.Context) error {
		srv := &http.Server{
			Addr:              net.JoinHostPort(host, fmt.Sprintf("%d", port)),
			Handler:           handler,
			ReadHeaderTimeout: time.Minute,
		}

		errc := make(chan error, 1)
		go func() { errc <- srv.ListenAndServe() }()
		log.Println("started", logkey.Host, host, logkey.Port, port)
		defer log.Println("shutdown")

		var err error
		select {
		case err = <-errc:
		case <-ctx.Done():
			err = srv.Shutdown(context.Background())
		}
		if err != nil {
			return fmt.Errorf("srv.ListenAndServe: %w", err)
		}

		return nil
	}
}