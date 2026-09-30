package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-study/29-todo/middleware"
	"go-study/29-todo/model"
	"go-study/29-todo/utils"
)

type CreateTodoReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  *int   `json:"status"`
}

type UpdateTodoReq struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *int    `json:"status"`
}

func todoVO(t model.Todo) gin.H {
	return gin.H{
		"id":          t.ID,
		"uid":         t.UID,
		"title":       t.Title,
		"content":     t.Content,
		"status":      t.Status,
		"create_time": t.CreateTime,
		"update_time": t.UpdateTime,
	}
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.FailApp(c, utils.BadRequest("id 必须是正整数"))
		return 0, false
	}
	return uint(id), true
}

// Create POST /todo  uid 来自 JWT，不能由前端指定
func Create(c *gin.Context) {
	var req CreateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailApp(c, utils.BadRequest("title 必填"))
		return
	}
	status := model.TodoStatusTodo
	if req.Status != nil {
		status = *req.Status
	}
	t := model.Todo{
		UID:     middleware.UID(c),
		Title:   req.Title,
		Content: req.Content,
		Status:  status,
	}
	if err := model.CreateTodo(&t); err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.Created(c, "创建成功", gin.H{"todo": todoVO(t)})
}

// List GET /todo/list?page=1&size=10
func List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}

	list, total, err := model.GetTodoListByUID(middleware.UID(c), page, size)
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	vos := make([]gin.H, 0, len(list))
	for _, t := range list {
		vos = append(vos, todoVO(t))
	}
	utils.OK(c, "获取列表成功", gin.H{
		"page": page, "size": size, "total": total, "list": vos,
	})
}

func Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	t, err := model.GetTodoByIDAndUID(id, middleware.UID(c))
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.TodoNotFound())
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.OK(c, "获取详情成功", gin.H{"todo": todoVO(t)})
}

func Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateTodoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailApp(c, utils.BadRequest("JSON 格式错误"))
		return
	}
	fields := map[string]any{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Content != nil {
		fields["content"] = *req.Content
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if len(fields) == 0 {
		utils.FailApp(c, utils.BadRequest("没有要更新的字段"))
		return
	}

	if err := model.UpdateTodo(id, middleware.UID(c), fields); err != nil {
		if model.IsNotFound(err) {
			utils.FailApp(c, utils.TodoNotFound())
			return
		}
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	t, _ := model.GetTodoByIDAndUID(id, middleware.UID(c))
	utils.OK(c, "更新成功", gin.H{"todo": todoVO(t)})
}

func Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := model.DeleteTodo(id, middleware.UID(c)); err != nil {
		if model.IsNotFound(err) {
			utils.FailApp(c, utils.TodoNotFound())
			return
		}
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.OK(c, "删除成功", nil)
}
