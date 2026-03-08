package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	sekairegion "Haruki-Command-Parser/internal/sekai_region"
	"errors"
	"fmt"
	"strings"
)

func (sekaiHandlers) CardDetailHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/card-detail", "/卡面", "/详情", "/查卡",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleCard, "card-detail"), nil
		},
	}
}

func (sekaiHandlers) CardListHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/查牌", "/查卡牌", "/卡牌列表", "/card", "/cards", "/pjsk card", "/pjsk member",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleCard, "card-list"), nil
		},
	}
}
func (sekaiHandlers) CardBoxHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/查箱", "/查框", "/卡牌一览", "/卡面一览", "/卡一览", "/box", "/card-box", "/pjsk box",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleCard, "card-box"), nil
		},
	}
}

// TODO
func (sekaiHandlers) CharaAliasHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk chara alias",
				"/角色别名", "/查角色别名",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			if args == "" {
				return nil, errors.New("请输入要查询的角色名或别名")
			}

			// TODO: 迁移 get_cid_by_nickname / get_nicknames_by_chara_id 逻辑
			return nil, fmt.Errorf("TODO: 角色别名查询未实现，query=%q", args)

		},
	}
}
func (sekaiHandlers) CardImgHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk card img",
				"/查卡面", "/卡面",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			if args == "" {
				return nil, errors.New("请输入要查询的卡牌")
			}

			// TODO: 迁移 search_single_card 逻辑
			// TODO: 迁移 only_has_after_training / has_after_training / get_card_image 逻辑
			return nil, fmt.Errorf("TODO: 卡面查询未实现，query=%q", args)
		},
	}
}

func (sekaiHandlers) CardStoryHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk card story",
				"/卡牌剧情", "/卡面剧情", "/卡剧情", "/卡牌故事", "/卡面故事", "/卡故事",
			},
			Disabled: true,
		},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			if args == "" {
				return nil, errors.New("请输入要查询的卡牌")
			}

			refresh := false
			save := true
			if strings.Contains(args, "refresh") {
				args = strings.TrimSpace(strings.ReplaceAll(args, "refresh", ""))
				refresh = true
			}

			model := ""
			if strings.Contains(args, "model:") {
				parts := strings.SplitN(args, "model:", 2)
				args = strings.TrimSpace(parts[0])
				model = strings.TrimSpace(parts[1])
				refresh = true
				save = false
			}

			// TODO: model: 仅超级用户可指定（check_superuser）
			// TODO: 默认模型读取 get_model_preset("sekai.story_summary.card")
			// TODO: 迁移 search_single_card + block_region + get_card_story_summary 逻辑
			// TODO: 按返回类型决定回复图片还是折叠文本
			return nil, fmt.Errorf(
				"TODO: 卡牌剧情查询未实现，query=%q, refresh=%t, save=%t, model=%q",
				args, refresh, save, model,
			)
		},
	}
}

func (sekaiHandlers) BoxHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk box",
				"/卡牌一览", "/卡面一览", "/卡一览",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())

			showID := false
			if strings.Contains(args, "id") {
				showID = true
				args = strings.TrimSpace(strings.ReplaceAll(args, "id", ""))
			}

			showBox := false
			if strings.Contains(args, "box") {
				showBox = true
				args = strings.TrimSpace(strings.ReplaceAll(args, "box", ""))
			}

			useAfterTraining := true
			if strings.Contains(args, "before") {
				useAfterTraining = false
				args = strings.TrimSpace(strings.ReplaceAll(args, "before", ""))
			}

			// TODO: 迁移 search_multi_cards(..., contain_leak=false) 并校验剩余参数
			// TODO: 迁移 compose_box_image(ctx, ctx.user_id, cards, showID, showBox, useAfterTraining)
			return nil, fmt.Errorf(
				"TODO: 卡牌一览查询未实现，query=%q, showID=%t, showBox=%t, useAfterTraining=%t",
				args, showID, showBox, useAfterTraining,
			)
		},
	}
}
