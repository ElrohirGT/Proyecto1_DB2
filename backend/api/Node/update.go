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

type UpdateNodeRequest struct {
	NodeType    string         `json:"NodeType"`     
	Identifier  map[string]any `json:"Identifier"`   
	Properties  map[string]any `json:"Properties"`   
}

func NewUpdateNodeHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use PUT"))
			return
		}
	
		var req UpdateNodeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid JSON: %s", err.Error())))
			return
		}

		if req.NodeType == "" || len(req.Identifier) == 0 || len(req.Properties) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request - `NodeType`, `Identifier` y `Properties` son requeridos"))
			return
		}

		var matchBuilder strings.Builder
		matchBuilder.WriteString("MATCH (n:")
		matchBuilder.WriteString(req.NodeType)
		matchBuilder.WriteString(" {")

		params := make(map[string]any)
		i := 1
		for key, value := range req.Identifier {
			matchBuilder.WriteString(key)
			matchBuilder.WriteString(": $")
			matchBuilder.WriteString("id_" + key)
			params["id_"+key] = value

			if i < len(req.Identifier) {
				matchBuilder.WriteString(", ")
			}
			i++
		}
		matchBuilder.WriteString("}) ")

		var setBuilder strings.Builder
		setBuilder.WriteString("SET ")
		i = 1
		for key, value := range req.Properties {
			setBuilder.WriteString("n.")
			setBuilder.WriteString(key)
			setBuilder.WriteString(" = $")
			setBuilder.WriteString("prop_" + key)
			params["prop_"+key] = value

			if i < len(req.Properties) {
				setBuilder.WriteString(", ")
			}
			i++
		}

		query := matchBuilder.String() + setBuilder.String() + " RETURN n"

		log.Info().Str("query", query).Msg("⏳ Ejecutando consulta de actualización en Neo4j...")
		_, err = neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("❌ Error actualizando en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Update Error: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ NODE UPDATED SUCCESSFULLY"))
	}
}
