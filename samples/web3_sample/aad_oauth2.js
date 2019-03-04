const {authCodeGrant, deviceCodeGrant, clientCredentialsGrant, printToken, refreshAccessToken} = require("./aad_oauth2_fn");
const Web3 = require('web3');

let opt = {
    authorityHostUrl: "https://login.microsoftonline.com",
    tenant: "<tenant-id>",
    clientId: "<client-id>",
    clientSecret: "<client-secret>", // required in client credentials flow
    resource: "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca",
    redirectUri: "<redirect_uri>", // required in auth code flow
};

(async function (opt, method) {
    const nodeUri = '<node_uri>';

    const web3 = new Web3();

    // auth code flow/device code flow use fixed settings to work
    if (method === "aadauthcode" || method === "aaddevice") {
        Object.assign(opt, {
            authorityHostUrl: "https://login.microsoftonline.com",
            tenant: "microsoft.onmicrosoft.com",
            clientId: "a8196997-9cc1-4d8a-8966-ed763e15c7e1",
            clientSecret: null,
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

    const providerFromToken = function (tok) {
        // generate an authorization header with retrieved access_token
        const header = {name: "Authorization", value: "Bearer " + tok.accessToken};
        return new web3.providers.HttpProvider(nodeUri, 0, "", "", [header]);
    };

    web3.setProvider(providerFromToken(tok));

    // refresh access_token if needed
    const expire = false;
    if (tok && tok.refreshToken && expire) {
        try {
            tok = await refreshAccessToken(tok.refreshToken, opt);
            printToken(tok);
            web3.setProvider(providerFromToken(tok));
        } catch (err) {
            console.error(err);
        }
    }

    // your code

})(opt, "aaddevice");