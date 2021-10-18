package controller

import (
	"errors"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Id       int64
}

type Driver struct {
	Username string
	Id       int64
}

// missing: split route into legs
type RouteInfo struct {
	Origin      Geolocation   `json:"Origin" binding:"required"`
	Destination Geolocation   `json:"Destination" binding:"required"`
	Duration    time.Duration `json:"Duration" binding:"required"`
	Cost        Money         `json:"Cost" binding:"required"`
}
type Money int64

type RideState string

const (
	RideCreated   = RideState("created")
	RideAccepted  = RideState("accepted")
	RideCancelled = RideState("cancelled")
)

var rideStates = map[string]RideState{
	"created":   RideCreated,
	"accepted":  RideAccepted,
	"cancelled": RideCancelled,
}

func ParseRideState(e string) (RideState, error) {
	_, found := rideStates[e]
	if !found {
		return RideState(""), errors.New("key not found")
	}
	return RideState(e), nil
}

type Ride struct {
	Id        int64
	User      int64
	Driver    int64
	CreatedAt time.Time
	State     RideState
	Route     RouteInfo
}

type Geolocation struct {
	X int `json:"X"`
	Y int `json:"Y"`
}
type UserService interface {
	GetUser(username string) User
	GetUserGeolocation(username string) Geolocation
}

type DriverService interface {
	GetDriver(username string) Driver
	GetDriverGeolocation(username string) Geolocation
	NotifyRide(driverId int64, ride Ride) error
}

type DriverServiceImpl struct{}

func (ds *DriverServiceImpl) GetDriver(username string) Driver {
	var y Driver
	return y
}
func (ds *DriverServiceImpl) GetDriverGeolocation(username string) Geolocation {
	var y Geolocation
	return y
}
func (ds *DriverServiceImpl) NotifyRide(driverId int64, ride Ride) error {
	log.Printf("Driver id %d has ride %v", driverId, ride)
	return nil
}

type MapService interface {
	FetchRouteInformation(params FetchRouteParams) RouteInfo
	GetAvailableDriversAround(loc Geolocation) []LocatedDriver
}

type RidesService interface {
	GetPendingRides() ([]Ride, error)
	AskForRide(userId int64, routeInfo RouteInfo) (int64, error)
	AcceptRide(rideId int64, driverId int64) error
	RejectRide(rideId int64)
	CancelRide(rideId int64)
	CompleteRide(rideId int64)
}

type RidesRepository struct {
	greatestId int64
	data       map[int64]Ride
}

func (rr *RidesRepository) CreateRide(ride Ride) (int64, error) {
	rr.greatestId = rr.greatestId + 1
	ride.Id = rr.greatestId
	rr.data[ride.Id] = ride
	return ride.Id, nil
}

type RidesServiceImpl struct {
	RidesRepository *RidesRepository
}

func NewRidesRepository() *RidesRepository {
	return &RidesRepository{greatestId: 0, data: make(map[int64]Ride)}
}

func (rr *RidesRepository) GetPendingRides() ([]Ride, error) {
	var res []Ride
	for _, v := range rr.data {
		if v.State == RideCreated {
			res = append(res, v)
		}
	}
	return res, nil
}

func (rr *RidesRepository) AcceptRide(rideId int64, driverId int64) error {
	v, exists := rr.data[rideId]
	if !exists {
		return errors.New("ride does not exist")
	}
	if v.Driver != 0 || v.State != RideCreated {
		return errors.New("driver is assigned or ride was already accepted")
	}
	v.Driver = driverId
	v.State = RideAccepted
	rr.data[rideId] = v
	return nil
}

func (rs *RidesServiceImpl) GetPendingRides() ([]Ride, error) {
	return rs.RidesRepository.GetPendingRides()
}

func (rs *RidesServiceImpl) AcceptRide(rideId int64, driverId int64) error {
	return rs.RidesRepository.AcceptRide(rideId, driverId)
}
func (rs *RidesServiceImpl) RejectRide(rideId int64)   { panic("not implemented!") } // TODO
func (rs *RidesServiceImpl) CancelRide(rideId int64)   { panic("not implemented!") } // TODO
func (rs *RidesServiceImpl) CompleteRide(rideId int64) { panic("not implemented!") } // TODO

func (rs *RidesServiceImpl) AskForRide(userId int64, routInfo RouteInfo) (int64, error) {
	rideId, err := rs.RidesRepository.CreateRide(
		Ride{
			User: userId, State: RideCreated, CreatedAt: time.Now(), Route: routInfo,
		},
	)
	if err != nil {
		return 0, err
	}
	return rideId, nil
}

type NuberController struct {
	UserService   UserService
	DriverService DriverService
	RidesService  RidesService
	MapService    MapService
}

type FetchRouteParams struct {
	Origin      Geolocation `json:"Origin" bindind:"required"`
	Destination Geolocation `json:"Destination" bindind:"required"`
	Time        time.Time   `json:"Time"`
}

func (ctl *NuberController) GetRouteCost(ctx *gin.Context, route RouteInfo) Money {
	return Money(math.Max(7, route.Duration.Minutes()))
}

func (ctl *NuberController) AskForRide(ctx *gin.Context) {
	params := new(RouteInfo)
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(400, gin.H{"error": "BadRequest"})
		return
	}
	userId := ctx.GetHeader("USER-TOKEN") // FIXME get userId from token
	if userId == "" {
		ctx.JSON(400, gin.H{"error": "userId does not exist"})
		return
	}
	userIdd, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		panic("cu")
	}
	rideId, err := ctl.RidesService.AskForRide(userIdd, *params)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "could not handle request"})
		return
	}

	ctx.JSON(200, gin.H{"RideId": rideId})
}

func (ctl *NuberController) AcceptRide(ctx *gin.Context) {
	type AcceptRideParams struct {
		RideId int64 `json:"RideId"`
	}
	params := new(AcceptRideParams)
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Println(err)
		ctx.JSON(400, gin.H{"error": "BadRequest"})
		return
	}
	driverId := ctx.GetHeader("USER-TOKEN") // FIXME get driverId from token
	if driverId == "" {
		ctx.JSON(400, gin.H{"error": "userId does not exist"})
		return
	}
	driverIdd, err := strconv.ParseInt(driverId, 10, 64)
	if err != nil {
		panic("cu")
	}
	err = ctl.RidesService.AcceptRide(params.RideId, driverIdd)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "could not accept ride"})
		return
	}
	ctx.JSON(200, gin.H{"status": "OK"})

}

type MapServiceImpl struct {
}

func (m *MapServiceImpl) FetchRouteInformation(params FetchRouteParams) RouteInfo {
	var timE time.Time
	if timE == params.Time {
		timE = time.Now()
	}
	ri := RouteInfo{
		Origin:      params.Origin,
		Destination: params.Destination,
		Duration:    time.Duration(6000),
	}
	return ri
}

type LocatedDriver struct {
	Username string
	ID       int64
	loc      Geolocation
}

func (m *MapServiceImpl) GetAvailableDriversAround(loc Geolocation) []LocatedDriver {
	return []LocatedDriver{
		{"joser", 64, Geolocation{10, 10}},
		{"Manguito", 67, Geolocation{16, 17}},
	}
}

func (ctl *NuberController) FetchRoute(ctx *gin.Context) {
	params := new(FetchRouteParams)
	if err := ctx.ShouldBindJSON(params); err != nil {
		log.Println(err)
		ctx.JSON(400, gin.H{"error": "BadRequest"})
		return
	}
	ri := ctl.MapService.FetchRouteInformation(*params)
	cost := ctl.GetRouteCost(ctx, ri) // FIXME pricing might be a separate microservice
	ri.Cost = cost
	ctx.JSON(200, ri)
}
