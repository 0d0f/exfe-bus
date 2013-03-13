package mail

import (
	"bytes"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"model"
	"net/mail"
	"testing"
)

func TestFindAddress(t *testing.T) {
	type Test struct {
		mails   []string
		pattern string
		domain  string
		ok      bool
		result  string
	}
	var tests = []Test{
		{[]string{"1@domain.com", "2@domain.com"}, "x", "exfe.com", false, ""},
		{[]string{"x@exfe.com", "1@domain.com", "2@domain.com"}, "x", "exfe.com", true, "[]"},
		{[]string{"x+c123@exfe.com", "x@exfe.com", "1@domain.com", "2@domain.com"}, "x\\+c([0-9]+)", "exfe.com", true, "[123]"},
	}

	for i, test := range tests {
		var addresses []*mail.Address
		for _, addr := range test.mails {
			addresses = append(addresses, &mail.Address{"", addr})
		}
		ok, match := findAddress(test.pattern, addresses)
		assert.Equal(t, ok, test.ok, fmt.Sprintf("test %d", i))
		if !test.ok {
			continue
		}
		assert.Equal(t, fmt.Sprintf("%v", match), test.result, fmt.Sprintf("test %d", i))
	}
}

func TestGetContent(t *testing.T) {
	buf1 := bytes.NewBufferString(`Delivered-To: x@0d0f.com
Received: by 10.229.56.222 with SMTP id z30csp34666qcg;
        Fri, 22 Feb 2013 01:39:34 -0800 (PST)
X-Received: by 10.50.157.163 with SMTP id wn3mr580297igb.89.1361525974497;
        Fri, 22 Feb 2013 01:39:34 -0800 (PST)
Return-Path: <googollee@gmail.com>
Received: from mail-ia0-x231.google.com (ia-in-x0231.1e100.net [2607:f8b0:4001:c02::231])
        by mx.google.com with ESMTPS id de2si1921845igb.47.2013.02.22.01.39.33
        (version=TLSv1 cipher=ECDHE-RSA-RC4-SHA bits=128/128);
        Fri, 22 Feb 2013 01:39:33 -0800 (PST)
Received-SPF: pass (google.com: domain of googollee@gmail.com designates 2607:f8b0:4001:c02::231 as permitted sender) client-ip=2607:f8b0:4001:c02::231;
Authentication-Results: mx.google.com;
       spf=pass (google.com: domain of googollee@gmail.com designates 2607:f8b0:4001:c02::231 as permitted sender) smtp.mail=googollee@gmail.com;
       dkim=pass header.i=@gmail.com
Received: by mail-ia0-f177.google.com with SMTP id o25so375608iad.22
        for <x@0d0f.com>; Fri, 22 Feb 2013 01:39:33 -0800 (PST)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
        d=gmail.com; s=20120113;
        h=x-received:mime-version:from:date:message-id:subject:to:cc
         :content-type;
        bh=JbKKJ00a5ssFJRhFuzXquEL+VcROcZWLhlPWr331pz0=;
        b=JqJRYeQUq09oqjFSbwwMxIHZr5FVCFLHOtg9Tof2pXe9OLmv47l3c416w+sHj0Bfwr
         l/7qyLRh9w4w2RUoNqnrp1Me6OwaPbcOB9INOfIS0Vh7zXfm72nWMR3k8fO+2bAHPhry
         pST5mXMlQnWdKwCNCcesfuS43oMkbnOlgtow6x9Dop5bXBJ6VNr6H3fSBhUuQ4u5Bmju
         kFYii7Yed7mVK7suHAaeRFp5ZnB2R1F4/k/IThaRtIzjGe5YrPygUzwiEOBD62PC7zlb
         HFuFS2UXrX5kdi/9ZJ8Kgv3EJFghfHNhjGMfcwrYcGNzse5wkQ/XGPYDevHS2faNuD8f
         uFFA==
X-Received: by 10.50.88.228 with SMTP id bj4mr569047igb.85.1361525972908; Fri,
 22 Feb 2013 01:39:32 -0800 (PST)
MIME-Version: 1.0
Received: by 10.42.18.199 with HTTP; Fri, 22 Feb 2013 01:39:11 -0800 (PST)
From: Googol Lee <googollee@gmail.com>
Date: Fri, 22 Feb 2013 17:39:11 +0800
Message-ID: <CAOf82vPBU=a7cO5TfPqWZRP8ZXPCy9Dc8n8e1HpeCOky5c5Yng@mail.gmail.com>
Subject: test
To: =?UTF-8?B?W0RFVl0gRVhGRSDCt1jCtw==?= <x@0d0f.com>
Cc: =?UTF-8?B?R29vZ29sIExlZSAtIEdvb2dsZee6r+eIt+S7rO+8gemTgeihgOecn+axieWtkO+8ge+8gQ==?= <googollee@gmail.com>
Content-Type: multipart/alternative; boundary=e89a8f3ba2bff9d1f904d64cf745

--e89a8f3ba2bff9d1f904d64cf745
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: base64

Y2MNCg0KLS0gDQrmlrDnmoTnkIborrrku47lsJHmlbDkurrnmoTkuLvlvKDliLDkuIDnu5/lpKnk
uIvvvIzlubbkuI3mmK/lm6DkuLrov5nkuKrnkIborrror7TmnI3kuobliKvkurrmipvlvIPml6fo
p4LngrnvvIzogIzmmK/lm6DkuLrkuIDku6PkurrnmoTpgJ3ljrvjgIINCg==
--e89a8f3ba2bff9d1f904d64cf745
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: base64

PGRpdiBkaXI9Imx0ciI+Y2M8YnIgY2xlYXI9ImFsbCI+PGRpdj48YnI+PC9kaXY+LS0gPGJyPuaW
sOeahOeQhuiuuuS7juWwkeaVsOS6uueahOS4u+W8oOWIsOS4gOe7n+WkqeS4i++8jOW5tuS4jeaY
r+WboOS4uui/meS4queQhuiuuuivtOacjeS6huWIq+S6uuaKm+W8g+aXp+ingueCue+8jOiAjOaY
r+WboOS4uuS4gOS7o+S6uueahOmAneWOu+OAgg0KPC9kaXY+DQo=
--e89a8f3ba2bff9d1f904d64cf745--

`)
	buf2 := bytes.NewBufferString(`Delivered-To: x@0d0f.com
Received: by 10.229.56.222 with SMTP id z30csp34666qcg;
        Fri, 22 Feb 2013 01:39:34 -0800 (PST)
X-Received: by 10.50.157.163 with SMTP id wn3mr580297igb.89.1361525974497;
        Fri, 22 Feb 2013 01:39:34 -0800 (PST)
Return-Path: <googollee@gmail.com>
Received: from mail-ia0-x231.google.com (ia-in-x0231.1e100.net [2607:f8b0:4001:c02::231])
        by mx.google.com with ESMTPS id de2si1921845igb.47.2013.02.22.01.39.33
        (version=TLSv1 cipher=ECDHE-RSA-RC4-SHA bits=128/128);
        Fri, 22 Feb 2013 01:39:33 -0800 (PST)
Received-SPF: pass (google.com: domain of googollee@gmail.com designates 2607:f8b0:4001:c02::231 as permitted sender) client-ip=2607:f8b0:4001:c02::231;
Authentication-Results: mx.google.com;
       spf=pass (google.com: domain of googollee@gmail.com designates 2607:f8b0:4001:c02::231 as permitted sender) smtp.mail=googollee@gmail.com;
       dkim=pass header.i=@gmail.com
Received: by mail-ia0-f177.google.com with SMTP id o25so375608iad.22
        for <x@0d0f.com>; Fri, 22 Feb 2013 01:39:33 -0800 (PST)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
        d=gmail.com; s=20120113;
        h=x-received:mime-version:from:date:message-id:subject:to:cc
         :content-type;
        bh=JbKKJ00a5ssFJRhFuzXquEL+VcROcZWLhlPWr331pz0=;
        b=JqJRYeQUq09oqjFSbwwMxIHZr5FVCFLHOtg9Tof2pXe9OLmv47l3c416w+sHj0Bfwr
         l/7qyLRh9w4w2RUoNqnrp1Me6OwaPbcOB9INOfIS0Vh7zXfm72nWMR3k8fO+2bAHPhry
         pST5mXMlQnWdKwCNCcesfuS43oMkbnOlgtow6x9Dop5bXBJ6VNr6H3fSBhUuQ4u5Bmju
         kFYii7Yed7mVK7suHAaeRFp5ZnB2R1F4/k/IThaRtIzjGe5YrPygUzwiEOBD62PC7zlb
         HFuFS2UXrX5kdi/9ZJ8Kgv3EJFghfHNhjGMfcwrYcGNzse5wkQ/XGPYDevHS2faNuD8f
         uFFA==
X-Received: by 10.50.88.228 with SMTP id bj4mr569047igb.85.1361525972908; Fri,
 22 Feb 2013 01:39:32 -0800 (PST)
References: <CABTp9f9aWfbh9xzHM6DZ6v9jfJL1+RCBHN3U7otSEHx+QsRUAA@mail.gmail.com>
	<5A846D28-194D-42C5-9B56-FB4B58622FB8@gmail.com>
	<CAOf82vPhowLHNZFjEVzYDz3spgPHJ31bi78xHp5jM2C45Bok9g@mail.gmail.com>
	<CAOf82vPDJsT-wqVyzOMBOXygOcbZ_u2Oqj87sUNHTwKbHyAOng@mail.gmail.com>
	<CA+wfOabc5aJbAqnvp0WwTY8iyaRkjUuJz5GUNi6YVb=pt8D=KA@mail.gmail.com>
	<CA+wfOaYpi47wZf8-cjSCHO2=L=2K-jJEEeZ=Ld-eUagKWatpWQ@mail.gmail.com>
	<CALsf4ScfeuNtdn4qFQLn7Cj0NK5s8aPzj=0MWUiR2ASAFvhDkw@mail.gmail.com>
	<CAEWKm+kUVvepoq9j=LCGbpnfZi1h5wo31ebpz6R1Mo=zLFLRTw@mail.gmail.com>
	<0F3F126D-B5F3-4B29-B048-95973ED822EC@gmail.com>
	<CAKib_jBJDBwy12GCaVrr3bP8ELTypODWu+xOSjFgZpUtNyfqQQ@mail.gmail.com>
	<CABTp9f_u+-SC1CoeXTYDio96OxLXYKKf0LQLS4YVU7vdg9NV6A@mail.gmail.com>
	<3C6F4630-44A8-4B5E-900F-447E41C2224B@gmail.com>
	<CA+wfOab3nzg_iuQ5=4Ybun8+_3HBkEBGCQ66bVkaeJR97=Njrg@mail.gmail.com>
	<CAKib_jA+06MGRMhi_G2DF1QGO=5Vu0iRGtYzVPMTjN4mG9V8pw@mail.gmail.com>
	<CAKib_jAYLChNff4Gp1L7T4G64qYswmFonb8DS8arT=f2rFRtng@mail.gmail.com>
	<CABTp9f8pTJeJBorZrU7wq=UbPxhzT0PYTt00g4nbN_-k6FbNMw@mail.gmail.com>
	<DB295B32-900D-45B2-98D7-981D2B832A8A@gmail.com>
	<CA+wfOaZBWUP7cNVO6P7t3fO9e_Mr3VKaSzob=FGQLFOyH4wbLA@mail.gmail.com>
	<C9BFFE30-2C92-4BA6-88D5-BC24186F828D@gmail.com>
	<CAOf82vPV57SgJ0Gr4Wz8q1nBZAffkVzcTz0csH96k=NnZhjDJw@mail.gmail.com>
	<4321185B-34EF-43D0-A6B5-9638E3DB41FE@gmail.com>
	<CAOf82vM9361tEcCB2G_rkYy8a1yJRu6EHZvGj1MB7M-xb1JOMw@mail.gmail.com>
	<A345842B-7C38-4551-A94C-382E53A95748@gmail.com>
	<CAOf82vNbm-YGY8VqfgZDCi2rFia1HU6TYThVZ09CgydQuywagA@mail.gmail.com>
MIME-Version: 1.0
Received: by 10.42.18.199 with HTTP; Fri, 22 Feb 2013 01:39:11 -0800 (PST)
From: Googol Lee <googollee@gmail.com>
Date: Fri, 22 Feb 2013 17:39:11 +0800
Message-ID: <CAOf82vPBU=a7cO5TfPqWZRP8ZXPCy9Dc8n8e1HpeCOky5c5Yng@mail.gmail.com>
Subject: test
To: =?UTF-8?B?W0RFVl0gRVhGRSDCt1jCtw==?= <x@0d0f.com>
Cc: =?UTF-8?B?R29vZ29sIExlZSAtIEdvb2dsZee6r+eIt+S7rO+8gemTgeihgOecn+axieWtkO+8ge+8gQ==?= <googollee@gmail.com>
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: base64

PGRpdiBkaXI9Imx0ciI+Y2M8YnIgY2xlYXI9ImFsbCI+PGRpdj48YnI+PC9kaXY+LS0gPGJyPuaW
sOeahOeQhuiuuuS7juWwkeaVsOS6uueahOS4u+W8oOWIsOS4gOe7n+WkqeS4i++8jOW5tuS4jeaY
r+WboOS4uui/meS4queQhuiuuuivtOacjeS6huWIq+S6uuaKm+W8g+aXp+ingueCue+8jOiAjOaY
r+WboOS4uuS4gOS7o+S6uueahOmAneWOu+OAgg0KPC9kaXY+DQo=

`)
	mail1, err := mail.ReadMessage(buf1)
	if err != nil {
		t.Fatal(err)
	}
	mail2, err := mail.ReadMessage(buf2)
	if err != nil {
		t.Fatal(err)
	}
	type Test struct {
		mail    *mail.Message
		ok      bool
		content string
	}
	var tests = []Test{
		{mail1, true, "cc"},
		{mail2, true, "cc"},
	}
	config := new(model.Config)
	log, err := logger.New(logger.Stderr, "bot")
	if err != nil {
		t.Fatal(err)
	}
	config.Log = log
	worker, err := New(config, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, test := range tests {
		content, err := worker.getContent(test.mail)
		if !test.ok {
			if err == nil {
				t.Errorf("test %d should not ok", i)
			}
			continue
		} else if err != nil {
			t.Errorf("test %d should ok: %s", i, err)
			continue
		}
		assert.Equal(t, content, test.content, fmt.Sprintf("test %d", i))
	}
}
