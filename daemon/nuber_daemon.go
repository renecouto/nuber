package daemon

import (
	"time"

	"github.com/renecouto/nuber/controller"
)

type NuberDaemon struct {
	UserService   controller.UserService
	DriverService controller.DriverService
	RidesService  controller.RidesService
	MapService    controller.MapService
}

func (nd *NuberDaemon) Run() {
	d, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}
	for {
		nd.loop()
		time.Sleep(d)
	}
}
func (nd *NuberDaemon) loop() {
	rides, err := nd.RidesService.GetPendingRides()
	if err != nil {
		panic(err)
	}
	for _, ride := range rides {
		candidateDrivers := nd.MapService.GetAvailableDriversAround(ride.Route.Origin)
		for _, c := range candidateDrivers {
			err := nd.DriverService.NotifyRide(c.ID, ride)
			if err != nil {
				panic(err)
			}
		}
	}
}
