package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var Version string

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func setupRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/exit", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Exiting the Program",
		})
		time.Sleep(2 * time.Second)
		os.Exit(0)
	})
	router.POST("/login", LoginHandler)
	router.POST("/adddl", AddDownloadHandler)
	router.POST("/canceldl", CancelDownloadHandler)
	router.GET("/dlinfo/:gid", GetDownloadInfoHandler)
}

func setupLoggingToFile(logfile string) {
	log.Println("Setting logToFile: ", logfile)
	os.Remove(logfile)
	handle, err := GetLogFileHandle(logfile)
	if err != nil {
		log.Println("Cannot open log file: ", err.Error())
	} else {
		log.SetOutput(io.MultiWriter(os.Stdout, handle))
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, handle)
}

func callback(c *cli.Context) error {
	ip := c.String("ip")
	port := c.String("port")
	apikey := c.String("apikey")
	logfile := c.String("logfile")
	uds := c.String("uds")
	if logfile != "" {
		setupLoggingToFile(logfile)
	}
	if apikey == "" {
		return fmt.Errorf("No mega.nz api key provided, exiting.")
	}
	megaClient = getMegaClient(apikey)
	debug := c.Bool("debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	setupRoutes(r)
	if uds != "" {
		if fileExists(fmt.Sprintf("%s", uds)) {
			log.Printf("Unix Socket/File unix:/%s already exists, Trying to remove it\n", uds)
			e := os.Remove(fmt.Sprintf("%s", uds))
			if e != nil {
				log.Fatal(e)
			}
		}
		log.Printf("Listening and serving HTTP on unix:/%s\n", uds)
		return r.RunUnix(fmt.Sprintf("%s", uds))
	} else {
		log.Printf("Serving on TCP %s:%s\n", ip, port)
		return r.Run(fmt.Sprintf("%s:%s", ip, port))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "MegaSDK-REST"
	app.Usage = "A web server encapsulating the downloading functionality of megasdk written in Go."
	app.Authors = []*cli.Author{
		{Name: "JaskaranSM"},
	}
	app.Action = callback
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: "6090",
			Usage: "port to listen on",
		},
		&cli.StringFlag{
			Name:  "ip",
			Value: "",
			Usage: "ip to listen on, by default webserver listens on localhost",
		},
		&cli.StringFlag{
			Name:  "apikey",
			Value: "",
			Usage: "API Key for MegaSDK. (mandatory)",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "run webserver in debug mode",
		},
		&cli.StringFlag{
			Name:  "logfile",
			Value: "",
			Usage: "log to file provided",
		},
		&cli.StringFlag{
			Name: "uds",
			Value: "",
			Usage: "To start the server on a Unix Domain Socket",
		},
	}
	if Version == "" {
		app.Version = "0.1.0"
	} else {
		app.Version = Version
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
