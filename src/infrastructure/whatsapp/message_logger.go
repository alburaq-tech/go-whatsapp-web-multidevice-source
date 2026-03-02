package whatsapp

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aldinokemal/go-whatsapp-web-multidevice/config"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	msgLogger     *logrus.Logger
	msgLoggerOnce sync.Once
)

func getMsgLogger() *logrus.Logger {
	msgLoggerOnce.Do(func() {
		msgLogger = logrus.New()
		msgLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "time",
			},
		})
		msgLogger.SetLevel(logrus.InfoLevel)

		logPath := config.LogMessagesPath
		if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
			logrus.Warnf("Failed to create message log directory: %v", err)
			msgLogger.SetOutput(os.Stdout)
			return
		}

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Warnf("Failed to open message log file %s: %v", logPath, err)
			msgLogger.SetOutput(os.Stdout)
			return
		}

		msgLogger.SetOutput(f)
	})
	return msgLogger
}

// LogIncomingMessage writes a structured JSON log entry for an incoming message.
func LogIncomingMessage(evt *events.Message, client *whatsmeow.Client) {
	deviceID := ""
	if client != nil && client.Store != nil && client.Store.ID != nil {
		deviceID = client.Store.ID.ToNonAD().String()
	}

	msgType := string(evt.Info.Type)
	if msgType == "" {
		msgType = detectIncomingType(evt.Message)
	}

	fields := logrus.Fields{
		"device_id":  deviceID,
		"message_id": evt.Info.ID,
		"from":       evt.Info.Sender.ToNonAD().String(),
		"to":         evt.Info.Chat.ToNonAD().String(),
		"is_from_me": evt.Info.IsFromMe,
		"type":       msgType,
		"timestamp":  evt.Info.Timestamp.UTC().Format(time.RFC3339),
	}
	if built := utils.BuildEventMessage(evt); built.Text != "" {
		fields["text"] = built.Text
	}

	getMsgLogger().WithFields(fields).Info("message")
}

// LogOutgoingMessage writes a structured JSON log entry for an outgoing message.
func LogOutgoingMessage(msgID, senderJID string, recipient types.JID, msg *waE2E.Message, content string, ts time.Time) {
	fields := logrus.Fields{
		"device_id":  senderJID,
		"message_id": msgID,
		"from":       senderJID,
		"to":         recipient.ToNonAD().String(),
		"is_from_me": true,
		"type":       detectOutgoingType(msg, content),
		"timestamp":  ts.UTC().Format(time.RFC3339),
	}
	if content != "" {
		fields["text"] = content
	}

	getMsgLogger().WithFields(fields).Info("message")
}

func detectIncomingType(msg *waE2E.Message) string {
	if msg == nil {
		return "unknown"
	}
	if msg.GetImageMessage() != nil {
		return "image"
	}
	if msg.GetVideoMessage() != nil {
		return "video"
	}
	if msg.GetAudioMessage() != nil {
		return "audio"
	}
	if msg.GetDocumentMessage() != nil {
		return "document"
	}
	if msg.GetStickerMessage() != nil {
		return "sticker"
	}
	if msg.GetLocationMessage() != nil {
		return "location"
	}
	if msg.GetContactMessage() != nil {
		return "contact"
	}
	if msg.GetReactionMessage() != nil {
		return "reaction"
	}
	if msg.GetPollCreationMessage() != nil {
		return "poll"
	}
	return "text"
}

func detectOutgoingType(msg *waE2E.Message, content string) string {
	if msg == nil {
		if content != "" {
			return "text"
		}
		return "unknown"
	}
	if msg.GetImageMessage() != nil {
		return "image"
	}
	if msg.GetVideoMessage() != nil {
		return "video"
	}
	if msg.GetAudioMessage() != nil {
		return "audio"
	}
	if msg.GetDocumentMessage() != nil {
		return "document"
	}
	if msg.GetStickerMessage() != nil {
		return "sticker"
	}
	if msg.GetLocationMessage() != nil {
		return "location"
	}
	if msg.GetContactMessage() != nil {
		return "contact"
	}
	if msg.GetPollCreationMessage() != nil {
		return "poll"
	}
	return "text"
}
