import {rental} from "./proto_gen/rental/rental_pb";
import {CoolCar} from "./request";

export namespace ProfileService{
    export function getProfile():Promise<rental.v1.IProfile>{
        return CoolCar.SendRequestWithAuthRetry({
            method: "GET",
            path: "/v1/profile",
            respMarshaller: rental.v1.Profile.fromObject
        })
    }

    export function submitProfile(req:rental.v1.IIdentity): Promise<rental.v1.IProfile>{
        return CoolCar.SendRequestWithAuthRetry({
            method: "POST",
            path: "/v1/profile",
            data: req,
            respMarshaller: rental.v1.Profile.fromObject
        })
    }

    export function clearProfile():Promise<rental.v1.IProfile> {
        return CoolCar.SendRequestWithAuthRetry({
            method: "DELETE",
            path: "/v1/profile",
            respMarshaller: rental.v1.Profile.fromObject
        })
    }
}