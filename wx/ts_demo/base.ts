
// Tuple
let x:[string, number, boolean]
x = ['hello', 10, false]
console.log(x)
console.log(x[0], x[1], x[2])

console.log("=====Enum 枚举类型============")
// Enum 枚举类型
enum Color {
    Red,
    Green=3, 
    Blue=4,
    Yellow,
    Gray="gray",
}
let c:Color = Color.Red
let y:Color = Color.Yellow
console.log(c, y, Color.Gray)

console.log("=====unkown 类型============")
//unkown 类型
let notSure:unknown=4
console.log(typeof notSure)
notSure=false
console.log(typeof notSure)

console.log("=====any 类型============")
// unkown 与 any的区别 https://juejin.cn/post/7021676475434663966
// any 类型
let nSoure:any = 4

// 它允许你在编译时可选择地包含或移除类型检查
console.log(nSoure.ifItExists())
console.log(typeof nSoure, nSoure)
nSoure = "maybe a string instead"
nSoure = false
console.log(typeof nSoure, nSoure)



// void类型

console.log("=====void类型============")
function warnUser(): void {
    console.log("This is my warning message")
}




