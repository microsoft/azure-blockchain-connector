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

function authCodeGrant(opt, port) {
    return new Promise((resolve, reject) => {
        const app = createApplication();
        const redirectPath = url.parse(opt.redirectUri).pathname;
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
                    reject(resp);
                    return;
                }
                resolve(resp);
            })
    })
}

module.exports = {
    authCodeGrant,
    clientCredentialsGrant,
};

(async function () {

    const opt = {
        clientId: "<client-id>",
        clientSecret: "<client-secret>",
        authorityHostUrl: "https://login.microsoftonline.com",
        tenant: "<tenant-id>",
        redirectUri: "<redirect_uri>",
        resource: "<resource>",
    };

    let tok;
    try {
        // authorization code grant
        tok = await authCodeGrant(opt, "3100");

        // client credentials grant
        tok = await clientCredentialsGrant(opt);
    } catch (err) {
        console.error(err)
    }

    console.log("Access:", tok.accessToken);
    console.log("Expires in:", tok.expiresIn);
    console.log("Expires on:", tok.expiresOn);

})();
