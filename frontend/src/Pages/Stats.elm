module Pages.Stats exposing (..)

import Html.Styled exposing (Html, div, h1, text, ul, li, section)
import Html.Styled.Attributes exposing (class)
import Http
import Api.Endpoint exposing (getStats, getStatsResponseDecoder, GetStatsResponse, request)
import Models.Product exposing (Product)
import Models.Provider exposing (Provider)
import Models.PurchasedProduct exposing (PurchasedProduct)
import Utils exposing (StyledDocument)
import Browser.Navigation as Nav
import Routing exposing (Route(..))
import Url exposing (Url)


-- MODEL

type alias Model =
    { topProducts : List Product
    , topProviders : List Provider
    , topPurchasedProducts : List PurchasedProduct
    }

init : ( Model, Cmd Msg )
init =
    ( { topProducts = []
      , topProviders = []
      , topPurchasedProducts = []
      }
    , fetchStats
    )


-- UPDATE

type Msg
    = FetchStats
    | StatsReceived (Result Http.Error GetStatsResponse)

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchStats ->
            ( model, fetchStats )

        StatsReceived (Ok stats) ->
            let 
                _ = Debug.log "Datos recibidos del backend" stats
            in
            ( { model | 
                topProducts = stats.topProducts
              , topProviders = stats.topProviders
              , topPurchasedProducts = stats.topPurchasedProducts
              }
            , Cmd.none
            )

        StatsReceived (Err err) ->
            let 
                _ = Debug.log "Error en la petición HTTP" err
            in
            ( model, Cmd.none )


-- HTTP REQUEST

fetchStats : Cmd Msg
fetchStats =
    request
        { url = getStats
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , headers = []
        , body = Http.emptyBody
        , expect = Http.expectJson StatsReceived getStatsResponseDecoder
        }


-- VIEW

view : Model -> StyledDocument Msg
view model =
    let _ = Debug.log "Model en view" model in
    { title = "Statistics"
    , body =
        [ section []
            [ h1 [] [ text "Top 3 Productos con Mejor Rating" ]
            , ul [] (List.map viewProduct model.topProducts)
            , h1 [] [ text "Top 5 Proveedores Favoritos" ]
            , ul [] (List.map viewProvider model.topProviders)
            , h1 [] [ text "Top 10 Productos Más Comprados" ]
            , ul [] (List.map viewPurchasedProduct model.topPurchasedProducts)
            ]
        ]
    }

viewProduct : Product -> Html Msg
viewProduct product =
    li [] [ text (product.name ++ " - Rating: " ++ String.fromFloat product.averageRating) ]

viewProvider : Provider -> Html Msg
viewProvider provider =
    li [] [ text (provider.name ++ " - Popularidad: " ++ String.fromInt provider.popularity) ]

viewPurchasedProduct : PurchasedProduct -> Html Msg
viewPurchasedProduct product =
    li [] [ text (product.productName ++ " - Compras: " ++ String.fromInt product.purchases) ]

-- MAIN

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none