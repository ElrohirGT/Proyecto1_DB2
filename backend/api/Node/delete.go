package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

type reqBody struct {
	NodeType   string
	Properties map[string]any
}

func NewDeleteNodeHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		decoder := json.NewDecoder(r.Body)
		var body reqBody

		err := decoder.Decode(&body)
		if err != nil {
			msg := fmt.Sprintf("BAD REQUEST - INVALID BODY! `%s`", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(msg))
			return
		}

		// MATCH (n:$NodeType {$nodeName: $nodeValue})
		// DETACH DELETE n
		// RETURN n
		// LIMIT 1
		b := strings.Builder{}
		b.WriteString("MATCH (n:")
		b.WriteString(body.NodeType)
		b.WriteString(" {")

		i := 1
		max_i := len(body.Properties)
		for property := range body.Properties {
			b.WriteString(property)
			b.WriteString(": ")
			b.WriteRune('$')
			b.WriteString(property)

			if i != max_i {
				b.WriteRune(',')
			}
		}
		b.WriteString("}) DETACH DELETE n LIMIT 1 RETURN n")

		query := b.String()
		log.Info().Str("query", query).Msg("Querying DB...")
		result, err := neo4j.ExecuteQuery(
			ctx, *client, query,
			body.Properties,
			neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"),
		)

		if err != nil {
			log.Fatal().Err(err).Msg("Error querying DB!")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("INTERNAL SERVER ERROR - QUERY ERROR `%s`", err.Error())
			w.Write([]byte(msg))
		}
		nodeCount := len(result.Records)
		log.Info().Int("recordCount", nodeCount).Msg("Done!")

		if nodeCount == 0 {
			log.Fatal().Err(err).Msg("No node found!")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 ERROR - NOT FOUND"))
			return
		}

		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)

		err = enc.Encode(result.Records[0])
		if err != nil {
			log.Fatal().Err(err).Interface("row", result.Records[0]).Msg("Error encoding row!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - INTERNAL SERVER ERROR"))
			return
		}

		w.Write(buff.Bytes())
	}
}
