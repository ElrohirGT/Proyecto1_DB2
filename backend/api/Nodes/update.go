package node

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

type UpdateNodeRequest struct {
	Target           utils.Neo4JObject
	UpdateProperties utils.Neo4JObjectProperties
}

func NewUpdateNodesHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
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
			w.Write(fmt.Appendf(nil, "400 Bad Request - Invalid JSON: %s", err.Error()))
			return
		}

		var b strings.Builder
		b.WriteString("MATCH ")
		req.Target.AppendAsNeo4JMatch(&b, []string{"(", ")"}, "n1")

		b.WriteString("SET ")
		for prop := range req.UpdateProperties {
			b.WriteString(prop)
			b.WriteString(" = ")
			b.WriteString("$new_")
			b.WriteString(prop)
		}

		params := make(map[string]any)
		for prop, val := range req.UpdateProperties {
			params["new_"+prop] = val
		}
		for prop, val := range req.Target.Properties {
			params["n1_"+prop] = val
		}

		b.WriteString(" RETURN n1")

		query := b.String()
		log.Info().Str("query", query).Msg("⏳ Ejecutando consulta de actualización en Neo4j...")
		_, err = neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("❌ Error actualizando en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "500 Internal Server Error - Neo4j Update Error: %s", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ NODE UPDATED SUCCESSFULLY"))
	}
}
