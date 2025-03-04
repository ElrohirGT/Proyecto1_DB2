module Pages.Trace exposing (..)

import Api.Endpoint exposing (GetHistoryResponse, getHistory, getHistoryResponseDecoder, request)
import Css exposing (alignItems, backgroundColor, bold, border, center, color, column, displayFlex, flexDirection, fontWeight, height, justifyContent, none, padding2, px, row, textAlign, textDecoration, verticalAlign, width, zero)
import Dict
import Html.Styled exposing (a, button, div, h1, h2, input, li, ol, p, pre, text)
import Html.Styled.Attributes exposing (css, placeholder, type_, value)
import Html.Styled.Events exposing (onClick, onInput)
import Http
import Json.Decode
import Routing exposing (goToHome)
import Theme exposing (colors)
import Utils exposing (StyledDocument)



-- MODEL


type alias Producers =
    List String


type alias HistoryState =
    Result Http.Error Producers


type alias APIResponse =
    Result Http.Error GetHistoryResponse


type alias Model =
    { productId : String
    , isLoading : Bool
    , history : Maybe HistoryState
    }


exampleResponse : String
exampleResponse =
    "{\"Values\":[{\"Nodes\":[{\"Id\":4152,\"ElementId\":\"4:5eaea7dc-9320-449f-993e-45d993464520:4152\",\"Labels\":[\"Provider\"],\"Props\":{\"created_at\":\"05/08/2000\",\"email\":\"ksenchenkoi8@issuu.com\",\"id\":\"657\",\"name\":\"Feeney-Ward\",\"owner\":\"Kit Senchenko\"}},{\"Id\":2776,\"ElementId\":\"4:5eaea7dc-9320-449f-993e-45d993464520:2776\",\"Labels\":[\"Product\"],\"Props\":{\"brand\":\"Vertex\",\"category\":\"Electronics\",\"id\":\"5\",\"name\":\"Smartphone\",\"weight\":\"3.24\"}}],\"Relationships\":[{\"Id\":1155175503443791928,\"ElementId\":\"5:5eaea7dc-9320-449f-993e-45d993464520:1155175503443791928\",\"StartId\":4152,\"StartElementId\":\"4:5eaea7dc-9320-449f-993e-45d993464520:4152\",\"EndId\":2776,\"EndElementId\":\"4:5eaea7dc-9320-449f-993e-45d993464520:2776\",\"Type\":\"PRODUCES\",\"Props\":{\"max_quantity\":31,\"since\":{},\"speed\":2}}]}],\"Keys\":[\"p1\"]}"


init : ( Model, Cmd Msg )
init =
    ( Model "5" False Nothing, Cmd.none )



-- UPDATE


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
            -- Local Parsing
            -- Currently doesn't work because of emulating Http.error
            -- let
            --     response =
            --         Decode.decodeString getHistoryResponseDecoder exampleResponse
            -- in
            -- ( { model | history = Just (Ok response) }, Cmd.none )
            -- API Parsing
            ( { model | isLoading = True }
            , request
                { url = getHistory model.productId
                , method = "GET"
                , timeout = Nothing
                , tracker = Nothing
                , headers = []
                , body = Http.emptyBody
                , expect = Http.expectJson GotHistory getHistoryResponseDecoder
                }
            )

        ProductIdChanged newValue ->
            ( { model | productId = newValue }, Cmd.none )

        GotHistory apiResponse ->
            let
                producers : HistoryState
                producers =
                    let
                        responseToProducerMapper : GetHistoryResponse -> Producers
                        responseToProducerMapper historyResponse =
                            historyResponse.values
                                |> List.map
                                    (\value ->
                                        value.nodes
                                            |> List.filterMap
                                                (\n ->
                                                    if List.any (\l -> l == "Provider") n.labels then
                                                        Dict.get "name" n.props
                                                            |> Maybe.andThen
                                                                (\nameValue ->
                                                                    nameValue
                                                                        |> Json.Decode.decodeValue Json.Decode.string
                                                                        |> Result.toMaybe
                                                                )

                                                    else
                                                        Nothing
                                                )
                                    )
                                |> List.concat
                    in
                    Result.map responseToProducerMapper apiResponse
            in
            ( { model | isLoading = False, history = Just producers }, Cmd.none )



-- VIEW


view : Model -> StyledDocument Msg
view model =
    { title = "Trace!"
    , body =
        let
            header =
                [ h1
                    [ css
                        [ color (Css.hex colors.secondary)
                        , textAlign center
                        ]
                    ]
                    [ text "¡Bienvenido al Trazador de Productos!" ]
                , div
                    [ css
                        [ displayFlex
                        , flexDirection row
                        , Css.property "gap" "0.5rem"
                        ]
                    ]
                    [ input
                        [ value model.productId
                        , onInput ProductIdChanged
                        , type_ "number"
                        , placeholder "ID del producto:"
                        ]
                        []
                    , button
                        [ onClick SearchClicked
                        , css Theme.btnStyles
                        ]
                        [ text "Buscar" ]
                    ]
                , a
                    [ goToHome
                    , css Theme.btnStyles
                    ]
                    [ text "Regresar a Inicio" ]
                ]

            mainContent =
                if model.isLoading then
                    [ h1 [ css [ color (Css.hex colors.warning) ] ] [ text "Cargando..." ] ]

                else
                    case model.history of
                        Nothing ->
                            header

                        Just response ->
                            case response of
                                Ok val ->
                                    header
                                        ++ [ div []
                                                [ h2
                                                    [ css
                                                        [ color (Css.hex colors.secondary)
                                                        ]
                                                    ]
                                                    [ text "Los proveedores son:" ]
                                                , ol [] (List.map (\v -> li [ css [ color (Css.hex "#ffffff") ] ] [ text v ]) val)
                                                ]
                                           ]

                                Err error ->
                                    let
                                        _ =
                                            Debug.log "HTTP ERROR: " error
                                    in
                                    case error of
                                        Http.BadBody debugError ->
                                            header
                                                ++ [ pre
                                                        [ css
                                                            [ color (Css.hex colors.error)
                                                            , fontWeight bold
                                                            ]
                                                        ]
                                                        [ text debugError ]
                                                   ]

                                        Http.BadStatus status ->
                                            if status == 404 then
                                                header
                                                    ++ [ div
                                                            [ css
                                                                [ color (Css.hex colors.error)
                                                                , fontWeight bold
                                                                ]
                                                            ]
                                                            [ text
                                                                (String.join " "
                                                                    [ "No product with ID:"
                                                                    , model.productId
                                                                    , "exists!"
                                                                    ]
                                                                )
                                                            ]
                                                       ]

                                            else
                                                header
                                                    ++ [ div
                                                            [ css
                                                                [ color (Css.hex colors.error)
                                                                , fontWeight bold
                                                                ]
                                                            ]
                                                            [ text "An error occurred while trying to get the product history!" ]
                                                       ]

                                        _ ->
                                            header
                                                ++ [ div [ css [ color (Css.hex colors.error) ] ] [ text "An error occurred while trying to get the product history!" ]
                                                   ]
        in
        [ div
            [ css
                (Theme.divBackgroundStyles
                    ++ [ flexDirection column
                       , Css.property "gap" "1rem"
                       ]
                )
            ]
            mainContent
        ]
    }
