module Api.Endpoint exposing (GetHistoryResponse, getHistory, getHistoryResponseDecoder, request)

import Http
import Json.Decode as Decode exposing (list)
import Json.Decode.Pipeline exposing (required)
import Models.Node exposing (Node, nodeDecoder)
import Models.Relation exposing (Relation, relationDecoder)
import Url.Builder exposing (string)


{-| Http.request, except it takes an Endpoint instead of a Url.
-}
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


{-| Get a URL to the localhost API.

This is not publicly exposed, because we want to make sure the only way to get one of these URLs is from this module.

-}
type Endpoint
    = Endpoint String


unwrap : Endpoint -> String
unwrap (Endpoint str) =
    str


url : List String -> List Url.Builder.QueryParameter -> Endpoint
url paths queryParams =
    Url.Builder.crossOrigin "http://localhost:8080" paths queryParams
        |> Endpoint



-- ENDPOINTS
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


getHistoryResponseDecoder : Decode.Decoder GetHistoryResponse
getHistoryResponseDecoder =
    Decode.succeed GetHistoryResponse
        |> required "Values"
            (list
                (Decode.succeed (\nodes -> \relationships -> { nodes = nodes, relationships = relationships })
                    |> required "Nodes" (list nodeDecoder)
                    |> required "Relationships" (list relationDecoder)
                )
            )
