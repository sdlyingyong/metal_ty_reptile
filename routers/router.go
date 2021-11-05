// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego_reptile_ty/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	//some times not delete exist router file and make exception
	//beego.Include(&controllers.CrawlController{})

	beego.Router("/reptile/sync/go", &controllers.CrawlController{}, "get:CrawlGoBlobSync")
	beego.Router("/reptile/go", &controllers.CrawlController{}, "get:CrawlGoBlob")

	beego.Router("/migrate/blob", &controllers.CrawlController{}, "get:MigrateBlob")
	beego.Router("/migrate/sync/blob", &controllers.CrawlController{}, "get:MigrateBlobSync")

	//demo
	//chan
	beego.Router("/demo/chan", &controllers.DemoController{}, "get:ChanDemo")
}
