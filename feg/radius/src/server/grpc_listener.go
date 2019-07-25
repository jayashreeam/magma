/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package server

import (
	"context"
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/counters"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/protos"
	"fbc/cwf/radius/session"
	"fmt"
	"math/rand"
	"net"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPCListener listens to gRpc
type GRPCListener struct {
	GrpcServer    *grpc.Server
	Server        *Server
	Config        config.ListenerConfig
	Modules       []Module
	HandleRequest modules.Middleware
	dupDropped    uint32
	ready         chan bool
}

// GetModules override
func (l *GRPCListener) GetModules() []Module {
	return l.Modules
}

// SetModules override
func (l *GRPCListener) SetModules(m []Module) {
	l.Modules = m
}

// AppendModule override
func (l *GRPCListener) AppendModule(m *Module) {
	l.Modules = append(l.Modules, *m)
}

// GetConfig override
func (l *GRPCListener) GetConfig() config.ListenerConfig {
	return l.Config
}

// SetHandleRequest override
func (l *GRPCListener) SetHandleRequest(hr modules.Middleware) {
	l.HandleRequest = hr
}

// Init override
func (l *GRPCListener) Init(
	server *Server,
	serverConfig config.ServerConfig,
	listenerConfig config.ListenerConfig,
) error {
	if server == nil {
		return errors.New("cannot initialize GRPC listener with null server")
	}

	l.ready = make(chan bool, 1)
	l.Server = server
	return nil
}

// ListenAndServe override
func (l *GRPCListener) ListenAndServe() error {
	// Start listenning
	listenAddress := fmt.Sprintf(":%d", l.GetConfig().Port)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		l.ready <- false
		return errors.New("grpc listener: failed to open tcp connection" + listenAddress)
	}

	// Start serving
	l.GrpcServer = grpc.NewServer()
	protos.RegisterAuthorizationServer(l.GrpcServer, &authorizationServer{Listener: l})
	go func() {
		l.GrpcServer.Serve(lis)
	}()

	// Signal listener is ready
	go func() {
		l.ready <- true
	}()
	return nil
}

// GetHandleRequest override
func (l *GRPCListener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// Shutdown override
func (l *GRPCListener) Shutdown(ctx context.Context) error {
	return nil
}

// GetDupDropped override
func (l *GRPCListener) GetDupDropped() *uint32 {
	return &l.dupDropped
}

// Ready override
func (l *GRPCListener) Ready() chan bool {
	return l.ready
}

// SetConfig override
func (l *GRPCListener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}

type authorizationServer struct {
	Listener *GRPCListener
}

func (s *authorizationServer) Change(ctx context.Context, request *protos.ChangeRequest) (*protos.CoaResponse, error) {
	// Convert to RADIUS request
	req := radius.Request{
		Packet: &radius.Packet{
			Code:   radius.CodeDisconnectRequest,
			Secret: []byte(s.Listener.Server.config.Secret),
		},
	}

	// Handle RADIUS request
	return s.handleCoaRequest(request.Ctx, &req)
}

func (s *authorizationServer) Disconnect(ctx context.Context, request *protos.DisconnectRequest) (*protos.CoaResponse, error) {
	// Convert to RADIUS request
	req := radius.Request{
		Packet: &radius.Packet{
			Code:   radius.CodeDisconnectRequest,
			Secret: []byte(s.Listener.Server.config.Secret),
		},
	}

	// Handle RADIUS request
	return s.handleCoaRequest(request.Ctx, &req)
}

func (s *authorizationServer) handleCoaRequest(ctx *protos.Context, request *radius.Request) (*protos.CoaResponse, error) {
	if ctx == nil {
		return nil, errors.New("cannot handle a request without context")
	}

	if request == nil {
		return nil, errors.New("cannot handle a nil request")
	}

	// Get session ID from the request, if exists, and setup correlation ID
	srv := s.Listener.Server
	var correlationField = zap.Uint32("correlation", rand.Uint32())
	requestContext := modules.RequestContext{
		RequestID: correlationField.Integer,
		Logger:    srv.loggerFactory.Bg().With(correlationField),
		SessionID: ctx.SessionId,
		SessionStorage: session.NewSessionStorage(
			srv.multiSessionStorage,
			ctx.SessionId,
		),
	}

	// Load state, read CoA identifier and persist the state again
	state, err := requestContext.SessionStorage.Get()
	if err != nil {
		return nil, err
	}

	// Add Acct-Session-Id attribute
	request.Attributes = radius.Attributes{}
	request.Set(rfc2866.AcctSessionID_Type, radius.Attribute(state.AcctSessionID))

	// Set Identifier
	request.Identifier = state.NextCoAIdentifier
	state.NextCoAIdentifier = (state.NextCoAIdentifier + 1) % 0xFF

	// Handle
	counter := counters.NewOperation("handle_grpc").Start()
	res, err := s.Listener.HandleRequest(&requestContext, request)
	if err != nil {
		requestContext.Logger.Error("failed to handle request", zap.Error(err))
		counter.Failure("grpc_handle_error")
		return nil, err
	}
	if res == nil {
		requestContext.Logger.Error("got nil response")
		counter.Failure("grpc_nil_response")
		return nil, err
	}
	counter.Success()

	// Persist state
	err = srv.multiSessionStorage.Set(ctx.SessionId, *state)
	if err != nil {
		return nil, err
	}

	// Convert response to CoA response
	return &protos.CoaResponse{
		CoaResponseType: convertCoaCode(res.Code),
		Ctx:             ctx,
	}, nil
}

func convertCoaCode(code radius.Code) protos.CoaResponseCoaResponseTypeEnum {
	if code == radius.CodeCoAACK || code == radius.CodeDisconnectACK {
		return protos.CoaResponse_ACK
	}
	return protos.CoaResponse_NAK
}
