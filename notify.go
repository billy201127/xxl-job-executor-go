package xxl

import (
	"context"

	"gomod.pri/golib/notify"
)

func Alert(webhook, secret, content string, isAtAll bool) {
	nf, _ := notify.NewNotification(notify.NotificationConfig{
		Type: notify.DingTalk,
		Config: notify.Config{
			Webhook: webhook,
			Secret:  secret,
		},
	})

	if isAtAll {
		_ = nf.SendText(context.Background(), content, notify.AtAll())
	} else {
		_ = nf.SendText(context.Background(), content)
	}
}
