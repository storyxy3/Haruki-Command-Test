package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	sekairegion "Haruki-Command-Parser/internal/sekai_region"
	"fmt"
	"strconv"
	"strings"
)

func (sekaiHandlers) ScoreControlHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/分数", "/查分数", "/pjsk score", "/score control",
				"/控分",
			},
		},
		Regions:    []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		PrefixArgs: []string{"wl"},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			parts := strings.SplitN(args, " ", 2)
			if len(parts) == 0 {
				return nil, fmt.Errorf("使用方式:\n%s 活动pt 歌曲名(可选)", ctx.originalTriggerCmd)
			}
			targetPT, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil || targetPT <= 0 {
				return nil, fmt.Errorf("使用方式:\n%s 活动pt 歌曲名(可选)", ctx.originalTriggerCmd)
			}
			// query := ""
			// if len(parts) > 1 {
			// 	query = strings.TrimSpace(parts[1])
			// }
			return makeResolvedCmd(ctx, parser.ModuleScore, "score-control"), nil
		},
	}
}

func (sekaiHandlers) CustomRoomScoreControlHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk custom room score", "/custom room score",
				"/自定义房间控分", "/自定义房控分", "/自定义控分",
				"/自定义房间分数", "/自定义分数",
			},
		},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			targetPT, err := strconv.Atoi(args)
			if err != nil || targetPT <= 0 {
				return nil, fmt.Errorf("使用方式: %s 目标PT", ctx.originalTriggerCmd)
			}
			return makeResolvedCmd(ctx, parser.ModuleScore, "score-custom-room"), nil
		},
	}
}

func (sekaiHandlers) MusicMetaHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk music meta", "/music meta",
				"/歌曲meta", "/曲目meta",
			},
			Priority: 1,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			segments := strings.Split(strings.ReplaceAll(args, "/", "|"), "|")
			clean := make([]string, 0, len(segments))
			for _, seg := range segments {
				seg = strings.TrimSpace(seg)
				if seg != "" {
					clean = append(clean, seg)
				}
			}
			if len(clean) == 0 {
				return nil, fmt.Errorf("请至少提供一个歌曲ID或名称")
			}
			if len(clean) > 3 {
				return nil, fmt.Errorf("一次最多进行3首歌曲的比较")
			}
			return makeResolvedCmd(ctx, parser.ModuleScore, "score-music-meta"), nil
		},
	}
}

func (sekaiHandlers) MusicBoardHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk music board", "/music board",
				"/歌曲排行", "/歌曲比较", "/歌曲排名", "/曲目榜",
			},
			Priority: 1,
		},
		Regions: []*sekairegion.SekaiRegion{sekairegion.GetRegionById("jp")},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleScore, "score-music-board"), nil
		},
	}
}
