module Theme exposing (..)

import Css exposing (alignItems, backgroundColor, border, center, color, displayFlex, height, justifyContent, none, padding2, px, textDecoration, width, zero)


colors : { background : String, error : String, secondary : String, warning : String, accent : String }
colors =
    { background = "#352f3d"
    , error = "#cf1500"
    , secondary = "#fc8020"
    , warning = "#f8a523"
    , accent = "#8f9044"
    }


divBackgroundStyles : List Css.Style
divBackgroundStyles =
    [ width (Css.vw 100)
    , height (Css.vh 100)
    , backgroundColor (Css.hex colors.background)
    , displayFlex
    , justifyContent center
    , alignItems center
    ]


btnStyles : List Css.Style
btnStyles =
    [ textDecoration none
    , padding2 (px 10) (px 20)
    , border zero
    , color (Css.hex "#ffffff")
    , backgroundColor (Css.hex colors.accent)
    ]
