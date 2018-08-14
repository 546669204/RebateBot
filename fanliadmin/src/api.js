import axios from 'axios'
import qs from 'qs'
axios.defaults.baseURL = "/"

function post(url,data){
    return axios({
        url,
        method:"POST",
        data:qs.stringify(data),
        headers:{'Content-Type': 'application/x-www-form-urlencoded'}
    });
}
function postjson(url,data){
    return axios({
        url,
        method:"POST",
        data,
    });
}
function get(url,data){
    return axios({
        url,
        method:"GET",
        params:data,
    });
}


export default {
    getService(data) {
        return get("/api/getService", data)
    },
    getUserData(data) {
        return get("/api/getUserData", data)
    },
    getOrderData(data) {
        return get("/api/getOrderData", data)
    },
    getWithdrawData(data) {
        return get("/api/getWithdrawData", data)
    },
    withdrawPay(data) {
        return post("/api/withdrawPay", data)
    },
    getTemplateData(data) {
        return get("/api/getTemplateData", data)
    },
    setTemplateData(data) {
        return post("/api/setTemplateData", data)
    },
    tblogin(data) {
        return post("/api/tblogin", data)
    },
    tbchecklogin(data) {
        return post("/api/tbchecklogin", data)
    },
    welogin(data) {
        return post("/api/welogin", data)
    },
    wechecklogin(data) {
        return post("/api/wechecklogin", data)
    },
    reboot(data) {
        return post("/api/reboot", data)
    },
}