package sekai

import (
	"Haruki-Command-Parser/cmd/handler"
	"Haruki-Command-Parser/internal/parser"
	"errors"
	"strings"
)

const MUSIC_SEARCH_HELP = `请输入要查询的曲目，支持以下查询方式:
1. 直接使用曲目名称或别名
2. 曲目ID: id123
3. 曲目负数索引: 例如 -1 表示最新的曲目，-1leak 则会包含未公开的曲目
4. 活动id: event123
5. 箱活: ick1`

// TODO: 占位，并不是真的要在这里画
func generate_music_chart(ctx SekaiHandlerContext, musicInfo *parser.MusicQueryInfo, refresh bool) (interface{}, error) {
	return nil, nil
}

func (sekaiHandlers) ChartHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		BaseCommandHandler: handler.BaseCommandHandler{
			Commands: []string{
				"/pjsk chart",
				"/谱面查询", "/铺面查询", "/谱面预览", "/铺面预览", "/谱面", "/铺面", "/查谱面", "/查铺面", "/查谱",
				"/技能预览",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			query := strings.TrimSpace(ctx.GetArgs())
			if query == "" {
				return nil, errors.New(MUSIC_SEARCH_HELP)
			}
			refresh := false
			if strings.Contains(query, "refresh") {
				refresh = true
				query = strings.Replace(query, "refresh", "", 1)
			}
			// TODO: 临时的，必须改掉
			musicParser := parser.NewMusicParser(map[string]int{"虾": 76})
			musicInfo, err := musicParser.Parse(query)
			if err != nil {
				return nil, err
			}

			return generate_music_chart(ctx, musicInfo, refresh)
		},
	}
}
