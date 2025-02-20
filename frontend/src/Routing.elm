module Routing exposing (..)

import Url exposing (Url)
import Url.Parser as P exposing ((</>), Parser, s)


type Route
    = Home
    | Trazability { productId : Int }
    | Statistics
    | Report
    | NotFound


routeParser : String -> Parser (Route -> c) c
routeParser basePath =
    P.oneOf
        [ P.map Home
            (s basePath </> P.top)
        , P.map (\id -> Trazability { productId = id }) (s basePath </> s "trace" </> P.int)
        , P.map Statistics (s basePath </> s "stats")
        , P.map Report (s basePath </> s "report")
        ]


parseUrl : String -> Url -> Route
parseUrl basePath url =
    let
        parsedUrl =
            P.parse (routeParser basePath) url
    in
    case parsedUrl of
        Just a ->
            a

        Nothing ->
            NotFound
