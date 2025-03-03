package relation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrohirGT/Proyecto1_DB2/api/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

type DeleteManyRelationsRequest struct {
	OriginNode      utils.Neo4JObject `json:"OriginNode"`
	DestinationNode utils.Neo4JObject `json:"DestinationNode"`
	Relation        utils.Neo4JObject `json:"Relation"`
	Limit           int               `json:"Limit,omitempty"`
}

type DeleteManyResponse struct {
	DeletedCount int `json:"DeletedCount"`
}

func NewDeleteManyRelationsHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use DELETE"))
			return
		}

		var req DeleteManyRelationsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid JSON: %s", err.Error())))
			return
		}

		if req.OriginNode.Category == "" || req.DestinationNode.Category == "" || req.Relation.Category == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request - `OriginNode`, `DestinationNode`, y `Relation` son requeridos"))
			return
		}

		var queryBuilder strings.Builder
		queryBuilder.WriteString("MATCH (n1:")
		queryBuilder.WriteString(req.OriginNode.Category)
		queryBuilder.WriteString(" {")

		i := 0
		for key := range req.OriginNode.Properties {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(key)
			queryBuilder.WriteString(": $n1_")
			queryBuilder.WriteString(key)
			i++
		}
		queryBuilder.WriteString("})-[r:")
		queryBuilder.WriteString(req.Relation.Category)
		queryBuilder.WriteString("]->(n2:")
		queryBuilder.WriteString(req.DestinationNode.Category)
		queryBuilder.WriteString(" {")

		i = 0
		for key := range req.DestinationNode.Properties {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(key)
			queryBuilder.WriteString(": $n2_")
			queryBuilder.WriteString(key)
			i++
		}

		deleteLimit := ""
		if req.Limit > 0 {
			deleteLimit = fmt.Sprintf(" LIMIT %d ", req.Limit)
		}

		// Aquí cambiamos DELETE r por DELETE r RETURN properties(r) AS deletedRelation
		queryBuilder.WriteString("}) WITH r ")
		queryBuilder.WriteString(deleteLimit)
		queryBuilder.WriteString("DELETE r")

		query := queryBuilder.String()

		params := make(map[string]any)

		for property, val := range req.OriginNode.Properties {
			params["n1_"+property] = val
		}
		for property, val := range req.Relation.Properties {
			params["r_"+property] = val
		}
		for property, val := range req.DestinationNode.Properties {
			params["n2_"+property] = val
		}

		log.Info().Str("query", query).Msg("Buscando y eliminando relación en Neo4j...")
		log.Debug().Msg(fmt.Sprintf("%v", params))

		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error al eliminar la relación en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Delete Error: %s", err.Error())))
			return
		}

		response := DeleteResponse{
			DeletedCount: result.Summary.Counters().RelationshipsDeleted(),
		}

		log.Info().Interface("Deleted Relation", response).Msg("Relación eliminada correctamente")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
