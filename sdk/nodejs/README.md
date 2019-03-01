# node-abc-provider

After require or import this module, you can use `ABCProvider` as a web3 provider to create an web3 instance.

```javascript
const web3 = new Web3(new ABCProvider({
    method: ABCMethods.AADDevice,
    remote: 'samplenode.blockchain.azure.com:3200',
    tenantId: 'microsoft.onmicrosoft.com',
    clientId: null,
    clientSecret: null
}));

web3.eth.net.isListening()
    .then(() => console.log('state: connected'))
    .catch(e => console.error(e));

// Output: 
// state: connected
```

In the behind, this module wraps the HttpProvider from web3.js 1.0.0beta. When catches an error, it will request a new access token with the given configuration. Now there is no advanced error handling inside, so you need to take care the status of the connection when use.

This module currently does not support `aadauthcode` method, as the `adal-node`'s auth code flow support requires a port listening, which is not recommended now. As an alternative, you may want to use `aaddevice`.