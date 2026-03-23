package port

import "context"

type LinkConferenceProvider interface {
	RequestLink(ctx context.Context) (string,error)
}