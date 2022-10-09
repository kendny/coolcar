/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-29 19:23:16
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-10-09 23:46:39
 * @FilePath: /coolcar/wx/miniprogram/pages/lock/lock.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// pages/lock/lock.ts

const shareLocationKey = "share_location"


Page({

  /**
   * 页面的初始数据
   */
  data: {
    shareLocation: false,
    avatarURL: '',
  },

  /**
   * 生命周期函数--监听页面加载
   */
  async onLoad() {
    const userInfo = await getApp<IAppOption>().globalData.userInfo
    this.setData({
      // @ts-ignore
      avatarURL: userInfo.avatarUrl,
      shareLocation: wx.getStorageSync(shareLocationKey) || false,
    })
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

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
      success: loc => {
        console.log('starting a trip：', {
          location: {
            latitude: loc.latitude,
            longitude: loc.longitude
          },
          avatarUrl: this.data.shareLocation ? this.data.avatarURL : ''
        })
        // 模拟开锁跳转页面
        setTimeout(() => {
          wx.redirectTo({
            url: '/pages/driving/driving',
            complete: () => {
              wx.hideLoading()
            }
          })
        }, 2000)

      },
      fail: () =>{
        wx.showToast({
          icon: "none",
          title: "请前往设置页授权位置信息"
        })
      }
    })

    wx.showLoading({
      title: "开锁中...",
      mask: true
    })


  }

})