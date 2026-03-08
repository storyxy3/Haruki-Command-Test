package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	sekairegion "Haruki-Command-Parser/internal/sekai_region"
	"fmt"
	"strings"
)

var multiEventCmds = []string{"/pjsk events", "/pjsk_events", "/events", "/活动列表", "/活动一览"}
var singleEventCmds = []string{"/pjsk event", "/pjsk_event", "/event", "/活动", "/查活动"}

func (sekaiHandlers) EventHandle() SekaiCommandHandler {
	// commands := make([]string, 0, len(singleEventCmds)+len(multiEventCmds))
	// commands = append(commands, singleEventCmds...)
	// commands = append(commands, multiEventCmds...)
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk events", "/pjsk_events", "/events", "/活动列表", "/活动一览",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleEvent, "event-list"), nil
		},
	}
}
func (sekaiHandlers) EventDetailHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/活动", "/查活动", "/event",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleEvent, "event-detail"), nil
		},
	}
}

// TODO
func (sekaiHandlers) EventStoryHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk event story", "/pjsk_event_story",
				"/活动剧情", "/活动故事", "/活动总结",
			},
			Disabled: true,
		},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			refresh := false
			save := true
			if strings.Contains(args, "refresh") {
				refresh = true
				args = strings.TrimSpace(strings.ReplaceAll(args, "refresh", ""))
			}
			model := ""
			if strings.Contains(args, "model:") {
				parts := strings.SplitN(args, "model:", 2)
				args = strings.TrimSpace(parts[0])
				model = strings.TrimSpace(parts[1])
				refresh = true
				save = false
			}

			// TODO: 默认模型 get_model_preset("sekai.story_summary.event")
			// TODO: model: 权限校验（仅超级用户）
			// TODO: 迁移 parse_search_single_event_args/get_current_event + block_region
			// TODO: 迁移 get_event_story_summary 并按返回类型回复
			return nil, fmt.Errorf(
				"TODO: 活动剧情未实现，query=%q, refresh=%t, save=%t, model=%q",
				args, refresh, save, model,
			)
		},
	}
}

func (sekaiHandlers) SendBoostHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk send boost", "/pjsk_send_boost", "/pjsk grant boost", "/pjsk_grant_boost",
				"/自动送火", "/送火",
			},
			Disabled: true,
		},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			// TODO: 迁移 send_boost(ctx, ctx.user_id)
			return nil, fmt.Errorf("TODO: 自动送火未实现，user_id=%s", ctx.GetUserId())
		},
	}
}

func (sekaiHandlers) EventRecordHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk event record", "/pjsk_event_record",
				"/活动记录", "/冲榜记录",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			// TODO: 迁移 compose_event_record_image(ctx, ctx.user_id) 回图逻辑
			return nil, fmt.Errorf("TODO: 活动记录未实现，user_id=%s", ctx.GetUserId())
		},
	}
}
