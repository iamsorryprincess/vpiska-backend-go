package v1

import (
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) eventHandler(ctx socketContext, body []byte) {
	prefix, message, isFound := strings.Cut(string(body), "/")
	if !isFound {
		return
	}
	switch prefix {
	case "chatMessage":
		{
			err := h.services.Events.SendChatMessage(ctx, service.ChatMessageInput{
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
	h.services.Publisher.Unsubscribe(ctx.EventID, subscriber)
	err := h.services.Events.RemoveUserInfo(ctx, ctx.EventID, ctx.UserID)

	if err != nil {
		if domain.IsInternalError(err) {
			h.logger.LogError(err)
		}
	}
}

func (h *Handler) rangeHandler(ctx socketContext, body []byte) {
}

func (h *Handler) rangeCloseHandler(ctx socketContext, subscriber service.Subscriber) {

}
