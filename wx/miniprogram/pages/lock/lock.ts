/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-29 19:23:16
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-10-09 23:46:39
 * @FilePath: /coolcar/wx/miniprogram/pages/lock/lock.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// pages/lock/lock.ts

import {routing} from "../../utils/routings";
import {IAppOption} from "../../appoption";
import {TripService} from "../../service/trip"

const shareLocationKey = "share_location"


Page({

    /**
     * 页面的初始数据
     */
    carID: "",
    data: {
        shareLocation: false,
        avatarURL: '',
    },

    /**
     * 生命周期函数--监听页面加载
     */
    // 如果含有多个参数： opt:Record<'trip_id|is_vip', string>
    async onLoad(opt: Record<'car_id', string>) { // Record<'trip_id', string> 范型，对类型进行扩展和保护
        const o: routing.LockOpts = opt
        this.carID = o.car_id
        const userInfo = await getApp<IAppOption>().globalData.userInfo
        this.setData({
            avatarURL: userInfo.avatarUrl,
            shareLocation: wx.getStorageSync(shareLocationKey) || false,
        })
    },

    //  获取用户信息
    onGetUserInfo(e: any) {
        console.log('onGetUserInfo:===', e)
        const userInfo: WechatMiniprogram.UserInfo = e.detail.userInfo
        if (userInfo) {
            //  会报 resolveUserInfo 不存在，需要自己定义
            getApp<IAppOption>().resolveUserInfo(userInfo)
            this.setData({
                shareLocation: true,
            })
            wx.setStorageSync(shareLocationKey, true)
        }
    },

    // 是否分享
    onShareLocation(e: any) {
        this.data.shareLocation = e.detail.value
        wx.setStorageSync(shareLocationKey, this.data.shareLocation)
    },

    onUnlockTap() {
        // 点击开锁之前，需要获取用户的位置
        wx.getLocation({
            type: 'gcj02',
            success: async loc => {
                console.log('starting a trip：', {
                    location: {
                        latitude: loc.latitude,
                        longitude: loc.longitude
                    },
                    avatarUrl: this.data.shareLocation ? this.data.avatarURL : ''
                })
                // 不存在carID就不开锁
                if (!this.carID) {
                    console.error("no carID specified!")
                    return;
                }

                const trip = await TripService.CreateTrip({
                    start: loc,
                    carId: this.carID,
                })

                if (!trip.id) {
                    console.error("no tripID in response!", trip)
                    return
                }
                wx.showLoading({
                    title: "开锁中...",
                    mask: true,
                })
                // 模拟开锁跳转页面
                setTimeout(() => {
                    wx.redirectTo({
                        url: routing.driving({
                            trip_id: trip.id!, // 会报类型不对，后面加! 或前面做一次判断
                        }),
                        complete: () => {
                            wx.hideLoading()
                        }
                    })
                }, 2000)

            },
            fail: () => {
                wx.showToast({
                    icon: "none",
                    title: "请前往设置页授权位置信息"
                })
            }
        })
    }

})