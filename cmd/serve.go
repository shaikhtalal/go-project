package cmd

import (
	"context"
	"net/http"
	"strings"

	"bcdpkg.in/go-project/pkg/hasura"
	"bcdpkg.in/go-project/pkg/param"
	"bcdpkg.in/go-project/pkg/router"
	"github.com/coreos/go-oidc"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/cenk/backoff"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
)

const (
	backoffRetries = 10
)

var (
	listenAddr     string
	serviceName    string
	corsHosts      []string
	dbDsn          string
	hasuraEndpoint string
	oidcIssuer     string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "REST server",
	Long:  `Use serve command to run the REST API call's`,
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	corsHosts = viper.GetStringSlice("cors-host")
	dbDsn = viper.GetString("db-dsn")
	serviceName = viper.GetString("service-name")
	listenAddr = viper.GetString("listen-addr")
	hasuraEndpoint = viper.GetString("hasura-endpoint")

	p := &param.Param{
		ServiceName:    serviceName,
		HTTPListenAddr: listenAddr,
		Client: param.Database{
			DSN: dbDsn,
		},
		CorsHosts: corsHosts,
		Log:       logrusEntry,
		OIDC: param.OIDC{
			Addr: oidcIssuer,
		},
	}

	return serve(p)

}

// postgres://postgres:postgres@localhost:5432/kuverit?sslmode=disable
func init() {
	serveCmd.Flags().StringSliceVar(&corsHosts, "cors-host", []string{"*"}, "CORS allowed hosts, comma separated")
	serveCmd.Flags().StringVarP(&dbDsn, "db-dsn", "", "postgres://postgres:postgres@localhost:5432/kuverit?sslmode=disable", "db DSN")
	serveCmd.Flags().StringVarP(&serviceName, "service-name", "", "go-project", "service name")
	serveCmd.Flags().StringVarP(&listenAddr, "listen-addr", "", ":9111", "listen address")
	serveCmd.Flags().StringVarP(&hasuraEndpoint, "hasura-endpoint", "", "http://localhost:8080/v1/graphql", "hasura URL to connect backend with hasura")
	serveCmd.Flags().StringVarP(&oidcIssuer, "oidc-issuer", "", "http://127.0.0.1:4444/", "oidc issuer url")
	RootCmd.AddCommand(serveCmd)

	viper.BindPFlags(serveCmd.Flags())
}

func serve(p *param.Param) error {
	var err error

	if p.OIDC.Client, err = connectOIDC(p.OIDC.Addr); err != nil {
		return err
	}

	p.Client.Db, err = connectPostgres(p.Client.DSN)
	if err != nil {
		return err
	}

	p.GraphQLClient, err = hasura.HasuraClient(hasuraEndpoint, p.Log)
	if err != nil {
		return err
	}
	chErr := make(chan error, 0)

	go func() {
		logrus.Infof("listening on %s", p.HTTPListenAddr)
		chErr <- serveHTTP(p)
	}()

	return <-chErr
}

func serveHTTP(p *param.Param) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Options{
		AllowedOrigins: p.CorsHosts,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"authorization", "*"},
		ExposedHeaders:   []string{"authorization", "*"},
		AllowCredentials: true,
	}))

	param.Inject(r, p)

	router.HandleHTTP(r)

	return errors.Wrap(r.Run(p.HTTPListenAddr), "unable start server")
}

func connectPostgres(dbDSN string) (db *gorm.DB, err error) {
	boff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), backoffRetries)
	u, err := dburl.Parse(dbDSN)
	if err != nil {
		return nil, errors.Wrap(err, "invalid db dsn "+dbDSN)
	}
	err = backoff.Retry(func() error {
	connect:
		db, err = gorm.Open(u.Driver, u.DSN)
		if err != nil {
			err = createDatabase(u.Driver, u.DSN, u.Path)
			if err != nil {
				return err
			}
			goto connect
		}
		return nil
	}, boff)
	return db, err
}

func createDatabase(driver, dsn, path string) error {
	cdb, err := sqlx.Open(driver, "postgres://postgres:postgres@localhost:5432?sslmode=disable")
	dbname := strings.TrimPrefix(path, "/")
	_, err = cdb.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		cdb.Close()
		return err
	}
	cdb.Close()
	logrus.Infof("Created empty database %s, as it did not exist", dbname)
	return nil
}

func connectOIDC(addr string) (op *oidc.Provider, err error) {
	boff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), backoffRetries)

	err = backoff.Retry(func() error {
		if op, err = oidc.NewProvider(context.Background(), addr); err != nil {
			logrus.Errorf("unable to connect to oidc: %v, retrying...", err)
		}
		return errors.Wrap(err, "unable to connect to oidc")
	}, boff)

	return op, err
}
