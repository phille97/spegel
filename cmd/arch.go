package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/phille97/spegel/discovery"
	"github.com/phille97/spegel/proxy"
)

// archCmd represents the arch command
var archCmd = &cobra.Command{
	Use:   "arch",
	Short: "pacman cache mirror",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listenPort, err := cmd.Flags().GetUint16("port")
		if err != nil {
			log.Fatal(err)
		}

		cacheDir, err := cmd.Flags().GetString("cache")
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		service := "_arch-spegel-service._tcp"

		handler := proxy.NewProxy()

		r := mux.NewRouter()
		r.HandleFunc("/proxy/{path:.*}", handler.HandleProxy)
		r.PathPrefix("/get/").Handler(http.StripPrefix("/get/", http.FileServer(http.Dir(cacheDir))))

		loggedRouter := handlers.LoggingHandler(os.Stdout, r)

		go func() {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), loggedRouter); err != http.ErrServerClosed {
				log.Fatalf("HTTP server ListenAndServe: %v", err)
			}
		}()

		server, err := discovery.NewServer(service, listenPort)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			if err := server.Register(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		for {
			nodes, err := discovery.Discover(service)
			if err != nil {
				log.Fatal(err)
			}
			handler.Update(*nodes)

			log.Printf("Found %d nodes", len(handler.Nodes))

			tc := time.After(60 * time.Second)
			select {
			case <-sig:
				os.Exit(0)
			case <-tc:
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(archCmd)

	archCmd.Flags().Uint16P("port", "p", 5151, "port to listen")
	archCmd.Flags().StringP("cache", "c", "/var/cache/pacman/pkg", "pacman cache directory")
}
