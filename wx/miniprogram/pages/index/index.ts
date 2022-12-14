// index.ts
// 获取应用实例
import {IAppOption} from "../../appoption";
import {routing} from "../../utils/routings";
import {ProfileService} from "../../service/profile";
import {rental} from "../../service/proto_gen/rental/rental_pb";
import {TripService} from "../../service/trip";

const app = getApp<IAppOption>()

interface Marker {
    iconPath: string
    id: number
    latitude: number
    longitude: number
    width: number
    height: number
}

// const defaultAvatar = '/resources/car.png'
const initialLat = 29.761267625855936
const initialLng = 121.87264654736123

Page({
    isPageShowing: false,
    data: {
        avatarURL: '',
        setting: {
            skew: 0,
            rotate: 0,
            showLocation: true,
            showScale: true,
            subKey: '',
            layerStyle: -1,
            enableZoom: true,
            enableScroll: true,
            enableRotate: false,
            showCompass: false,
            enable3D: false,
            enableOverlooking: false,
            enableSatellite: false,
            enableTraffic: false,
        },
        location: {
            latitude: initialLat,
            longitude: initialLng,
        },
        scale: 16,
        markers: [] as Marker[],
    },

    async onLoad() {
        // const userInfo:WechatMiniprogram.UserInfo = await getApp<IAppOption>().globalData.userInfo
        // this.setData({
        //    // @ts-ignore
        //   avatarURL: userInfo.avatarUrl,
        // })
        app.globalData.userInfo.then(userInfo => {
            this.setData({
                userInfo,
                // hasUserInfo:true
            })
            console.log('userInfo:===', userInfo)
        })

        // @ts-ignore
        if (wx.getUserProfile) {
            this.setData({
                canIUseGetUserProfile: true
            })
        }
        this.setData({
            markers: [
                {
                    iconPath: '/resources/car.png',
                    id: 0,
                    latitude: 23.0999994,
                    longitude: 113.324520,
                    width: 50,
                    height: 50,
                },
                {
                    iconPath: '/resources/car.png',
                    id: 1,
                    latitude: 23.0999994,
                    longitude: 114.324520,
                    width: 50,
                    height: 50,
                }
            ]
        })
    },

    onShow() {
        this.isPageShowing = true
    },
    onHide() {
        this.isPageShowing = false
    },

    async onScanClicked() {
        //  查看是否有行程
        const trips = await TripService.GetTrips(rental.v1.TripStatus.IN_PROGRESS)
        if ((trips.trips?.length || 0) > 0) {
            await this.selectComponent('#tripModal').showModal()
            wx.navigateTo({
                url: routing.driving({
                    trip_id: trips.trips![0].id!,
                }),
            })
            return
        }
        wx.scanCode({
            success: async res => {
                console.log(res)
                // TODO... get car id from scan result
                const carID = "car123"
                const lockURL = routing.lock({
                    car_id: carID
                })
                // 获取认证信息
                const prof = await ProfileService.getProfile()
                if (prof.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
                    wx.navigateTo({
                        url: lockURL,
                    })
                } else {
                    await this.selectComponent("#licModal").showModal()
                    wx.navigateTo({
                        url: routing.register({
                            redirectURL: lockURL
                        })
                    })
                }

                wx.navigateTo({
                    url: '/pages/register/register'
                })
            },
            fail: console.error,
        })
    },
    // 事件处理函数
    bindViewTap() {
        wx.navigateTo({
            url: '../logs/logs',
        })
    },
    onMyTripsTap() {
        // 查看个人行程
        wx.navigateTo({
            url: "/pages/mytrips/mytrips"
        })
    },
    onMyLocationTap() {
        wx.getLocation({
            type: 'gcj02',
            success: res => {
                console.log('onMyLocationTap:==', res)
                this.setData({
                    location: {
                        latitude: res.latitude,
                        longitude: res.longitude
                    }
                })
            },
            fail: (error) => {
                console.log("getLocation:fail", error)
                wx.showToast({
                    icon: 'none',
                    title: "请前往设置页授权",
                })
            }
        })
    },

    moveCars() {
        console.log("开始移动：==")
        const map = wx.createMapContext("map")
        const dest = {
            latitude: 23.0999994,
            longitude: 113.234520,
        }
        const moveCar = () => {
            dest.latitude += 1,
                dest.longitude += 1
            map.translateMarker({
                destination: {
                    latitude: dest.latitude,
                    longitude: dest.longitude,
                },
                markerId: 0,
                autoRotate: false,
                rotate: 0,
                duration: 5000,
                animationEnd: () => {
                    if (this.isPageShowing) {
                        moveCar()
                    }
                }
            })
        }
        moveCar()
    },
    getUserProfile() {
        // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认，开发者妥善保管用户快速填写的头像昵称，避免重复弹窗
        wx.getUserProfile({
            desc: '展示用户信息', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
            success: (res) => {
                console.log(res)
                this.setData({
                    userInfo: res.userInfo,
                    hasUserInfo: true
                })
            }
        })
    },
    getUserInfo(e: any) {
        // 不推荐使用getUserInfo获取用户信息，预计自2021年4月13日起，getUserInfo将不再弹出弹窗，并直接返回匿名的用户个人信息
        // e 为any 类型，后面所获取的值的类型将不受到保护，解决方法，如果能确定所获取值的类型，则可以手动限定类型
        const userInfo: WechatMiniprogram.UserInfo = e.detail.userInfo;
        // app.globalData.userInfo = userInfo;
        app.resolveUserInfo(userInfo)
        console.log(e)
        // this.setData({
        //   userInfo: userInfo,
        //   hasUserInfo: true
        // })
    }
})
