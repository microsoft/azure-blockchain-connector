# Microsoft.Blockchain.JsonRpc

`AzureBlockchainRpcClient` is a wrapper of `Nethereum.JsonRpc.Client`.
You can use this with `Nethereum.Web3` to access a transaction node.
Same as the connector, `remote` is your transaction node URI, `tenantId` is your directory ID. 

```c#
var web3 = new Web3(new AzureBlockchainRpcClient(remote, tenantId));
```
For target `net472`, the client will use method `aadauthcode` (popup a webview) to request access.
For target `netstandard2_0`, which does not support webview, method `aaddevice` will be used.

You can explicitly specify to use `aaddevice` by setting the`useDeviceFlow` flag to `true`.
```c#
var web3 = new Web3(new AzureBlockchainRpcClient(remote, tenantId, useDeviceFlow));
```

For method `aadclient`, you can pass `clientId` and `clientSecret` to the client flow override.
```c#
var web3 = new Web3(new AzureBlockchainRpcClient(remote, tenantId, clientId, clientSecret)
```