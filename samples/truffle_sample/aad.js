const {AuthenticationContext} = require("adal-node");
const createApplication = require('express');
const crypto = require('crypto');
const url = require('url');

// This is a sample code using the adal-node library to show how azure-blockchain-connector works in the background.
// The most part is simply wraps the adal-node library and you can check the last part to see how the grant flow.
// To use this code, you may want to run `npm install adal-node express` to get deps.

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

function printToken(tok) {
    console.log("Access:", tok.accessToken);
    console.log("Expires in:", tok.expiresIn);
    console.log("Expires on:", tok.expiresOn);
    console.log()
}

module.exports = {
    authCodeGrant,
    clientCredentialsGrant,
    refreshAccessToken
};

(async function () {

    let opt = {
        authorityHostUrl: "https://login.microsoftonline.com",
        tenant: "<tenant-id>",
        clientId: "<client-id>",
        clientSecret: "<client-secret>", // required in client credentials flow
        resource: "<resource>",
        redirectUri: "<redirect_uri>", // required in auth code flow
    };

    // here you may choose a method
    const method = "aadauthcode";

    let tok;
    try {
        switch (method) {
            case "aadauthcode":
                // authorization code grant
                // for azure-blockchain-connector, it uses fixed settings for auth code mode to work
                Object.assign(opt, {
                    authorityHostUrl: "https://login.microsoftonline.com",
                    tenant: "microsoft.onmicrosoft.com",
                    clientId: "a8196997-9cc1-4d8a-8966-ed763e15c7e1",
                    clientSecret: null,
                    resource: "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca",
                    redirectUri: "http://localhost:3100/_callback"
                });
                tok = await authCodeGrant(opt);
                break;
            case "aadclient":
                // client credentials grant
                tok = await clientCredentialsGrant(opt);
        }
    } catch (err) {
        console.error(err)
    }

    if (tok && tok.refreshToken) {
        // refresh access_token when expires
        // Please note that setInterval and setTimeout is not accurate in JavaScript.
        const schedule = (timeoutFn, fn) => setTimeout(() => {
            fn();
            schedule(timeoutFn, fn);
        }, timeoutFn());

        schedule(() => tok.expiresIn * 1000, async () => {
            try {
                tok = await refreshAccessToken(tok.refreshToken, opt);
                printToken(tok);
            } catch (err) {
                console.log(err);
            }
        });
    }

    printToken(tok);

})();
