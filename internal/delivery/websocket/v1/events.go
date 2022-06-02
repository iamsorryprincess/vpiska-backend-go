package v1

import (
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) upgradeEventConnection(writer http.ResponseWriter, request *http.Request) {
	params := strings.Split(request.URL.RawQuery, "&")
	eventId := getQueryValue("eventId", params)
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

	ctx, conn, err := h.upgradeConnection(writer, request)

	if err != nil {
		return
	}

	ctx.EventID = eventId
	sub := newSubscriber(h.logger, conn)
	h.publisher.Subscribe(eventId, sub)
	err = h.events.AddUserInfo(request.Context(), service.AddUserInfoInput{
		EventID: eventId,
		UserID:  ctx.UserID,
	})

	if err != nil {
		h.publisher.Unsubscribe(eventId, sub)
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

	go h.readMessages(ctx, conn, sub, h.eventHandler, h.eventCloseHandler)
}

func (h *Handler) eventHandler(ctx socketContext, body []byte) {
	prefix, message, isFound := strings.Cut(string(body), "/")
	if !isFound {
		return
	}
	switch prefix {
	case "chatMessage":
		{
			err := h.events.SendChatMessage(ctx, service.ChatMessageInput{
				EventID:     ctx.EventID,
				UserID:      ctx.UserID,
				UserName:    ctx.UserName,
				UserImageID: ctx.UserImageID,
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

func (h *Handler) eventCloseHandler(ctx socketContext, subscriber service.Subscriber) {
	h.publisher.Unsubscribe(ctx.EventID, subscriber)
	err := h.events.RemoveUserInfo(ctx, ctx.EventID, ctx.UserID)

	if err != nil {
		if domain.IsInternalError(err) {
			h.logger.LogError(err)
		}
	}
}
