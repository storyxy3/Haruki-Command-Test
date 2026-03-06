package handler

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"
)

var commandHandlerTree handlerTreeNode
var maxDepth int

const DefaultPriority = 100

// TODO: 返回的结果统一封装，具体怎么样之后再说
func Dispatch(ctx context.Context, event Event) (interface{}, error) {
	handlerContext := &HandlerContext{
		Context:     ctx,
		MessageType: event.MessageType,
		Message:     event.Message,
		Event:       event,
		MessageId:   event.MessageId,
		UserId:      event.UserId,
		SenderName:  event.SenderName,
		GroupId:     event.GroupId,
	}
	prefix := make([]rune, 0, min(maxDepth, len(event.Message)))
	triggerCmd, handler := commandHandlerTree.Get(prefix, []rune(event.Message), nil, nil)
	if handler == nil {
		return nil, nil
	}
	handlerContext.ArgText = strings.TrimSpace(event.Message[len(string(triggerCmd)):])
	handlerContext.TriggerCmd = clearCmd(triggerCmd)
	return handler.Handle(handlerContext)
}

func clearCmd(cmd []rune) string {
	result := make([]rune, 0, len(cmd))
	for _, c := range cmd {
		if IsCommandSeg(c) {
			continue
		}
		result = append(result, c)
	}
	return string(result)
}

type CommandHandler interface {
	GetCommands() []string               // 获取该处理器负责的命令列表
	GetPriority() int                    // 处理器的优先级
	GetHelper() string                   // 帮助文本
	Handle(Context) (interface{}, error) // 处理方法，TODO: 返回的结构暂未定义
}

type BaseCommandHandler struct {
	Commands   []string
	Priority   int
	Helper     string
	handleFunc func(Context) (interface{}, error)
}

func (h *BaseCommandHandler) GetCommands() []string {
	return h.Commands
}
func (h *BaseCommandHandler) GetPriority() int {
	return h.Priority
}
func (h *BaseCommandHandler) GetHelper() string {
	return h.Helper
}
func (b *BaseCommandHandler) Handle(ctx Context) (interface{}, error) {
	if b.handleFunc != nil {
		return b.handleFunc(ctx)
	}
	cmdName := "未定义"
	if len(b.Commands) > 0 {
		cmdName = b.Commands[0]
	}
	return nil, fmt.Errorf("命令处理器 %s 没有处理方法", cmdName)
}

// 将处理器添加到树
func RegisterCommandHandler(handler CommandHandler) {
	for _, command := range handler.GetCommands() {
		runeCommand := []rune(command)
		prefix := make([]rune, 0, len(runeCommand))
		commandHandlerTree.Add(prefix, []rune(command), handler)
	}
}

func IsCommandSeg(r rune) bool {
	switch r {
	case ' ':
		return true
	case '_':
		return true
	case '-':
		return true
	case '.':
		return true
	default:
		return false
	}
}

// 指令树结构，用于解析指令
type handlerTreeNode struct {
	priority int                       //指令优先级
	depth    int                       // 当前的指令树深度
	handler  CommandHandler            // 指令处理器
	children map[rune]*handlerTreeNode // 子节点
}

// 将指令处理器添加到指令解析树中
// 将指令的每一个字符加入到树中，如果是指令分隔符，就跳过它
// 如果指令走到了最后一个字符，而当前树节点没有处理器，将它注册到节点上，否则报一个警告
func (t *handlerTreeNode) Add(prefix, command []rune, handler CommandHandler) {
	// 查找下一个字符
	var nextR rune
	for i, r := range command {
		if IsCommandSeg(r) {
			continue
		}
		nextR = unicode.ToLower(r)
		prefix = append(prefix, r)
		command = command[i+1:]
		break
	}
	// 如果没有下一个字符，就将handler添加到当前节点
	if nextR == 0 {
		handlerPriority := handler.GetPriority()
		if t.handler != nil {
			fmt.Fprintf(os.Stderr, "前缀 %s 已被注册，已有的优先级：%d，待注册优先级：%d\n", string(prefix), t.priority, handlerPriority)
			// 如果待注册的handler优先级不为0，且当前handler优先级为0或大于待注册handler（优先级数值越低，优先级越高，但是0为最低优先）
			if handlerPriority > 0 && (handlerPriority < t.priority || t.priority == 0) {
				t.priority = handlerPriority
				t.depth = len(prefix)
				t.handler = handler
				fmt.Fprintf(os.Stderr, "待注册的指令解析器优先级更高，替换已有的解析器\n")
			}
			return
		}
		t.priority = handlerPriority
		t.depth = len(prefix)
		if t.depth > maxDepth {
			maxDepth = t.depth
		}
		t.handler = handler
		return
	}
	// 有下一个字符，就将处理器添加到子节点
	if t.children == nil {
		t.children = make(map[rune]*handlerTreeNode)
	}
	child := t.children[nextR]
	if child == nil {
		child = &handlerTreeNode{}
	}
	child.Add(prefix, command, handler)
	t.children[nextR] = child
}

// 获取指令的处理器，可能为空
// 将指令按字符在指令树中查找
func (t *handlerTreeNode) Get(prefix, command, handledPrefix []rune, handler CommandHandler) ([]rune, CommandHandler) {
	handlerPriority := 0
	if handler != nil {
		handlerPriority = handler.GetPriority()
	}
	// 如果传入的解析器为空或者优先级不高于当前的handler，将handler替换为当前的（如果有）
	if t.handler != nil &&
		((t.priority > 0 && t.priority <= handlerPriority) ||
			handler == nil) {
		handledPrefix = slices.Clone(prefix)
		handler = t.handler
	}
	// 查找命令的下一个字符
	var nextR rune
	for i, r := range command {
		prefix = append(prefix, r)
		if IsCommandSeg(r) {
			continue
		}
		nextR = unicode.ToLower(r)
		command = command[i+1:]
		break
	}
	// 如果没有下一个字符，返回当前高优先级节点
	if nextR == 0 {
		return handledPrefix, handler
	}
	// 如果没有子节点，返回当前高优先级节点
	if t.children == nil {
		return handledPrefix, handler
	}
	// 如果下一个字符在子节点中，在子节点中查找
	if child := t.children[nextR]; child != nil {
		return child.Get(prefix, command, handledPrefix, handler)
	}
	// 如果不在子节点中，返回当前高优先级节点
	return handledPrefix, handler
}
