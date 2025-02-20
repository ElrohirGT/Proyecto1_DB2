module Utils exposing (..)

import Html


type alias StyledDocument msg =
    { title : String
    , body : List (Html.Html msg)
    }
