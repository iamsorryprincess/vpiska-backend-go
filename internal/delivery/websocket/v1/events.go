package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) upgradeEventConnection(writer http.ResponseWriter, request *http.Request) {
	eventId := request.URL.Query().Get("eventId")
	isValid, err := validateId(eventId)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, conn, err := h.upgradeConnection(writer, request)

	if err != nil {
		return
	}

	userInfo.EventID = eventId
	ch := make(chan []byte)
	sub := newSubscriber(ch)
	h.publisher.Subscribe(eventId, sub)
	websocketContext := &readContext{context: newSocketContext(request.Context()), conn: conn, subscriber: sub, userData: userInfo}

	go readMessages(websocketContext, h.eventHandler, h.eventCloseHandler)
	go writeMessages(h.pingPeriod, conn, ch)

	err = h.events.AddUserInfo(request.Context(), service.AddUserInfoInput{
		EventID: eventId,
		UserID:  userInfo.UserID,
	})

	if err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			h.logger.LogError(closeErr)
		}
		if domain.IsInternalError(err) {
			h.logger.LogError(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) eventHandler(ctx context.Context, userData userData, body []byte) {
	prefix, message, isFound := strings.Cut(string(body), "/")
	if !isFound {
		return
	}
	switch prefix {
	case "chatMessage":
		{
			err := h.events.SendChatMessage(ctx, service.ChatMessageInput{
				EventID:     userData.EventID,
				UserID:      userData.UserID,
				UserName:    userData.UserName,
				UserImageID: userData.UserImageID,
				Message:     message,
			})

			if err != nil {
				if domain.IsInternalError(err) {
					h.logger.LogError(err)
				}
			}
		}
	default:
		h.logger.LogWarning("unknown route")
	}
}

func (h *Handler) eventCloseHandler(ctx context.Context, userData userData, subscriber service.Subscriber) {
	h.publisher.Unsubscribe(userData.EventID, subscriber)
	err := h.events.RemoveUserInfo(ctx, userData.EventID, userData.UserID)

	if err != nil {
		if domain.IsInternalError(err) {
			h.logger.LogError(err)
		}
	}
}
