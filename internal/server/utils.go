package server

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gosimple/slug"
	"github.com/mubeng/mubeng/internal/metrics"
)

// Stop stops the server and all gateways (if any).
func Stop(ctx context.Context) {
	if len(handler.Gateways) > 0 {
		handler.mu.Lock()
		defer handler.mu.Unlock()

		for _, gateway := range handler.Gateways {
			_ = gateway.Close(ctx)
		}
	}

	if metricsServer != nil {
		_ = metricsServer.Shutdown(ctx)
	}

	_ = server.Shutdown(ctx)
}

func interrupt(sig chan os.Signal) {
	<-sig
	log.Warn("Interrupted. Exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	Stop(ctx)
}

func watch(w *fsnotify.Watcher) {
	for {
		select {
		case event := <-w.Events:
			if event.Op == 2 {
				log.Info("Proxy file has changed, reloading...")

				err := handler.Options.ProxyManager.Reload()
				if err != nil {
					log.Fatal(err)
				}

				if metricsEnabled {
					metrics.ProxyPoolSize.Set(float64(handler.Options.ProxyManager.Count()))
				}
			}
		case err := <-w.Errors:
			log.Fatal(err)
		}
	}
}

func getGatewayKey(baseURL, region string) string {
	return slug.Make(baseURL + "-" + region)
}
