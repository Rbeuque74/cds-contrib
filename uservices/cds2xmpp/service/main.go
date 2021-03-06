package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VERSION ...
const VERSION = "0.1.0"

var mainCmd = &cobra.Command{
	Use: "cds2xmpp",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("cds")
		viper.AutomaticEnv()

		switch viper.GetString("log_level") {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "error":
			log.SetLevel(log.WarnLevel)
			gin.SetMode(gin.ReleaseMode)
		default:
			log.SetLevel(log.DebugLevel)
		}

		router := gin.New()
		router.Use(gin.Recovery())

		router.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE",
			RequestHeaders:  "Origin, Authorization, Content-Type, Accept",
			MaxAge:          50 * time.Second,
			Credentials:     true,
			ValidateHeaders: false,
		}))

		router.GET("/mon/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"version": VERSION})
		})

		s := &http.Server{
			Addr:           ":" + viper.GetString("listen_port"),
			Handler:        router,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Infof("Running cds2xmpp on %s", viper.GetString("listen_port"))

		if err := born(); err != nil {
			log.Fatalf("Error while initialize cds bot: %s", err)
		}
		go do()

		helloWorld()

		if err := s.ListenAndServe(); err != nil {
			log.Errorf("Error while running ListenAndServe: %s", err.Error())
		}
	},
}

func init() {
	flags := mainCmd.Flags()

	flags.String("log-level", "", "Log Level : debug, info or warn")
	viper.BindPFlag("log_level", flags.Lookup("log-level"))

	flags.String("listen-port", "8085", "Listen Port")
	viper.BindPFlag("listen_port", flags.Lookup("listen-port"))

	flags.String("event-kafka-broker-addresses", "", "Ex: --event-kafka-broker-addresses=host:port,host2:port2")
	viper.BindPFlag("event_kafka_broker_addresses", flags.Lookup("event-kafka-broker-addresses"))

	flags.String("event-kafka-topic", "", "Ex: --kafka-topic=your-kafka-topic")
	viper.BindPFlag("event_kafka_topic", flags.Lookup("event-kafka-topic"))

	flags.String("event-kafka-user", "", "Ex: --kafka-user=your-kafka-user")
	viper.BindPFlag("event_kafka_user", flags.Lookup("event-kafka-user"))

	flags.String("event-kafka-password", "", "Ex: --kafka-password=your-kafka-password")
	viper.BindPFlag("event_kafka_password", flags.Lookup("event-kafka-password"))

	flags.String("event-kafka-group", "", "Ex: --kafka-group=your-kafka-group")
	viper.BindPFlag("event_kafka_group", flags.Lookup("event-kafka-group"))

	flags.String("xmpp-server", "", "XMPP Server")
	viper.BindPFlag("xmpp_server", flags.Lookup("xmpp-server"))

	flags.String("xmpp-bot-jid", "cds@localhost", "XMPP Bot JID")
	viper.BindPFlag("xmpp_bot_jid", flags.Lookup("xmpp-bot-jid"))

	flags.String("xmpp-bot-password", "", "XMPP Bot Password")
	viper.BindPFlag("xmpp_bot_password", flags.Lookup("xmpp-bot-password"))

	flags.String("xmpp-hello-world", "", "Sending Hello World message to this jabber id")
	viper.BindPFlag("xmpp_hello_world", flags.Lookup("xmpp-hello-world"))

	flags.Bool("xmpp-debug", false, "XMPP Debug")
	viper.BindPFlag("xmpp_debug", flags.Lookup("xmpp-debug"))

	flags.Bool("xmpp-notls", true, "XMPP No TLS")
	viper.BindPFlag("xmpp_notls", flags.Lookup("xmpp-notls"))

	flags.Bool("xmpp-starttls", false, "XMPP Start TLS")
	viper.BindPFlag("xmpp_starttls", flags.Lookup("xmpp-starttls"))

	flags.Bool("xmpp-session", true, "XMPP Session")
	viper.BindPFlag("xmpp_session", flags.Lookup("xmpp-session"))

	flags.Bool("xmpp-insecure-skip-verify", true, "XMPP InsecureSkipVerify")
	viper.BindPFlag("xmpp_insecure_skip_verify", flags.Lookup("xmpp-insecure-skip-verify"))

	flags.String("xmpp-default-hostname", "", "Default Hostname for user, enter your.jabber.net for @your.jabber.net")
	viper.BindPFlag("xmpp_default_hostname", flags.Lookup("xmpp-default-hostname"))
}

func main() {
	mainCmd.Execute()
}
