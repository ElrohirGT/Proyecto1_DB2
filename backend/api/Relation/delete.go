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

type DeleteRelationRequest struct {
	OriginNode      utils.Neo4JObject `json:"OriginNode"`
	DestinationNode utils.Neo4JObject `json:"DestinationNode"`
	Relation        utils.Neo4JObject `json:"Relation"`
}

func NewDeleteRelationHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodDelete {
			http.Error(w, "405 Method Not Allowed - Use DELETE", http.StatusMethodNotAllowed)
			return
		}

		var req DeleteRelationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 Bad Request - Invalid JSON: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if req.OriginNode.Category == "" || req.DestinationNode.Category == "" || req.Relation.Category == "" {
			http.Error(w, "400 Bad Request - OriginNode, DestinationNode, y Relation son requeridos", http.StatusBadRequest)
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
		queryBuilder.WriteString("}) DELETE r RETURN COUNT(r) AS deletedCount")

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

		log.Info().Str("query", query).Interface("params", params).Msg("🔍 Ejecutando DELETE en Neo4j...")

		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			log.Error().Err(err).Msg("Error al eliminar la relación en Neo4j")
			http.Error(w, fmt.Sprintf("500 Internal Server Error - Neo4j Delete Error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if len(result.Records) == 0 {
			log.Warn().Msg("No se encontró la relación para eliminar")
			http.Error(w, "404 ERROR - No relation found to delete", http.StatusNotFound)
			return
		}

		record := result.Records[0]

		deletedCountInterface, found := record.Get("deletedCount")
		if !found {
			log.Warn().Msg("No se encontró la relación para eliminar")
			http.Error(w, "404 ERROR - No relation found to delete", http.StatusNotFound)
			return
		}

		deletedCount, ok := deletedCountInterface.(int64)
		if !ok {
			log.Error().Msg("Error al convertir deletedCount a int64")
			http.Error(w, "500 Internal Server Error - Error en la conversión de datos", http.StatusInternalServerError)
			return
		}

		if deletedCount == 0 {
			log.Warn().Msg("No se encontró la relación para eliminar")
			http.Error(w, "404 ERROR - No relation found to delete", http.StatusNotFound)
			return
		}

		if deletedCount == 0 {
			log.Warn().Msg("No se encontró la relación para eliminar")
			http.Error(w, "404 ERROR - No relation found to delete", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(req)
	}
}
