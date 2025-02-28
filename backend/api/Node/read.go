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

type GetNodeGeneric struct {
	NodeType   string
	Properties map[string]any
}

func NewReadNodeHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		url_queries := r.URL.Query()
		nodeType := url_queries.Get("NodeType")

		if nodeType == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BAD REQUEST - NO `NodeType` query!"))
			return
		}

		properties := url_queries.Get("Properties")
		if properties == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BAD REQUEST - NO `Properties` query!"))
			return
		}

		decoder := json.NewDecoder(strings.NewReader(properties))
		var nodeProperties map[string]any

		err := decoder.Decode(&nodeProperties)
		if err != nil {
			msg := fmt.Sprintf("BAD REQUEST - INVALID `Properties` array! `%s`", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(msg))
			return
		}

		// MATCH (n:$NodeType {$nodeName: $nodeValue})
		// RETURN n
		b := strings.Builder{}
		b.WriteString("MATCH (n:")
		b.WriteString(nodeType)
		b.WriteString(" {")

		i := 1
		max_i := len(nodeProperties)
		for property := range nodeProperties {
			b.WriteString(property)
			b.WriteString(": ")
			b.WriteRune('$')
			b.WriteString(property)

			if i != max_i {
				b.WriteRune(',')
			}
		}
		b.WriteString("}) RETURN n LIMIT 1")

		query := b.String()
		// query = "MATCH (n) -[]->() RETURN n"
		log.Info().Str("query", query).Msg("Querying DB...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, nodeProperties, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Fatal().Err(err).Msg("Error querying DB!")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("INTERNAL SERVER ERROR - QUERY ERROR `%s`", err.Error())
			w.Write([]byte(msg))
		}
		log.Info().Int("recordCount", len(result.Records)).Msg("Done!")

		for _, record := range result.Records {
			log.Info().Interface("row", record).Msg("Node found!")
		}

		w.Write([]byte("OK"))
	}
}
