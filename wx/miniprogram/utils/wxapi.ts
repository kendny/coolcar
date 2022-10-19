/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2021-06-13 13:29:41
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-10-09 22:53:37
 * @FilePath: /coolcar/wx/miniprogram/utils/wxapi.ts
 * @Description: 这是默认设置, 请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export function getSetting(): Promise<WechatMiniprogram.GetSettingSuccessCallbackResult> {
  return new Promise((resolve, reject) => {
    wx.getSetting({
      success: resolve, // 等价 res => resolve(res)
      fail: reject, // 等价 err => resolve(err)
    })
  })
}

export function getUserInfo(): Promise<WechatMiniprogram.GetUserInfoSuccessCallbackResult> {
  return new Promise((resolve, reject) => {
    wx.getUserInfo({
      success: resolve,
      fail: reject,
    })
  })
}