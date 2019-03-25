package db

import (
	"log"
	"strings"

	"bitbucket.org/danielper/util"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// ConnectArango returns a connection to a arango database
func ConnectArango() (driver.Client, error) {
	arangoHost := util.GetEnvOrDefault("ARANGO_HOST", "localhost")
	arangoUser := util.GetEnvOrDefault("ARANGO_USER", "arango")
	arangoPass := util.GetEnvOrDefault("ARANGO_PASS", "arango")

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: strings.Split(arangoHost, ","),
	})
	if err != nil {
		return nil, err
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangoUser, arangoPass),
	})
	if err != nil {
		return nil, err
	}

	log.Println("ArangoDB connected.")

	return c, nil
}
