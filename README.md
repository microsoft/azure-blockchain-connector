# Azure Blockchain Connector

ABC is a proxy for you to access blockchain safely. With this project, you can connect to nodes with basic authentication or Azure Active Directory. To use it, you just need to provide a local address, your target node address and your authentication information. Then, you can access the nodes using your client applications via the local address.

## Quickstart

To build this project, you need to install golang 1.11 and tdm-gcc.

```
git clone https://github.com/Microsoft/azure-blockchain-connector.git
cd .\azure-blockchain-connector\
go build .\cmd\abc

.\abc <your parameters here>
```
You can run [./build.cmd](./build.cmd) to build for windows, macOS and linux at the same time. The output executables will be located at `./build`. Docker build is also supported but it is not recommended. You may want to check the content of `./.dockerfile`.

For authentication, this project supports basic authentication and several Azure Active Directory OAuth2 interfaces. To use the proxy, you need to supply parameters of the basic settings (local, remote addresses, secure settings), and then choose an authentication method and provide corresponding information. 

### Basic Parameters

**-remote** *string, **required***: The host of the blockchain node you want to access.

**-local** *string*: The local listen host. Your client(e.g geth) need to connect to this address and requests sent to this host will be redirected to the remote host. If not provided, `localhost:3100` will be used by default. 

**-insecure** *bool*: Indicating whether it should skip certificate verification. Sending authentication information via an insecure HTTP connection is dangerous, please check Basic Authentication and OAuth2 protocols.

**-cert** *string*: The certification file(PEM)'s path.

**-method** *string*: the authentication method you want to use. Available options are `basic` ,  `aadauthcode`, `aaddevice`, `aadclient`, for basic auth, AAD auth code flow, AAD device code flow, AAD client credentials flow respectively. Please refer to the following sections for detail method parameters. Default: `basic`.

### Basic Authentication

Basic auth adds an authentication header to the users' requests using their user-ID and password. If you choose basic auth, the following two args are required.

**-username** *string*: basic auth username field.

**-password** *string*: basic auth password field.

```shell
# Example

.\abc -remote="samplenode.blockchain.azure.com:3200" -username="alice" -password="123"
```

### Azure Active Directory

This proxy supports several AAD OAuth2 authentication flows. **-tenant-id** is always a required argument. **-client-id** and **-client-secret** are sometimes required for specified methods (In fact, this is because we pre-populate certain methods with fixed values that should never change, so you don't need to supply them). 

Now auth code flow and client credentials flow are supported. In auth code flow, you should provide **-tenant-id** and other args are not required. For client credentials flow, you must also provide **-client-id** and **-client-secret** values.

**-tenant-id** *string*: required field id of the Azure Active Directory your Azure Blockchain Member belongs to. You can use both the domain string and the GUID string.

**-client-id** *string*: required in the client credentials flow, specifies the AAD application you want to use to access.

**-client-secret** *string*: a secret value for the application. required in client credentials flow to indicate the user is the owner of the AAD application.

**-webview** *bool*: an optional arg for auth code flow. In Windows, its default value is true, the proxy will popup a webview window to ask the user to select their account to grant authentications. Otherwise, the proxy will listen to a specified host to receive a credential callback(auth code) from a authentication server.

**-authcode-addr** *string*: an optional arg for auth code flow. When popping-up window is not supported, or **-webview** is false, the proxy will listen to the host specified by this arg to receive auth code. The default value is `localhost:3100`. It only works when corresponding values are set in the Azure Portal. If using AAD, it will print the first pair of the access_token and refresh_token.

```shell
# Example

.\abc -remote="samplenode.blockchain.azure.com:3200" -method="aadauthcode" -tenant-id="<your_directory_id>"

.\abc -remote="samplenode.blockchain.azure.com:3200" -method="aaddevice" -tenant-id="<your_directory_id>"

.\abc -remote="samplenode.blockchain.azure.com:3200" -method="aadclient" -tenant-id="<your_directory_id>" -client-id="12345678-abcd-efgh-ijkl-1234567890ab"
-client-secret="q@w#e%r^t&y*u(i)o_p"

```

### Logging

**-whenlog** *string*: configuration about in what cases logs should be prited. Alternatives are "always", "onNon200" and "onError". Default "is onError".

- **onError**: print log only for those who raise exceptions.
- **onNon200**: print log for those who have a non-200 response, or those who raise exceptions.
- **always**: print log for every request

**-whatlog** *string*: configuration about what information should be included in logs. Alternatives are "basic" and "detailed". Default is "basic".

- **basic**: print the request's method and URI, and the response status code (and the exception message, if exception raised) in the log.
- **detailed**: print the request's method, URI and body, and the response status code and body (and the exception message, if exception raised) in the log.

**-debugmode** *bool*: open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.

# Access Nodes without ABC

This project is only an out-of-the-box tool for users to connect to nodes conveniently. You can also update your own workflows to support interacting with a transaction nodes. In general, it is to add basic auth support and OAuth2 support. For some oauth2 methods, some pre-specified settings should be used.

To add basic auth support, add an authentication header with base64-encoded username:password pair to all requests.
```
Authentication: Basic base64(<username>:<password>)
```
To add AAD support, use an OAuth2 grant flow to retrieve a token. Then append the token with a bearer authentication header. We support auth code flow, device flow and client credentials flow. The former two are three-legged flows, you should use specified OAuth2 settings to let the user logging in. Or you may register self-managed AAD application to use the client credentials flow.
```
Authentication: Bearer <access_token>
```
You can find sample code in [/samples](samples), which includes samples for [web3.js](samples/web3_sample), [truffle](samples/truffle_sample) and [Nethereum](samples/nethereum_sample). You can also get the specified settings mentioned above from these samples.


# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.