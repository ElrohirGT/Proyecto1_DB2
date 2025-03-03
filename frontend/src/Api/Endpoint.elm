module Api.Endpoint exposing (GetHistoryResponse, getHistory, getHistoryResponseDecoder, GetStatsResponse, getStats, getStatsResponseDecoder, request)

import Http
import Json.Decode exposing (Decoder, field, list, map3, succeed)
import Json.Decode.Pipeline exposing (required)
import Models.Node exposing (Node, nodeDecoder)
import Models.Relation exposing (Relation, relationDecoder)
import Models.Product exposing (Product, productDecoder)
import Models.Provider exposing (Provider, providerDecoder)
import Models.PurchasedProduct exposing (PurchasedProduct, purchasedProductDecoder)
import Url.Builder exposing (string)


-- HTTP Request Helper

request :
    { body : Http.Body
    , expect : Http.Expect msg
    , headers : List Http.Header
    , method : String
    , timeout : Maybe Float
    , url : Endpoint
    , tracker : Maybe String
    }
    -> Cmd msg
request config =
    Http.request
        { body = config.body
        , expect = config.expect
        , headers = config.headers
        , method = config.method
        , timeout = config.timeout
        , url = unwrap config.url
        , tracker = config.tracker
        }


-- TYPES

type Endpoint
    = Endpoint String

unwrap : Endpoint -> String
unwrap (Endpoint str) =
    str

url : List String -> List Url.Builder.QueryParameter -> Endpoint
url paths queryParams =
    Url.Builder.crossOrigin "http://localhost:8080" paths queryParams
        |> Endpoint


-- GET /history endpoint

getHistory : String -> Endpoint
getHistory productId =
    url [ "history" ] [ string "ProductId" productId ]

type alias GetHistoryResponse =
    { values :
        List
            { nodes : List Node
            , relationships : List Relation
            }
    }

getHistoryResponseDecoder : Decoder GetHistoryResponse
getHistoryResponseDecoder =
    succeed GetHistoryResponse
        |> required "Values"
            (list
                (succeed (\nodes relationships -> { nodes = nodes, relationships = relationships })
                    |> required "Nodes" (list nodeDecoder)
                    |> required "Relationships" (list relationDecoder)
                )
            )


-- GET /stats endpoint

getStats : Endpoint
getStats =
    url [ "statistics" ] []

type alias GetStatsResponse =
    { topProducts : List Product
    , topProviders : List Provider
    , topPurchasedProducts : List PurchasedProduct
    }

getStatsResponseDecoder : Decoder GetStatsResponse
getStatsResponseDecoder =
    map3 GetStatsResponse
        (field "top_products" (list productDecoder))
        (field "top_providers" (list providerDecoder))
        (field "top_purchased_products" (list purchasedProductDecoder))
