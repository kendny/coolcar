package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgtJXT3aWFNJ73Z3qGxWp
4DHLnnHM++dMfTfKN3ExRG5p86VIJGSBIfKY8DV/EvfBSmFM8x/Mo1H3hb1gjYZX
RreePf9P2wRE88nIkl/rztkkjAn+u6OdvzjY8UoxKRF3P8jWPLnCoOGIqK71yzDD
58QM5F1QnSwjctnahHbUZ8MNCPq0LzLj7iW99Njywjqxw3saopG/PC1Skr4pCObw
WIuepIeQqc90RdZUCuyw4L7UW3Wf8V3jj/3ROUrogLPLFtjHjV7DZIvg6ZQj9h6l
Jpa5Dd9xJlcLSnMu3hZE+6vBwhqxwjZNjZnWDtTdNZ4mtVbB5/g782NhSPZnUIRu
kwIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Errorf("parse public key failed: %v\n", err)
	}

	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	//tkn := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY1NDQzMTMsImlhdCI6MTY2NjUzNzExMywiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjM1M2QxMDJiNmU0MzU4YjdjMzFlMGVjIn0.ZUr0b1x7G4qj01far2KA1eO1sHqA6GI8I0zpxItJ364NoUBh354jJ8IeJRr1Ls-udkxNJDGzljv9jX_WB8AQXr-9gFaBBXBNH2myOXegUcV3-q-AU1N9ImHTHYR1JONwJ6uUThiYsJ71EomweViXXBq9KQtZ7mM3PoxiFi3Jmlw-PNCWm0JfsdbEyj1M0QcWIYJB1DehZ21BRQr7TKLEYcSyXc-ragcONmrh06bneMKEV9ld_2cp1v2I-9apS9NUnxwkxK3rzWX39N7oVpU4gD5kx32OqpBxeOye925PidPmO2B9JYGgliRins-hXxiWw5EA4B9qmFcVhbbpU-IWbA"
	////固定时间
	//jwt.TimeFunc = func() time.Time {
	//	return time.Unix(1666537213, 0)
	//}
	//accountID, err := v.Verify(tkn)
	//
	//if err != nil {
	//	t.Errorf("verifition failed: %v\n", err)
	//}
	//
	//want := "6353d102b6e4358b7c31e0ec"
	//if accountID != want {
	//	t.Errorf("wrong account ID. want %q, got %q", want, accountID)
	//}

	//	 表格化测试
	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name: "valid_token",
			tkn:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY1NDQzMTMsImlhdCI6MTY2NjUzNzExMywiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjM1M2QxMDJiNmU0MzU4YjdjMzFlMGVjIn0.ZUr0b1x7G4qj01far2KA1eO1sHqA6GI8I0zpxItJ364NoUBh354jJ8IeJRr1Ls-udkxNJDGzljv9jX_WB8AQXr-9gFaBBXBNH2myOXegUcV3-q-AU1N9ImHTHYR1JONwJ6uUThiYsJ71EomweViXXBq9KQtZ7mM3PoxiFi3Jmlw-PNCWm0JfsdbEyj1M0QcWIYJB1DehZ21BRQr7TKLEYcSyXc-ragcONmrh06bneMKEV9ld_2cp1v2I-9apS9NUnxwkxK3rzWX39N7oVpU4gD5kx32OqpBxeOye925PidPmO2B9JYGgliRins-hXxiWw5EA4B9qmFcVhbbpU-IWbA",
			now:  time.Unix(1666537213, 0),
			want: "6353d102b6e4358b7c31e0ec",
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY1NDQzMTMsImlhdCI6MTY2NjUzNzExMywiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjM1M2QxMDJiNmU0MzU4YjdjMzFlMGVjIn0.ZUr0b1x7G4qj01far2KA1eO1sHqA6GI8I0zpxItJ364NoUBh354jJ8IeJRr1Ls-udkxNJDGzljv9jX_WB8AQXr-9gFaBBXBNH2myOXegUcV3-q-AU1N9ImHTHYR1JONwJ6uUThiYsJ71EomweViXXBq9KQtZ7mM3PoxiFi3Jmlw-PNCWm0JfsdbEyj1M0QcWIYJB1DehZ21BRQr7TKLEYcSyXc-ragcONmrh06bneMKEV9ld_2cp1v2I-9apS9NUnxwkxK3rzWX39N7oVpU4gD5kx32OqpBxeOye925PidPmO2B9JYGgliRins-hXxiWw5EA4B9qmFcVhbbpU-IWbA",
			now:     time.Unix(1666555513, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1666537213, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "11eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY1NDQzMTMsImlhdCI6MTY2NjUzNzExMywiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjM1M2QxMDJiNmU0MzU4YjdjMzFlMGVjIn0.ZUr0b1x7G4qj01far2KA1eO1sHqA6GI8I0zpxItJ364NoUBh354jJ8IeJRr1Ls-udkxNJDGzljv9jX_WB8AQXr-9gFaBBXBNH2myOXegUcV3-q-AU1N9ImHTHYR1JONwJ6uUThiYsJ71EomweViXXBq9KQtZ7mM3PoxiFi3Jmlw-PNCWm0JfsdbEyj1M0QcWIYJB1DehZ21BRQr7TKLEYcSyXc-ragcONmrh06bneMKEV9ld_2cp1v2I-9apS9NUnxwkxK3rzWX39N7oVpU4gD5kx32OqpBxeOye925PidPmO2B9JYGgliRins-hXxiWw5EA4B9qmFcVhbbpU-IWbA",
			now:     time.Unix(1666537213, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)
			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v\n", err)
			}

			if c.wantErr && err == nil {
				t.Errorf("want error; got no error;\n")
			}

			if accountID != c.want {
				t.Errorf("wrong account ID. want %q, got %q", c.want, accountID)
			}

		})
	}

}
