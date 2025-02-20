module Pages.Home exposing (..)

import Html exposing (h1, text)
import Utils exposing (StyledDocument)


type alias Model =
    {}


init : ( Model, Cmd msg )
init =
    ( {}, Cmd.none )


type Msg
    = None


subscriptions : a -> Sub msg
subscriptions model =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )


view : Model -> StyledDocument msg
view model =
    { title = "Home!"
    , body = [ h1 [] [ text "In Home!" ] ]
    }
