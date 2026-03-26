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
	settingSvc := grpcserver.NewSettingServiceServer()
	memberSvc := grpcserver.NewMembershipServiceServer()
	inviteSvc := grpcserver.NewInvitationServiceServer()
	notifySvc := grpcserver.NewNotificationServiceServer()
	commentSvc := grpcserver.NewCommentServiceServer()

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
	mux.Handle(pb.AuthService_QueryTopNavDisplay_FullMethodName, connect.NewUnaryHandler(
		pb.AuthService_QueryTopNavDisplay_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.QueryTopNavDisplayRequest]) (*connect.Response[pb.QueryTopNavDisplayResponse], error) {
			resp, err := authSvc.QueryTopNavDisplay(ctx, req.Msg)
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
	mux.Handle(pb.DocumentService_RecordRecentDocument_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_RecordRecentDocument_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.RecordRecentDocumentRequest]) (*connect.Response[pb.RecordRecentDocumentResponse], error) {
			resp, err := docSvc.RecordRecentDocument(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))
	mux.Handle(pb.DocumentService_ListRecentDocuments_FullMethodName, connect.NewUnaryHandler(
		pb.DocumentService_ListRecentDocuments_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListRecentDocumentsRequest]) (*connect.Response[pb.ListRecentDocumentsResponse], error) {
			resp, err := docSvc.ListRecentDocuments(ctx, req.Msg)
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

	mux.Handle(pb.SettingService_GetSettings_FullMethodName, connect.NewUnaryHandler(
		pb.SettingService_GetSettings_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetSettingsRequest]) (*connect.Response[pb.GetSettingsResponse], error) {
			resp, err := settingSvc.GetSettings(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.SettingService_UpdateSettings_FullMethodName, connect.NewUnaryHandler(
		pb.SettingService_UpdateSettings_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateSettingsRequest]) (*connect.Response[pb.UpdateSettingsResponse], error) {
			resp, err := settingSvc.UpdateSettings(ctx, req.Msg)
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

	mux.Handle(pb.MembershipService_ListUsers_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_ListUsers_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
			resp, err := memberSvc.ListUsers(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_UpdateUser_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_UpdateUser_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
			resp, err := memberSvc.UpdateUser(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_ListUserGroupMembers_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_ListUserGroupMembers_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListUserGroupMembersRequest]) (*connect.Response[pb.ListUserGroupMembersResponse], error) {
			resp, err := memberSvc.ListUserGroupMembers(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_ListOrgNodes_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_ListOrgNodes_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListOrgNodesRequest]) (*connect.Response[pb.ListOrgNodesResponse], error) {
			resp, err := memberSvc.ListOrgNodes(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_GetOrgNode_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_GetOrgNode_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.GetOrgNodeRequest]) (*connect.Response[pb.GetOrgNodeResponse], error) {
			resp, err := memberSvc.GetOrgNode(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_CreateOrgNode_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_CreateOrgNode_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateOrgNodeRequest]) (*connect.Response[pb.CreateOrgNodeResponse], error) {
			resp, err := memberSvc.CreateOrgNode(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_UpdateOrgNode_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_UpdateOrgNode_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateOrgNodeRequest]) (*connect.Response[pb.UpdateOrgNodeResponse], error) {
			resp, err := memberSvc.UpdateOrgNode(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_DeleteOrgNode_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_DeleteOrgNode_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.DeleteOrgNodeRequest]) (*connect.Response[pb.DeleteOrgNodeResponse], error) {
			resp, err := memberSvc.DeleteOrgNode(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_MoveOrgNode_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_MoveOrgNode_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.MoveOrgNodeRequest]) (*connect.Response[pb.MoveOrgNodeResponse], error) {
			resp, err := memberSvc.MoveOrgNode(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.MembershipService_ListOrgNodeMembers_FullMethodName, connect.NewUnaryHandler(
		pb.MembershipService_ListOrgNodeMembers_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListOrgNodeMembersRequest]) (*connect.Response[pb.ListOrgNodeMembersResponse], error) {
			resp, err := memberSvc.ListOrgNodeMembers(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.InvitationService_CreateInvitation_FullMethodName, connect.NewUnaryHandler(
		pb.InvitationService_CreateInvitation_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateInvitationRequest]) (*connect.Response[pb.CreateInvitationResponse], error) {
			resp, err := inviteSvc.CreateInvitation(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.InvitationService_ListInvitations_FullMethodName, connect.NewUnaryHandler(
		pb.InvitationService_ListInvitations_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListInvitationsRequest]) (*connect.Response[pb.ListInvitationsResponse], error) {
			resp, err := inviteSvc.ListInvitations(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.InvitationService_UpdateInvitation_FullMethodName, connect.NewUnaryHandler(
		pb.InvitationService_UpdateInvitation_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateInvitationRequest]) (*connect.Response[pb.UpdateInvitationResponse], error) {
			resp, err := inviteSvc.UpdateInvitation(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.InvitationService_DeleteInvitation_FullMethodName, connect.NewUnaryHandler(
		pb.InvitationService_DeleteInvitation_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.DeleteInvitationRequest]) (*connect.Response[pb.DeleteInvitationResponse], error) {
			resp, err := inviteSvc.DeleteInvitation(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.InvitationService_ListInvitationRecords_FullMethodName, connect.NewUnaryHandler(
		pb.InvitationService_ListInvitationRecords_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListInvitationRecordsRequest]) (*connect.Response[pb.ListInvitationRecordsResponse], error) {
			resp, err := inviteSvc.ListInvitationRecords(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.NotificationService_CreateNotification_FullMethodName, connect.NewUnaryHandler(
		pb.NotificationService_CreateNotification_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateNotificationRequest]) (*connect.Response[pb.CreateNotificationResponse], error) {
			resp, err := notifySvc.CreateNotification(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.NotificationService_ListNotifications_FullMethodName, connect.NewUnaryHandler(
		pb.NotificationService_ListNotifications_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListNotificationsRequest]) (*connect.Response[pb.ListNotificationsResponse], error) {
			resp, err := notifySvc.ListNotifications(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.NotificationService_MarkNotificationsRead_FullMethodName, connect.NewUnaryHandler(
		pb.NotificationService_MarkNotificationsRead_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.MarkNotificationsReadRequest]) (*connect.Response[pb.MarkNotificationsReadResponse], error) {
			resp, err := notifySvc.MarkNotificationsRead(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.NotificationService_MarkAllNotificationsRead_FullMethodName, connect.NewUnaryHandler(
		pb.NotificationService_MarkAllNotificationsRead_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.MarkAllNotificationsReadRequest]) (*connect.Response[pb.MarkAllNotificationsReadResponse], error) {
			resp, err := notifySvc.MarkAllNotificationsRead(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.CommentService_CreateDocumentComment_FullMethodName, connect.NewUnaryHandler(
		pb.CommentService_CreateDocumentComment_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.CreateDocumentCommentRequest]) (*connect.Response[pb.CreateDocumentCommentResponse], error) {
			resp, err := commentSvc.CreateDocumentComment(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.CommentService_ListDocumentComments_FullMethodName, connect.NewUnaryHandler(
		pb.CommentService_ListDocumentComments_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.ListDocumentCommentsRequest]) (*connect.Response[pb.ListDocumentCommentsResponse], error) {
			resp, err := commentSvc.ListDocumentComments(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.CommentService_DeleteDocumentComment_FullMethodName, connect.NewUnaryHandler(
		pb.CommentService_DeleteDocumentComment_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.DeleteDocumentCommentRequest]) (*connect.Response[pb.DeleteDocumentCommentResponse], error) {
			resp, err := commentSvc.DeleteDocumentComment(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(resp), nil
		},
		opts...,
	))

	mux.Handle(pb.CommentService_UpdateDocumentComment_FullMethodName, connect.NewUnaryHandler(
		pb.CommentService_UpdateDocumentComment_FullMethodName,
		func(ctx context.Context, req *connect.Request[pb.UpdateDocumentCommentRequest]) (*connect.Response[pb.UpdateDocumentCommentResponse], error) {
			resp, err := commentSvc.UpdateDocumentComment(ctx, req.Msg)
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
