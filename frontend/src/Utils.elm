module Utils exposing (..)

import Html


type alias StyledDocument msg =
    { title : String
    , body : List (Html.Html msg)
    }


mapMsg : (a -> b) -> StyledDocument a -> StyledDocument b
mapMsg mapper doc =
    { title = doc.title
    , body = List.map (Html.map mapper) doc.body
    }
