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
	NodeType   string         `json:"NodeType"`
	Properties map[string]any `json:"Properties"`
}

func NewCreateNodeHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use POST"))
			return
		}

		var req NodeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid JSON: %s", err.Error())))
			return
		}

		if req.NodeType == "" || len(req.Properties) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request - `NodeType` y `Properties` son requeridos"))
			return
		}

		var queryBuilder strings.Builder
		queryBuilder.WriteString("CREATE (n:")
		queryBuilder.WriteString(req.NodeType)
		queryBuilder.WriteString(" {")

		params := make(map[string]any)
		i := 1
		for key, value := range req.Properties {
			queryBuilder.WriteString(key)
			queryBuilder.WriteString(": $")
			queryBuilder.WriteString(key)
			params[key] = value

			if i < len(req.Properties) {
				queryBuilder.WriteString(", ")
			}
			i++
		}
		queryBuilder.WriteString("}) RETURN n")

		query := queryBuilder.String()

		log.Info().Str("query", query).Msg("Ejecutando consulta en Neo4j...")
		_, err = neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error insertando en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Insert Error: %s", err.Error())))
			return
		}

		//Success message
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("✅ NODE CREATED SUCCESSFULLY"))
	}
}
