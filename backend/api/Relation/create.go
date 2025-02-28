package relation


import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

type CreateRelationRequest struct {
	OriginNode      neo4JObject `json:"OriginNode"`
	DestinationNode neo4JObject `json:"DestinationNode"`
	Relation        neo4JObject `json:"Relation"`
}

func NewCreateRelationHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use POST"))
			return
		}

		var req CreateRelationRequest
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
		queryBuilder.WriteString("}), (n2:")
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
		queryBuilder.WriteString("}) MERGE (n1)-[r:")
		queryBuilder.WriteString(req.Relation.Category)
		queryBuilder.WriteString(" {")

		
		i = 0
		for key := range req.Relation.Properties {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(key)
			queryBuilder.WriteString(": $r_")
			queryBuilder.WriteString(key)
			i++
		}
		queryBuilder.WriteString("}]->(n2) RETURN r")

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

		log.Info().Str("query", query).Msg("Creando relación en Neo4j...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		result, err = neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error al crear la relación en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Create Error: %s", err.Error())))
			return
		}

		if len(result.Records) == 0 {
			log.Warn().Msg(" No se encontró la relación creada")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - Relationship was not created"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("RELATION CREATED SUCCESSFULLY"))

	}
}
