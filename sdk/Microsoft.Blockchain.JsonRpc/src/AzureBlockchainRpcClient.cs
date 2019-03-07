// <copyright file="AzureBlockchainRpcClient.cs" company="Microsoft">
// Copyright (c) Microsoft Corporation. All rights reserved.
// </copyright>

namespace Microsoft.Blockchain.JsonRpc
{
    using System;
    using System.Net.Http.Headers;
    using System.Threading.Tasks;
    using Microsoft.IdentityModel.Clients.ActiveDirectory;
    using Nethereum.JsonRpc.Client;

    public class AzureBlockchainRpcClient : IClient
    {
        private readonly string remote;
        private readonly string tenantId;
        private readonly string method;
        private readonly string clientId;
        private readonly string clientSecret;

        private IClient client;
        private AuthenticationResult tok;
        private Task requestAccessTask;

#if NET472
        public AzureBlockchainRpcClient(string remote, string tenantId, bool useDeviceFlow = false)
        {
            this.remote = remote;
            this.tenantId = tenantId;
            this.method = useDeviceFlow ? AzureBlockchainMethodsConstants.AadDevice : AzureBlockchainMethodsConstants.AadAuthCode;
            this.clientId = "a8196997-9cc1-4d8a-8966-ed763e15c7e1";
            this.UpdateClient();
        }
#elif NETSTANDARD2_0
        public AzureBlockchainRpcClient(string remote, string tenantId)
        {
            this.remote = remote;
            this.tenantId = tenantId;
            this.method = AzureBlockchainMethodsConstants.AadDevice;
            this.clientId = "a8196997-9cc1-4d8a-8966-ed763e15c7e1";
            this.UpdateClient();
        }
#endif

        public AzureBlockchainRpcClient(string remote, string tenantId, string clientId, string clientSecret)
        {
            this.remote = remote;
            this.tenantId = tenantId;
            this.method = AzureBlockchainMethodsConstants.AadClient;
            this.clientId = clientId;
            this.clientSecret = clientSecret;
            this.UpdateClient();
        }

        public RequestInterceptor OverridingRequestInterceptor { get; set; }

        public async Task<T> SendRequestAsync<T>(
            RpcRequest request,
            string route = null)
        {
            try
            {
                return await this.client.SendRequestAsync<T>(request, route).ConfigureAwait(false);
            }
            catch (RpcClientUnknownException)
            {
                await this.RequestAccess().ConfigureAwait(false);
                this.UpdateClient();
                return await this.client.SendRequestAsync<T>(request, route).ConfigureAwait(false);
            }
        }

        public Task<T> SendRequestAsync<T>(string method, string route = null, params object[] paramList)
        {
            return this.client.SendRequestAsync<T>(method, route, paramList);
        }

        public Task SendRequestAsync(RpcRequest request, string route = null)
        {
            return this.client.SendRequestAsync(request, route);
        }

        public Task SendRequestAsync(string method, string route = null, params object[] paramList)
        {
            return this.client.SendRequestAsync(method, route, paramList);
        }

        private void UpdateClient()
        {
            if (this.tok == null)
            {
                this.client = new RpcClient(new Uri($"https://{this.remote}"));
            }
            else
            {
                var authHeader = new AuthenticationHeaderValue(this.tok.AccessTokenType, this.tok.AccessToken);
                this.client = new RpcClient(new Uri($"https://{this.remote}"), authHeader);
            }
        }

        private async Task RequestAccess()
        {
            if (this.requestAccessTask == null)
            {
                this.requestAccessTask = Task.Run(async () => { this.tok = await this.RetrieveToken().ConfigureAwait(false); });
                await this.requestAccessTask.ConfigureAwait(false);
                this.requestAccessTask = null;
            }
            else
            {
                await this.requestAccessTask.ConfigureAwait(false);
            }
        }

        private async Task<AuthenticationResult> RetrieveToken()
        {
            var authority = $"https://login.microsoftonline.com/{this.tenantId}";
            var resource = "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca";

            var ctx = new AuthenticationContext(authority);
            AuthenticationResult result;

            try
            {
                result = await ctx.AcquireTokenSilentAsync(resource, this.clientId).ConfigureAwait(false);
            }
            catch (Exception)
            {
                switch (this.method)
                {
                    case AzureBlockchainMethodsConstants.AadAuthCode:
#if NET472
                        result = await ctx.AcquireTokenAsync(
                            resource,
                            this.clientId,
                            new Uri("urn:ietf:wg:oauth:2.0:oob"),
                            new PlatformParameters(PromptBehavior.SelectAccount)).ConfigureAwait(false);
#elif NETSTANDARD2_0
                        result = await ctx.AcquireTokenAsync(
                            resource,
                            this.clientId,
                            new Uri("urn:ietf:wg:oauth:2.0:oob"),
                            new PlatformParameters()).ConfigureAwait(false);
#endif
                        break;
                    case AzureBlockchainMethodsConstants.AadDevice:
                        var codeResult = await ctx.AcquireDeviceCodeAsync(resource, this.clientId).ConfigureAwait(false);
                        Console.WriteLine("Open: " + codeResult.VerificationUrl);
                        Console.WriteLine("Enter: " + codeResult.UserCode);
                        result = await ctx.AcquireTokenByDeviceCodeAsync(codeResult).ConfigureAwait(false);
                        break;
                    case AzureBlockchainMethodsConstants.AadClient:
                        var clientCredential = new ClientCredential(this.clientId, this.clientSecret);
                        result = await ctx.AcquireTokenAsync(resource, clientCredential).ConfigureAwait(false);
                        break;
                    default:
                        throw new Exception("method not found");
                }
            }

            return result;
        }
    }
}