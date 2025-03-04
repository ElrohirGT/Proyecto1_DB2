module Pages.Home exposing (..)

import Css exposing (center, color, column, displayFlex, flexDirection, justifyContent, minWidth, row, spaceAround, textAlign)
import Html.Styled exposing (a, div, h1, text)
import Html.Styled.Attributes exposing (css)
import Routing exposing (goToStats, goToTrace)
import Theme exposing (colors)
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
        [ div [ css (Theme.divBackgroundStyles ++ [ flexDirection column, Css.property "gap" "2rem" ]) ]
            [ h1
                [ css
                    [ textAlign center
                    , color (Css.hex colors.secondary)
                    ]
                ]
                [ text "¡Bienvenido al Inicio!" ]
            , div
                [ css
                    [ displayFlex
                    , justifyContent spaceAround
                    , minWidth (Css.vw 50)
                    ]
                ]
                [ a [ goToTrace, css Theme.btnStyles ] [ text "Trazabilidad" ]
                , a [ goToStats, css Theme.btnStyles ] [ text "Estadísticas" ]
                ]
            ]
        ]
    }
