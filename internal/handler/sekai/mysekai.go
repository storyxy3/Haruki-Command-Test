package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	"fmt"
	"strconv"
	"strings"
)

func (sekaiHandlers) MysekaiResourceHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/mysekai-resource", "/mysekai资源", "/烤森资源", "/msmap", "/msa",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-resource"), nil
		},
	}
}

func (sekaiHandlers) MysekaiTalkListHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/mysekai-talk-list", "/mysekai对话列表", "/烤森对话列表",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-talk-list"), nil
		},
	}
}
func (sekaiHandlers) MysekaiFixtureListHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/mysekai-fixture-list", "/mysekai家具列表", "/烤森家具列表", "/msf",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-fixture-list"), nil
		},
	}
}

func (sekaiHandlers) MysekaiFurnitureHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk mysekai furniture", "/pjsk mysekai fixture",
				"/msf", "/mysekai 家具", "/家具列表",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-fixture-detail"), nil
		},
	}
}

func (sekaiHandlers) MysekaiDoorUpgradeHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/mysekai-door-upgrade", "/mysekai大门升级", "/烤森大门升级", "/msg", "/msgate",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-door-upgrade"), nil
		},
	}
}
func (sekaiHandlers) MysekaiMusicRecordHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/mysekai-music-record", "/mysekai唱片", "/烤森唱片", "/msm", "/mss",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMysekai, "mysekai-music-record"), nil
		},
	}
}

// TODO
func (sekaiHandlers) MysekaiBlueprintHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk mysekai blueprint", "/mysekai blueprint",
				"/msb", "/mysekai 蓝图",
			},
			Disabled: true,
		},
		// TODO: refer 中限制 regions=get_regions(RegionAttributes.MYSEKAI)
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			showID := strings.Contains(args, "id")
			showAllTalks := strings.Contains(args, "all")
			// TODO: 迁移 unit/cid 解析 + compose_mysekai_fixture_list_image/compose_mysekai_talk_list_image 回图逻辑
			return nil, fmt.Errorf("TODO: mysekai蓝图查询未实现，args=%q, show_id=%t, show_all_talks=%t", args, showID, showAllTalks)
		},
	}
}
func (sekaiHandlers) MysekaiPhotoHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk mysekai photo", "/pjsk mysekai picture",
				"/msp", "/mysekai 照片",
			},
			Disabled: true,
		},
		// TODO: refer 中限制 regions=get_regions(RegionAttributes.MYSEKAI)
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			seq, err := strconv.Atoi(args)
			if err != nil {
				return nil, fmt.Errorf("请输入正确的照片编号（从1或-1开始）")
			}
			// TODO: 迁移群限制校验 + get_mysekai_photo_and_time + 回图逻辑
			return nil, fmt.Errorf("TODO: mysekai照片下载未实现，seq=%d", seq)
		},
	}
}

func (sekaiHandlers) CheckMysekaiDataHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk check mysekai data",
				"/pjsk烤森抓包数据", "/pjsk烤森抓包", "/烤森抓包", "/烤森抓包数据",
				"/msd",
			},
			Disabled: true,
		},
		// TODO: refer 中限制 regions=get_regions(RegionAttributes.MYSEKAI)
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			// TODO: 迁移 at用户解析 + get_player_bind_id + get_mysekai_upload_time + 文本组装逻辑
			return nil, fmt.Errorf("TODO: 烤森抓包状态未实现，user_id=%s", ctx.GetUserId())
		},
	}
}

func (sekaiHandlers) MSRChangeBindHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/msr换绑",
			},
			Disabled: true,
		},
		// TODO: refer 中限制 regions=get_regions(RegionAttributes.BD_MYSEKAI)
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			force := strings.Contains(args, "force")
			// TODO: force 需要 superuser 权限
			// TODO: 迁移参数校验 + update_bd_msr_limit_uid 调用
			return nil, fmt.Errorf("TODO: msr换绑未实现，force=%t, args=%q", force, args)
		},
	}
}
