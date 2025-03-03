module Models.Provider exposing (Provider, providerDecoder)

import Json.Decode exposing (Decoder, field, int, string)

type alias Provider =
    { name : String
    , popularity : Int
    }

providerDecoder : Decoder Provider
providerDecoder =
    Json.Decode.map2 Provider
        (field "name" string)
        (field "popularity" int)