using System;
using System.Net.Http.Headers;
using System.Threading.Tasks;
using Microsoft.IdentityModel.Clients.ActiveDirectory;
using Nethereum.JsonRpc.Client;

namespace dotnet
{
    public class AadConfig
    {
        public string AuthorityHost = "https://login.microsoftonline.com";
        public string Tenant = "<tenant>";
        public string ClientId = "<client_id>";
        public string ClientSecret = "<client_secret>";
        public string Resource = "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca";
        public string RedirectUri = "<redirect_uri>";
    }

    public static class AbcMethods
    {
        public static readonly string AadAuthCode = "aadauthcode";
        public static readonly string AadDevice = "aaddevice";
        public static readonly string AadClient = "aadclient";

        public static readonly string Default = AadAuthCode;
    }

    public class AbcConfig
    {
        public string Remote;
        public string TenantId;
        public string Method;
        public string ClientId;
        public string ClientSecret;

        public string Host => $"https://{Remote}";

        private void CheckValue(string s, string name)
        {
            if (string.IsNullOrEmpty(TenantId))
            {
                throw new Exception($"AbcConfig: {name} should not be null or empty");
            }
        }

        public AbcConfig Check()
        {
            if (string.IsNullOrEmpty(Method))
            {
                Method = AbcMethods.Default;
            }

            CheckValue(TenantId, "TenantId");
            if (Method == AbcMethods.AadClient)
            {
                CheckValue(ClientId, "ClientId");
                CheckValue(ClientSecret, "ClientSecret");
            }

            return this;
        }
    }

    public class AbcRpcClient : IClient
    {
        public AbcConfig Config;
        public IClient Client;
        public AuthenticationResult Tok;

        public AbcRpcClient(AbcConfig config)
        {
            Config = config.Check();
            UpdateClient();
        }

        private void UpdateClient()
        {
            if (Tok == null)
            {
                Client = new RpcClient(new Uri(Config.Host));
            }
            else
            {
                var authHeader = new AuthenticationHeaderValue(Tok.AccessTokenType, Tok.AccessToken);
                Client = new RpcClient(new Uri(Config.Host), authHeader);
            }
        }

        public async Task<T> SendRequestAsync<T>(RpcRequest request, string route = null)
        {
            try
            {
                return await Client.SendRequestAsync<T>(request, route);
            }
            catch (RpcClientUnknownException)
            {
                await RequestAccess();
                UpdateClient();
                return await Client.SendRequestAsync<T>(request, route);
            }
        }

        public Task<T> SendRequestAsync<T>(string method, string route = null, params object[] paramList)
        {
            return Client.SendRequestAsync<T>(method, route, paramList);
        }

        public Task SendRequestAsync(RpcRequest request, string route = null)
        {
            return Client.SendRequestAsync(request, route);
        }

        public Task SendRequestAsync(string method, string route = null, params object[] paramList)
        {
            return Client.SendRequestAsync(method, route, paramList);
        }

        public RequestInterceptor OverridingRequestInterceptor { get; set; }


        private static async Task<AuthenticationResult> RetrieveToken(string method, AadConfig config)
        {
            // aadauthcode and aaddevice method use fixed settings
            if (method == "aadauthcode" || method == "aaddevice")
            {
                config.ClientId = "a8196997-9cc1-4d8a-8966-ed763e15c7e1";
                config.ClientSecret = null;
                config.RedirectUri = "urn:ietf:wg:oauth:2.0:oob";
            }

            var authority = config.AuthorityHost + "/" + config.Tenant;

            var ctx = new AuthenticationContext(authority);
            AuthenticationResult result;

            try
            {
                result = await ctx.AcquireTokenSilentAsync(config.Resource, config.ClientId);
            }
            catch (Exception)
            {
                switch (method)
                {
                    case "aadauthcode":
                        result = await ctx.AcquireTokenAsync(config.Resource, config.ClientId,
                            new Uri(config.RedirectUri),
                            new PlatformParameters(PromptBehavior.SelectAccount));
                        break;
                    case "aaddevice":
                        var codeResult = await ctx.AcquireDeviceCodeAsync(config.Resource, config.ClientId);
                        Console.WriteLine("Open: " + codeResult.VerificationUrl);
                        Console.WriteLine("Enter: " + codeResult.UserCode);
                        result = await ctx.AcquireTokenByDeviceCodeAsync(codeResult);
                        break;
                    case "aadclient":
                        var clientCredential = new ClientCredential(config.ClientId, config.ClientSecret);
                        result = await ctx.AcquireTokenAsync(config.Resource, clientCredential);
                        break;
                    default:
                        throw new Exception("method not found");
                }
            }

            return result;
        }

        private Task _requestAccessTask;

        private async Task RequestAccess()
        {
            if (_requestAccessTask == null)
            {
                _requestAccessTask = Task.Run(async () =>
                {
                    Tok = await RetrieveToken(Config.Method, new AadConfig
                    {
                        AuthorityHost = "https://login.microsoftonline.com",
                        Tenant = Config.TenantId,
                        ClientId = Config.ClientId,
                        ClientSecret = Config.ClientSecret
                    });
                });
                await _requestAccessTask;
                _requestAccessTask = null;
            }
            else
            {
                await _requestAccessTask;
            }
        }
    }
}