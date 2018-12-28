import {TokenResponse} from "adal-node";

const AuthenticationContext = require("adal-node");
const createApplication = require('express');
const crypto = require('crypto');


class OAuth2Options {
    clientId: string;
    clientSecret?: string;
    resource?: string;
    authorityHostUrl: string;
    tenant: string;
    redirectUri: string;

    get authorityUrl(): string {
        return this.authorityHostUrl + '/' + this.tenant;
    }

    get authorizationUrlTemplate() {
        return 'https://login.microsoftonline.com/' +
            this.tenant +
            '/oauth2/authorize?response_type=code&client_id=' +
            this.clientId +
            '&redirect_uri=' +
            this.redirectUri +
            '&state=<state>&resource=' +
            this.resource;
    }

    authorizationUrl(stateToken: string): string {
        return this.authorizationUrlTemplate.replace('<state>', stateToken);
    }
}

class AzureADOAuth2Service {
    private tokenResp: TokenResponse;

    constructor(public opt: OAuth2Options) {
    }

    get header() {
        if (!this.tokenResp.accessToken) {
            console.error("header: no accessToken has been fetched.")
        }
        return {
            name: "Authorization",
            value: "Bearer " + this.tokenResp.accessToken
        }
    }

    authCodeGrant(listenPath: string) {
        return new Promise((resolve, reject) => {
            const app = createApplication();
            app.get('/auth', (req, res) => {
                crypto.randomBytes(48, (ex, buf) => {
                    const token = buf.toString('base64').replace(/\//g, '_').replace(/\+/g, '-');
                    res.cookie('authstate', token);

                    const authorizationUrl = this.opt.authorizationUrl(token);
                    res.redirect(authorizationUrl);
                });
            });

            app.get('/getAToken', (req, res) => {
                if (req.cookies.authstate !== req.query.state) {
                    res.send('error: state does not match');
                }

                const ctx = new AuthenticationContext(this.opt.authorityUrl);
                ctx.acquireTokenWithAuthorizationCode(
                    req.query.code,
                    this.opt.redirectUri,
                    this.opt.resource,
                    this.opt.clientId,
                    this.opt.clientSecret,
                    (err, resp) => {
                        let errorMessage = '';
                        if (err) {
                            errorMessage = 'error: ' + err.message + '\n';
                        }
                        errorMessage += 'response: ' + JSON.stringify(resp);
                        res.send(errorMessage);
                        app.close();
                        if (err) {
                            reject(resp);
                            return;
                        }
                        this.tokenResp = resp;
                        resolve(resp);
                    }
                );
            });
            console.log("Authorize: " + listenPath + '/auth');
            app.listen(listenPath)
        });
    }

    clientCredentialsGrant() {
        return new Promise((resolve, reject) => {
            const ctx = new AuthenticationContext(this.opt.authorityUrl);
            ctx.acquireTokenWithClientCredentials(this.opt.resource, this.opt.clientId, this.opt.clientSecret,
                (err, resp) => {
                    if (err) {
                        reject(resp);
                        return;
                    }
                    this.tokenResp = resp;
                    resolve(resp);
                })
        })
    }
}

module.exports = {
    OAuth2Options,
    AzureADOAuth2Service,
};