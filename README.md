# u2fcli

u2fcli is a tool designed to handing registering, signing, and verify U2F tokens on the command line.

## Install

`go get github.com/mdp/u2fcli`

### Usage


#### Registration

Choose a challenge and App ID to send to the U2F device

```
[mdp u2fcli]$ u2fcli reg --challenge complexChallengeGoesHere --appid https://mdp.im

Registering, press the button on your U2F device
Registered Data: BQQsElODRKZR9flGJaIaIcJ5uvGFPR-eET5Qn3zVmvN8ZTkFkNfVq6kvCdex1l5ObwBRdoHBjJ__9YV8QoMmMcP5QI0CSTcV0d8q7bx2UWxwskb7X0Z_8hjws-XWHwHDjJVJiFOoGpKu_QO1gW3-I7vHxu4F2E31QsgBhx0v1xL9gXswggJEMIIBLqADAgECAgRVYr6gMAsGCSqGSIb3DQEBCzAuMSwwKgYDVQQDEyNZdWJpY28gVTJGIFJvb3QgQ0EgU2VyaWFsIDQ1NzIwMDYzMTAgFw0xNDA4MDEwMDAwMDBaGA8yMDUwMDkwNDAwMDAwMFowKjEoMCYGA1UEAwwfWXViaWNvIFUyRiBFRSBTZXJpYWwgMTQzMjUzNDY4ODBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEszH3c9gUS5mVy-RYVRfhdYOqR2I2lcvoWsSCyAGfLJuUZ64EWw5m8TGy6jJDyR_aYC4xjz_F2NKnq65yvRQwmjOzA5MCIGCSsGAQQBgsQKAgQVMS4zLjYuMS40LjEuNDE0ODIuMS41MBMGCysGAQQBguUcAgEBBAQDAgUgMAsGCSqGSIb3DQEBCwOCAQEArBbZs262s6m3bXWUs09Z9Pc-28n96yk162tFHKv0HSXT5xYU10cmBMpypXjjI-23YARoXwXn0bm-BdtulED6xc_JMqbK-uhSmXcu2wJ4ICA81BQdPutvaizpnjlXgDJjq6uNbsSAp98IStLLp7fW13yUw-vAsWb5YFfK9f46Yx6iakM3YqNvvs9M9EUJYl_VrxBJqnyLx2iaZlnpr13o8NcsKIJRdMUOBqt_ageQg3ttsyq_3LyoNcu7CQ7x8NmeCGm_6eVnZMQjDmwFdymwEN4OxfnM5MkcKCYhjqgIGruWkVHsFnJa8qjZXneVvKoiepuUQyDEJ2GcqvhU2YKY1zBEAiBhRCTflfJIhFb3k_Rkm3oT6uHcWKWuJUS1IJmCLYNvCAIgZ95Ojyj1cVSenQGcQUuOicnaClx7x_z_WhCeUwHARwU
Public Key: BCwSU4NEplH1-UYlohohwnm68YU9H54RPlCffNWa83xlOQWQ19WrqS8J17HWXk5vAFF2gcGMn__1hXxCgyYxw_k
Key Handle: jQJJNxXR3yrtvHZRbHCyRvtfRn_yGPCz5dYfAcOMlUmIU6gakq79A7WBbf4ju8fG7gXYTfVCyAGHHS_XEv2Bew
```

#### Signing

Using the `Key Handle` from above, we can ask the U2F token to sign a challenge

```
[mdp u2fcli]$ u2fcli sig --challenge anotherChallengeGoesHere --appid https://mdp.im --keyhandle jQJJNxXR3yrtvHZRbHCyRvtfRn_yGPCz5dYfAcOMlUmIU6gakq79A7WBbf4ju8fG7gXYTfVCyAGHHS_XEv2Bew
Authenticating, press the button on your U2F device

Counter: 31
Signature: MEQCICYcvvB7x_lzj2pFZ1IIH80kp9GlHjm0NE3nVrI6ZYvAAiB4ilXXEIUPRGhzgETi_wi3ICryV6ePnuOsjIJqCF1grQ
Raw Response: AQAAAB8wRAIgJhy-8HvH-XOPakVnUggfzSSn0aUeObQ0TedWsjpli8ACIHiKVdcQhQ9EaHOAROL_CLcgKvJXp4-e46yMgmoIXWCt
```

#### Verifying

Finally, we can verify the `Raw Response` from above by providing it to u2fcli along with the `Public Key` we recieved at registration.

```
[mdp u2fcli]$ u2fcli ver --challenge complexChallengeGoesHere --appid https://mdp.im --publickey BCwSU4NEplH1-UYlohohwnm68YU9H54RPlCffNWa83xlOQWQ19WrqS8J17HWXk5vAFF2gcGMn__1hXxCgyYxw_k --signature AQAAAB8wRAIgJhy-8HvH-XOPakVnUggfzSSn0aUeObQ0TedWsjpli8ACIHiKVdcQhQ9EaHOAROL_CLcgKvJXp4-e46yMgmoIXWCt
Signature verified
```
