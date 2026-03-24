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
	memberSvc := grpcserver.NewMembershipServiceServer()

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

	mux.Handle(pb.AuthService_QuerySideBarDisplay_FullMethodName, connect.NewUnaryHandler(
		pb.AuthService_QuerySideBarDisplay_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.QuerySideBarDisplayRequest]) (*connect.Response[pb.QuerySideBarDisplayResponse], error) {
			resp, err := authSvc.QuerySideBarDisplay(ctx, req.Msg)
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

	mux.Handle(pb.DocumentService_GetDocumentYjsSnapshot_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_GetDocumentYjsSnapshot_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetDocumentYjsSnapshotRequest]) (*connect.Response[pb.GetDocumentYjsSnapshotResponse], error) {
			resp, err := docSvc.GetDocumentYjsSnapshot(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_SaveDocumentYjsSnapshot_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_SaveDocumentYjsSnapshot_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.SaveDocumentYjsSnapshotRequest]) (*connect.Response[pb.SaveDocumentYjsSnapshotResponse], error) {
			resp, err := docSvc.SaveDocumentYjsSnapshot(ctx, req.Msg)
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

	mux.Handle(pb.DocumentService_UpdateDocumentMeta_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_UpdateDocumentMeta_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateDocumentMetaRequest]) (*connect.Response[pb.UpdateDocumentMetaResponse], error) {
			resp, err := docSvc.UpdateDocumentMeta(ctx, req.Msg)
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

	mux.Handle(pb.KnowledgeBaseService_CreateKnowledgeBase_FullMethodName, connect.NewUnaryHandler(
		pb.KnowledgeBaseService_CreateKnowledgeBase_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateKnowledgeBaseRequest]) (*connect.Response[pb.CreateKnowledgeBaseResponse], error) {
			resp, err := kbSvc.CreateKnowledgeBase(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.KnowledgeBaseService_ListRecentKnowledgeBases_FullMethodName, connect.NewUnaryHandler(
		pb.KnowledgeBaseService_ListRecentKnowledgeBases_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListRecentKnowledgeBasesRequest]) (*connect.Response[pb.ListRecentKnowledgeBasesResponse], error) {
			resp, err := kbSvc.ListRecentKnowledgeBases(ctx, req.Msg)
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

	mux.Handle(pb.MembershipService_ListUserGroups_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_ListUserGroups_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListUserGroupsRequest]) (*connect.Response[pb.ListUserGroupsResponse], error) {
			resp, err := memberSvc.ListUserGroups(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_GetUserGroup_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_GetUserGroup_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetUserGroupRequest]) (*connect.Response[pb.GetUserGroupResponse], error) {
			resp, err := memberSvc.GetUserGroup(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_CreateUserGroup_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_CreateUserGroup_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateUserGroupRequest]) (*connect.Response[pb.CreateUserGroupResponse], error) {
			resp, err := memberSvc.CreateUserGroup(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_UpdateUserGroup_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_UpdateUserGroup_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateUserGroupRequest]) (*connect.Response[pb.UpdateUserGroupResponse], error) {
			resp, err := memberSvc.UpdateUserGroup(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_DeleteUserGroup_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_DeleteUserGroup_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.DeleteUserGroupRequest]) (*connect.Response[pb.DeleteUserGroupResponse], error) {
			resp, err := memberSvc.DeleteUserGroup(ctx, req.Msg)
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

	mux.Handle(pb.DocumentService_CreateTemplate_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_CreateTemplate_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateTemplateRequest]) (*connect.Response[pb.CreateTemplateResponse], error) {
			resp, err := docSvc.CreateTemplate(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_GetTemplate_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_GetTemplate_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetTemplateRequest]) (*connect.Response[pb.GetTemplateResponse], error) {
			resp, err := docSvc.GetTemplate(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_ListTemplates_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_ListTemplates_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListTemplatesRequest]) (*connect.Response[pb.ListTemplatesResponse], error) {
			resp, err := docSvc.ListTemplates(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.DocumentService_ListRecentTemplates_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_ListRecentTemplates_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListRecentTemplatesRequest]) (*connect.Response[pb.ListRecentTemplatesResponse], error) {
			resp, err := docSvc.ListRecentTemplates(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))
}
