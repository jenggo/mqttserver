package main

import (
	"crypto/tls"
	"mqttserver/config"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	rds "github.com/go-redis/redis/v8"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}).With().Timestamp().Logger()
	if err := cleanenv.ReadConfig(config.FileConfig, &config.Load); err != nil {
		log.Fatal().Msg(err.Error())
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	authRules := &auth.Ledger{
		Auth: auth.AuthRules{
			{Username: auth.RString(config.Load.Mqtt.Username), Password: auth.RString(config.Load.Mqtt.Password), Allow: true},
		},
		ACL: auth.ACLRules{
			{
				Username: auth.RString(config.Load.Mqtt.Username),
				Filters: auth.Filters{
					"#": auth.ReadWrite,
				},
			},
			{
				Filters: auth.Filters{
					"#": auth.Deny,
				},
			},
		},
	}

	cert, err := getTls()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	server := mqtt.New(nil)
	server.Log = &log.Logger

	if err := server.AddHook(new(auth.Hook), &auth.Options{Ledger: authRules}); err != nil {
		log.Fatal().Msg(err.Error())
	}

	addr := net.JoinHostPort(config.Load.Redis.Host, config.Load.Redis.Port)

	if err := server.AddHook(new(redis.Hook), &redis.Options{
		Options: &rds.Options{
			Addr:     addr,
			Password: config.Load.Redis.Password,
			DB:       config.Load.Redis.DB,
		},
	}); err != nil {
		log.Fatal().Msg(err.Error())
	}

	tcp := listeners.NewTCP(config.Load.Mqtt.ID, config.Load.Mqtt.Listen, &listeners.Config{
		TLSConfig: cert,
	})
	if err := server.AddListener(tcp); err != nil {
		log.Fatal().Msg(err.Error())
	}

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	<-done
	server.Log.Warn().Msg("stopping mqttserver...")
	server.Close()
}

func getTls() (cfg *tls.Config, err error) {
	var cert tls.Certificate

	if config.Load.TLS.UseFile {
		cert, err = tls.LoadX509KeyPair(config.Load.TLS.CertFile, config.Load.TLS.KeyFile)
	} else {
		cert, err = tls.X509KeyPair([]byte(config.Load.TLS.Cert), []byte(config.Load.TLS.Key))
	}

	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return tlsConfig, nil
}
