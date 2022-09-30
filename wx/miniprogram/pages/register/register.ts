/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-29 08:06:04
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-29 19:25:04
 * @FilePath: /coolcar/wx/miniprogram/pages/register/register.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// pages/register/register.ts
Page({

  /**
   * 页面的初始数据
   */
  data: {
    licNo: '',
    name: '',
    birthDate: '2000-01-01',
    genderIndex: 0,
    genders: ['未知', '男', '女', '其他'],
    licImgURL: '',
    state: 'UNSUBMITTED' as 'UNSUBMITTED' |  'PENDING' | 'VERIFIED' ,
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad() {

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

  /***
   * 
   * 上传照片的点击事件
  */
  onUploadLic(){
    console.log("onUploadLic===")
    wx.chooseImage({
      success: (res) => {
        console.log("chooseImage:==", res)
        if(res.tempFilePaths.length > 0){
          this.setData({
            licImgURL: res.tempFilePaths[0]
          })

          // todo... upload image
          setTimeout(() => {
            this.setData({
              licNo: "1213232",
              name: "kendny",
              genderIndex: 1,
              birthDate: "2017-01-01",
              state: 'UNSUBMITTED'
            })
          }, 1000)
        }
      }
    })
  },

  onGenderChange(e:any) {
    console.log("onGenderChange:==", e)
    this.setData({
      genderIndex: e.detail.value,
    })
  },

  onBirthDateChange(e:any) {
    console.log("BirthDateChange:==", e)
    this.setData({
      birthDate: e.detail.value,
    })
  },

  onSubmit(){
    // TODO: submit the form to server
    this.setData({
      state: 'PENDING'
    })

    // 模拟审核
    setTimeout(() => {
      this.onLicVerified()
    }, 3000)
  },

  onResubmit() {
    this.setData({
      state: 'UNSUBMITTED',
      licImgURL: ''
    })
  },

  onLicVerified(){
    // 模拟审核
    this.setData({
      state: 'VERIFIED'
    })

    // 审核通过，跳转到开锁页面
    wx.navigateTo({
      url: '/pages/lock/lock'
    })
  },

})