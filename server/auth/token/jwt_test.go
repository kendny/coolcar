package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAgtJXT3aWFNJ73Z3qGxWp4DHLnnHM++dMfTfKN3ExRG5p86VI
JGSBIfKY8DV/EvfBSmFM8x/Mo1H3hb1gjYZXRreePf9P2wRE88nIkl/rztkkjAn+
u6OdvzjY8UoxKRF3P8jWPLnCoOGIqK71yzDD58QM5F1QnSwjctnahHbUZ8MNCPq0
LzLj7iW99Njywjqxw3saopG/PC1Skr4pCObwWIuepIeQqc90RdZUCuyw4L7UW3Wf
8V3jj/3ROUrogLPLFtjHjV7DZIvg6ZQj9h6lJpa5Dd9xJlcLSnMu3hZE+6vBwhqx
wjZNjZnWDtTdNZ4mtVbB5/g782NhSPZnUIRukwIDAQABAoIBAFy8MGpC/OarwmR6
Aa7Lb41enyGlOBLFhqIo/y7IDY2d23t8iaRKnoNrDmEQ+xB+EkXgrYSW4SBRBW6H
rT8FoS6fEQyPxJLT9vD99DFWz7jkkpS2rR1kQXmBGcAJwMmE+Qx9GOkZIU+cMJyq
0dAEtIrBngXv5CpWVJagudc2ySD3xhJXwLzWINBd7Qd/KZ5cWIhKFRIQT3aG2N8D
R88SPaMzTZ1dFTcUAgFuBgY6GevSMVVHd9GnmvkKrKgEbE6PG8llzq/jpZ0fp0My
pDFzoO7clOFigSB1y7vfCzvo5DNbsc0u5wZkbBzx9j96YxjSE1TkEs4SbKd/FKp1
E417F8kCgYEA7ax5f9ruL7BxLZbWO5P3AR6NRjIfaGvNE5+CPnOUnK2PT0UMFc8S
Q3ykM5XtEpdsl2zVIJZMXqJMrk5h6DFGqqtWe2GxDX7IS8ko1iEe7e11QSPrEhKC
1k5lwDC6fIDKsenEMiGA5XmndxCFLX/oYDAH5OaRioWxcpNcLNWS5B0CgYEAjOis
8Moh8ijT3Hm8/axdwrnjMVFEK2MHqxrdUa1A1/0K15lkwGPkuORlNNPuppdTrOb+
/Cln3Wl9/WduC1TQL73l9Z/0BBCYbYpMQupPEeUxQHe7PGqySfhlDcPHqqJHk+0e
l/oXbfdyOVaCJenujfYgkk4Bp5IZQ5FUFH++vm8CgYEAtQfQK6CYF97u9eiUGSEk
3Mdml/cJkUG7HJ08WVEz+vr/00MmZ14n7Rt8/oXN44FPBy+wTnfsh5Bbk9DEJlWS
G+ERqDzK3RqaeY8o+aCUrGlYDFvNayCKY62nNvJmuPpoaYdDF2QJh8BX3ArLRdLN
Lqam/KTaaoZWmIzBgqzSi8kCgYA/rAwpqkz4jfZeFCyQPEkJ6tF8wYoaSN94M+V6
ON2qO8+gaNIcFYeO/LW2z2VxpcPLx33FGDi01ix2SzxwplyEljzJZwPuqkkWhn27
ZwFfxr8gsHnM7TGvNy0CsUsSEc5iS62dYwcfS8cznaGl3DVNtMA3HESnId+EprDl
qmvYxQKBgQDKG/cZlRuS+xlubqCmidqfl9dSu4CwMLBOzLamin9Xy4xSlDCfAuk4
yzsb+L+gGZpGl+ij/YICe0eqBf+yAbmLl1LcJ5y+nk/SM7d1VC2z4Ox16haVDxiu
z9y/FQlNNI1gQgy0FKdh9dmTD2hH8j5QA9Hlk974DxkZK1Hzq8FYtg==
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v\n", err)
		panic(err)
	}

	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		//return time.Now()
		return time.Unix(1666537113, 0)
	}
	tkn, err := g.GenerateToken("6353d102b6e4358b7c31e0ec", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}
	fmt.Println(tkn)
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY1NDQzMTMsImlhdCI6MTY2NjUzNzExMywiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjM1M2QxMDJiNmU0MzU4YjdjMzFlMGVjIn0.ZUr0b1x7G4qj01far2KA1eO1sHqA6GI8I0zpxItJ364NoUBh354jJ8IeJRr1Ls-udkxNJDGzljv9jX_WB8AQXr-9gFaBBXBNH2myOXegUcV3-q-AU1N9ImHTHYR1JONwJ6uUThiYsJ71EomweViXXBq9KQtZ7mM3PoxiFi3Jmlw-PNCWm0JfsdbEyj1M0QcWIYJB1DehZ21BRQr7TKLEYcSyXc-ragcONmrh06bneMKEV9ld_2cp1v2I-9apS9NUnxwkxK3rzWX39N7oVpU4gD5kx32OqpBxeOye925PidPmO2B9JYGgliRins-hXxiWw5EA4B9qmFcVhbbpU-IWbA"
	if tkn != want {
		t.Errorf("worng token gennerated. want:%q; got: %q", want, tkn)
	}
}
