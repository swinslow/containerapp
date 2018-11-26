import axios from 'axios';

class TokenManager {
    constructor(apiroot, onFetchToken, onFetchLoginInfo) {
        this.apiroot = apiroot;
        this.onFetchToken = onFetchToken;
        this.onFetchLoginInfo = onFetchLoginInfo;

        // bind handlers
        this.fetchToken = this.fetchToken.bind(this);
        this.fetchLoginInfo = this.fetchLoginInfo.bind(this);
    }

    fetchToken(email) {
        // need to send request as form encoded data, not as JSON
        var params = new URLSearchParams();
        params.append('email', email);
        axios.post(this.apiroot + "/oauth/getToken", params)
        .then(res => {
            const token = res.data.token;
            this.onFetchToken(token);
        })
        .catch(err => {
            // FIXME should probably take a callback from
            // FIXME caller for what to do on error
            this.onFetchToken(null);
            console.log("error: " + err);
        });
    }

    fetchLoginInfo(token) {
        const url = this.apiroot + "/landing";
        const config = {
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        }
        
        // FIXME should also check for "unauthorized" returns,
        // FIXME e.g. if token is invalid or indicates caller is not logged in
        axios.get(url, config)
        .then(res => {
            console.log("res.status: " + res.status)
            if (res.status === 200) {
                let myself = {
                    isKnownUser: (res.data.id !== 0),
                    id: res.data.id,
                    name: res.data.name,
                    email: res.data.email,
                    isAdmin: res.data.is_admin
                };
                this.onFetchLoginInfo(myself);
            } else {
                this.onFetchLoginInfo(null);
            }
        })
        .catch(err => {
            this.onFetchLoginInfo(null);
            console.log("error: " + err)
        });
    }
}

export default TokenManager;
