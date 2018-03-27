# u2fcli

u2fcli is a tool designed to handle registering, signing, and verifying U2F tokens on the command line.

## Install

`go get github.com/mdp/u2fcli`

## Usage


### Registration

Choose a challenge and App ID to send to the U2F device

```
[mdp u2fcli]$ u2fcli reg --challenge complexChallengeGoesHere \
  --appid https://mdp.im
Registering, press the button on your U2F device

{
  "KeyHandle": "0JGeJ3MhvDzK_YjKhK4VkPOegGn0x3wxJENJ8J1JanozbSr8Elz2KRcARLh2sF__l_Vof2xiydPw6CEicpzs0A",
  "PublicKey": "BPQPBz7NV3LwksVwjbGdn7ODP5omKHt8CetrHnDZeUUxmFChHcKuYNHgLm0HdtsSD6p7cjrFZdb9mNOLg3huRcI",
  "RegisteredData": "0JGeJ3MhvDzK_YjKhK4VkPOegGn0x3wxJENJ8J1JanozbSr8Elz2KRcARLh2sF__l_Vof2xiydPw6CEicpzs0DCCAkQwggEuoAMCAQICBFVivqAwCwYJKoZIhvcNAQELMC4xLDAqBgNVBAMTI1l1YmljbyBVMkYgUm9vdCBDQSBTZXJpYWwgNDU3MjAwNjMxMCAXDTE0MDgwMTAwMDAwMFoYDzIwNTAwOTA0MDAwMDAwWjAqMSgwJgYDVQQDDB9ZdWJpY28gVTJGIEVFIFNlcmlhbCAxNDMyNTM0Njg4MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESzMfdz2BRLmZXL5FhVF-F1g6pHYjaVy-haxILIAZ8sm5RnrgRbDmbxMbLqMkPJH9pgLjGPP8XY0qerrnK9FDCaM7MDkwIgYJKwYBBAGCxAoCBBUxLjMuNi4xLjQuMS40MTQ4Mi4xLjUwEwYLKwYBBAGC5RwCAQEEBAMCBSAwCwYJKoZIhvcNAQELA4IBAQCsFtmzbrazqbdtdZSzT1n09z7byf3rKTXra0Ucq_QdJdPnFhTXRyYEynKleOMj7bdgBGhfBefRub4F226UQPrFz8kypsr66FKZdy7bAnggIDzUFB0-629qLOmeOVeAMmOrq41uxICn3whK0sunt9bXfJTD68CxZvlgV8r1_jpjHqJqQzdio2--z0z0RQliX9WvEEmqfIvHaJpmWemvXejw1ywoglF0xQ4Gq39qB5CDe22zKr_cvKg1y7sJDvHw2Z4Iab_p5WdkxCMObAV3KbAQ3g7F-czkyRwoJiGOqAgau5aRUewWclryqNled5W8qiJ6m5RDIMQnYZyq-FTZgpjXMEUCIEwCqGbDrEYu0F2vVl6IC_5u3M3WVZGm3A6efdh55j1aAiEAooYy3e8q5eakgpXC8FPy0VdGphzV6sjbpuExuJdCqlk"
}
```

### Signing

Using the `Key Handle` from above, we can ask the U2F token to sign a challenge

```
[mdp u2fcli]$ u2fcli sig --appid https://mdp.im \
  --challenge anotherChallenge \
  --keyhandle 0JGeJ3MhvDzK_YjKhK4VkPOegGn0x3wxJENJ8J1JanozbSr8Elz2KRcARLh2sF__l_Vof2xiydPw6CEicpzs0A
Authenticating, press the button on your U2F device

{
  "Counter": 33,
  "Signature": "AQAAACEwRQIgetEQfx2p2SB7ch2JtvDYxjqTekMfZuDPjrJ0deNTXysCIQD5LehJ4gXf1vpJ37_XWefnSkRzwfwZ3Uffq7jWWTZYkw"
}
```

### Verifying

Finally, we can verify the `Signature` from above by providing it to u2fcli along with the `PublicKey` we recieved at registration.

```
[mdp u2fcli]$ u2fcli ver --appid https://mdp.im \
  --challenge anotherChallenge \
  --publickey BPQPBz7NV3LwksVwjbGdn7ODP5omKHt8CetrHnDZeUUxmFChHcKuYNHgLm0HdtsSD6p7cjrFZdb9mNOLg3huRcI \
  --signature AQAAACEwRQIgetEQfx2p2SB7ch2JtvDYxjqTekMfZuDPjrJ0deNTXysCIQD5LehJ4gXf1vpJ37_XWefnSkRzwfwZ3Uffq7jWWTZYkw

Signature verified
```


## Credit

This tool wouldn't be possible without the work of https://github.com/flynn/u2f which
handles all the interactions with the hardware and USB HID devices

