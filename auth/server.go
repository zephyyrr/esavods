package main

import (
    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/tylerb/graceful"
    "os"
    "path"
    "time"
    "errors"
    "net/http"
    "strconv"
)

type AuthKey []byte

func (key AuthKey) String() string {
    return string(key)
}

type KeyData struct {
    AuthKey
    Authenticater string
    Valid bool
    Expires time.Time
}

var (
    db *bolt.DB
)

func main () {
    var err error
    db, err = bolt.Open(path.Join(DBFolder, "AuthKeys.bolt"), 0600, nil)
    if err != nil {
        os.Exit(2)
    }
    server := echo.New()

    server.Use(echo.WrapMiddleware(Protected))

    server.Get("/key/:key", CheckKey)
    server.Post("/authenticate", Authenticate)

    s := standard.New(":" + strconv.Itoa(Port))
	s.SetHandler(server)
    graceful.ListenAndServe(s.Server, 3*time.Second)
}

func CheckKey(c echo.Context) error {
    // Fetch key from DB
    key := c.Param("key")
    return db.View(func (tx *bolt.Tx) (err error) {
        metadata := KeyData{
            AuthKey: []byte(key),
        }
        b := tx.Bucket([]byte(key))
        metadata.Authenticater = string(b.Get([]byte("Authenticater")))
        metadata.Expires, err = time.Parse(time.RFC3339, string(b.Get( []byte("Expires"))))
        if err != nil {
            return err
        }
        // Set validity on key
        metadata.Valid = metadata.Expires.Before(time.Now())

        // Return key to requester.
        return c.JSON(200, metadata)
    })
}

func Authenticate(c echo.Context) error {
    // Generate new key
    // Store in DB
    // Return to requester
    return c.JSON(http.StatusNotImplemented, "ab2458fecb221")
}

//Configuration
var (
	DebugMode    bool   //Activate various debugging features.
	Secret       string //Shared secret between the cluster.
	Port         int    //Port to listen on.
	DBFolder     string //Folder with database files.
)

// Reads the configuration of the server.
// Currently uses environment variables.
func ReadConfiguration() {
    var err error
	DebugMode = os.Getenv("AUTH_DEBUG") == "true"
    Secret = os.Getenv("ESAVODS_SECRET")
	DBFolder = os.Getenv("AUTH_DB_FOLDER")
    Port, err = strconv.Atoi(os.Getenv("AUTH_PORT"))
    if err != nil {
        os.Exit(1)
    }
}

var NotAuthorizedError = errors.New("Not Authorized.")
func Protected(c echo.Context) error {
    if c.Request().Header().Get("ESAVods-Shared-Key") != Secret {
        c.Error(NotAuthorizedError)
        return c.JSON(http.StatusUnauthorized, "The Eastmost Penninsula holds the secret.")
    }
    return nil
}
