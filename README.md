# Blockchain Connector

This is a proxy for you to connect to the blockchain safely. We provide three ways for you to use it: 1.Compile our source code (implement with golang) and run it. 2.Run it in Docker Container, with the Dockerfile we presented in this branch. 3.Run with the binary release (not available yet).

## Run with the Source Code

1. Get the environment of golang ready.

2. 
```bash
go get github.com/Microsoft/azure-blockchain-connector/tree/basic_auth
```

3. In the project directory, run
```bash
go build -o <outputFile> HttpProxy.go
```

where \<outputFile\> is the name of the executable file you want to specify with, "HttpProxy" if not specified.

4. 
```bash
./<outputFile>  <parameters>
```
where \<parameters\> is the parameters needed by the program. To know about the parameters in detail, see the section "Parameters". 

## Run with the Dockerfile

1. Get the environment for Docker.

2. Clone or download this repo to anywhere you like (If you clone it, please make sure to check out to the brunch "basic_auth").

3. Run in the project directory
```bash
docker build -t <image_name> .
```
where \<image_name\> is the name of the output image.

4. 
```bash
docker run --net=host --name=<container_name> <image_name> <parameters>
```
where <container_name> is the name of the container you are to run, and \<parameters\> are the parameters of the entrypoint (which essentially are the parameters of the golang code). To know about the parameters in detail, see the section "Parameters".

## Parameters

1. parameters:

   - **username** *string*

   ​        The username you want to login with. (no default)

   - **password** *string*

   ​        The password you want to login with. (no default)

   - **local** *string*

   ​        Local address to bind to. (default "1234")

   - **remote** *string*

           The host you want to send to. (no default)

2. example:

   ```bash
   ./HttpProxy -username user -password 12345 -port 1111
   ```

3. parameters for developer:

If you want to use this connector to test your environment, but you only have a self-signed SSL certificate, the parameters below may help you:

   - **cert** *string*

            the CA cert of your self-signed certificate.

   - **insecure** *bool*

            indicate if it should skip certificate verifications.

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
