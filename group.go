package ddpportal

import (
	"bitbucket.org/symbol-it/voltiger/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/ksuid"
)

const (
	MemberAlreadyInGroup = utils.Error("Member already in group")
)

type Group struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Members []string `json:"members,omitempty"`
}

func NewGroup(name string, member string) Group {
	ID := ksuid.New()
	group := Group{ID: ID.String(), Name: name}
	group.Members = []string{member}

	return group
}

func (g *Group) AddMember(member string) error {
	for _, m := range g.Members {
		if m == member {
			return MemberAlreadyInGroup
		}
	}

	g.Members = append(g.Members, member)
	spew.Dump(g)
	return nil
}
