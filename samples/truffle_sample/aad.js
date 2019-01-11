// This is a sample code using the adal-node library to show how azure-blockchain-connector works in the background.
// The most part is simply wrapping the adal-node library for a clear view.
// To use this code, you may want to run `npm install adal-node express` to get deps.
// You may also find the code in web3_sample.

const {AuthenticationContext} = require("adal-node");
const createApplication = require('express');
const crypto = require('crypto');
const url = require('url');

function getAuthorityUrl(opt) {
    return opt.authorityHostUrl + '/' + opt.tenant;
}

function getAuthUrl(opt, state) {
    return 'https://login.microsoftonline.com/' +
        opt.tenant +
        '/oauth2/authorize?response_type=code&client_id=' +
        opt.clientId +
        '&redirect_uri=' +
        opt.redirectUri +
        `&state=${state}&resource=` +
        opt.resource +
        '&prompt=select_account';
}

function authCodeGrant(opt) {
    return new Promise((resolve, reject) => {
        const app = createApplication();
        const u = url.parse(opt.redirectUri);
        const redirectPath = u.pathname;
        const port = u.port;
        let server;
        let stateToken;
        app.get('/_click_to_auth', (req, res) => {
            crypto.randomBytes(48, (ex, buf) => {
                stateToken = buf.toString('base64').replace(/\//g, '_').replace(/\+/g, '-');
                res.redirect(getAuthUrl(opt, stateToken));
            });
        });

        app.get(redirectPath, (req, res) => {
            if (stateToken !== req.query.state) {
                res.send('error: state does not match');
            }
            const ctx = new AuthenticationContext(getAuthorityUrl(opt));
            ctx.acquireTokenWithAuthorizationCode(
                req.query.code,
                opt.redirectUri,
                opt.resource,
                opt.clientId,
                opt.clientSecret,
                (err, resp) => {
                    let errorMessage = '';
                    if (err) {
                        errorMessage = 'error: ' + err.message + '\n';
                    }
                    errorMessage += 'response: ' + JSON.stringify(resp);
                    res.send(errorMessage);
                    if (server) {
                        server.close();
                    }
                    if (err) {
                        reject(resp);
                        return;
                    }
                    resolve(resp);
                }
            );
        });
        console.log("Authorize: http://localhost:" + port + '/_click_to_auth');
        server = app.listen(port)
    });
}

function clientCredentialsGrant(opt) {
    return new Promise((resolve, reject) => {
        const ctx = new AuthenticationContext(getAuthorityUrl(opt));
        ctx.acquireTokenWithClientCredentials(opt.resource, opt.clientId, opt.clientSecret,
            (err, resp) => {
                if (err) {
                    reject(err);
                    return;
                }
                resolve(resp);
            })
    })
}

function deviceCodeGrant(opt, language) {
    return new Promise((resolve, reject) => {
        const ctx = new AuthenticationContext(getAuthorityUrl(opt));
        ctx.acquireUserCode(opt.resource, opt.clientId, language, (err, userCodeInfo) => {
            if (err) {
                reject(err);
                return;
            }
            console.log("Open:", userCodeInfo.verificationUrl);
            console.log("Enter:", userCodeInfo.userCode);
            ctx.acquireTokenWithDeviceCode(opt.resource, opt.clientId, userCodeInfo,
                (err, resp) => {
                    if (err) {
                        reject(err);
                        return;
                    }
                    resolve(resp);
                })
        })
    })
}

function refreshAccessToken(refreshToken, opt) {
    return new Promise((resolve, reject) => {
        const ctx = new AuthenticationContext(getAuthorityUrl(opt));
        ctx.acquireTokenWithRefreshToken(refreshToken, opt.clientId, opt.clientSecret, opt.resource,
            (err, resp) => {
                if (err) {
                    reject(resp);
                    return;
                }
                resolve(resp);
            })
    })

}

// scheduleRefreshAccessToken shows a way to refresh token on expire
// Please note that setInterval and setTimeout is not accurate in JavaScript.
// A better way is to handle send()/sendAsync() errors, but in that case you might have to write your own provider.
function scheduleRefreshAccessToken(tok, opt, callback) {
    if (tok && tok.refreshToken) {
        const schedule = (timeoutFn, fn) => {
            let timeout;
            try {
                timeout = timeoutFn();
            } catch (e) {
                console.error(e);
            }
            setTimeout(() => {
                fn();
                schedule(timeoutFn, fn);
            }, timeout)
        };

        schedule(() => tok.expiresIn * 999, async () => {
            try {
                tok = await refreshAccessToken(tok.refreshToken, opt);
                callback(tok)
            } catch (err) {
                console.log(err);
            }
        });
    }
}

function printToken(tok) {
    console.log("Access:", tok.accessToken);
    console.log("Expires in:", tok.expiresIn);
    console.log("Expires on:", tok.expiresOn);
    console.log()
}

async function retrieveToken(nodeUri, opt, method) {
    // auth code flow/device code flow use fixed settings to work
    if (method === "aadauthcode" || method === "aaddevice") {
        Object.assign(opt, {
            authorityHostUrl: "https://login.microsoftonline.com",
            tenant: "microsoft.onmicrosoft.com",
            clientId: "a8196997-9cc1-4d8a-8966-ed763e15c7e1",
            clientSecret: null,
            resource: "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca",
            redirectUri: "http://localhost:3100/_callback"
        });
    }

    let tok;
    // use a method to retrieve
    try {
        switch (method) {
            case "aadauthcode":
                tok = await authCodeGrant(opt);
                break;
            case "aaddevice":
                tok = await deviceCodeGrant(opt);
                break;
            case "aadclient":
                tok = await clientCredentialsGrant(opt);
                break;
        }
    } catch (err) {
        console.error(err)
    }
    printToken(tok);

    if (!tok) {
        console.error("no token");
        return;
    }
    scheduleRefreshAccessToken(tok, opt, token => {
        tok = token;
        printToken(tok);
    });

}


module.exports = {
    authCodeGrant,
    deviceCodeGrant,
    clientCredentialsGrant,
    refreshAccessToken,
    scheduleRefreshAccessToken,
    printToken,
    retrieveToken
};

(function () {
    let opt = {
        authorityHostUrl: "https://login.microsoftonline.com",
        tenant: "<tenant-id>",
        clientId: "<client-id>",
        clientSecret: "<client-secret>", // required in client credentials flow
        resource: "<resource>",
        redirectUri: "<redirect_uri>", // required in auth code flow
    };
    retrieveToken("<node_uri>", opt, "aaddevice").then();
})();
