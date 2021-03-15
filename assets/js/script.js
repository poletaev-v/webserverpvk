window.addEventListener("load", function() {
    const http = axios.create({
        headers: { 'Cache-Control': 'no-cache' }
    });
    const ti = 1 * 1000; // Milliseconds
    const ttl = 10 * 1000;
    const app = new Vue({
        el: '#app',
        data: {
            message: "",
            lastTime: 0,
            ttl: ttl,
            objData: {},
            availableFields: [],
            settings: {},
            cntItems: 0,
        },
        mounted: function() {
            this.getSettings()
            this.reloadData()
        },
        methods: {
            getSettings: function() {
                http
                    .get('http://192.168.100.4:8080/data/settings.json')
                    .then(response => {
                        if (Object.keys(response.data).length) {
                            this.settings = response.data
                            this.availableFields = Object.keys(response.data)
                        }
                    })
            },
            reloadData: function() {
                setInterval(() => {
                    http
                        .get('http://192.168.100.4:8080/data/data.json')
                        .then(response => {
                            if (Number(response.data.timeStamp) && this.lastTime < Number(response.data.timeStamp)) {
                                this.lastTime = Number(response.data.timeStamp)
                            }

                            if (Number(response.data.timeStamp) && (+new Date()) < (Number(response.data.timeStamp) + this.ttl)) {
                                if (this.availableFields.length) {
                                    for (let key in response.data) {
                                        if (this.availableFields.includes(key)) {
                                            this.objData[this.settings[key]] = response.data[key]
                                        }
                                    }
                                    this.cntItems = Object.keys(this.objData).length
                                }
                            } else {
                                this.message = "Ведется весовой и габаритный контроль"
                                this.objData = {}
                                this.cntItems = 0
                            }
                            // console.log(this.objData)
                            // console.log("Call then!", (+new Date()), (Number(response.data.timeStamp) + this.ttl), (+new Date()) < (Number(response.data.timeStamp) + this.ttl))
                        })
                }, ti);
            },
        }
    });
});