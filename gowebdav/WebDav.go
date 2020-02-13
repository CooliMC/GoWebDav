package gowebdav

import (
	"flag"
	"fmt"
	"github.com/abbot/go-http-auth"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
)

var dirFlag, sqlAddress, sqlUsername, sqlPassword, sqlDatabase string
var httpPort, httpsPort, sqlPort int
var httpEnabled, httpsEnabled, authEnabled, authDigest bool

func init() {
	flag.StringVar(&dirFlag, "root_dir", "./media", "Directory to server from. Default is media.")
	flag.IntVar(&httpPort, "port_http", 80, "Port to server HTTP.")
	flag.IntVar(&httpsPort, "port_https", 443, "Port to server HTTPS.")
	flag.BoolVar(&httpEnabled, "http_enabled", true, "Server HTTP. Default true.")
	flag.BoolVar(&httpsEnabled, "https_enabled", false, "Server HTTPS. Default false.")
	flag.BoolVar(&authEnabled, "auth_enabled", true, "Authentication enabled. Default true.")
	flag.BoolVar(&authDigest, "auth_digest", false, "Digest Authentication. Default Basic.")

	flag.StringVar(&sqlAddress, "sql_address", "127.0.0.1", "SQL-Server address.")
	flag.IntVar(&sqlPort, "sql_port", 3306, "SQL-Server port.")
	flag.StringVar(&sqlUsername, "sql_username", "root", "SQL-Server username for this application.")
	flag.StringVar(&sqlPassword, "sql_password", "", "SQL-Server user-password for this application.")
	flag.StringVar(&sqlDatabase, "sql_database", "gowebdav", "SQL-Server database for this application.")

	flag.Parse()
}

//Main WebDAVServer Execute() function
func Execute() {
	dirFlag = "D://MediaTest//"
	sqlAddress = "192.168.2.200"
	sqlPort = 43306
	sqlPassword = "my-secret-pw"

	srv := &webdav.Handler{
		FileSystem: DynamicFileSystem{webdav.Dir(dirFlag)},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}

	sqlServer, err := setupDatabase()
	if err != nil {
		log.Printf("SQL-Database Error: %s", err.Error())
		return
	}

	if authEnabled {
		if authDigest {
			//Add the authentication "Handler" that checks the user credentials
			authenticator := getDigestAuth(sqlServer)

			//Add the above created authentication handler as a pre-WebDAV-authentication-layer
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if username, authinfo := authenticator.CheckAuth(r); username == "" {
					authenticator.RequireAuth(w, r)
				} else {
					rN := r.WithContext(context.WithValue(r.Context(), "username", username))
					ar := auth.AuthenticatedRequest{Request: *rN, Username: username}
					ar.Header.Set(auth.AuthUsernameHeader, ar.Username)

					if authinfo != nil {
						w.Header().Set(authenticator.Headers.V().AuthInfo, *authinfo)
					}
					srv.ServeHTTP(w, &ar.Request)
				}
			})
		} else {
			//Add the authentication "Handler" that checks the user credentials
			authenticator := getBasicAuth(sqlServer)

			//Add the above created authentication handler as a pre-WebDAV-authentication-layer
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if username := authenticator.CheckAuth(r); username == "" {
					authenticator.RequireAuth(w, r)
				} else {
					rN := r.WithContext(context.WithValue(r.Context(), "username", username))
					ar := &auth.AuthenticatedRequest{Request: *rN, Username: username}
					ar.Header.Set(auth.AuthUsernameHeader, ar.Username)

					srv.ServeHTTP(w, &ar.Request)
				}
			})
		}
	} else {
		http.HandleFunc("/", srv.ServeHTTP)
	}

	if httpsEnabled {
		if _, err := os.Stat("./cert.pem"); err != nil {
			fmt.Println("[x] No cert.pem in current directory. Please provide a valid cert")
			return
		}
		if _, er := os.Stat("./key.pem"); er != nil {
			fmt.Println("[x] No key.pem in current directory. Please provide a valid cert")
			return
		}

		if httpEnabled {
			go http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), "cert.pem", "key.pem", nil)
		} else if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), "cert.pem", "key.pem", nil); err != nil {
			log.Fatalf("Error with WebDAV server: %v", err)
		}
	}

	if httpEnabled {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatalf("Error with WebDAV server: %v", err)
		}
	}

}

func getBasicAuth(sqlServer *DatabaseConnection) *auth.BasicAuth {
	return &auth.BasicAuth{ Realm: "WebDAV", Secrets: func(user, realm string) string {
		if pwd, err := sqlServer.getUserPassword(user); err == nil {
			return pwd
		} else {
			return ""
		}
	}}
}

func getDigestAuth(sqlServer *DatabaseConnection) *auth.DigestAuth {
	return auth.NewDigestAuthenticator("WebDAV", func(user, realm string) string {
		if pwd, err := sqlServer.getUserPassword(user); err == nil {
			return pwd
		} else {
			return ""
		}
	})
}

/**
 * The DynamicFileSystem struct implements an overridden
 * FileSystem struct with context based username parsing.
 *
 * @author	CooliMC
 * @version	1.0
 * @since	2019-07-01
 */

type DynamicFileSystem struct{
	webdav.FileSystem
}

func (d DynamicFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Test MKDIR: | %s |", fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	//return webdav.Dir(d).Mkdir(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), perm)
	/*if strings.HasPrefix(name, "/C") {
		return d.FileSystem.Mkdir(ctx, fmt.Sprintf("/user/%s%s", "admin", name), perm)
	} else {
		return d.FileSystem.Mkdir(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), perm)
	}*/
	return d.FileSystem.Mkdir(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), perm)
}

func (d DynamicFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("Test OpenFile: | %s |", fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	//return webdav.Dir(d).OpenFile(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), flag, perm)
	/*if strings.HasPrefix(name, "/C") {
		return d.FileSystem.OpenFile(ctx, fmt.Sprintf("/user/%s%s", "admin", name), flag, perm)
	} else {
		return d.FileSystem.OpenFile(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), flag, perm)
	}*/
	return d.FileSystem.OpenFile(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name), flag, perm)
}

func (d DynamicFileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("Test RemoveAll: | %s |", fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	//return webdav.Dir(d).RemoveAll(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	/*if strings.HasPrefix(name, "/C") {
		if name == "/C" {
			_ = d.FileSystem.RemoveAll(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
		}
		return d.FileSystem.RemoveAll(ctx, fmt.Sprintf("/user/%s%s", "admin", name))
	} else {
		return d.FileSystem.RemoveAll(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	}*/
	return d.FileSystem.RemoveAll(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
}

func (d DynamicFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("Test Rename: | %s | %s |", fmt.Sprintf("/user/%s%s", ctx.Value("username"), oldName), fmt.Sprintf("/user/%s%s", ctx.Value("username"), newName))
	//return webdav.Dir(d).Rename(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), oldName), fmt.Sprintf("/user/%s%s", ctx.Value("username"), newName))
	/*if strings.HasPrefix(oldName, "/C") {
		_ = d.FileSystem.Rename(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), oldName), fmt.Sprintf("/user/%s%s", ctx.Value("username"), newName))
		return d.FileSystem.Rename(ctx, fmt.Sprintf("/user/%s%s", "admin", oldName), fmt.Sprintf("/user/%s%s", "admin", newName))
	} else {
		return d.FileSystem.Rename(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), oldName), fmt.Sprintf("/user/%s%s", ctx.Value("username"), newName))
	}*/
	return d.FileSystem.Rename(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), oldName), fmt.Sprintf("/user/%s%s", ctx.Value("username"), newName))
}

func (d DynamicFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Test Stat: | %s |", fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	//return webdav.Dir(d).Stat(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	/*if strings.HasPrefix(name, "/C") {
		return d.FileSystem.Stat(ctx, fmt.Sprintf("/user/%s%s", "admin", name))
	} else {
		return d.FileSystem.Stat(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
	}*/
	return d.FileSystem.Stat(ctx, fmt.Sprintf("/user/%s%s", ctx.Value("username"), name))
}

