package ddpportal_test

import (
	"testing"

	ddpportal "gitlab.adeo.com/ddp-portal-api"
)

func Test_NewGroup(t *testing.T) {
	group := ddpportal.NewGroup("test", "2000xxxx")

	if group.ID == "" {
		t.Error("ID group can't be empty")
	}

	if group.Name != "test" {
		t.Error("Group name must be test")
	}

	if len(group.Members) != 1 {
		t.Error("Group must have one member")
	}

	var found bool
	for _, m := range group.Members {
		if m == "2000xxxx" {
			found = true
		}
	}

	if !found {
		t.Error("2000xxxx must be member of group")
	}
}

func Test_GroupAddMember_Ok(t *testing.T) {
	group := ddpportal.NewGroup("test", "2000xxxx")

	err := group.AddMember("2000yyyy")

	if err != nil {
		t.Error("2000yyyy must be added to group", err)
	}

	if len(group.Members) != 2 {
		t.Error("group must have 2 members")
	}

	var found bool
	for _, m := range group.Members {
		if m == "2000yyyy" {
			found = true
		}
	}

	if !found {
		t.Error("2000yyyy must be member of group")
	}

}

func Test_GroupAddMember_Error(t *testing.T) {
	group := ddpportal.NewGroup("test", "2000xxxx")

	err := group.AddMember("2000xxxx")

	if err == nil {
		t.Error("2000xxxx cannot be added to group ")
	}

}
