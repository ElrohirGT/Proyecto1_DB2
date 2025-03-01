package relation

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

type neo4JObject struct {
	Category   string
	Properties map[string]any
}

type reqBody struct {
	OriginNode      neo4JObject
	DestinationNode neo4JObject
	Relation        neo4JObject
	NewProperties   map[string]any
}

func (self *neo4JObject) appendAsNeo4JMatch(b *strings.Builder, limits []string, queryId string) {
	if len(limits) != 2 {
		panic("The should only be two limits! Example []string {\"[\", \"]\"}")
	}

	b.WriteString(limits[0])
	b.WriteString(queryId)
	b.WriteRune(':')
	b.WriteString(self.Category)

	propertiesCount := len(self.Properties)
	if propertiesCount > 0 {
		b.WriteString(" {")

		i := 1
		for property := range self.Properties {
			b.WriteString(property)
			b.WriteString(": ")
			b.WriteRune('$')
			b.WriteString(queryId)
			b.WriteRune('_')
			b.WriteString(property)

			if i != propertiesCount {
				b.WriteRune(',')
			}
			i++
		}
		b.WriteString("}")
	}
	b.WriteString(limits[1])
}

func NewUpdateRelationHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
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

		// MATCH (n:$NodeType {$nodeName: $nodeValue}) -[r:$RelationType {$relName: $relValue}]-> (n2:$NodeType {$nodeName: $nodeValue})
		// SET r.$NewProperty = $NewValue
		// RETURN r
		// LIMIT 1
		b := &strings.Builder{}
		b.WriteString("MATCH ")
		body.OriginNode.appendAsNeo4JMatch(b, []string{"(", ")"}, "n1")
		b.WriteString(" -")
		body.Relation.appendAsNeo4JMatch(b, []string{"[", "]"}, "r")
		b.WriteString("-> ")
		body.DestinationNode.appendAsNeo4JMatch(b, []string{"(", ")"}, "n2")

		for propertyName := range body.NewProperties {
			b.WriteString("SET ")
			b.WriteString("r.")
			b.WriteString(propertyName)
			b.WriteString(" = ")
			b.WriteString("$new_")
			b.WriteString(propertyName)
			b.WriteRune(' ')
		}
		b.WriteString(" RETURN r LIMIT 1")

		query := b.String()
		params := make(map[string]any)

		for property, val := range body.OriginNode.Properties {
			key := strings.Join([]string{"n1_", property}, "")
			params[key] = val
		}

		for property, val := range body.Relation.Properties {
			key := strings.Join([]string{"r_", property}, "")
			params[key] = val
		}

		for property, val := range body.DestinationNode.Properties {
			key := strings.Join([]string{"n2_", property}, "")
			params[key] = val
		}

		for property, val := range body.NewProperties {
			key := strings.Join([]string{"new_", property}, "")
			params[key] = val
		}

		log.Info().Str("query", query).Msg("Querying DB...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error querying DB!")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("INTERNAL SERVER ERROR - QUERY ERROR `%s`", err.Error())
			w.Write([]byte(msg))
			return
		}
		nodeCount := len(result.Records)
		log.Info().Int("recordCount", nodeCount).Msg("Done!")

		if nodeCount == 0 {
			log.Error().Err(err).Msg("No row found!")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 ERROR - NOT FOUND"))
			return
		}

		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)

		err = enc.Encode(result.Records[0])
		if err != nil {
			log.Error().Err(err).Interface("row", result.Records[0]).Msg("Error encoding row!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 ERROR - INTERNAL SERVER ERROR"))
			return
		}

		w.Write(buff.Bytes())
	}
}
