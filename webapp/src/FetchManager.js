import axios from 'axios';

class FetchManager {
    constructor() {

        // bind handlers
        this.getWithToken = this.getWithToken.bind(this);
        this.postWithToken = this.postWithToken.bind(this);
    }

    // getWithToken makes API GET calls with a JWT token included.
    getWithToken(url, payload, token, success, failure) {
        const config = {
            url: url,
            data: payload,
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        }
        
        // FIXME should also check for "unauthorized" returns,
        // FIXME e.g. if token is invalid or indicates caller is not logged in
        axios.get(config)
        .then(res => {
            success(res.data)
        })
        .catch(err => {
            failure(err)
        });
    }

    // postWithToken makes API POST calls with a JWT token included.
    postWithToken(url, payload, token, success, failure) {
        const config = {
            url: url,
            data: payload,
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        }
        
        // FIXME should also check for "unauthorized" returns,
        // FIXME e.g. if token is invalid or indicates caller is not logged in
        axios.post(config)
        .then(res => {
            success(res.data)
        })
        .catch(err => {
            failure(err)
        });
    }
}

export default FetchManager;
