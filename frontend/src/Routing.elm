module Routing exposing (..)

import Html exposing (Attribute)
import Html.Attributes exposing (href)
import Url exposing (Url)
import Url.Parser as P exposing ((</>), Parser, s)


type Route
    = Home
    | Trace { productId : Int }
    | Statistics
    | Report
    | NotFound


routeParser : Parser (Route -> c) c
routeParser =
    P.oneOf
        [ P.map Home
            P.top
        , P.map (\id -> Trace { productId = id }) (s "trace" </> P.int)
        , P.map Statistics (s "stats")
        , P.map Report (s "report")
        ]


parseUrl : Url -> Route
parseUrl url =
    let
        parsedUrl =
            P.parse routeParser url
    in
    case parsedUrl of
        Just a ->
            a

        Nothing ->
            NotFound


goToTrace : Int -> Attribute msg
goToTrace productId =
    href (String.join "/" [ "/trace", String.fromInt productId ])
