# web3.js Connection Sample

This is a web3.js sample to connect our transaction nodes without the blockchain connector. 

## Quickstart

To run the sample code, you need to firstly install some dependencies.
```
npm install web3 adal-node express
```

You may want to have a look at `basic_auth.js` and `aad_oauth2.js`. In `basic_auth.js`, the method to add basic authentication is also documented in the web3.js official repository. Azure Active Directory support(`aad_oauth2.js`), is implemented using the Azure official `adal-node` library. In this sample, oauth2 methods are wrapped into promises for intuition. 
Following is a quick walkthrough.

### Basic Authentication
Use a provider with basic auth information in URI.
```javascript
const provider = new web3.providers.HttpProvider(`http://${username}:${password}@${nodeUri}`);
web3.setProvider(provider);
```

### Azure Active Directory OAuth2
Invoke an auth grant method from `adal-node` to acquire a token, generate a auth header, put it into a provider. Refresh when token expires is similar.
```javascript
const tok = await retrieveToken(config);
const header = {name: "Authorization", value: "Bearer " + tok.accessToken};
const provider = web3.providers.HttpProvider(nodeUri, 0, "", "", [header]);
web3.setProvider(provider);
```
For example, `retrieveToken` can be thw following, which is a acquireToken method from `adal-node` library. See `aad_oauth2_fn.js` for more details.
```javascript
new Promise((resolve, reject) => {
    const ctx = new AuthenticationContext(authorityUrl);
    ctx.acquireTokenWithClientCredentials(resource, clientId, clientSecret,
        (err, resp) => {
            if (err) {
                reject(err);
                    return;
            }
            resolve(resp);
        })
    })
```