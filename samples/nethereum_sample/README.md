# Nethereum Connection Sample

This is a Nethereum sample to connect our transaction nodes without the blockchain connector. 

## Quickstart

This sample uses two packages from Nethereum and ADAL.
```
dotnet add package Nethereum.Web3
dotnet add package Microsoft.IdentityModel.Clients.ActiveDirectory
```

In `Program.cs`, we demonstrate how to configure basic authentication and AAD OAuth2 to create a `Nethereum.Web3` object. Then you can use the instance to access to a target transaction node.

### Basic Authentication
To configure basic auth, you need to encode your username and password into a base64 string, and then put it into a header. `Nethereum.Web3` supports setting an auth header in `RpcClient`'s constructor.
```c#
var authValue = Convert.ToBase64String(Encoding.UTF8.GetBytes($"{username}:{password}"));
var authHeader = new AuthenticationHeaderValue("Basic", authValue);
var web3 = Web3(new RpcClient(new Uri(nodeUri), authHeader, null, null));
```

### Azure Active Directory OAuth2
In short, using AAD OAuth2 is to retrieve a token via the ADAL library, and then set it in requests' headers. This is similar to basic auth, but you should notice to handle exceptions for the token retrieve process.
```c#
var tok = await retrieveToken(config);
var authHeader = new AuthenticationHeaderValue("Bearer", tok.AccessToken);
var web3 = Web3(new RpcClient(new Uri(nodeUri), authHeader, null, null));
```
For example, the `retrieveToken("aadauthcode")` above can be equivalent to:
```c#
ctx.AcquireTokenAsync(resource, clientId, new Uri(redirectUri), new PlatformParameters(PromptBehavior.SelectAccount));
```