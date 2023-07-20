package userinfo

//{
//    "iss": "http://127.0.0.1:5556/dex",
//    "sub": "Cg0wLTM4NS0yODA4OS0wEgRtb2Nr",
//    "aud": "app1",
//    "exp": 1689069658,
//    "iat": 1688983258,
//    "at_hash": "igWP_QRsD5PjeRCOx9eyUA",
//    "email": "kilgore@kilgore.trout",
//    "email_verified": true,
//    "groups": [
//        "authors"
//    ],
//    "name": "Kilgore Trout"
//}

type Userinfo struct {
	Subject  string `json:"sub"`
	Audience string `json:"aud"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
