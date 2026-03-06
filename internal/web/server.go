// Package web 提供 Web 管理界面服务
// 使用 Gin 框架构建 RESTful API 和管理页面
package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/platform"
	"github.com/yahao333/get_jobs/internal/storage"
)

// Server Web 服务器
type Server struct {
	router *gin.Engine
	port   int
}

// NewServer 创建 Web 服务器
func NewServer(port int) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.New("index").Parse(htmlTemplate)))

	// 静态文件
	r.Static("/static", "./static")

	server := &Server{
		router: r,
		port:   port,
	}

	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 首页
	s.router.GET("/", s.handleIndex)

	// API 路由
	api := s.router.Group("/api")
	{
		// 岗位相关
		api.GET("/jobs", s.handleGetJobs)
		api.POST("/jobs/search", s.handleSearchJobs)
		api.POST("/jobs/deliver", s.handleDeliverJob)
		api.DELETE("/jobs/:id", s.handleDeleteJob)

		// 黑名单相关
		api.GET("/blacklist", s.handleGetBlacklist)
		api.POST("/blacklist", s.handleAddBlacklist)
		api.DELETE("/blacklist/:id", s.handleDeleteBlacklist)

		// 投递记录
		api.GET("/deliveries", s.handleGetDeliveries)

		// 平台相关
		api.GET("/platforms", s.handleGetPlatforms)
		api.GET("/status", s.handleGetStatus)

		// 配置相关
		api.GET("/config", s.handleGetConfig)
		api.POST("/config", s.handleUpdateConfig)
	}
}

// handleIndex 首页
func (s *Server) handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "求职自动化工具",
		"port":  s.port,
	})
}

// handleGetJobs 获取岗位列表
func (s *Server) handleGetJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var jobs []storage.BossData
	var err error

	if status != "" {
		err = storage.Where(&jobs, "delivery_status = ?", status)
	} else {
		err = storage.Where(&jobs, "1 = 1")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(jobs) {
		jobs = []storage.BossData{}
	} else {
		if end > len(jobs) {
			end = len(jobs)
		}
		jobs = jobs[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":      jobs,
		"page":      page,
		"page_size": pageSize,
		"total":     len(jobs),
	})
}

// handleSearchJobs 搜索岗位
func (s *Server) handleSearchJobs(c *gin.Context) {
	var req struct {
		Platform string `json:"platform"`
		Keyword  string `json:"keyword"`
		CityCode string `json:"city_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 调用平台搜索功能
	c.JSON(http.StatusOK, gin.H{
		"message": "搜索功能待实现",
		"params":  req,
	})
}

// handleDeliverJob 投递简历
func (s *Server) handleDeliverJob(c *gin.Context) {
	var req struct {
		JobID   int64  `json:"job_id"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 调用投递功能
	c.JSON(http.StatusOK, gin.H{
		"message": "投递功能待实现",
		"job_id":  req.JobID,
	})
}

// handleDeleteJob 删除岗位
func (s *Server) handleDeleteJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := storage.Delete(&storage.BossData{}, "id = ?", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// handleGetBlacklist 获取黑名单
func (s *Server) handleGetBlacklist(c *gin.Context) {
	blacklistType := c.Query("type")

	var blacklists []storage.Blacklist
	var err error

	if blacklistType != "" {
		err = storage.Where(&blacklists, "type = ?", blacklistType)
	} else {
		err = storage.Where(&blacklists, "1 = 1")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blacklist": blacklists})
}

// handleAddBlacklist 添加黑名单
func (s *Server) handleAddBlacklist(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blacklist := storage.Blacklist{
		Keyword: req.Keyword,
		Type:    req.Type,
		Source:  "manual",
	}

	if err := storage.Create(&blacklist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// handleDeleteBlacklist 删除黑名单
func (s *Server) handleDeleteBlacklist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := storage.Delete(&storage.Blacklist{}, "id = ?", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// handleGetDeliveries 获取投递记录
func (s *Server) handleGetDeliveries(c *gin.Context) {
	var records []storage.DeliveryRecord
	if err := storage.Where(&records, "1 = 1"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deliveries": records})
}

// handleGetPlatforms 获取支持的平台
func (s *Server) handleGetPlatforms(c *gin.Context) {
	platforms := platform.AllPlatforms
	c.JSON(http.StatusOK, gin.H{"platforms": platforms})
}

// handleGetStatus 获取系统状态
func (s *Server) handleGetStatus(c *gin.Context) {
	// 统计岗位数量
	totalJobs, _ := storage.Count(&storage.BossData{}, "1 = 1")
	deliveredJobs, _ := storage.Count(&storage.BossData{}, "delivery_status = ?", "delivered")
	pendingJobs, _ := storage.Count(&storage.BossData{}, "delivery_status = ?", "pending")

	// 统计黑名单
	totalBlacklist, _ := storage.Count(&storage.Blacklist{}, "1 = 1")

	// 统计投递记录
	totalDeliveries, _ := storage.Count(&storage.DeliveryRecord{}, "1 = 1")

	c.JSON(http.StatusOK, gin.H{
		"total_jobs":       totalJobs,
		"delivered_jobs":   deliveredJobs,
		"pending_jobs":     pendingJobs,
		"total_blacklist":  totalBlacklist,
		"total_deliveries": totalDeliveries,
	})
}

// handleGetConfig 获取配置
func (s *Server) handleGetConfig(c *gin.Context) {
	configs := map[string]interface{}{
		"search.keywords":          config.GetStringSlice("search.keywords"),
		"search.city_codes":        config.Get("search.city_codes"),
		"delivery.daily_limit":     config.GetInt("delivery.daily_limit"),
		"delivery.send_img_resume": config.GetBool("delivery.send_img_resume"),
		"ai.enable":                config.GetBool("ai.enable"),
		"filter.filter_dead_hr":    config.GetBool("filter.filter_dead_hr"),
	}

	c.JSON(http.StatusOK, gin.H{"config": configs})
}

// handleUpdateConfig 更新配置
func (s *Server) handleUpdateConfig(c *gin.Context) {
	// TODO: 实现配置更新
	c.JSON(http.StatusOK, gin.H{"message": "配置更新功能待实现"})
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	config.Info("Web 管理界面启动: http://localhost", addr)
	return s.router.Run(addr)
}

// 简单的 HTML 模板
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { background: #fff; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header h1 { color: #333; font-size: 24px; }
        .stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 15px; margin-bottom: 20px; }
        .stat-card { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stat-card h3 { color: #666; font-size: 14px; margin-bottom: 10px; }
        .stat-card .value { color: #333; font-size: 28px; font-weight: bold; }
        .card { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .card h2 { color: #333; font-size: 18px; margin-bottom: 15px; padding-bottom: 10px; border-bottom: 1px solid #eee; }
        .btn { padding: 8px 16px; background: #007AFF; color: #fff; border: none; border-radius: 4px; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #FF3B30; }
        .btn-danger:hover { background: #d32f2f; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f9f9f9; color: #666; font-weight: 500; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .status-pending { background: #FFF3E0; color: #F57C00; }
        .status-delivered { background: #E8F5E9; color: #388E3C; }
        .status-failed { background: #FFEBEE; color: #D32F2F; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.title}}</h1>
            <p>AI 驱动的求职自动化工具 | 管理端口: {{.port}}</p>
        </div>

        <div class="stats">
            <div class="stat-card">
                <h3>总岗位数</h3>
                <div class="value" id="totalJobs">-</div>
            </div>
            <div class="stat-card">
                <h3>已投递</h3>
                <div class="value" id="deliveredJobs">-</div>
            </div>
            <div class="stat-card">
                <h3>待投递</h3>
                <div class="value" id="pendingJobs">-</div>
            </div>
            <div class="stat-card">
                <h3>黑名单</h3>
                <div class="value" id="blacklistCount">-</div>
            </div>
        </div>

        <div class="card">
            <h2>岗位管理</h2>
            <button class="btn" onclick="searchJobs()">搜索新岗位</button>
            <table>
                <thead>
                    <tr>
                        <th>公司</th>
                        <th>职位</th>
                        <th>薪资</th>
                        <th>地点</th>
                        <th>状态</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody id="jobsTable">
                    <tr><td colspan="6">加载中...</td></tr>
                </tbody>
            </table>
        </div>

        <div class="card">
            <h2>黑名单管理</h2>
            <button class="btn" onclick="addBlacklist()">添加黑名单</button>
            <table>
                <thead>
                    <tr>
                        <th>关键词</th>
                        <th>类型</th>
                        <th>来源</th>
                        <th>添加时间</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody id="blacklistTable">
                    <tr><td colspan="5">加载中...</td></tr>
                </tbody>
            </table>
        </div>
    </div>

    <script>
        // 加载数据
        async function loadData() {
            try {
                // 加载状态
                const statusRes = await fetch('/api/status');
                const status = await statusRes.json();
                document.getElementById('totalJobs').textContent = status.total_jobs || 0;
                document.getElementById('deliveredJobs').textContent = status.delivered_jobs || 0;
                document.getElementById('pendingJobs').textContent = status.pending_jobs || 0;
                document.getElementById('blacklistCount').textContent = status.total_blacklist || 0;

                // 加载岗位
                const jobsRes = await fetch('/api/jobs?page_size=10');
                const jobsData = await jobsRes.json();
                renderJobs(jobsData.jobs || []);

                // 加载黑名单
                const blRes = await fetch('/api/blacklist');
                const blData = await blRes.json();
                renderBlacklist(blData.blacklist || []);
            } catch(e) {
                console.error(e);
            }
        }

        function renderJobs(jobs) {
            const tbody = document.getElementById('jobsTable');
            if (jobs.length === 0) {
                tbody.innerHTML = '<tr><td colspan="6">暂无数据</td></tr>';
                return;
            }
            tbody.innerHTML = jobs.map(job => {
                const companyName = job.company_name || '';
                const jobName = job.job_name || '';
                const salary = job.salary || '';
                const location = job.location || '';
                const status = job.delivery_status || '';
                const statusClass = 'status status-' + status;
                return '<tr>'
                    + '<td>' + companyName + '</td>'
                    + '<td>' + jobName + '</td>'
                    + '<td>' + salary + '</td>'
                    + '<td>' + location + '</td>'
                    + '<td><span class="' + statusClass + '">' + status + '</span></td>'
                    + '<td>'
                    + '<button class="btn" onclick="deliverJob(' + job.id + ')">投递</button>'
                    + '<button class="btn btn-danger" onclick="deleteJob(' + job.id + ')">删除</button>'
                    + '</td>'
                    + '</tr>';
            }).join('');
        }

        function renderBlacklist(blacklist) {
            const tbody = document.getElementById('blacklistTable');
            if (blacklist.length === 0) {
                tbody.innerHTML = '<tr><td colspan="5">暂无数据</td></tr>';
                return;
            }
            tbody.innerHTML = blacklist.map(bl => {
                const keyword = bl.keyword || '';
                const type = bl.type || '';
                const source = bl.source || '';
                const createdAt = bl.created_at || '';
                return '<tr>'
                    + '<td>' + keyword + '</td>'
                    + '<td>' + type + '</td>'
                    + '<td>' + source + '</td>'
                    + '<td>' + createdAt + '</td>'
                    + '<td><button class="btn btn-danger" onclick="deleteBlacklist(' + bl.id + ')">删除</button></td>'
                    + '</tr>';
            }).join('');
        }

        async function searchJobs() {
            alert('搜索功能待实现');
        }

        async function deliverJob(id) {
            alert('投递功能待实现');
        }

        async function deleteJob(id) {
            if (!confirm('确定删除?')) return;
            await fetch('/api/jobs/' + id, { method: 'DELETE' });
            loadData();
        }

        async function addBlacklist() {
            const keyword = prompt('请输入黑名单关键词:');
            const type = prompt('请输入类型 (company/hr/job):');
            if (!keyword || !type) return;
            await fetch('/api/blacklist', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({keyword, type})
            });
            loadData();
        }

        async function deleteBlacklist(id) {
            if (!confirm('确定删除?')) return;
            await fetch('/api/blacklist/' + id, { method: 'DELETE' });
            loadData();
        }

        loadData();
    </script>
</body>
</html>
`
