package api

import (
	"net/http"

	node "github.com/ElrohirGT/Proyecto1_DB2/api/Node"
	"github.com/ElrohirGT/Proyecto1_DB2/api/health"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Api struct {
	dbClient *neo4j.DriverWithContext

	// Handlers: endpoint functions
	CheckHealthHandler http.HandlerFunc
	CreateUserHandler  http.HandlerFunc

	// CRUD

	ReadNodeHandler   http.HandlerFunc
	DeleteNodeHandler http.HandlerFunc

	// CRUD (multiple)

	// FUNC REQUIREMENTS

}

func NewApi(
	mongoClient *neo4j.DriverWithContext,
) *Api {

	return &Api{
		dbClient: mongoClient,

		CheckHealthHandler: health.CheckHealthHandler,

		ReadNodeHandler:   node.NewReadNodeHandler(mongoClient),
		DeleteNodeHandler: node.NewDeleteNodeHandler(mongoClient),
	}
}
