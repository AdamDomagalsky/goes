package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/AdamDomagalsky/goes/bank/token"
	"github.com/AdamDomagalsky/goes/bank/util"
	"google.golang.org/grpc/metadata"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata with credentials")
	}

	values := md.Get(util.AuthorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization token")
	}

	authorizationHeader := values[0]
	fields := strings.Fields(authorizationHeader)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != util.AuthorizationTypeKey {
		return nil, fmt.Errorf("invalid authorization header type: %s", authorizationType)
	}
	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid authorization token: %w", err)
	}

	return payload, nil
}
