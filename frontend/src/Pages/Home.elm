module Pages.Home exposing (..)

import Html exposing (a, div, h1, text)
import Routing exposing (goToTrace)
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
    { title = "Home!"
    , body =
        [ div []
            [ h1 [] [ text "In Home!" ]
            , a [ goToTrace 5 ] [ text "Go to Trace!" ]
            ]
        ]
    }
