package cmd

import (
	"encoding/base64"
	"testing"
)

var verifyAppID = "https://mdp.im"
var verifyChallenge = "complexChallengeGoesHere"
var verifyPubKey, _ = base64.RawURLEncoding.DecodeString("BCwSU4NEplH1-UYlohohwnm68YU9H54RPlCffNWa83xlOQWQ19WrqS8J17HWXk5vAFF2gcGMn__1hXxCgyYxw_k")
var verifySignature, _ = base64.RawURLEncoding.DecodeString("AQAAAB8wRAIgJhy-8HvH-XOPakVnUggfzSSn0aUeObQ0TedWsjpli8ACIHiKVdcQhQ9EaHOAROL_CLcgKvJXp4-e46yMgmoIXWCt")

func TestVerify(t *testing.T) {
	if err := verify(verifyAppID, verifyChallenge, verifySignature, verifyPubKey); err != nil {
		t.Error(err)
	}
}

func TestVerifyAppID(t *testing.T) {
	if err := verify("https://mdp.io", verifyChallenge, verifySignature, verifyPubKey); err == nil {
		t.Error("Should fail with the incorrect appID")
	}
}
func TestVerifyChallenge(t *testing.T) {
	if err := verify(verifyAppID, "notTheRightChallenge", verifySignature, verifyPubKey); err == nil {
		t.Error("Should fail with the incorrect challlenge")
	}
}
