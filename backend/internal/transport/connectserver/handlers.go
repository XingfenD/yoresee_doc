package connectserver

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/XingfenD/yoresee_doc/internal/transport/grpcserver"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

func registerHandlers(mux *http.ServeMux, opts []connect.HandlerOption) {
	authSvc := grpcserver.NewAuthServiceServer()
	docSvc := grpcserver.NewDocumentServiceServer()
	kbSvc := grpcserver.NewKnowledgeBaseServiceServer()
	sysSvc := grpcserver.NewSystemServiceServer()

	mux.Handle(pb.AuthService_Login_FullMethodName, connect.NewUnaryHandler(
		pb.AuthService_Login_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.AuthLoginRequest]) (*connect.Response[pb.AuthLoginResponse], error) {
			resp, err := authSvc.Login(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.AuthService_Register_FullMethodName, connect.NewUnaryHandler(
		pb.AuthService_Register_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.AuthRegisterRequest]) (*connect.Response[pb.AuthRegisterResponse], error) {
			resp, err := authSvc.Register(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_ListDocuments_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_ListDocuments_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListDocumentsRequest]) (*connect.Response[pb.ListDocumentsResponse], error) {
			resp, err := docSvc.ListDocuments(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_GetDocumentContent_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_GetDocumentContent_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetDocumentContentRequest]) (*connect.Response[pb.GetDocumentContentResponse], error) {
			resp, err := docSvc.GetDocumentContent(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_GetOwnDocuments_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_GetOwnDocuments_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetOwnDocumentsRequest]) (*connect.Response[pb.GetOwnDocumentsResponse], error) {
			resp, err := docSvc.GetOwnDocuments(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_CreateDocument_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_CreateDocument_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateDocumentRequest]) (*connect.Response[pb.CreateDocumentResponse], error) {
			resp, err := docSvc.CreateDocument(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.KnowledgeBaseService_ListKnowledgeBases_FullMethodName, connect.NewUnaryHandler(
		pb.KnowledgeBaseService_ListKnowledgeBases_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListKnowledgeBasesRequest]) (*connect.Response[pb.ListKnowledgeBasesResponse], error) {
			resp, err := kbSvc.ListKnowledgeBases(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.KnowledgeBaseService_GetKnowledgeBase_FullMethodName, connect.NewUnaryHandler(
		pb.KnowledgeBaseService_GetKnowledgeBase_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetKnowledgeBaseRequest]) (*connect.Response[pb.GetKnowledgeBaseResponse], error) {
			resp, err := kbSvc.GetKnowledgeBase(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.SystemService_Health_FullMethodName, connect.NewUnaryHandler(
		pb.SystemService_Health_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.HealthRequest]) (*connect.Response[pb.HealthResponse], error) {
			resp, err := sysSvc.Health(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.SystemService_SystemInfo_FullMethodName, connect.NewUnaryHandler(
		pb.SystemService_SystemInfo_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.SystemInfoRequest]) (*connect.Response[pb.SystemInfoResponse], error) {
			resp, err := sysSvc.SystemInfo(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_UpdateDocument_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_UpdateDocument_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateDocumentRequest]) (*connect.Response[pb.UpdateDocumentResponse], error) {
			resp, err := docSvc.UpdateDocument(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))
}
