package relationproperties

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElrohirGT/Proyecto1_DB2/api/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

type CreateRelationPropertiesRequest struct {
	OriginNode      utils.Neo4JObject `json:"OriginNode"`
	DestinationNode utils.Neo4JObject `json:"DestinationNode"`
	Relation        utils.Neo4JObject `json:"Relation"`
	Properties      map[string]any
}

type CreateRelationPropertiesResponse struct {
	DeletedCount int `json:"DeletedCount"`
}

func NewCreateRelationPropertiesHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var req CreateRelationPropertiesRequest
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
		queryBuilder.WriteString("}) ")

		for prop := range req.Properties {
			queryBuilder.WriteString(" SET ")
			queryBuilder.WriteString("r.")
			queryBuilder.WriteString(prop)
			queryBuilder.WriteString(" = ")
			queryBuilder.WriteString("$new_")
			queryBuilder.WriteString(prop)
		}

		queryBuilder.WriteString(" RETURN r")

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
		for prop, val := range req.Properties {
			params["new_"+prop] = val
		}

		log.Info().Str("query", query).Msg("Creando/Actualizando propiedades de nodos...")
		log.Debug().Msg(fmt.Sprintf("%v", params))

		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error al eliminar la relación en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Delete Error: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)
		err = enc.Encode(result.Records)
		if err != nil {
			log.Error().Err(err).Interface("row", result.Records[0]).Msg("Error encoding row!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - INTERNAL SERVER ERROR"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(buff.Bytes())
	}
}
