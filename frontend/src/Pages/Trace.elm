module Pages.Trace exposing (..)

import Api.Endpoint exposing (getHistory, request)
import Html.Styled exposing (button, div, h1, input, text)
import Html.Styled.Attributes exposing (value)
import Html.Styled.Events exposing (onClick, onInput)
import Http
import Utils exposing (StyledDocument)


type alias APIResponse =
    Result Http.Error String


type alias Model =
    { productId : String
    , isLoading : Bool
    , history : Maybe APIResponse
    }


init : ( Model, Cmd Msg )
init =
    ( Model "5" False Nothing, Cmd.none )


type Msg
    = SearchClicked
    | GotHistory APIResponse
    | ProductIdChanged String


subscriptions : a -> Sub Msg
subscriptions model =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SearchClicked ->
            ( { model | isLoading = True }
            , request
                { url = getHistory model.productId
                , method = "GET"
                , timeout = Nothing
                , tracker = Nothing
                , headers = []
                , body = Http.emptyBody
                , expect = Http.expectString GotHistory
                }
            )

        ProductIdChanged newValue ->
            ( { model | productId = newValue }, Cmd.none )

        GotHistory response ->
            ( { model | isLoading = False, history = Just response }, Cmd.none )


view : Model -> StyledDocument Msg
view model =
    { title = "Trace!"
    , body =
        if model.isLoading then
            [ h1 [] [ text "Cargando..." ]
            ]

        else
            let
                basicHeader =
                    [ h1 [] [ text "Welcome to the product tracer!" ]
                    , div []
                        [ input
                            [ value model.productId
                            , onInput ProductIdChanged
                            ]
                            []
                        , button [ onClick SearchClicked ] [ text "Search" ]
                        ]
                    ]
            in
            case model.history of
                Nothing ->
                    basicHeader

                Just response ->
                    case response of
                        Ok val ->
                            basicHeader
                                ++ [ div []
                                        [ Html.Styled.p [] [ text val ]
                                        ]
                                   ]

                        Err error ->
                            let
                                _ =
                                    Debug.log "HTTP ERROR: " error
                            in
                            basicHeader
                                ++ [ div [] [ text "An error occurred while trying to get the product history!" ]
                                   ]
    }
