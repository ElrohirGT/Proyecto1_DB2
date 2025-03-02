module Models.Relation exposing (..)

import Dict exposing (Dict)
import Json.Decode as Decode exposing (dict, int, string)
import Json.Decode.Pipeline exposing (required)


type alias Relation =
    { startId : Int
    , endId : Int
    , relType : String
    , props : Dict String Decode.Value
    }


relationDecoder : Decode.Decoder Relation
relationDecoder =
    Decode.succeed Relation
        |> required "StartId" int
        |> required "EndId" int
        |> required "Type" string
        |> required "Props" (dict Decode.value)
