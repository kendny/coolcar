import { formatDuration, formatFee } from '../../utils/format';
// pages/driving/driving.ts

const centPerSec = 0.4
Page({
  timer: undefined as number|undefined,
  tripID: '',

  /**
   * 页面的初始数据
   */
  data: {
    location: {
      latitude: 32.39,
      longitude: 118.46,
    },
    scale: 12,
    elapsed: '00:00:00',
    fee: '0.00',
    // markers: [
    //   {
    //     iconPath: "/resources/car.png",
    //     id: 0,
    //     latitude: initialLat,
    //     longitude: initialLng,
    //     width: 20,
    //     height: 20,
    //   },
    // ],
  },

  setupLocationUpdator(){
    wx.startLocationUpdate({
      fail: console.error
    })
    wx.onLocationChange(loc=>{
      this.setData({
        location: {
          latitude: loc.latitude,
          longitude: loc.longitude
        }
      })
    })
  },

  setupTimer(){
    let elapsedSec = 0;
    let cents = 0
    this.timer = setInterval(() => {
      elapsedSec++
      cents += centPerSec
      this.setData({
        elapsed: formatDuration(elapsedSec),
        fee: formatFee(cents)
      })
    }, 1000)
  },
  
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad() {
    this.setupLocationUpdator()
    this.setupTimer()
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
    // 关闭监听实时位置变化，前后台都停止消息接收
    wx.stopLocationUpdate()
    if (this.timer) {
      clearInterval(this.timer)
    }
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

  
})