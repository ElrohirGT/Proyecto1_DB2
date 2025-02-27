package api

import (
	"net/http"

	node "github.com/ElrohirGT/Proyecto1_DB2/api/Node"
	relation "github.com/ElrohirGT/Proyecto1_DB2/api/Relation"
	"github.com/ElrohirGT/Proyecto1_DB2/api/health"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Api struct {
	dbClient *neo4j.DriverWithContext

	// Handlers: endpoint functions
	CheckHealthHandler http.HandlerFunc
	CreateUserHandler  http.HandlerFunc

	// CRUD
	CreateNodeHandler http.HandlerFunc
	ReadNodeHandler   http.HandlerFunc
	UpdateNodeHandler http.HandlerFunc
	DeleteNodeHandler http.HandlerFunc

	ReadRelationHandler  http.HandlerFunc
	UpdateRelationHandler http.HandlerFunc

	// CRUD (multiple)

	// FUNC REQUIREMENTS

}

func NewApi(
	mongoClient *neo4j.DriverWithContext,
) *Api {

	return &Api{
		dbClient: mongoClient,

		CheckHealthHandler: health.CheckHealthHandler,
		CreateNodeHandler: node.NewCreateNodeHandler(mongoClient),
		ReadNodeHandler:       node.NewReadNodeHandler(mongoClient),
		UpdateNodeHandler:	 node.NewUpdateNodeHandler(mongoClient),
		DeleteNodeHandler:     node.NewDeleteNodeHandler(mongoClient),
		ReadRelationHandler:	relation.NewReadRelationHandler(mongoClient),
		UpdateRelationHandler: relation.NewUpdateRelationHandler(mongoClient),
	}
}
