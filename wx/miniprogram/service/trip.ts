import {rental} from "./proto_gen/rental/rental_pb";
import {CoolCar} from "./request";

export namespace TripService {
    export function CreateTrip(req:rental.v1.ICreateTripRequest):Promise<rental.v1.ICreateTripResponse> {
        return CoolCar.SendRequestWithAuthRetry({
            method:"POST",
            path: "/v1/trip",
            data: req,
            respMarshaller: rental.v1.CreateTripResponse.fromObject
        })
    }
}