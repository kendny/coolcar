import {auth} from "./proto_gen/auth/auth_pb";
import camelcaseKeys from "camelcase-keys";

export namespace CoolCar {
    const serverAddr = "http://localhost:8080"
    const AUTH_ERR = "AUTH_ERR"

    // 用户认证 信息
    const authData = {
        token: "",
        expiryMs: 0
    }

    // 封装接口 REQ 对  RequestOption 进行限制
    interface RequestOption<REQ, RES> {
        method: "GET" | "PUT" | "POST" | "DELETE"
        path: string
        data?: REQ // ? data可选项
        respMarshaller: (r: object) => RES
    }

    export interface AuthOption {
        attachAuthHeader: boolean // 是否加上token
        retryOnAuthError: boolean // 重试
    }

    export async function SendRequestWithAuthRetry<REQ, RES>(o: RequestOption<REQ, RES>, a?: AuthOption): Promise<RES> {
        const authOpt = a || {
            attachAuthHeader: true,
            retryOnAuthError: true
        }

        try {
            await login() // 重新登陆
            return sendRequest(o, authOpt)
        } catch (err) {
            if (err === AUTH_ERR && authOpt.retryOnAuthError) {
                // 重置
                authData.token = ""
                authData.expiryMs = 0
                return SendRequestWithAuthRetry(o, {
                    attachAuthHeader: true,
                    retryOnAuthError: false
                })
            } else {
                throw err
            }
        }
    }

    export async function login() {
        const reqTimeMs = Date.now()
        // 如果token有效，不登陆，直接发业务请求
        if (authData.token && authData.expiryMs >= reqTimeMs) {
            return
        }
        const wxResp = await wxLogin()
        const resp = await sendRequest<auth.v1.ILoginRequest, auth.v1.ILoginResponse>({
            method: "POST",
            path: "/v1/auth/login",
            data: {
                code: wxResp.code,
            },
            respMarshaller: auth.v1.LoginResponse.fromObject,
        }, {
            attachAuthHeader: false,
            retryOnAuthError: false
        })

        authData.token = resp.accessToken! // ！一定有token
        authData.expiryMs = reqTimeMs + resp.expiresIn! * 1000
    }

    function sendRequest<REQ, RES>(o: RequestOption<REQ, RES>, a: AuthOption): Promise<RES> {
        const authOpt = a || {
            attachAuthHeader: true,
        }
        // 将微信的wx.request 转成Promise
        return new Promise((resolve, reject) => {
            //请求头中添加过期时间
            const header: Record<string, any> = {}
            if (authOpt.attachAuthHeader) {
                if (authData.token && authData.expiryMs >= Date.now()) {
                    header.authorization = "Bearer " + authData.token
                } else {
                    // token过期 抛出登陆异常的错误
                    reject(AUTH_ERR)
                    return
                }
            }

            wx.request({
                url: serverAddr + o.path,
                method: o.method,
                // @ts-ignore
                // todo...?
                data: o.data,
                header,
                success: res => {
                    // 对请求结果的处理
                    if (res.statusCode === 401) {
                        reject(AUTH_ERR)
                    } else if (res.statusCode >= 400) {
                        reject(res)
                    } else {
                        resolve(o.respMarshaller(camelcaseKeys(res.data as object, {
                            deep: true
                        })))
                    }

                },
                fail: reject,
            })
        })
    }

    function wxLogin(): Promise<WechatMiniprogram.LoginSuccessCallbackResult> {
        return new Promise((resolve, reject) => {
            wx.login({
                success: resolve,
                fail: reject
            })
        })
    }

}