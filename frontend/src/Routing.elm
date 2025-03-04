module Routing exposing (..)

import Html.Styled exposing (Attribute)
import Html.Styled.Attributes exposing (href)
import Url exposing (Url)
import Url.Parser as P exposing ((</>), Parser, s)


type Route
    = Home
    | Trace
    | Statistics
    | Report
    | NotFound


routeParser : Parser (Route -> c) c
routeParser =
    P.oneOf
        [ P.map Home
            P.top
        , P.map Trace (s "trace")
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


goToTrace : Attribute msg
goToTrace =
    href "/trace"


goToStats : Attribute msg
goToStats =
    href "/stats"


goToHome : Attribute msg
goToHome =
    href "/"
