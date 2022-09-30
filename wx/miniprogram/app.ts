/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-02 08:29:18
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-29 20:11:32
 * @FilePath: /coolcar/wx/miniprogram/app.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// app.ts

import { getSetting, getUserInfo } from "./utils/wxapi"
let resolveUserInfo: (value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void
let rejectUserInfo: (reason?: any) => void


App<IAppOption>({
  globalData: {
    userInfo: new Promise((resolve, reject) => {
      resolveUserInfo = resolve // 等价 res => resolve(res)
      rejectUserInfo = reject  // 等价 error => reject(error)
    })
  },
  async onLaunch() {
    // 登录
    // Coolcar.login()

    // 获取用户信息
    try {
      const setting = await getSetting()
      if (setting.authSetting['scope.userInfo']) {
        const userInfoRes = await getUserInfo()
        resolveUserInfo(userInfoRes.userInfo)
      }
    } catch (err) {
      rejectUserInfo(err)
    }
  },
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo) {
    resolveUserInfo(userInfo)
  }
})