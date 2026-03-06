package handler

import "context"

type MessageType string

const (
	MessageTypePrivate MessageType = "private"
	MessageTypeGroup   MessageType = "group"
)

// 可能还有其它需要的字段，先这样吧
type Event struct {
	MessageType MessageType
	Message     string
	MessageId   string
	UserId      string
	SenderName  string
	GroupId     string
}

type Context interface {
	context.Context              // 继承go的上下文
	GetTriggerCmd() string       // 获取触发的命令
	GetArgs() string             // 获取命令参数
	GetMessageType() MessageType // 获取消息类型
	GetMessage() string          // 获取原始文本消息
	GetEvent() Event             // 获取原始事件对象
	GetMessageId() string        // 获取消息ID
	GetUserId() string           // 获取发送者ID
	GetSenderName() string       // 获取发送者名称
	GetGroupId() string          // 获取群号
}

type HandlerContext struct {
	context.Context        // 继承go的上下文
	TriggerCmd      string // 触发的命令
	ArgText         string // 命令参数\
	MessageType     MessageType
	Message         string // 原始文本消息
	Event           Event  // 原始事件对象
	MessageId       string // 消息ID
	UserId          string // 发送者ID
	SenderName      string // 发送者名称
	GroupId         string // 群号
}

func (h *HandlerContext) GetTriggerCmd() string {
	return h.TriggerCmd
}
func (h *HandlerContext) GetArgs() string {
	return h.ArgText
}
func (h *HandlerContext) GetMessageType() MessageType {
	return h.MessageType
}
func (h *HandlerContext) GetMessage() string {
	return h.Message
}
func (h *HandlerContext) GetEvent() Event {
	return h.Event
}
func (h *HandlerContext) GetMessageId() string {
	return h.MessageId
}

func (h *HandlerContext) GetUserId() string {
	return h.UserId
}
func (h *HandlerContext) GetSenderName() string {
	return h.SenderName
}
func (h *HandlerContext) GetGroupId() string {
	return h.GroupId
}
