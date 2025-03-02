package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)


type NodeRequest struct {
	NodeType   string         json:"NodeType"
	Properties map[string]any json:"Properties"
}

func NewCreateNodeHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var req NodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid JSON: %s", err.Error())))
			return
		}

		if req.NodeType == "" || len(req.Properties) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request - NodeType y Properties son requeridos"))
			return
		}

		var queryBuilder strings.Builder
		queryBuilder.WriteString("CREATE (n:")
		queryBuilder.WriteString(req.NodeType)
		queryBuilder.WriteString(" {")

		params := make(map[string]any)
		i := 0
		for key, value := range req.Properties {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(key)
			queryBuilder.WriteString(": $")
			queryBuilder.WriteString(key)
			params[key] = value
			i++
		}
		queryBuilder.WriteString("}) RETURN n")

		query := queryBuilder.String()

		log.Info().Str("query", query).Msg("Creando nodo en Neo4j...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error al crear el nodo en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Create Error: %s", err.Error())))
			return
		}

		if len(result.Records) == 0 {
			log.Warn().Msg(" No se pudo crear el nodo")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - Node was not created"))
			return
		}
		
		record := result.Records[0]
		createdNode, found := record.Get("n") 
		
		if !found {
			log.Warn().Msg("⚠ No se pudo recuperar el nodo creado")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - Created node not found"))
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdNode)
	}
}