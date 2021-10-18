package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/renecouto/nuber/controller"
	"github.com/renecouto/nuber/daemon"
)

func runWebsite(ctl *controller.NuberController) {

	r := gin.Default()
	// r.LoadHTMLGlob("web/templates/**/*")
	// var itemsRepo models.ItemsRepository

	// if setup_data {
	// 	SetupData(itemsRepo)
	// }

	bindApiRoutes(r, ctl)
	r.Run()
}

func bindApiRoutes(r *gin.Engine, ctl *controller.NuberController) {

	r.POST("/routes/_fetch", ctl.FetchRoute)
	r.POST("/rides/_ask", ctl.AskForRide)
	r.POST("/rides/_accept", ctl.AcceptRide)
	// r.POST("/rides/_reject", ctl.RejectRide)
}

func runDaemon(d *daemon.NuberDaemon) {
	d.Run()
}

func seconds(t int) time.Duration {
	d, err := time.ParseDuration(fmt.Sprintf("%ds", t))
	if err != nil {
		panic(err)
	}
	return d
}

func main() {
	var program = flag.String("program", "both", "web for website, daemon for daemon, both for both")
	var userService controller.UserService
	driverService := controller.DriverServiceImpl{}
	ridesRepo := controller.NewRidesRepository()
	ridesService := controller.RidesServiceImpl{RidesRepository: ridesRepo}
	mapService := controller.MapServiceImpl{}
	ctl := controller.NuberController{
		UserService:   userService,
		DriverService: &driverService,
		RidesService:  &ridesService,
		MapService:    &mapService,
	}
	d := daemon.NuberDaemon{
		UserService:   userService,
		DriverService: &driverService,
		RidesService:  &ridesService,
		MapService:    &mapService,
	}
	switch *program {
	case "web":
		runWebsite(&ctl)
		break
	case "daemon":
		runDaemon(&d)
		break
	case "both":
		go runDaemon(&d)
		go runWebsite(&ctl)
		for {
			log.Println("sleeping")
			time.Sleep(seconds(20))
		}

	default:
		panic("invalid program choice")

	}

}
