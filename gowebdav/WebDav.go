package gowebdav

import (
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
)

var dirFlag, sqlAddress, sqlUsername, sqlPassword, sqlDatabase string
var httpPort, httpsPort, sqlPort int
var serveSecure, authDigest bool

func init() {
	dirFlag = *flag.String("root_dir", "./media", "Directory to server from. Default is media.")
	httpPort = *flag.Int("port_http", 80, "Port to server HTTP.")
	httpsPort = *flag.Int("port_https", 443, "Port to server HTTPS.")
	serveSecure = *flag.Bool("https_only", false, "Server HTTPS. Default false.")
	authDigest = *flag.Bool("auth_digest", false, "Digest Authentication. Default Basic.")

	sqlAddress = *flag.String("sql_address", "127.0.0.1", "SQL-Server address.")
	sqlPort = *flag.Int("sql_port", 3306, "SQL-Server port.")
	sqlUsername = *flag.String("sql_username", "root", "SQL-Server username for this application.")
	sqlPassword = *flag.String("sql_password", "", "SQL-Server user-password for this application.")
	sqlDatabase = *flag.String("sql_database", "gowebdav", "SQL-Server database for this application.")

	flag.Parse()
}

func Execute() {
	srv := &webdav.Handler{
		FileSystem: webdav.Dir(dirFlag),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}

	sqlAddress = "192.168.56.1"
	sqlAddress = "127.0.0.1"
	sqlPassword = "my-secret-pw"
	sqlPassword = ""

	sqlServer, err := setupDatabase()
	if err != nil {
		println("FEHLER")
	}

	sqlServer.getUserPassword("coolimc")


	http.Handle("/", srv)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}


}