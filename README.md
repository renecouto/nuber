# Nuber - Not UBER
## Goals
- Deeper experience with golang and microservices
- Get insights into ETL and online can mix
- Use grpc (maybe graphql)

## Non-goals
- prod-like environment
- beauty

## Architecture
- Users
- Drivers
- Payments
- Map
- Routes
- Rides
- Backend For Frontend

## Main App Flow
- User searches ride
    - Routes.FetchRoutes(source, destination, time) (cost, estimated time)
- User Asks for ride
    - Rides.AskForRide(ride_id, route)
    - Rides service replies with a 200 and asynchronously sends notifications to drivers
        - Fetch active drivers close to the user by geolocation, increasing the radius as needed
            - which service? drivers? map?
        - Driver can accept or decline the ride
    - User waits until a driver accepts the ride
- User waits for driver to arrive
- Driver starts the ride
- Driver or user completes the ride
- Payment is debited to drivers account

## TODOS
- make daemon not call the service methods
    - use the api?
- should drivers be users?
- where do i get the info that drivers are in a ride? propagte to maps and drivers apis?
- users and drivers crud


## Commands
- Rides.AskForRide {
    pickup location
    destination
    cost
}
- Rides.CancelRide
- Rides.CompleteRide
- Users.SendCurrentPosition
- Rides.OfferRideToDriver
- Drivers.AcceptRide
- Drivers.CancelRide
- Drivers.CompleteRide
- Drivers.StartWorking
- Drivers.StopWorking
- Users.SendMessage
- Drivers.SendMessage
- Drivers.SendCurrentPosition

## Sagas/Queries
- link ride with driver
- fetch driver estimated time to arrive
- fetch ride entire path
- fetch ride distance, cost, estimated time
- update leg cost, estimated time

## Displays
- map with current positions
- search for address

## Work-arounds
- Since I don't want to waste too much time learning about maps, I'll implement my own
### Maps
- first let's make it a grid of size 16x16, all streets connecting to each other
- we can make the sections real time vary with some random number
- if we have never gone through a specific leg, we can assume it has a leg time equal to the speed limit (?)


