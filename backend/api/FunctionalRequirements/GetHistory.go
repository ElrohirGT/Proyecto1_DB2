package functionalrequirements

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

func NewGetHistoryHandler(client *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		queries := r.URL.Query()
		productId := queries.Get("ProductId")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		if productId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BAD REQUEST - NO `ProductId` query!"))
			return
		}

		query := `MATCH p1=(:Provider)-[:PRODUCES]->(:Product {id: $id})
RETURN p1
UNION
MATCH p1=(:Provider)-[:PRODUCES]->(:Material)<-[:NEEDS]-(:Product {id: $id})
RETURN p1`

		params := make(map[string]any)
		params["id"] = productId

		log.Info().Str("query", query).Msg("Querying DB...")
		result, err := neo4j.ExecuteQuery(ctx, *client, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

		if err != nil {
			log.Error().Err(err).Msg("Error querying DB!")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("INTERNAL SERVER ERROR - QUERY ERROR `%s`", err.Error())
			w.Write([]byte(msg))
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
