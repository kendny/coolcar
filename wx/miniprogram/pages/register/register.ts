/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-29 08:06:04
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-10-09 23:26:55
 * @FilePath: /coolcar/wx/miniprogram/pages/register/register.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// pages/register/register.ts
import {routing} from "../../utils/routings";
import {rental} from "../../service/proto_gen/rental/rental_pb";
import {formatData} from "../../utils/format";
import {ProfileService} from "../../service/profile";

Page({

  /**
   * 页面的初始数据
   */
  rediretURL:'',
  profileRefresher:0,

  data: {
    licNo: '',
    name: '',
    birthDate: '2000-01-01',
    genderIndex: 0,
    genders: ['未知', '男', '女', '其他'],
    licImgURL: '',
    state:  rental.v1.IdentityStatus[rental.v1.IdentityStatus.UNSUBMITTED] ,
  },

  // 渲染profile
  renderProfile: function (p: rental.v1.IProfile) {
    this.setData({
      licNo: p.identity?.licNumber || "",
      name: p.identity?.name || "",
      genderIndex: p.identity?.gender || 0,
      birthDate: formatData(p.identity?.birthDateMillis || 0),
      state: rental.v1.IdentityStatus[p.identityStatus || 0],
    })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(opt:Record<'redirect', string>) {
    const o:routing.RegisterOpts = opt
    if(o.redirect) {
      this.rediretURL = decodeURIComponent(o.redirect)
    }
    ProfileService.getProfile().then(p => this.renderProfile(p))
    console.log(o, o.redirect)
  },


  /***
   * 
   * 上传照片的点击事件
  */
  onUploadLic(){
    wx.chooseImage({
      success: (res) => {
        console.log("chooseImage:==", res)
        if(res.tempFilePaths.length > 0){
          this.setData({
            licImgURL: res.tempFilePaths[0]
          })

          const data = wx.getFileSystemManager().readFileSync(res.tempFilePaths[0])
          wx.request({
            method: "PUT",
            url: "https://wuhan-1259722894.cos.ap-shanghai.myqcloud.com/account_2/639446f1ce6ee5f76eab5220?q-sign-algorithm=sha1&q-ak=AKIDbkfNr78vUq32pOhoiQxHMDpDPPESeicR&q-sign-time=1670661873%3B1670662873&q-key-time=1670661873%3B1670662873&q-header-list=host&q-url-param-list=&q-signature=bcd7f962d5130b9ad5fc6e8d32219f259f1b24ae",
            data,
            success: console.log,
            fail: console.error,
          })

          // todo... upload image
          // 此处 wx.uploadFile 上传文件有坑点
          // setTimeout(() => {
          //   this.setData({
          //     licNo: "1213232",
          //     name: "kendny",
          //     genderIndex: 1,
          //     birthDate: "2017-01-01",
          //     state: 'UNSUBMITTED'
          //   })
          // }, 1000)
        }
      }
    })
  },

  onGenderChange(e:any) {
    console.log("onGenderChange:==", e)
    this.setData({
      genderIndex: parseInt(e.detail.value),
    })
  },

  onBirthDateChange(e:any) {
    console.log("BirthDateChange:==", e)
    this.setData({
      birthDate: e.detail.value,
    })
  },

  // 清除掉轮询
  clearProfileRefresher() {
    if(this.profileRefresher) {
      clearInterval(this.profileRefresher)
      this.profileRefresher = 0
    }
  },
  onLicVerified() {
    this.setData({
      state: rental.v1.IdentityStatus[rental.v1.IdentityStatus.VERIFIED]
    })
    if(this.rediretURL) {
      wx.redirectTo({
        url: this.rediretURL,
      })
    }
  },

  scheduleProfileRefresher(){
    // 轮询获取profile更新 状态
    let i = 0;
    this.profileRefresher = setInterval(() => {
      console.log("profileRefresher:=", ++i)
      ProfileService.getProfile().then(p => {
        this.renderProfile(p)
        if(p.identityStatus != rental.v1.IdentityStatus.PENDING) {
          this.clearProfileRefresher()
        }
        if(p.identityStatus == rental.v1.IdentityStatus.VERIFIED) {
          this.onLicVerified()
        }

      })
    }, 1000)
  },
  onSubmit(){
    ProfileService.submitProfile({
      licNumber: this.data.licNo,
      name: this.data.name,
      gender: this.data.genderIndex,
      birthDateMillis: Date.parse(this.data.birthDate)
    }).then(p => {
      this.renderProfile(p)
      this.scheduleProfileRefresher()
    })
  },

  onResubmit() {
    // 重新提交
    ProfileService.clearProfile().then(p => this.renderProfile(p))
  },
  onUnload(): void | Promise<void> {
    this.clearProfileRefresher()
  }

})