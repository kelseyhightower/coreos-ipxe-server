// Copyright 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	config()

	// Register images static file server.
	http.Handle("/images/", http.StripPrefix("/images/",
		http.FileServer(http.Dir(filepath.Join(dataDir, "images")))))

	// Register cloud configs static file server.
	http.Handle("/configs/", http.StripPrefix("/configs/",
		http.FileServer(http.Dir(filepath.Join(dataDir, "configs")))))

	// Register the iPXE boot script server.
	http.HandleFunc("/", ipxeBootScriptServer)

	// Start the iPXE Boot Server.
	fmt.Println("Starting CoreOS iPXE Server...")
	fmt.Printf("Listening on %s\n", hostPort)
	fmt.Printf("Advertised URL %s\n", baseUrl)
	fmt.Printf("Data directory: %s\n", dataDir)
	log.Fatal(http.ListenAndServe(hostPort, nil))
}
