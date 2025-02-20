package db_client

import (
	"context"

	"github.com/ElrohirGT/Proyecto1_DB2/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// NewDriver initializes and returns a Neo4j driver along with its close function.
func NewDriver(config *config.DatabaseConfig) (*neo4j.DriverWithContext, func(context.Context) error, error) {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(config.DBUri, neo4j.BasicAuth(config.DBUser, config.DBPassword, ""))
	if err != nil {
		return nil, nil, err
	}

	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx) // Ensure cleanup if verification fails
		return nil, nil, err
	}

	// Return the driver and its close function
	return &driver, driver.Close, nil
}
