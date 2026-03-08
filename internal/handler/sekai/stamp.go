package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	"fmt"
	"strconv"
	"strings"
)

func (sekaiHandlers) StampHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/贴纸", "/查贴纸", "/pjsk贴纸", "/pjsk stamp", "/pjsk bq", "/stamp",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleStamp, "stamp-list"), nil
		},
	}
}

// TODO
func (sekaiHandlers) StampMakeHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk stamp", "/pjsk bq",
				"/pjsk表情", "/pjsk表情制作",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			if args == "" {
				return nil, fmt.Errorf("使用方式\n查询某个角色: %s miku\n根据id查询: %s 123\n查询多个: %s 123 456\n制作表情: %s 123 文本",
					ctx.originalTriggerCmd, ctx.originalTriggerCmd, ctx.originalTriggerCmd, ctx.originalTriggerCmd)
			}
			format := "gif"
			if strings.Contains(args, "png") {
				format = "png"
				args = strings.TrimSpace(strings.ReplaceAll(args, "png", ""))
			}

			// TODO: 迁移 block_region + 参数解析(qtype=id/cid/id_text) + 获取/制作表情逻辑
			return nil, fmt.Errorf("TODO: 表情查询/制作未实现，query=%q, format=%s", args, format)
		},
	}
}

func (sekaiHandlers) RandStampHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk rand stamp", "/pjsk rand bq",
				"/pjsk随机表情", "/pjsk随机表情制作", "/随机表情",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			format := "gif"
			if strings.Contains(args, "png") {
				format = "png"
				args = strings.TrimSpace(strings.ReplaceAll(args, "png", ""))
			}
			// TODO: 迁移 block_region + 随机 sid 选择 + 指定角色/制作文本逻辑
			return nil, fmt.Errorf("TODO: 随机表情未实现，query=%q, format=%s", args, format)
		},
	}
}

func (sekaiHandlers) StampRefreshHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk stamp refresh", "/pjsk refresh stamp",
				"/pjsk表情刷新", "/pjsk刷新表情", "/pjsk刷新表情底图", "/pjsk表情刷新底图",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			sid, err := strconv.Atoi(args)
			if err != nil || sid < 0 {
				return nil, fmt.Errorf("使用方式: %s 123", ctx.originalTriggerCmd)
			}
			// TODO: 迁移 block_region + 删除旧cutout + ensure_stamp_maker_base_image + gif 回传逻辑
			return nil, fmt.Errorf("TODO: 刷新表情底图未实现，sid=%d", sid)
		},
	}
}

func (sekaiHandlers) StampRefreshBatchHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk stamp refresh batch",
				"/pjsk表情刷新批量",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			// TODO: 迁移 block_region + superuser 校验 + 批量刷新逻辑
			return nil, fmt.Errorf("TODO: 批量刷新表情底图未实现")
		},
	}
}

func (sekaiHandlers) StampBaseHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk stamp base",
				"/pjsk表情底图",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			gif := true
			if strings.Contains(args, "png") {
				gif = false
				args = strings.TrimSpace(strings.ReplaceAll(args, "png", ""))
			}
			sid, err := strconv.Atoi(args)
			if err != nil || sid < 0 {
				return nil, fmt.Errorf("使用方式: %s 123", ctx.originalTriggerCmd)
			}
			// TODO: 迁移 ensure_stamp_maker_base_image + gif/png 回传逻辑
			return nil, fmt.Errorf("TODO: 查看表情底图未实现，sid=%d, gif=%t", sid, gif)
		},
	}
}

func (sekaiHandlers) StampBaseDeleteHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk remove stamp base", "/pjsk del stamp base",
				"/pjsk删除表情底图",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			parts := strings.Fields(args)
			if len(parts) == 0 {
				return nil, fmt.Errorf("使用方式: %s 123 456", ctx.originalTriggerCmd)
			}
			sids := make([]int, 0, len(parts))
			for _, p := range parts {
				sid, err := strconv.Atoi(p)
				if err != nil {
					return nil, fmt.Errorf("使用方式: %s 123 456", ctx.originalTriggerCmd)
				}
				sids = append(sids, sid)
			}
			// TODO: 迁移 superuser 校验 + sid 存在性校验 + 删除 cutout 底图逻辑
			return nil, fmt.Errorf("TODO: 删除表情底图未实现，sids=%v", sids)
		},
	}
}
