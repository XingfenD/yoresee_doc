package dto

type TemplateContainer int

const (
	TemplateContainerOwn TemplateContainer = iota
	TemplateContainerKnowledgeBase
	TemplateContainerPublic
)

type CreateTemplateRequest struct {
	UserExternalID          string
	TargetContainer         TemplateContainer
	KnowledgeBaseExternalID *string
	TemplateContent         string
}
