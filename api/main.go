package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/tylerb/graceful"
	"github.com/unrolled/render"
	. "github.com/zephyyrr/esavods/api/models"
	"math/big"
	"os"
	"path"
	"time"
)

var (
	server *echo.Echo
	r      *render.Render
	log    *logrus.Logger
)

func main() {
	ReadConfiguration()

	log = &logrus.Logger{
		Out:       logrus.StandardLogger().Out,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.ErrorLevel,
	}

	if DebugMode {
		log.Level = logrus.DebugLevel
	}

	encoder = Json{}
	var err error
	db, err = OpenDatabase()
	if err != nil {
		server.Logger().Fatalf("Could not open database for read/write. Quitting.")
	}
	defer db.Close()

	server = echo.New()
	r = render.New(render.Options{
		Delims:     render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset:    "UTF-8",                     // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                        // Output human readable JSON.
		IndentXML:  true,                        // Output human readable XML.
		//PrefixJSON:                []byte(")]}',\n"),                                // Prefixes JSON responses with the given bytes.
		PrefixXML:                 []byte("<?xml version='1.0' encoding='UTF-8'?>"), // Prefixes XML responses with the given bytes.
		UnEscapeHTML:              true,                                             // Replace ensure '&<>' are output correctly (JSON only).
		StreamingJSON:             true,                                             // Streams the JSON response via json.Encoder.
		RequirePartials:           true,                                             // Return an error if a template is missing a partial used in a layout.
		DisableHTTPErrorRendering: true,                                             // Disables automatic rendering of http.StatusInternalServerError when an error o
	})
	server.SetHTTPErrorHandler(ErrorHandler)
	server.Use(echo.WrapMiddleware(echologrus(log)))
	setupAPI()
	server.Static("/static", StaticFolder)
	if DebugMode {
		server.File("/debug", path.Join(StaticFolder, "debug.html"))
	}
	server.File("/favicon.ico", path.Join(StaticFolder, "static/favicon.ico"))

	s := standard.New(":" + os.Getenv("API_PORT"))
	s.SetHandler(server)

	timelimit := 3 * time.Second
	if Https {
		graceful.ListenAndServeTLS(s.Server, CertFile, PrivKeyFile, timelimit)
	} else {
		graceful.ListenAndServe(s.Server, timelimit)
	}
}

//Configuration
var (
	DebugMode    bool   //Activate various debugging features.
	StaticFolder string //Folder for serving static content
	Port         int    //Port to listen on.
	DBFolder     string //Folder with database files.

	Https       bool   //Use HTTPS instead.
	CertFile    string //Certificate to use with HTTPS
	PrivKeyFile string //Private Key File to use with HTTPS
)

// Reads the configuration of the server.
// Currently uses environment variables.
func ReadConfiguration() {
	DebugMode = os.Getenv("API_DEBUG") == "true"

	DBFolder = os.Getenv("API_DB_FOLDER")
	StaticFolder = os.Getenv("API_STATIC_FOLDER")
	if StaticFolder == "" {
		StaticFolder = "static/"
	}

	Https = os.Getenv("API_HTTPS") == "true"
	CertFile = os.Getenv("API_CERTFILE")
	PrivKeyFile = os.Getenv("API_KEYFILE")
}

//Installs all API functions.
func setupAPI() {
	server.Get("/events", GetEvents)
	server.Get("/event/:name", GetEvent)
	server.Post("/events", PostEvents)
	server.Get("/runs", GetRuns)
	server.Post("/runs", PostRuns)

}

//Error logging function.
func echologrus(l *logrus.Logger) func(echo.Context) error {
	return func(c echo.Context) error {
		l.WithFields(logrus.Fields{
			"remote": c.Request().RemoteAddress(),
			"method": c.Request().Method(),
			"url":    c.Request().URL(),
		}).Debug()
		return nil
	}
}

func ErrorHandler(err error, c echo.Context) {

	switch e := err.(type) {
	case Error:
		log.WithFields(logrus.Fields{
			"remote":      c.Request().RemoteAddress(),
			"method":      c.Request().Method(),
			"url":         c.Request().URL(),
			"input_data":  c.ParamValues(),
			"output_data": e.Data,
			"error":       e.Internal,
		}).Error(e.Message)
		c.JSON(e.HttpStatus, e)
	case *echo.HTTPError:
		log.Error(err)
		c.JSON(e.Code, e.Message)
	}
}

func NewID() Id {
	max := big.NewInt(1)
	max.Lsh(max, 128)
	src, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("Error reading randomness: " + err.Error())
	}
	return Id(base64.URLEncoding.EncodeToString(src.Bytes()))
}
