# Blockchain Connector

This is a proxy for you to connect to the blockchain safely. We provide three ways for you to use it: 

    - Compile our source code (implement with golang) and run it.
    
    - Run it in Docker Container, with the Dockerfile we presented in this branch. 
    
    - Run with the binary release (not available yet).

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

            The username you want to login with (no default).

   - **password** *string*

            The password you want to login with (no default).

   - **local** *string*

            Local address to bind to (default "127.0.0.1:3100"). Note that local address can be non-local, but if so, you should make sure the connection from your computer to the "local address" is safe (e.g. the connection is in a LAN).

   - **remote** *string*

           The host you want to send to (no default).

    - **cert** *string*

            The CA cert of the remote.

    - **insecure** *bool*

            Indicating if it should skip certificate verifications.

    - **whenlog** *string*

            Configuration about in what cases logs should be prited. Alternatives are "always", "onNon200" and "onError". Default "is onError". See 4 for details.

    - **whatlog** *string*

            Configuration about what information should be included in logs. Alternatives are "basic" and "detailed". Default is "basic". See 5 for details.

    - **debugmode** *bool*

            Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.

2. example for users who run with our source code:

   ```bash
   ./<outputFile> -username user -password 12345 -remote https://microsoft.com/ -local 127.0.0.1:3100 -insecure -debugmode
   ```

3. examples for users who run with Docker
    ```bash
   docker run --net=host --name=<container_name> <image_name> -username user -password 12345 -remote https://microsoft.com/ -local 127.0.0.1:3100 -insecure
   ```

4. explainations about the alternatives for **whenlog**

    - **onError**

            Print log only for those who raise exceptions.

    - **onNon200** 

            Print log for those who have a non-200 response, or those who raise exceptions.

    - **always** 

            Print log for every request

5. explainations about the alternatives for **whatlog**

    - **basic**

            Print the request's method and URI, and the response status code (and the exception message, if exception raised) in the log.

    - **detailed** 

            Print the request's method, URI and body, and the response status code and body (and the exception message, if exception raised) in the log.

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
