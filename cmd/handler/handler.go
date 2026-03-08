package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

var commandHandlerTree handlerTreeNode
var treeMutex = &sync.RWMutex{}
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
	matched := MatchCommandHandler(event.Message)
	if matched.Handler == nil {
		return nil, nil
	}
	handlerContext.ArgText = strings.TrimSpace(string(matched.ArgText))
	handlerContext.TriggerCmd = matched.Command
	return matched.Handler.Handle(handlerContext)
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

// 定义统一的处理器接口
type CommandHandler interface {
	GetCommands() []string               // 获取该处理器负责的命令列表
	GetPriority() int                    // 处理器的优先级
	GetHelper() string                   // 帮助文本
	Handle(Context) (interface{}, error) // 处理方法，TODO: 返回的结构暂未定义
}

// 一个基本的处理器
type CommandHandlerBase struct {
	Commands   []string
	Priority   int
	Helper     string
	handleFunc func(Context) (interface{}, error)
}

func (h *CommandHandlerBase) GetCommands() []string {
	return h.Commands
}
func (h *CommandHandlerBase) GetPriority() int {
	return h.Priority
}
func (h *CommandHandlerBase) GetHelper() string {
	return h.Helper
}
func (b *CommandHandlerBase) Handle(ctx Context) (interface{}, error) {
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
	treeMutex.Lock()
	defer treeMutex.Unlock()
	for _, command := range handler.GetCommands() {
		runeCommand := []rune(command)
		prefix := make([]rune, 0, len(runeCommand))
		commandHandlerTree.Add(prefix, []rune(command), handler)
	}
}

// 从树中匹配处理器
func MatchCommandHandler(message string) matchedHandler {
	treeMutex.RLock()
	defer treeMutex.RUnlock()
	messageRune := []rune(message)
	matched := commandHandlerTree.Get(messageRune, 0, matchedHandler{})
	matched.ArgText = messageRune[matched.PrefixLength:]
	return matched
}

// 指令分隔符
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
	priority int                       // 指令优先级
	depth    int                       // 当前的指令树深度
	command  string                    // 当前指令树的指令
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
				log.Println("待注册的指令处理器优先级更高，替换已有的处理器")
			} else {
				// 否则返回，不添加到树
				return
			}
		}
		t.priority = handlerPriority
		t.command = string(prefix)
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

type matchedHandler struct {
	Command      string         // 触发的命令
	PrefixLength int            // 匹配到处理器的消息前缀长度
	ArgText      []rune         // 消息的参数（去掉匹配到的前缀之后的消息内容）
	Handler      CommandHandler // 处理器
}

// 获取指令的处理器，可能为空
// 将指令按字符在指令树中查找
func (t *handlerTreeNode) Get(command []rune, prefixLength int, macthed matchedHandler) matchedHandler {
	handlerPriority := 0
	if macthed.Handler != nil {
		handlerPriority = macthed.Handler.GetPriority()
	}
	// 如果传入的解析器为空或者优先级不高于当前的handler，将handler替换为当前的（如果有）
	if t.handler != nil &&
		((t.priority > 0 && t.priority <= handlerPriority) ||
			macthed.Handler == nil) {
		macthed.Command = t.command
		macthed.PrefixLength = prefixLength
		macthed.Handler = t.handler
	}
	// 查找命令的下一个字符
	var nextR rune
	for i, r := range command {
		prefixLength++
		if IsCommandSeg(r) {
			continue
		}
		nextR = unicode.ToLower(r)
		command = command[i+1:]
		break
	}
	// 如果没有下一个字符，返回当前高优先级节点
	if nextR == 0 {
		return macthed
	}
	// 如果没有子节点，返回当前高优先级节点
	if t.children == nil {
		return macthed
	}
	child := t.children[nextR]
	// 如果下一个字符不在子节点中，返回当前高优先级节点
	if child == nil {
		return macthed
	}
	// 如果在子节点中，在子节点中查找
	return child.Get(command, prefixLength, macthed)
}
