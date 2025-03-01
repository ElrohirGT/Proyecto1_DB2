package properties

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

type DeleteNodeRequest struct {
	Target           utils.Neo4JObject
	RemoveProperties []string
	Limit            *int
}

func NewDeleteNodePropertiesHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use DELETE"))
			return
		}

		var req DeleteNodeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(fmt.Appendf(nil, "400 Bad Request - Invalid JSON: %s", err.Error()))
			return
		}

		var b strings.Builder
		b.WriteString("MATCH ")
		req.Target.AppendAsNeo4JMatch(&b, []string{"(", ")"}, "n1")

		for _, prop := range req.RemoveProperties {
			b.WriteString(" REMOVE ")
			b.WriteString("n1.")
			b.WriteString(prop)
		}

		params := make(map[string]any)
		for prop, val := range req.Target.Properties {
			params["n1_"+prop] = val
		}

		b.WriteString(" RETURN n1 ")

		if req.Limit != nil {
			b.WriteString("LIMIT ")
			b.WriteString(fmt.Sprint(*req.Limit))
		}

		query := b.String()
		log.Info().Str("query", query).Msg("⏳ Ejecutando consulta de actualización en Neo4j...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("❌ Error actualizando en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "500 Internal Server Error - Neo4j Update Error: %s", err.Error()))
			return
		}

		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)

		err = enc.Encode(result.Records)
		if err != nil {
			log.Error().Err(err).Interface("array", result.Records[0]).Msg("Error encoding array!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - INTERNAL SERVER ERROR"))
			return
		}

		log.Info().Msg("✅ PROPERTIES REMOVED SUCCESSFULLY")
		w.Write(buff.Bytes())
	}
}
