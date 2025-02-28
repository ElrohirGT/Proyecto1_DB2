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

type ReadNeo4JObject struct {
	Category   string         `json:"Category"`
	Properties map[string]any `json:"Properties"`
}

type ReadRelationRequest struct {
	OriginNode      ReadNeo4JObject `json:"OriginNode"`
	DestinationNode ReadNeo4JObject `json:"DestinationNode"`
	Relation        ReadNeo4JObject `json:"Relation"`
}


func NewReadRelationHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method Not Allowed - Use GET"))
			return
		}

		queryParams := r.URL.Query()
		origin := queryParams.Get("OriginNode")
		destination := queryParams.Get("DestinationNode")
		relation := queryParams.Get("Relation")

		if origin == "" || destination == "" || relation == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request - Missing required query parameters"))
			return
		}

		var originNode, destinationNode, relationObj ReadNeo4JObject
		err := json.Unmarshal([]byte(origin), &originNode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid `OriginNode`: %s", err.Error())))
			return
		}

		err = json.Unmarshal([]byte(destination), &destinationNode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid `DestinationNode`: %s", err.Error())))
			return
		}

		err = json.Unmarshal([]byte(relation), &relationObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("400 Bad Request - Invalid `Relation`: %s", err.Error())))
			return
		}


		b := &strings.Builder{}
		b.WriteString("MATCH ")
		originNode.appendAsNode(b, "n1") 
		b.WriteString(" -")
		relationObj.appendAsRelation(b, "r")
		b.WriteString("-> ")
		destinationNode.appendAsNode(b, "n2")
		b.WriteString(" RETURN r, n1, n2 LIMIT 10")

		query := b.String()
		params := make(map[string]any)

		for property, val := range originNode.Properties {
			params["n1_"+property] = val
		}

		for property, val := range relationObj.Properties {
			params["r_"+property] = val
		}

		for property, val := range destinationNode.Properties {
			params["n2_"+property] = val
		}

		log.Info().Str("query", query).Msg("Ejecutando consulta de búsqueda en Neo4j...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error consultando en Neo4j")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 Internal Server Error - Neo4j Query Error: %s", err.Error())))
			return
		}

		if len(result.Records) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 ERROR - No relationships found"))
			return
		}

		var relationships []map[string]any
		for _, record := range result.Records {
			relData := map[string]any{}

			if rel, ok := record.Get("r"); ok {
				relData["relation"] = rel
			}
			if origin, ok := record.Get("n1"); ok {
				relData["origin"] = origin
			}
			if destination, ok := record.Get("n2"); ok {
				relData["destination"] = destination
			}

			relationships = append(relationships, relData)
		}

		response, err := json.Marshal(relationships)
		if err != nil {
			log.Error().Err(err).Msg("Error codificando la respuesta")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - INTERNAL SERVER ERROR"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func (self *ReadNeo4JObject) appendAsNode(b *strings.Builder, queryId string) {
	b.WriteRune('(')
	b.WriteString(queryId)
	b.WriteRune(':')
	b.WriteString(self.Category)

	if len(self.Properties) > 0 {
		b.WriteString(" {")
		i := 0
		for property := range self.Properties {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(property)
			b.WriteString(": $")
			b.WriteString(queryId)
			b.WriteString("_")
			b.WriteString(property)
			i++
		}
		b.WriteString("}")
	}
	b.WriteRune(')')
}

func (self *ReadNeo4JObject) appendAsRelation(b *strings.Builder, queryId string) {
	b.WriteRune('[')
	b.WriteString(queryId)
	b.WriteRune(':')
	b.WriteString(self.Category)

	if len(self.Properties) > 0 {
		b.WriteString(" {")
		i := 0
		for property := range self.Properties {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(property)
			b.WriteString(": $")
			b.WriteString(queryId)
			b.WriteString("_")
			b.WriteString(property)
			i++
		}
		b.WriteString("}")
	}
	b.WriteRune(']')
}
