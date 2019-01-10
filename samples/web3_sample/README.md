# web3.js Connection Sample

This is a web3.js sample to connect our transaction nodes without the blockchain connector. 

## Quickstart

To run the sample code, you need to firstly install some dependencies.
```
npm install web3 adal-node express
```

You may want to have a look at `basic_auth.js` and `aad_oauth2.js`. In `basic_auth.js`, the method to add basic authentication is also documented in the web3.js official repository. 

Azure Active Directory support(`aad_oauth2.js`), is implemented using the Azure official `adal-node` library. In this sample, oauth2 methods are wrapped into promises for intuition. As you can see, in short, it is fetching an access_token and then generate a new provider with corresponding auth header.
