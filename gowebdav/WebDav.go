package gowebdav

import (
	"flag"
	"fmt"
	"github.com/abbot/go-http-auth"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
)

var dirFlag, sqlAddress, sqlUsername, sqlPassword, sqlDatabase string
var httpPort, httpsPort, sqlPort int
var serveSecure, authEnabled, authDigest bool

func init() {
	dirFlag = *flag.String("root_dir", "./media", "Directory to server from. Default is media.")
	httpPort = *flag.Int("port_http", 80, "Port to server HTTP.")
	httpsPort = *flag.Int("port_https", 443, "Port to server HTTPS.")
	serveSecure = *flag.Bool("https_only", false, "Server HTTPS. Default false.")
	authEnabled = *flag.Bool("auth_enabled", true, "Authentication enabled. Default true.")
	authDigest = *flag.Bool("auth_digest", false, "Digest Authentication. Default Basic.")

	sqlAddress = *flag.String("sql_address", "127.0.0.1", "SQL-Server address.")
	sqlPort = *flag.Int("sql_port", 3306, "SQL-Server port.")
	sqlUsername = *flag.String("sql_username", "root", "SQL-Server username for this application.")
	sqlPassword = *flag.String("sql_password", "", "SQL-Server user-password for this application.")
	sqlDatabase = *flag.String("sql_database", "gowebdav", "SQL-Server database for this application.")

	flag.Parse()
}

func Execute() {
	dirFlag = "D://MediaTest//"

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


	sqlAddress = "192.168.2.200"
	sqlPort = 43306
	sqlPassword = "my-secret-pw"

	sqlServer, err := setupDatabase()
	if err != nil {
		log.Printf("SQL-Database Error: %s", err.Error())
		return
	}

	if authEnabled {
		if authDigest {
			authenticator := auth.NewDigestAuthenticator("WebDAV", func(user, realm string) string {
				if pwd, err := sqlServer.getUserPassword(user); err == nil {
					return pwd
				} else {
					return ""
				}
			})

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if username, authinfo := authenticator.CheckAuth(r); username == "" {
					authenticator.RequireAuth(w, r)
				} else {
					ar := auth.AuthenticatedRequest{Request: *r, Username: username}
					ar.Header.Set(auth.AuthUsernameHeader, ar.Username)
					if authinfo != nil {
						w.Header().Set(authenticator.Headers.V().AuthInfo, *authinfo)
					}
					srv.ServeHTTP(w, &ar.Request)
				}
			})
		} else {
			authenticator := &auth.BasicAuth{Realm: "WebDAV", Secrets: func(user, realm string) string {
				if pwd, err := sqlServer.getUserPassword(user); err == nil {
					return pwd
				} else {
					return ""
				}
			}}

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if username := authenticator.CheckAuth(r); username == "" {
					authenticator.RequireAuth(w, r)
				} else {
					ar := &auth.AuthenticatedRequest{Request: *r, Username: username}
					ar.Header.Set(auth.AuthUsernameHeader, ar.Username)
					srv.ServeHTTP(w, &ar.Request)
				}
			})
		}
	} else {
		http.HandleFunc("/", srv.ServeHTTP)
	}


	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}


}