package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankush-dev-eng/lanternfish"
	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/detection"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
	"github.com/ankush-dev-eng/lanternfish/pkg/server"
	"github.com/ankush-dev-eng/lanternfish/pkg/simulator"
)

const banner = `
      ___       _   _ _____ _____ ____  _   _ 
     / _ \     | | | | ____|_   _/ ___|| | | |
    | | | |    | |_| |  _|   | | \___ \| |_| |
    | |_| |    |  _  | |___  | |  ___) |  _  |
     \___/     |_| |_|_____| |_| |____/|_| |_|
        \
         \    .---.      LANTERNFISH v1.0.0
          \  /     \     Cyber Deception & Intruder Tracing
      (O)---'  o    |    Oceanic Field-Instrument Platform
             \     /     
              '---'
`

func main() {
	portFlag := flag.String("port", "", "Port to bind HTTP server (overrides PORT env)")
	flag.Parse()

	port := os.Getenv("PORT")
	if *portFlag != "" {
		port = *portFlag
	}
	if port == "" {
		port = "8080"
	}

	fmt.Print(banner)
	fmt.Printf("[*] Initializing LANTERNFISH Deception Architecture...\n")

	// 1. Generate enterprise network graph
	graph := network.GenerateEnterpriseTopology()
	fmt.Printf("[+] Generated %d simulated enterprise hosts across %d subnets\n", len(graph.Nodes), len(graph.Subnets))

	// 2. Initialize seeded decoys
	decoys := decoy.NewRegistry()
	fmt.Printf("[+] Armed %d decoy lures (canary files, honeytokens, decoy hosts)\n", len(decoys.GetAll()))

	// 3. Telemetry event bus
	bus := eventbus.NewBus(2000)

	// 4. Detection Engine
	detector := detection.NewEngine(graph, decoys)

	// 5. Simulation Engine
	simEngine := simulator.NewEngine(bus, decoys, detector)

	// 6. Sub-filesystem for embedded static UI
	staticFS, err := fs.Sub(lanternfish.WebFS, "web")
	if err != nil {
		log.Fatalf("[-] Failed to create static sub-filesystem: %v", err)
	}

	// 7. Initialize HTTP Server
	srv := server.New(server.Config{
		Graph:     graph,
		Decoys:    decoys,
		Bus:       bus,
		Detector:  detector,
		SimEngine: simEngine,
		StaticFS:  staticFS,
	})

	addr := ":" + port
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv.Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("[+] Field-Instrument Server running at: http://localhost:%s\n", port)
		fmt.Printf("[+] Live Telemetry SSE stream:         http://localhost:%s/api/stream\n", port)
		fmt.Printf("[+] Forensics Report (Markdown):       http://localhost:%s/api/report.md\n", port)
		fmt.Printf("[+] Health Probe:                      http://localhost:%s/healthz\n", port)
		fmt.Printf("[*] Operator console ready. Press Ctrl+C to terminate.\n\n")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[-] Server failed: %v", err)
		}
	}()

	<-stopChan
	fmt.Printf("\n[*] Gracefully stopping LANTERNFISH deception grid...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("[-] Server shutdown warning: %v\n", err)
	}
	fmt.Printf("[+] Deception grid disarmed cleanly. Farewell.\n")
}
