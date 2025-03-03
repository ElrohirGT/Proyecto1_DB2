module Models.PurchasedProduct exposing (PurchasedProduct, purchasedProductDecoder)

import Json.Decode exposing (Decoder, field, int, string)

type alias PurchasedProduct =
    { productName : String
    , productId : String
    , purchases : Int
    }

purchasedProductDecoder : Decoder PurchasedProduct
purchasedProductDecoder =
    Json.Decode.map3 PurchasedProduct
        (field "product_name" string)
        (field "product_id" string)
        (field "purchases" int)
