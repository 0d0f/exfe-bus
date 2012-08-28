package tokenmanager

import (
	"github.com/googollee/go-mysql"
	"testing"
	"time"
)

func TestTokenManager(t *testing.T) {
	client, err := mysql.DialTCP("127.0.0.1:3306", "root", "", "exfe_dev")
	if err != nil {
		panic(err)
	}
	mgr := New(client)
	resource := "resource"
	tk := "fjadsklfjkldasfdasiffjuoru21urjew"

	{
		_, err := mgr.GetResource(tk)
		if err == nil {
			t.Fatalf("get resource should failed")
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err == nil || ok {
			t.Errorf("tk(%s) verify with resource(%s) should failed", tk, resource)
		}
	}

	tk, err = mgr.GenerateToken(resource, time.Second)
	if err != nil {
		t.Fatalf("generate token failed: %s", err)
	}

	{
		r, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	time.Sleep(time.Second * 2)

	{
		r, err := mgr.GetResource(tk)
		if err != ExpiredError {
			t.Fatalf("get resource should expired")
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err == nil {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, time.Second)

	{
		r, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.ExpireToken(tk)

	{
		r, err := mgr.GetResource(tk)
		if err != ExpiredError {
			t.Fatalf("get resource should expired")
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err == nil {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, NeverExpire)

	{
		r, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	time.Sleep(time.Second * 2)

	{
		r, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	err = mgr.DeleteToken(tk)
	if err != nil {
		t.Errorf("delete fail: %s", err)
	}
}
