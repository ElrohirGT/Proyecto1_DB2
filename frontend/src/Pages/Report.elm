module Pages.Report exposing (..)

import Html exposing (h1, text)
import Utils exposing (StyledDocument)


type alias Model =
    {}


init : ( Model, Cmd Msg )
init =
    ( {}, Cmd.none )


type Msg
    = None


subscriptions : a -> Sub Msg
subscriptions model =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )


view : Model -> StyledDocument Msg
view model =
    { title = "Report!"
    , body = [ h1 [] [ text "In Report!" ] ]
    }
