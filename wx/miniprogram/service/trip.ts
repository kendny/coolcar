import {rental} from "./proto_gen/rental/rental_pb";
import {CoolCar} from "./request";

export namespace TripService {
    // 这个接口传参比较复杂，所以借用 ICreateTripRequest
    export function CreateTrip(req:rental.v1.ICreateTripRequest):Promise<rental.v1.ITripEntity> {
        return CoolCar.SendRequestWithAuthRetry({
            method:"POST",
            path: "/v1/trip",
            data: req,
            respMarshaller: rental.v1.TripEntity.fromObject
        })
    }

    export function GetTrip(id:string):Promise<rental.v1.ITrip>{
        return CoolCar.SendRequestWithAuthRetry({
            method: "GET",
            path: `/v1/trip/${encodeURIComponent(id)}`,
            respMarshaller: rental.v1.Trip.fromObject,
        })
    }

    export function GetTrips(s?:rental.v1.TripStatus):Promise<rental.v1.IGetTripsResponse> {
        return CoolCar.SendRequestWithAuthRetry({
            method: "GET",
            path: s? `/v1/trips?status=${s}`: "/v1/trips",
            respMarshaller: rental.v1.GetTripsResponse.fromObject
        })
    }

    export function UpdateTripPos(id:string, loc?:rental.v1.ILocation) {
        return updateTrip({
            id,
            current: loc
        })
    }

    export function FinishTrip(id:string) {
        return updateTrip({
            id,
            endTrip: true
        })
    }


    function updateTrip(r:rental.v1.IUpdateTripRequest): Promise<rental.v1.ITrip> {
        if(!r.id) {
            return Promise.reject("must specify id")
        }
        return CoolCar.SendRequestWithAuthRetry({
            method: "PUT",
            path: `/v1/trip/${encodeURIComponent(r.id)}`,
            data: r,
            respMarshaller: rental.v1.Trip.fromObject
        })
    }
}