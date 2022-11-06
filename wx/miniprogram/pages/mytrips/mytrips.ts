import {TripService} from "../../service/trip";


Page({    data: {
        promotionItems: [
            {
                img: 'https://img.mukewang.com/5f7301d80001fdee18720764.jpg',
                promotionID: 1,
            },
            {
                img: 'https://img.mukewang.com/5f6805710001326c18720764.jpg',
                promotionID: 2,
            },
            {
                img: 'https://img.mukewang.com/5f6173b400013d4718720764.jpg',
                promotionID: 3,
            },
            {
                img: 'https://img.mukewang.com/5f7141ad0001b36418720764.jpg',
                promotionID: 4,
            },
        ],
    },

    /**
     * 生命周期函数--监听页面加载
     */
    async onLoad() {
        // 获取所有行程
        //rental.v1.TripStatus.FINISHED
        const res = await TripService.GetTrips().then(console.log)
        console.log(res)
    },

})