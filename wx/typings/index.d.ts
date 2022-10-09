/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-02 08:29:18
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-10-10 00:42:21
 * @FilePath: /coolcar/wx/typings/index.d.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/// <reference path="./types/index.d.ts" />

interface IAppOption {
  globalData: {
    // userInfo?: WechatMiniprogram.UserInfo,
    userInfo: Promise<WechatMiniprogram.UserInfo>
  }
  userInfoReadyCallback?: WechatMiniprogram.GetUserInfoSuccessCallback,
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo): void;
}