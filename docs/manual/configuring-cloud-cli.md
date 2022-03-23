This section will introduce how to configure cloud-cli with credentials for accessing API7 CLOUD.

# Create Access Token
First, you need to log in to [API7 Cloud Web Console](https://console.api7.cloud) to create an **Access Token** associated with your account. 

**NOTE**: This Token will have the authority to operate the resources under your account, so please be sure to safekeeping.

It is recommended to configure a reasonable validity period for the Access Token. 

When the Access Token is leaked, you can revoke the token in API7 Cloud to avoid possible security risks.

# cloud-cli configure command
After executing the `cloud-cli configure` command, cloud-cli will prompt you to enter the Access Token.
````
$ cloud-cli configure
API7 Cloud Access Token: {PASTE YOUR ACCESS TOKEN HARE}
````
When the token you entered is verified to be correct and valid, cloud-cli will give the following prompt:
````
WARNING: your access token will expire at 2024-01-31T00:03:50+08:00
successfully configured api7 cloud access token, your account is jack@api7.ai
````

Enjoy!
