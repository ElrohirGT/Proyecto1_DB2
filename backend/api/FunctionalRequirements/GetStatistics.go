package functionalrequirements

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func GetStatisticsHandler(db *neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if db == nil {
			http.Error(w, "500 Internal Server Error - DB connection is nil", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		queryTopProducts := `
			MATCH (c:Consumer)-[r:RATES]->(p:Product)
			RETURN p.name AS name, AVG(r.rating) AS average_rating
			ORDER BY average_rating DESC
			LIMIT 3
		`

		queryTopProviders := `
			MATCH (p:Provider)<-[r:PREFERS]-(c:Retailer)
			RETURN p.name AS name, COUNT(r) AS popularity
			ORDER BY popularity DESC
			LIMIT 5
		`

		queryTopPurchasedProducts := `
			MATCH (c:Consumer)-[r:BUYS_FROM_RETAILER]->(retailer:Retailer),
			      (p:Product {id: r.productId})  
			RETURN p.name AS product_name, r.productId AS product_id, COUNT(r) AS purchases
			ORDER BY purchases DESC
			LIMIT 10
		`

		log.Info().Msg("Ejecutando consultas de estadísticas...")

		// 🏆 Ejecutar top productos
		resultTopProducts, err := neo4j.ExecuteQuery(ctx, *db, queryTopProducts, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			http.Error(w, fmt.Sprintf("500 Internal Server Error - Error en top productos: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var topProducts []map[string]any
		for _, record := range resultTopProducts.Records {
			name, _ := record.Get("name")
			averageRating, _ := record.Get("average_rating")

			topProducts = append(topProducts, map[string]any{
				"name":           name,
				"average_rating": averageRating,
			})
		}

		// 🔝 Ejecutar top proveedores
		resultTopProviders, err := neo4j.ExecuteQuery(ctx, *db, queryTopProviders, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			http.Error(w, fmt.Sprintf("500 Internal Server Error - Error en top proveedores: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var topProviders []map[string]any
		for _, record := range resultTopProviders.Records {
			name, _ := record.Get("name")
			popularity, _ := record.Get("popularity")

			topProviders = append(topProviders, map[string]any{
				"name":       name,
				"popularity": popularity,
			})
		}

		// 🛍️ Ejecutar top productos comprados
		resultTopPurchasedProducts, err := neo4j.ExecuteQuery(ctx, *db, queryTopPurchasedProducts, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			http.Error(w, fmt.Sprintf("500 Internal Server Error - Error en top productos comprados: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		var topPurchasedProducts []map[string]any
		for _, record := range resultTopPurchasedProducts.Records {
			productName, _ := record.Get("product_name")
			productID, _ := record.Get("product_id")
			purchases, _ := record.Get("purchases")

			topPurchasedProducts = append(topPurchasedProducts, map[string]any{
				"product_name": productName,
				"product_id":   productID,
				"purchases":    purchases,
			})
		}

		response := map[string]any{
			"top_products":           topProducts,
			"top_providers":          topProviders,
			"top_purchased_products": topPurchasedProducts,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
