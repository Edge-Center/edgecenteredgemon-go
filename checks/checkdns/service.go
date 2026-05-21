package checkdns

import (
	"github.com/Edge-Center/edgecenteredgemon-go/checks"
	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter"
)

type Service = checks.Service[Request, Response]

func New(r edgecenter.Requester) Service {
	return checks.NewService[Request, Response](r, "dns")
}
