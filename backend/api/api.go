package api

import (
	"net/http"

	functionalrequirements "github.com/ElrohirGT/Proyecto1_DB2/api/FunctionalRequirements"
	node "github.com/ElrohirGT/Proyecto1_DB2/api/Node"
	properties "github.com/ElrohirGT/Proyecto1_DB2/api/Properties"
	relation "github.com/ElrohirGT/Proyecto1_DB2/api/Relation"
	relationproperties "github.com/ElrohirGT/Proyecto1_DB2/api/RelationProperties"
	"github.com/ElrohirGT/Proyecto1_DB2/api/health"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Api struct {
	dbClient *neo4j.DriverWithContext

	// Handlers: endpoint functions
	CheckHealthHandler http.HandlerFunc
	CreateUserHandler  http.HandlerFunc

	// CRUD
	CreateNodeHandler      http.HandlerFunc
	ReadNodeHandler        http.HandlerFunc
	UpdateNodeHandler      http.HandlerFunc
	DeleteNodeHandler      http.HandlerFunc
	DeleteManyNodesHandler http.HandlerFunc

	CreateRelationHandler      http.HandlerFunc
	ReadRelationHandler        http.HandlerFunc
	UpdateRelationHandler      http.HandlerFunc
	DeleteRelationHandler      http.HandlerFunc
	DeleteManyRelationsHandler http.HandlerFunc

	// RELATION PROPERTIES
	CreateRelationPropertiesHandler http.HandlerFunc

	// CRUD (multiple)
	UpdatePropertiesHandler http.HandlerFunc
	DeletePropertiesHandler http.HandlerFunc

	// FUNC REQUIREMENTS
	GetProductHistoryHandler http.HandlerFunc
	GetStatisticsHandler     http.HandlerFunc
}

func NewApi(
	mongoClient *neo4j.DriverWithContext,
) *Api {

	return &Api{
		dbClient: mongoClient,

		CheckHealthHandler: health.CheckHealthHandler,

		CreateNodeHandler:      node.NewCreateNodeHandler(mongoClient),
		ReadNodeHandler:        node.NewReadNodeHandler(mongoClient),
		UpdateNodeHandler:      node.NewUpdateNodeHandler(mongoClient),
		DeleteNodeHandler:      node.NewDeleteNodeHandler(mongoClient),
		DeleteManyNodesHandler: node.NewDeleteManyNodesHandler(mongoClient),

		CreateRelationHandler:      relation.NewCreateRelationHandler(mongoClient),
		ReadRelationHandler:        relation.NewReadRelationHandler(mongoClient),
		UpdateRelationHandler:      relation.NewUpdateRelationHandler(mongoClient),
		DeleteRelationHandler:      relation.NewDeleteRelationHandler(mongoClient),
		DeleteManyRelationsHandler: relation.NewDeleteManyRelationsHandler(mongoClient),

		CreateRelationPropertiesHandler: relationproperties.NewCreateRelationPropertiesHandler(mongoClient),

		UpdatePropertiesHandler: properties.NewUpdatePropertiesHandler(mongoClient),
		DeletePropertiesHandler: properties.NewDeleteNodePropertiesHandler(mongoClient),

		GetProductHistoryHandler: functionalrequirements.NewGetHistoryHandler(mongoClient),
		GetStatisticsHandler:     functionalrequirements.GetStatisticsHandler(mongoClient),
	}

}
