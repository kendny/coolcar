/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-21 15:42:02
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-21 18:26:54
 * @FilePath: /coolcar/wx/ts_demo/greeter.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
var Student = /** @class */ (function () {
    function Student(firstName, middleInitial, lastName) {
        this.firstName = firstName;
        this.middleInitial = middleInitial;
        this.lastName = lastName;
        this.fullName = firstName + " " + middleInitial + " " + lastName;
    }
    return Student;
}());
function greeter(person) {
    return "Hello, " + person.firstName + " " + person.lastName;
}
var user = new Student("Jane", "M.", "User");
// document.body.textContent = greeter(user)
console.log(user);

// void类型

