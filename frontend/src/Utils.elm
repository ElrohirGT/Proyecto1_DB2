module Utils exposing (..)

import Browser
import Html.Styled


type alias StyledDocument msg =
    { title : String
    , body : List (Html.Styled.Html msg)
    }


mapMsg : (a -> b) -> StyledDocument a -> StyledDocument b
mapMsg mapper doc =
    { title = doc.title
    , body = List.map (Html.Styled.map mapper) doc.body
    }


toUnstyled : StyledDocument msg -> Browser.Document msg
toUnstyled doc =
    { title = doc.title
    , body = List.map Html.Styled.toUnstyled doc.body
    }
