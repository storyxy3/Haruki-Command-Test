package sekai

import (
	"Haruki-Command-Parser/internal/handler"
	"Haruki-Command-Parser/internal/parser"
	"fmt"
	"strconv"
	"strings"
)

func (sekaiHandlers) MusicDetailHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/查曲", "/查歌", "/查乐", "/查音乐", "/查询乐曲", "/查歌曲", "/歌曲", "/乐曲", "/song", "/music",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMusic, "music-detail"), nil
		},
	}
}
func (sekaiHandlers) MusicListHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/歌曲列表", "/歌曲一览", "/乐曲列表", "/乐曲一览", "/难度排行", "/定数表", "/歌曲定数", "/查乐曲", "/music-list",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMusic, "music-list"), nil
		},
	}
}
func (sekaiHandlers) MusicRewardsHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/曲目奖励", "/歌曲奖励", "/music rewards", "/music-rewards", "/pjsk music rewards",
				"/打歌奖励", "/歌曲挖矿", "/打歌挖矿",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMusic, "music-rewards"), nil
		},
	}
}

func (sekaiHandlers) MusicProgressHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/打歌进度", "/歌曲进度", "/打歌信息", "/pjsk进度", "/progress", "/music-progress",
			},
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			return makeResolvedCmd(ctx, parser.ModuleMusic, "music-progress"), nil
		},
	}
}

// TODO
func (sekaiHandlers) AliasDelHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk alias del", "/pjskalias del",
				"/删除歌曲别名", "/歌曲别名删除",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 MusicAliasDB.remove + 日志与批量反馈逻辑
			return nil, fmt.Errorf("TODO: 删除歌曲别名未实现，args=%q", args)
		},
	}
}
func (sekaiHandlers) SongHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk song", "/pjsk music", "/song", "/music",
				"/查曲", "/查歌", "/歌曲", "/查歌曲",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			query := strings.TrimSpace(ctx.GetArgs())
			if query == "" {
				return nil, fmt.Errorf("请输入要查询的歌曲名或ID")
			}
			// TODO: 迁移 leak 查询、多曲查询、单曲详情查询逻辑
			return nil, fmt.Errorf("TODO: 查曲未实现，query=%q", query)
		},
	}
}

func (sekaiHandlers) NoteNumHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk note num", "/pjsk note count",
				"/物量", "/查物量",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			noteCount, err := strconv.Atoi(args)
			if err != nil {
				return nil, fmt.Errorf("请输入物量数值")
			}
			// TODO: 迁移按 totalNoteCount 查询谱面并组装结果文本
			return nil, fmt.Errorf("TODO: 物量查询未实现，note_count=%d", noteCount)
		},
	}
}

func (sekaiHandlers) PlayProgressHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk progress",
				"/pjsk进度", "/打歌进度", "/歌曲进度", "/打歌信息",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 extract_diff + compose_play_progress_image 回图逻辑
			return nil, fmt.Errorf("TODO: 打歌进度未实现，args=%q", args)
		},
	}
}

func (sekaiHandlers) SyncMusicAliasHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/sync music alias", "/sma", "/同步歌曲别名",
			},
			Disabled: true,
		},
		// TODO: refer 中这里是 CmdHandler（非 SekaiCmdHandler），后续确认是否迁到通用 handler
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			// TODO: 迁移 superuser 校验 + block + sync_music_alias 流程
			return nil, fmt.Errorf("TODO: 同步歌曲别名未实现")
		},
	}
}

func (sekaiHandlers) BPMHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk bpm", "/查bpm", "/查BPM",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			query := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 search_music + get_chart_bpm + 封面与 BPM 文本拼接逻辑
			return nil, fmt.Errorf("TODO: BPM查询未实现，query=%q", query)
		},
	}
}

func (sekaiHandlers) MusicCoverHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk music cover",
				"/查曲绘", "/曲绘",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			query := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 search_music + 读取曲绘资源并回复逻辑
			return nil, fmt.Errorf("TODO: 曲绘查询未实现，query=%q", query)
		},
	}
}

func (sekaiHandlers) AliasSetHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk alias add", "/pjskalias add",
				"/添加歌曲别名", "/歌曲别名添加",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 search_music + MusicAliasDB.add + 日志与反馈逻辑
			return nil, fmt.Errorf("TODO: 添加歌曲别名未实现，args=%q", args)
		},
	}
}
func (sekaiHandlers) AliasHandle() SekaiCommandHandler {
	return SekaiCommandHandler{
		CommandHandlerBase: handler.CommandHandlerBase{
			Commands: []string{
				"/pjsk alias", "/music alias",
				"/歌曲别名", "/查歌曲别名",
			},
			Disabled: true,
		},
		handleFunc: func(ctx SekaiHandlerContext) (interface{}, error) {
			args := strings.TrimSpace(ctx.GetArgs())
			// TODO: 迁移 search_music + MusicAliasDB.get_aliases 逻辑
			return nil, fmt.Errorf("TODO: 查看歌曲别名未实现，args=%q", args)
		},
	}
}
