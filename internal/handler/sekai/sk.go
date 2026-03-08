package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	sekairegion "Haruki-Command-Parser/internal/sekai_region"
	"fmt"
	"strconv"
	"strings"
)

func (sekaiHandlers) SKLineHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/sk-line", "/sk线", "/榜线", "/pjsk sk line", "/pjsk board line", "/skl",
			},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-line"), nil
		},
	}
}

func (sekaiHandlers) SKQueryHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{Commands: []string{
			"/sk-query", "/sk查询", "/sk查分", "/pjsk sk board", "/pjsk board",
		},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-query"), nil
		},
	}
}

func (sekaiHandlers) SKSpeedHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk sk speed", "/pjsk board speed", "/时速", "/sks", "/skv", "/sk时速",
				"/sk-speed", "/sk时速", "/时速线", "/pjsk sk speed", "/pjsk board speed", "/sks", "/skv", "/sktime",
			},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-speed"), nil
		},
	}
}
func (sekaiHandlers) SKCheckRoomHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/sk-check-room", "/sk查房", "/查房", "/cf", "/pjsk查房", "/csb", "/冲水板", "/pjsk冲水板",
			},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-check-room"), nil
		},
	}
}

func (sekaiHandlers) SKPlayerTraceHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/sk-player-trace", "/sk玩家轨迹", "/玩家轨迹", "/ptr", "/pjsk玩家追踪", "/pjsk ptr",
			},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-player-trace"), nil
		},
	}
}

func (sekaiHandlers) SKRankTraceHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/sk-rank-trace", "/sk档线轨迹", "/档线轨迹", "/rtr", "/skt", "/sklt", "/sktl", "/pjsk追踪", "/pjsk sk追踪",
			},
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-rank-trace"), nil
		},
	}
}

func (sekaiHandlers) WinratePredictHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{Commands: []string{
			"/pjsk winrate predict", "/胜率预测", "/5v5预测", "/胜率", "/5v5胜率", "/预测胜率", "/预测5v5",
		}},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleSK, "sk-winrate"), nil
		},
	}
}

// TODO
func (sekaiHandlers) SKDailySpeedHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk sk daily speed", "/pjsk board daily speed", "/日速", "/skds", "/skdv", "/sk日速",
			},
			Disabled: true,
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs()) + ctx.prefixArg
			days := 1
			if v, err := strconv.Atoi(strings.TrimSpace(args)); err == nil {
				days = v
			}
			// TODO: 迁移 extract_wl_event + compose_sks_image(unit='d') 回图逻辑
			return nil, fmt.Errorf("TODO: 日速查询未实现，days=%d, query=%q", days, args)
		},
	}
}

func (sekaiHandlers) SKPredictHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk sk predict", "/pjsk board predict", "/sk预测", "/榜线预测", "/skp",
			},
			Disabled: true,
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs()) + ctx.prefixArg
			// TODO: 迁移 extract_wl_event + wl 单榜拦截 + compose_skp_image 回图逻辑
			return nil, fmt.Errorf("TODO: 榜线预测未实现，query=%q", args)
		},
	}
}

func (sekaiHandlers) SKBoardHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk sk board", "/pjsk board", "/sk",
			},
			Disabled: true,
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs()) + ctx.prefixArg
			// TODO: 迁移 extract_wl_event + parse_sk_query_params + compose_sk_image 回图逻辑
			return nil, fmt.Errorf("TODO: 指定榜线查询未实现，query=%q", args)
		},
	}
}

func (sekaiHandlers) CSBHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/csb", "/查水表", "/pjsk查水表", "/停车时间",
			},
			Disabled: true,
		},
		PrefixArgs: []string{"", "wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs()) + ctx.prefixArg
			// TODO: 迁移 extract_wl_event + parse_sk_query_params + compose_csb_image 回图逻辑
			return nil, fmt.Errorf("TODO: 查水表未实现，query=%q", args)
		},
	}
}
