module Models.Node exposing (..)

import Dict exposing (Dict)
import Json.Decode as Decode exposing (dict, int, list, string, value)
import Json.Decode.Pipeline exposing (required)


type alias Node =
    { id : Int
    , labels : List String
    , props : Dict String Decode.Value
    }


nodeDecoder : Decode.Decoder Node
nodeDecoder =
    Decode.succeed Node
        |> required "Id" int
        |> required "Labels" (list string)
        |> required "Props" (dict value)
