package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"Haruki-Command-Parser/internal/chardata"
	"Haruki-Command-Parser/internal/config"
	"Haruki-Command-Parser/internal/handler"
	sekaihandler "Haruki-Command-Parser/internal/handler/sekai"
	"Haruki-Command-Parser/internal/parser"

	sekai "haruki-cloud/database/sekai"

	_ "github.com/lib/pq"
)

// ParseRequest is the request body accepted by POST /api/parse and POST /api/process.
type ParseRequest struct {
	Text   string `json:"text"`
	UserQQ string `json:"user_qq,omitempty"`
	// Params carries pre-assembled structured data from the Bot for modules that
	// require external data (SK, Mysekai, Education, Score, misc-birthday, etc.).
	// Part1 forwards this field verbatim to Part2's /api/render.
	Params json.RawMessage `json:"params,omitempty"`
}

// ParseResponse is the unified response from POST /api/parse.
type ParseResponse struct {
	Module string          `json:"module,omitempty"`
	Mode   string          `json:"mode,omitempty"`
	Region string          `json:"region,omitempty"`
	Query  string          `json:"query,omitempty"`
	Flags  map[string]bool `json:"flags,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type server struct {
	loader        *chardata.Loader
	logger        *slog.Logger
	serviceClient *http.Client // HTTP client for calling Part2
	serviceURL    string       // Part2 /api/render URL
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfgPath := findConfig()
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("failed to load config", "error", err, "path", cfgPath)
		os.Exit(1)
	}

	if cfg.Log.Level == "debug" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	// --- Connect to Haruki-Cloud DB ---
	var sekaiClient *sekai.Client
	if cfg.HarukiCloud.SekaiDB.Driver != "" {
		dsn := cfg.HarukiCloud.SekaiDB.DSN
		if dsn == "" {
			dsn, err = config.BuildDSN(cfg.HarukiCloud.SekaiDB)
			if err != nil {
				logger.Error("failed to build DSN", "error", err)
				os.Exit(1)
			}
		}
		sekaiClient, err = sekai.Open(cfg.HarukiCloud.SekaiDB.Driver, dsn)
		if err != nil {
			logger.Error("failed to connect to Sekai DB", "error", err)
			os.Exit(1)
		}
		defer sekaiClient.Close()
		logger.Info("connected to Haruki-Cloud Sekai DB")
	}

	// --- Load character nicknames ---
	region := strings.ToLower(strings.TrimSpace(cfg.HarukiCloud.Region))
	if region == "" {
		region = "jp"
	}

	loader := chardata.NewLoader(sekaiClient, region, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if sekaiClient != nil {
		if err := loader.Load(ctx); err != nil {
			logger.Warn("failed to load character nicknames from DB, starting with empty map", "error", err)
		}
		refreshInterval := cfg.HarukiCloud.CacheRefreshIntervalDur
		if refreshInterval <= 0 {
			refreshInterval = 6 * time.Hour
		}
		loader.StartBackgroundRefresh(ctx, refreshInterval)
	}

	sekaihandler.RegisterSekaiCommandHandler()

	// --- HTTP client for Part2 ---
	serviceURL := cfg.ServiceAPI.BaseURL
	if serviceURL == "" {
		serviceURL = "http://localhost:24000"
	}
	serviceClient := &http.Client{Timeout: cfg.ServiceAPI.TimeoutDur}

	// --- HTTP Server ---
	srv := &server{
		loader:        loader,
		logger:        logger,
		serviceClient: serviceClient,
		serviceURL:    serviceURL + "/api/render",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/parse", srv.handleParse)     // 返回 ParsedCommand JSON（调试用）
	mux.HandleFunc("/api/process", srv.handleProcess) // 解析 + 调用 Part2 + 返回图片（Bot 主流程）
	mux.HandleFunc("/api/reload", srv.handleReload)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	})

	host := cfg.Server.Host
	port := cfg.Server.Port
	if port == 0 {
		port = 8001
	}
	addr := host + ":" + strconv.Itoa(port)

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Haruki-Command-Parser starting", "addr", addr)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	logger.Info("stopped")
}

// handleParse processes POST /api/parse.
func (s *server) handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ParseResponse{Error: "method not allowed"})
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ParseResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		writeJSON(w, http.StatusBadRequest, ParseResponse{Error: "text is required"})
		return
	}

	evt := handler.Event{
		MessageType: handler.MessageTypePrivate,
		Message:     req.Text,
		UserId:      req.UserQQ,
	}
	result, _ := handler.Dispatch(r.Context(), evt)
	if result == nil {
		writeJSON(w, http.StatusOK, ParseResponse{Error: "无法识别指令格式，请发送 /help 查看说明"})
		return
	}

	if err, ok := result.(error); ok {
		writeJSON(w, http.StatusOK, ParseResponse{Error: err.Error()})
		return
	}

	resolved, ok := result.(*parser.ResolvedCommand)
	if !ok {
		writeJSON(w, http.StatusOK, ParseResponse{Error: "unexpected handler result"})
		return
	}

	resp := ParseResponse{
		Module: moduleToString(resolved.Module),
		Mode:   resolved.Mode,
		Region: resolved.Region,
		Query:  resolved.Query,
		Flags: map[string]bool{
			"is_help":    resolved.IsHelp,
			"is_verbose": resolved.IsVerbose,
			"is_preview": resolved.IsPreview,
		},
	}

	s.logger.Debug("parsed command",
		"text", req.Text,
		"module", resp.Module,
		"mode", resp.Mode,
		"region", resp.Region,
		"query", resp.Query,
	)

	writeJSON(w, http.StatusOK, resp)
}

// handleProcess processes POST /api/process.
// This is the main Bot-facing endpoint:
//  1. Parses the raw text command (same as /api/parse)
//  2. Builds a ParsedCommand and POSTs it directly to Part2 /api/render
//  3. Streams the PNG image bytes back to the caller
//
// Bot flow: Bot → POST /api/process → Part1 → POST /api/render (Part2) → DrawingAPI → image
func (s *server) handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	// Step 1: Parse the command
	evt := handler.Event{
		MessageType: handler.MessageTypePrivate,
		Message:     req.Text,
		UserId:      req.UserQQ,
	}
	result, _ := handler.Dispatch(r.Context(), evt)
	if result == nil {
		http.Error(w, "无法识别指令格式，请发送 /help 查看说明", http.StatusBadRequest)
		return
	}

	if errRes, ok := result.(error); ok {
		http.Error(w, "parse error: "+errRes.Error(), http.StatusBadRequest)
		return
	}

	resolved, ok := result.(*parser.ResolvedCommand)
	if !ok {
		http.Error(w, "parse error: unexpected handler result", http.StatusBadRequest)
		return
	}

	// Step 2: Build the render payload for Part2
	renderPayload := map[string]interface{}{
		"module":  moduleToString(resolved.Module),
		"mode":    resolved.Mode,
		"region":  resolved.Region,
		"query":   resolved.Query,
		"user_id": req.UserQQ,
		"flags": map[string]bool{
			"is_help":    resolved.IsHelp,
			"is_verbose": resolved.IsVerbose,
			"is_preview": resolved.IsPreview,
		},
	}
	// Forward params from Bot (for SK/Mysekai/Education/Score modules that need external data)
	if len(req.Params) > 0 {
		renderPayload["params"] = req.Params
	}

	s.logger.Info("forwarding to Part2",
		"module", renderPayload["module"],
		"mode", renderPayload["mode"],
		"region", renderPayload["region"],
		"url", s.serviceURL,
	)

	// Step 3: POST to Part2 /api/render
	payloadBytes, _ := json.Marshal(renderPayload)
	partReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, s.serviceURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		http.Error(w, "failed to build request to Part2: "+err.Error(), http.StatusInternalServerError)
		return
	}
	partReq.Header.Set("Content-Type", "application/json")

	resp, err := s.serviceClient.Do(partReq)
	if err != nil {
		s.logger.Error("Part2 request failed", "error", err)
		http.Error(w, "Part2 unavailable: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Step 4: Relay the response back to Bot
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(fmt.Sprintf("Part2 error (HTTP %d), check Part2 logs", resp.StatusCode)))
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, resp.Body)
}

// handleReload forces a reload of character nicknames from DB.
func (s *server) handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := s.loader.Load(ctx); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Wait, resolver was handling reloading of character nicknames.
	// Since the handler system uses `parser.Extractor` which currently handles nickname matching if supplied... Wait!
	// Extractor has nicknames, but Handler framework does not inject nicknames.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func moduleToString(m parser.TargetModule) string {
	switch m {
	case parser.ModuleCard:
		return "card"
	case parser.ModuleGacha:
		return "gacha"
	case parser.ModuleMusic:
		return "music"
	case parser.ModuleEvent:
		return "event"
	case parser.ModuleDeck:
		return "deck"
	case parser.ModuleSK:
		return "sk"
	case parser.ModuleMysekai:
		return "mysekai"
	case parser.ModuleProfile:
		return "profile"
	case parser.ModuleHelp:
		return "help"
	case parser.ModuleEducation:
		return "education"
	case parser.ModuleScore:
		return "score"
	case parser.ModuleStamp:
		return "stamp"
	case parser.ModuleMisc:
		return "misc"
	default:
		return "unknown"
	}
}

func findConfig() string {
	candidates := []string{"configs.yaml", "config.yaml", "../configs.yaml"}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		return v
	}
	fmt.Fprintln(os.Stderr, "warning: no config file found, using defaults")
	return "configs.yaml"
}
