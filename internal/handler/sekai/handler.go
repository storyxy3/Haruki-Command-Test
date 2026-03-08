package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	sekairegion "Haruki-Command-Parser/internal/sekai_region"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
)

type SekaiHandlerContext struct {
	handler.HandlerContext
	region             *sekairegion.SekaiRegion // 区服
	originalTriggerCmd string                   // 原始触发命令，未去除区服前缀
	prefixArg          string                   // 额外前缀
	uidArg             string                   // UID参数 /u / uid / @
	flags              map[string]bool          // -verbose, -preview, -help 等开关
}

type SekaiCommandHandler struct {
	handler.CommandHandlerBase
	Regions    []*sekairegion.SekaiRegion
	PrefixArgs []string
	handleFunc func(SekaiHandlerContext) (interface{}, error)
}

func (s *SekaiHandlerContext) Region() *sekairegion.SekaiRegion {
	return s.region
}
func (s *SekaiHandlerContext) PrefixArg() string {
	return s.prefixArg
}
func (s *SekaiHandlerContext) Flags() map[string]bool {
	return s.flags
}
func (s *SekaiHandlerContext) SetArgs(args string) {
	s.ArgText = args
}

func (skh *SekaiCommandHandler) Handle(ctx handler.Context) (interface{}, error) {
	if skh.handleFunc == nil {
		cmdName := "未定义"
		if len(skh.Commands) > 0 {
			cmdName = skh.Commands[0]
		}
		return nil, fmt.Errorf("Sekai 命令处理器 %s 没有处理方法", cmdName)
	}

	// 处理指令区服前缀
	var cmdRegion *sekairegion.SekaiRegion
	originalTriggerCmd := ctx.GetTriggerCmd()
	triggerCmd := originalTriggerCmd
	for _, region := range skh.Regions {
		cmdRegionPrefix := fmt.Sprintf("/%s", region.Id())
		if strings.HasPrefix(triggerCmd, cmdRegionPrefix) {
			cmdRegion = region
			triggerCmd = strings.Replace(triggerCmd, cmdRegionPrefix, "/", 1)
			break
		}
	}
	// 处理前缀参数
	prefixArg := ""
	for _, prefix := range skh.PrefixArgs {
		cmdPrefix := fmt.Sprintf("/%s", prefix)
		if strings.HasPrefix(triggerCmd, cmdPrefix) {
			prefixArg = prefix
			triggerCmd = strings.Replace(triggerCmd, cmdPrefix, "/", 1)
			break
		}
	}
	// TODO: 如果没有指定区服，并且用户有默认区服，并且用户默认区服在可用区服列表中，则使用用户的默认区服

	// 如果没有指定区服，并且用户没有默认区服，则使用指令的默认区服
	if cmdRegion == nil && len(skh.Regions) > 0 {
		cmdRegion = skh.Regions[0]
	}
	// TODO: 处理账号参数等
	args := ctx.GetArgs()
	uidArg := ""

	// 提取通用 flags 和 region flag (不需要 nicknames 也能提取这些)
	ext := parser.NewExtractor(nil)
	flags := make(map[string]bool)

	regRes := ext.ExtractRegion(args)
	if regRes.Value != "" {
		if r := sekairegion.GetRegionById(regRes.Value); r != nil {
			cmdRegion = r
		}
	}
	args = regRes.Remaining

	verbRes := ext.ExtractVerbose(args)
	flags["is_verbose"] = verbRes.Value
	args = verbRes.Remaining

	preRes := ext.ExtractPreview(args)
	flags["is_preview"] = preRes.Value
	args = preRes.Remaining

	helpRes := ext.ExtractHelp(args)
	flags["is_help"] = helpRes.Value
	args = helpRes.Remaining

	skCtx := SekaiHandlerContext{
		HandlerContext: handler.HandlerContext{
			Context:     ctx,
			TriggerCmd:  triggerCmd,
			ArgText:     args,
			MessageType: ctx.GetMessageType(),
			Message:     ctx.GetMessage(),
			Event:       ctx.GetEvent(),
			MessageId:   ctx.GetMessageId(),
			UserId:      ctx.GetUserId(),
			SenderName:  ctx.GetSenderName(),
			GroupId:     ctx.GetGroupId(),
		},
		region:             cmdRegion,
		originalTriggerCmd: originalTriggerCmd,
		prefixArg:          prefixArg,
		uidArg:             uidArg, //
		flags:              flags,
	}
	return skh.handleFunc(skCtx)
}

var DefaultRegions = sekairegion.Regions

type sekaiHandlers struct{}

func RegisterSekaiCommandHandler() {
	handlersVal := reflect.ValueOf(sekaiHandlers{})
	handlersTyp := handlersVal.Type()
	configTyp := reflect.TypeOf(SekaiCommandHandler{})
	for i := 0; i < handlersVal.NumMethod(); i++ {
		methodVal := handlersVal.Method(i)
		methodTyp := methodVal.Type()
		methodName := handlersTyp.Method(i).Name
		//
		if methodTyp.NumIn() == 0 &&
			methodTyp.NumOut() == 1 &&
			methodTyp.Out(0) == configTyp {
			log.Printf("注册指令解析器：%s\n", methodName)
			results := methodVal.Call(nil)
			skHandler := results[0].Interface().(SekaiCommandHandler)

			if len(skHandler.Regions) == 0 {
				skHandler.Regions = DefaultRegions
			}
			if len(skHandler.PrefixArgs) == 0 {
				skHandler.PrefixArgs = []string{""}
			}
			allRegionCommands := make(map[string]bool, len(skHandler.Commands)*len(skHandler.Regions)*len(skHandler.PrefixArgs))
			for _, prefix := range skHandler.PrefixArgs {
				for _, region := range skHandler.Regions {
					for _, cmd := range skHandler.Commands {
						if strings.HasPrefix(cmd, fmt.Sprintf("/%s%s", region.Id(), prefix)) {
							fmt.Fprintf(os.Stderr, "指令 %s 本身包含了区服前缀！", cmd)
						}
						allRegionCommands[cmd] = true
						allRegionCommands[strings.Replace(cmd, "/", fmt.Sprintf("/%s", prefix), 1)] = true
						allRegionCommands[strings.Replace(cmd, "/", fmt.Sprintf("/%s%s", region.Id(), prefix), 1)] = true
					}
				}
			}
			// 去除重复指令
			skHandler.Commands = make([]string, 0, len(allRegionCommands))
			for cmd := range allRegionCommands {
				skHandler.Commands = append(skHandler.Commands, cmd)
			}
			slices.Sort(skHandler.Commands)
			// 默认优先级
			if skHandler.Priority == 0 {
				skHandler.Priority = handler.DefaultPriority
			}
			handler.RegisterCommandHandler(&skHandler)
		}
	}
}
