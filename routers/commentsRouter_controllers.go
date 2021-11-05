package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "MigrateBlob",
            Router: "/migrate/blob",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "CrawlBeeDoc",
            Router: "/reptile/beego",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "CrawlFish",
            Router: "/reptile/fish",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "CrawlGoBlob",
            Router: "/reptile/go",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "CrawlMovie",
            Router: "/reptile/movie",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"] = append(beego.GlobalControllerRouter["beego_reptile_ty/controllers:CrawlController"],
        beego.ControllerComments{
            Method: "CrawlGoBlobSync",
            Router: "/reptile/sync/go",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
