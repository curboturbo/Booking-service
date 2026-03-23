package conference

import (
	"context"
	"fmt"
	port "test-backend-1-curboturbo/internal/port/outbound"
	"time"
)

type linkConferenceService struct{}


func NewLinkConferenceService() port.LinkConferenceProvider{
	return &linkConferenceService{}
}


func (l *linkConferenceService) RequestLink(ctx context.Context) (string,error) {
    select {
    case <-ctx.Done():
		return "", fmt.Errorf("call interrupt")
    case <-time.After(2 * time.Second):
        return "https://leetcode.com/maevec", nil
    }
}