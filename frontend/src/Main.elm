module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Pages.Home as HomePage
import Pages.NotFound as NotFoundPage
import Pages.Report as ReportPage
import Pages.Stats as StatsPage
import Pages.Trace as TracePage
import Routing exposing (Route(..))
import Url exposing (Url)
import Utils


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
    | NotFound NotFoundPage.Model
    | Report ReportPage.Model
    | Stats StatsPage.Model
    | Trace TracePage.Model


type alias Model =
    { key : Nav.Key
    , model : PageModels
    }


init : a -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        route =
            Routing.parseUrl url

        ( innerModel, initialEffect ) =
            case route of
                Routing.Home ->
                    Tuple.mapBoth Home (Cmd.map HomeMsg) HomePage.init

                Routing.NotFound ->
                    Tuple.mapBoth NotFound (Cmd.map NotFoundMsg) NotFoundPage.init

                Routing.Report ->
                    Tuple.mapBoth Report (Cmd.map ReportMsg) ReportPage.init

                Routing.Statistics ->
                    Tuple.mapBoth Stats (Cmd.map StatsMsg) StatsPage.init

                Routing.Trace ->
                    Tuple.mapBoth Trace (Cmd.map TraceMsg) TracePage.init
    in
    ( Model navKey innerModel, initialEffect )


subscriptions : Model -> Sub Msg
subscriptions model =
    case model.model of
        Home inner ->
            Sub.map HomeMsg (HomePage.subscriptions inner)

        NotFound inner ->
            Sub.map NotFoundMsg (NotFoundPage.subscriptions inner)

        Report inner ->
            Sub.map ReportMsg (ReportPage.subscriptions inner)

        Stats inner ->
            Sub.map StatsMsg (StatsPage.subscriptions inner)

        Trace inner ->
            Sub.map TraceMsg (TracePage.subscriptions inner)


type Msg
    = UrlChanged Url
    | LinkClicked Browser.UrlRequest
    | HomeMsg HomePage.Msg
    | NotFoundMsg NotFoundPage.Msg
    | ReportMsg ReportPage.Msg
    | StatsMsg StatsPage.Msg
    | TraceMsg TracePage.Msg


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlChanged url ->
            init {} url model.key

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

                _ ->
                    ( model, Cmd.none )

        NotFoundMsg innerMsg ->
            case model.model of
                NotFound innerModel ->
                    let
                        ( newModel, effect ) =
                            NotFoundPage.update innerMsg innerModel
                                |> Tuple.mapBoth NotFound (Cmd.map NotFoundMsg)
                    in
                    ( { model | model = newModel }, effect )

                _ ->
                    ( model, Cmd.none )

        ReportMsg innerMsg ->
            case model.model of
                Report innerModel ->
                    let
                        ( newModel, effect ) =
                            ReportPage.update innerMsg innerModel
                                |> Tuple.mapBoth Report (Cmd.map ReportMsg)
                    in
                    ( { model | model = newModel }, effect )

                _ ->
                    ( model, Cmd.none )

        StatsMsg innerMsg ->
            case model.model of
                Stats innerModel ->
                    let
                        ( newModel, effect ) =
                            StatsPage.update innerMsg innerModel
                                |> Tuple.mapBoth Stats (Cmd.map StatsMsg)
                    in
                    ( { model | model = newModel }, effect )

                _ ->
                    ( model, Cmd.none )

        TraceMsg innerMsg ->
            case model.model of
                Trace innerModel ->
                    let
                        ( newModel, effect ) =
                            TracePage.update innerMsg innerModel
                                |> Tuple.mapBoth Trace (Cmd.map TraceMsg)
                    in
                    ( { model | model = newModel }, effect )

                _ ->
                    ( model, Cmd.none )


view : Model -> Browser.Document Msg
view model =
    let
        styledDoc =
            case model.model of
                Home inner ->
                    Utils.mapMsg HomeMsg (HomePage.view inner)

                NotFound inner ->
                    Utils.mapMsg NotFoundMsg (NotFoundPage.view inner)

                Report inner ->
                    Utils.mapMsg ReportMsg (ReportPage.view inner)

                Stats inner ->
                    Utils.mapMsg StatsMsg (StatsPage.view inner )



                Trace inner ->
                    Utils.mapMsg TraceMsg (TracePage.view inner)
    in
    Utils.toUnstyled styledDoc
