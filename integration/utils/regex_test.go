package utils

import (
	"testing"
)

const bodyBytes = `"token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcm5hbWUxIiwiYWRtaW4iOnRydWUsInRlbmFudCI6Inh4eDEiLCJpZCI6IjY3YzliZDUwYmE5Yjc5ODkxMGM2MDEzOCIsInJvbGUiOiJvd25lciIsImNsYWltcyI6InVzZXIiLCJleHAiOjE2NDM0NjY0MzV9.YxliGp0CSnWhzAUzGRDjYIn5I-OOL6wHh4y-OYIYOwHNiZx86hiS2vDtGLnWUyH33EngMhyBbBfSV0NOFlRWMcN-RusikHSxu3oiqe-lvgG9wM9D5D5GcL8N96__HFBtWDLe4a383RcXNt1ufsxbofALofcYS_PiE4GoVeC1WselnztCaOUNykWAVwcnT41SF2J8sO3zzbQNPRdJ-GRDBGYaaSn5zR8NrRtmAxyVTEOytS5Y88vL3rnVjCXmlICpp93gsKHVbH8NPxM98PCzjCWCye_gWlW--V1kLt8YJgIRPsjJ6Nv0vip3fDRMFklS5mZjR6Hhu5B0xDN9svjNPg","user":"username1","name":"user1","id":"67c9bd50ba9b798910c60138","tenant":"xxx1","role":"owner","email":"user1.email.com"}`

func TestGetFieldValueFromJSON(t *testing.T) {
	val, err := GetFieldValueFromJSON("token", bodyBytes)
	if err != nil {
		t.Errorf("error \n")
	}

	if val != "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcm5hbWUxIiwiYWRtaW4iOnRydWUsInRlbmFudCI6Inh4eDEiLCJpZCI6IjY3YzliZDUwYmE5Yjc5ODkxMGM2MDEzOCIsInJvbGUiOiJvd25lciIsImNsYWltcyI6InVzZXIiLCJleHAiOjE2NDM0NjY0MzV9.YxliGp0CSnWhzAUzGRDjYIn5I-OOL6wHh4y-OYIYOwHNiZx86hiS2vDtGLnWUyH33EngMhyBbBfSV0NOFlRWMcN-RusikHSxu3oiqe-lvgG9wM9D5D5GcL8N96__HFBtWDLe4a383RcXNt1ufsxbofALofcYS_PiE4GoVeC1WselnztCaOUNykWAVwcnT41SF2J8sO3zzbQNPRdJ-GRDBGYaaSn5zR8NrRtmAxyVTEOytS5Y88vL3rnVjCXmlICpp93gsKHVbH8NPxM98PCzjCWCye_gWlW--V1kLt8YJgIRPsjJ6Nv0vip3fDRMFklS5mZjR6Hhu5B0xDN9svjNPg" {
		t.Errorf("doest not match")
	}
}
