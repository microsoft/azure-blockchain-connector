# AzureBlockchainConnector.Web3

`AbcRpcClient` is a wrapper of `RpcClient` which accepts an `AbcConfig` instance. You can use this with `Nethereum.Web3` to access a transaction node.

```c#
var web3 = new Web3(new AbcRpcClient(new AbcConfig
    {
        Remote = "samplenode.blockchain.azure.com:3200",
        TenantId = "microsoft.onmicrosoft.com",
        Method = AbcMethods.AadAuthCode,
        ClientId = "",
        ClientSecret = ""
    }));
```