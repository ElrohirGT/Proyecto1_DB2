module Api.Endpoint exposing (getHistory, request)

import Http
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


getHistory : String -> Endpoint
getHistory productId =
    url [ "history" ] [ string "ProductId" productId ]
