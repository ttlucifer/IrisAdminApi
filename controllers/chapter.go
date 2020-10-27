package controllers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/transformer"
	"github.com/snowlyg/blog/validates"
	gf "github.com/snowlyg/gotransformer"
)

/**
* @api {get} /admin/chapters/:id 根据id获取分类信息
* @apiName 根据id获取分类信息
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 根据id获取分类信息
* @apiSampleRequest /admin/chapters/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetChapter(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	chapter, err := models.GetChapterById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, chapterTransform(chapter), "操作成功"))
}

/**
* @api {post} /admin/chapters/ 新建分类
* @apiName 新建分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 新建分类
* @apiSampleRequest /admin/chapters/
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func CreateChapter(ctx iris.Context) {
	chapter := new(models.Chapter)
	if err := ctx.ReadJSON(chapter); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*chapter)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(200, nil, e))
				return
			}
		}
	}

	err = chapter.CreateChapter()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(200, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if chapter.ID == 0 {
		_, _ = ctx.JSON(ApiResource(200, nil, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, chapterTransform(chapter), "操作成功"))
	}

}

/**
* @api {post} /admin/chapters/:id/update 更新分类
* @apiName 更新分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 更新分类
* @apiSampleRequest /admin/chapters/:id/update
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func UpdateChapter(ctx iris.Context) {
	aul := new(models.Chapter)

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(200, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	aul.ID = id
	err = models.UpdateChapterById(id, aul)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(200, nil, fmt.Sprintf("Error update chapter: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if aul.ID == 0 {
		_, _ = ctx.JSON(ApiResource(200, nil, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, chapterTransform(aul), "操作成功"))
	}

}

/**
* @api {delete} /admin/chapters/:id/delete 删除分类
* @apiName 删除分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 删除分类
* @apiSampleRequest /admin/chapters/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func DeleteChapter(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteChapterById(id)
	if err != nil {

		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /tts 获取所有的分类
* @apiName 获取所有的分类
* @apiGroup Chapters
* @apiVersion 1.0.0
* @apiDescription 获取所有的分类
* @apiSampleRequest /tts
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiChapter null
 */
func GetAllChapters(ctx iris.Context) {
	offset := libs.ParseInt(ctx.URLParam("offset"), 1)
	limit := libs.ParseInt(ctx.URLParam("limit"), 20)
	searchStr := ctx.FormValue("searchStr")
	docId := uint(libs.ParseInt(ctx.URLParam("docId"), 0))
	orderBy := ctx.FormValue("orderBy")
	fmt.Println(fmt.Sprintf("docId:%d", docId))

	chapters, err := models.GetAllChapters(docId, searchStr, orderBy, offset, limit)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(200, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, chaptersTransform(chapters), "操作成功"))
}

func chaptersTransform(chapters []*models.Chapter) []*transformer.Chapter {
	var rs []*transformer.Chapter
	for _, chapter := range chapters {
		r := chapterTransform(chapter)
		rs = append(rs, r)
	}
	return rs
}

func chapterTransform(chapter *models.Chapter) *transformer.Chapter {
	r := &transformer.Chapter{}
	g := gf.NewTransform(r, chapter, time.RFC3339)
	_ = g.Transformer()
	if chapter.Doc != nil {
		transform := docTransform(chapter.Doc)
		r.Doc = *transform
	}
	return r
}