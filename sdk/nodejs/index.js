const {authCodeGrant, deviceCodeGrant, clientCredentialsGrant, printToken, refreshAccessToken} = require("./methods");
const Web3 = require('web3');

const ABCMethods = Object.freeze({
    AADAuthCode: 'aadauthcode',
    AADDevice: 'aaddevice',
    AADClient: 'aadclient',
});

function checkEmpty(obj, ...names) {
    names.forEach(name => {
        const v = obj[name];
        if (typeof v != 'string' || v === '') {
            throw `config: ${name} should be a non-empty string`;
        }
    })
}

async function requestAccess(config, tok) {
    if (tok && tok.refreshToken) {
        tok = await refreshAccessToken(tok.refreshToken, config.adalOptions);
    } else {
        switch (config.method) {
            case ABCMethods.AADAuthCode:
                tok = await authCodeGrant(config.adalOptions);
                break;
            case ABCMethods.AADDevice:
                tok = await deviceCodeGrant(config.adalOptions);
                break;
            case ABCMethods.AADClient:
                tok = await clientCredentialsGrant(config.adalOptions);
                break;
        }
    }

    if (!tok || !tok.accessToken) {
        throw 'fetch: got no access token';
    }
    return tok
}

const DEFAULT_METHOD = ABCMethods.AADDevice;

class ABCProviderConfig {

    constructor(obj) {
        this.httpOptions = {};
        this.remote = null;
        this.tenantId = null;
        this.method = DEFAULT_METHOD;
        this.clientId = null;
        this.clientSecret = null;

        Object.assign(this, obj);
        if (Object.values(this).indexOf(this.method) < 0) {
            throw "config: method not available"
        }
        checkEmpty(this, "remote", "tenantId");
        if (this.method === ABCMethods.AADClient) {
            checkEmpty(this, "clientId", "clientSecret");
        }
    }

    get host() {
        return `https://${this.remote}`
    }

    // adalOptions returns the corresponding options object for adal-node methods
    get adalOptions() {
        let options = {
            authorityHostUrl: "https://login.microsoftonline.com",
            tenant: this.tenantId,
            clientId: this.clientId,
            resource: "5838b1ed-6c81-4c2f-8ca1-693600b4e2ca",
            clientSecret: this.clientSecret,
        };

        // auth code flow/device code flow use fixed settings to work
        if (this.method === "aadauthcode" || this.method === "aaddevice") {
            Object.assign(options, {
                authorityHostUrl: "https://login.microsoftonline.com",
                clientId: "a8196997-9cc1-4d8a-8966-ed763e15c7e1",
                clientSecret: null,
                redirectUri: "http://localhost:3100/_callback"
            });
        }
        return options
    }

}

class ABCProvider {

    constructor(config) {
        this.config = config instanceof ABCProviderConfig ? config : new ABCProviderConfig(config);
        if (this.config.method === ABCMethods.AADAuthCode) {
            throw 'error: method aadauthcode not supported in this sdk'
        }
        this.provider = null;
        this.tok = null;
        this.updateProviderPromise = null;
    }

    async _updateProvider() {
        if (!this.updateProviderPromise) {
            this.updateProviderPromise = new Promise(async (resolve) => {
                this.tok = await requestAccess(this.config, this.tok);
                let options = this.config.httpOptions || {};
                options.headers = options.headers || [];
                options = options.headers.filter(header => header.name !== 'Authorization');
                options.push({name: "Authorization", value: "Bearer " + this.tok.accessToken});
                this.provider = new Web3.providers.HttpProvider(this.config.host, options);
                resolve()
            }).then(() => {
                this.updateProviderPromise = null;
            });
        }
        return this.updateProviderPromise;
    }

    get host() {
        return this.provider.host;
    }

    get conneted() {
        // now HttpProvider.connected always returns true
        return this.provider && this.provider.connected;
    }

    async _request(fnName, args) {
        if (!this.provider) {
            await this._updateProvider();
        }
        const fn = this.provider[fnName];
        try {
            fn.apply(this.provider, args);
        } catch (e) {
            await this._updateProvider();
            fn.apply(this.provider, args);
        }
    }

    send() {
        return this._request('send', arguments);
    }

    sendBatch() {
        return this._request('sendBatch', arguments);
    }

    disconnect() {
        return this.provider.disconnect();
    }
}

module.exports = {
    ABCMethods,
    ABCProviderConfig,
    ABCProvider
};
