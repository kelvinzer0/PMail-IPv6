import {defineStore} from "pinia";
import {http} from "@/utils/axios";

const useGlobalStatusStore = defineStore('useGlobalStatusStore', {
    state() {
        return {
            userInfos: {},
            sidebarVisible: true // Add new state for sidebar visibility
        }
    },
    getters: {
        isLogin(state) {
            return Object.keys(state.userInfos).length > 0
        }
    },
    actions: {
        async init(callback) {
            let that = this
            try {
                const res = await http.post("/api/user/info", {});
                if (res.errorNo === 0) {
                    Object.assign(that.userInfos, res.data)
                    console.log("userInfos initialized:", that.userInfos);
                    callback();
                } else {
                    console.error("Failed to get user info:", res.errorMsg);
                    that.userInfos = {}; // Clear user info on error
                    callback();
                }
            } catch (error) {
                console.error("Error fetching user info:", error);
                that.userInfos = {}; // Clear user info on error
                callback();
            }
        },
        toggleSidebar() {
            this.sidebarVisible = !this.sidebarVisible;
        }
    }
})


export {useGlobalStatusStore};