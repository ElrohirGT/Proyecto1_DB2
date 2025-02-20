module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Html exposing (Html, h1, text)
import Pages.Home as HomePage
import Routing exposing (Route(..))
import Url exposing (Url)


main : Program String Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }


type PageModels
    = Home HomePage.Model


type alias Model =
    { key : Nav.Key
    , basePath : String
    , model : PageModels
    }


init : String -> Url -> Nav.Key -> ( Model, Cmd Msg )
init basePath url navKey =
    let
        route =
            Routing.parseUrl basePath url

        ( innerModel, initialEffect ) =
            case route of
                Routing.Home ->
                    Tuple.mapFirst Home HomePage.init

                _ ->
                    Tuple.mapFirst Home HomePage.init
    in
    ( Model navKey basePath innerModel, initialEffect )


subscriptions : Model -> Sub Msg
subscriptions model =
    case model.model of
        Home inner ->
            HomePage.subscriptions inner


type Msg
    = UrlChanged Url
    | LinkClicked Browser.UrlRequest
    | HomeMsg HomePage.Msg


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlChanged url ->
            init model.basePath url model.key

        LinkClicked req ->
            case req of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        HomeMsg innerMsg ->
            case model.model of
                Home innerModel ->
                    let
                        ( newModel, effect ) =
                            HomePage.update innerMsg innerModel
                                |> Tuple.mapBoth Home (Cmd.map HomeMsg)
                    in
                    ( { model | model = newModel }, effect )


view : Model -> Browser.Document Msg
view model =
    let
        styledDoc =
            case model.model of
                Home inner ->
                    HomePage.view inner
    in
    styledDoc
