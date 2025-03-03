module Models.Product exposing (Product, productDecoder)

import Json.Decode exposing (Decoder, field, float, string)

type alias Product =
    { name : String
    , averageRating : Float
    }

productDecoder : Decoder Product
productDecoder =
    Json.Decode.map2 Product
        (field "name" string)
        (field "average_rating" float)