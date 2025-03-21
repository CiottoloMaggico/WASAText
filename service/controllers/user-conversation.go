package controllers

import (
	apierrors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/api/filter"
	"github.com/ciottolomaggico/wasatext/service/controllers/translators"
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
)

type UserConversationController interface {
	GetUserConversations(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error)
	GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error)
}

type UserConversationControllerImpl struct {
	Model             models.UserConversationModel
	ConversationModel models.ConversationModel
	Filter            filter.Filter
}

func (controller UserConversationControllerImpl) GetUserConversations(requestIssuerUUID string, paginationPs pagination.PaginationParams) (pagination.PaginatedView, error) {
	queryParameters, err := database.NewQueryParameters(paginationPs, controller.Filter)
	if err != nil {
		return pagination.PaginatedView{}, apierrors.InvalidUrlParameters()
	}

	conversationsCount, cursor, err := controller.Model.Count(requestIssuerUUID, queryParameters)
	if err != nil {
		return pagination.PaginatedView{}, err
	}

	queryParameters.Cursor = cursor
	var conversations []models.UserConversation
	if conversationsCount > 0 {
		conversations, cursor, err = controller.Model.GetUserConversations(requestIssuerUUID, queryParameters)

		if err != nil {
			return pagination.PaginatedView{}, err
		}
	}

	paginationPs.Cursor = cursor
	return pagination.ToPaginatedView(paginationPs, conversationsCount, translators.UserConversationListToSummaryView(conversations))
}

func (controller UserConversationControllerImpl) GetUserConversation(requestIssuerUUID string, conversationId int64) (views.UserConversationView, error) {
	conversation, err := controller.Model.GetUserConversation(requestIssuerUUID, conversationId)
	if err != nil {
		return views.UserConversationView{}, translators.DBErrorToApiError(err)
	}

	participants, err := controller.ConversationModel.GetConversationParticipants(conversationId)
	if err != nil {
		return views.UserConversationView{}, err
	}

	return translators.UserConversationToView(*conversation, participants), nil
}
