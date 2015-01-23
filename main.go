// Copyright 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kelseyhightower/coreos-ipxe-server/config"
)

func main() {
	for _, s := range []string{"/images/", "/configs/", "/profiles/"} {
		// Register static file servers.
		http.Handle(s, http.StripPrefix(s,
			http.FileServer(http.Dir(filepath.Join(config.DataDir, s)))))
	}

	// Register the sshkey script server.
	http.HandleFunc("/keys", sshKeyServer)

	// Register the iPXE boot script server.
	http.HandleFunc("/", ipxeBootScriptServer)

	// Start the iPXE Boot Server.
	fmt.Println("Starting CoreOS iPXE Server...")
	fmt.Printf("Listening on %s\n", config.ListenAddr)

	if config.BaseUrl != "" {
		fmt.Printf("Advertised URL %s\n", config.BaseUrl)
	}
	fmt.Printf("Data directory: %s\n", config.DataDir)

	log.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}
