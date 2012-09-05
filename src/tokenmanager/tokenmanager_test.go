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
	mgr := New(client, "tokens")
	resource := "resource"
	tk := "fjadsklfjkldasfdasiffjuoru21urjew"

	{
		_, _, err := mgr.GetResource(tk)
		if err == nil {
			t.Fatalf("get resource should failed")
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err == nil || ok {
			t.Errorf("tk(%s) verify with resource(%s) should failed", tk, resource)
		}
	}

	tk, err = mgr.GenerateToken(resource, "", time.Second)
	if err != nil {
		t.Fatalf("generate token failed: %s", err)
	}

	{
		r, d, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
		if got, expect := d, ""; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, d, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := d, ""; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		err = mgr.UpdateData(tk, "abc")
		if err != nil {
			t.Errorf("tk(%s) update data failed: %s", tk, err)
		}

		ok, d, err = mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := d, "abc"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		err = mgr.UpdateData(tk, "123")
		if err != nil {
			t.Errorf("tk(%s) update data failed: %s", tk, err)
		}

		ok, d, err = mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
		if got, expect := d, "123"; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	time.Sleep(time.Second * 2)

	{
		r, _, err := mgr.GetResource(tk)
		if err != ExpiredError {
			t.Fatalf("get resource should expired")
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err == nil {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, time.Second)

	{
		r, _, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.ExpireToken(tk)

	{
		r, _, err := mgr.GetResource(tk)
		if err != ExpiredError {
			t.Fatalf("get resource should expired")
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err == nil {
			t.Errorf("tk(%s) verify with resource(%s) should expired", tk, resource)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	mgr.RefreshToken(tk, NeverExpire)

	{
		r, _, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
		if err != nil {
			t.Errorf("tk(%s) verify with resource(%s) failed: %s", tk, resource, err)
		}
		if !ok {
			t.Errorf("tk(%s) should verify with resource(%s), but not", tk, resource)
		}
	}

	time.Sleep(time.Second * 2)

	{
		r, _, err := mgr.GetResource(tk)
		if err != nil {
			t.Fatalf("get resource failed: %s", err)
		}
		if got, expect := r, resource; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}

		ok, _, err := mgr.VerifyToken(tk, resource)
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
