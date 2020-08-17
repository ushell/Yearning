// Copyright 2019 HenryYee.
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"Yearning-go/src/handle"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"github.com/cookieY/yee"
	"github.com/cookieY/yee/middleware"
	"net/http"
)

func SuperManageDB() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		if lib.SuperAuth(c, "db") {
			c.Next()
			return nil
		}
		return c.JSON(http.StatusForbidden, "非法越权操作！")
	}
}

func SuperManageUser() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		if lib.SuperAuth(c, "user") {
			c.Next()
			return
		}
		return c.JSON(http.StatusForbidden, "非法越权操作！")

	}
}

func SuperManageGroup() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		user, _ := lib.JwtParse(c)
		if user == "admin" {
			c.Next()
			return
		}
		return c.JSON(http.StatusForbidden, "非法越权操作！")
	}
}

func AuditGroup() yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		_, rule := lib.JwtParse(c)
		if rule == "admin" || rule == "perform" {
			c.Next()
			return
		}
		return c.JSON(http.StatusForbidden, "非法越权操作！")
	}
}

func AddRouter(e *yee.Core) {
	e.GET("/", func(c yee.Context) error {
		return c.HTMLTml(http.StatusOK, "./dist/index.html")
	})
	e.POST("/login", handle.UserGeneralLogin)
	e.POST("/register", handle.UserRegister)
	e.GET("/fetch", handle.UserReqSwitch)
	e.POST("/ldap", handle.UserLdapLogin)

	r := e.Group("/api/v2", middleware.JWTWithConfig(middleware.JwtConfig{SigningKey: []byte(model.JWT)}))
	r.POST("/dash/initMenu", handle.DashInit)
	r.GET("/dash/pie", handle.DashPie)
	r.GET("/dash/axis", handle.DashAxis)
	r.GET("/dash/count", handle.DashCount)
	r.PUT("/dash/userinfo", handle.DashUserInfo)
	r.PUT("/dash/stmt", handle.DashStmt)

	r.POST("/user/password_reset", handle.ChangePassword)
	r.POST("/user/mail_reset", handle.ChangeMail)
	r.PUT("/user/order", handle.GeneralFetchMyOrder)
	r.GET("/fetch/sql", handle.GeneralFetchSQLInfo)
	r.GET("/fetch/idc", handle.GeneralIDC)
	r.GET("/fetch/source/:idc/:xxx", handle.GeneralSource)
	r.POST("/fetch/marge", handle.SuperUserRuleMarge)
	r.GET("/fetch/base/:source", handle.GeneralBase)
	r.PUT("/fetch/table", handle.GeneralTable)
	r.PUT("/fetch/tableinfo", handle.GeneralTableInfo)
	r.PUT("/fetch/test", handle.GeneralSQLTest)
	r.GET("/fetch/detail", handle.GeneralOrderDetailList)
	r.GET("/fetch/roll", handle.GeneralOrderDetailRollSQL)
	r.POST("/fetch/rollorder", handle.RollBackSQLOrder)
	r.GET("/fetch/undo", handle.GeneralFetchUndo)
	r.PUT("/query/status", handle.FetchQueryStatus)
	r.POST("/query/refer", handle.ReferQueryOrder)
	r.PUT("/query/fetchbase", handle.FetchQueryDatabaseInfo)
	r.GET("/query/fetchtable/:t/:source", handle.FetchQueryTableInfo)
	r.GET("/query/tableinfo/:base/:table/:source", handle.FetchQueryTableStruct)
	r.POST("/query", handle.FetchQueryResults)
	r.DELETE("/query/undo", handle.UndoQueryOrder)
	r.PUT("/query/merge", handle.GeneralMergeDDL)
	r.POST("/sql/refer", handle.SQLReferToOrder)
	r.GET("/board", handle.GeneralFetchBoard)

	audit := r.Group("/audit", AuditGroup())
	audit.POST("/refer/perform", handle.MulitAuditOrder)
	audit.PUT("", handle.FetchAuditOrder)
	audit.GET("/sql", handle.FetchOrderSQL)
	audit.GET("/kill/:work_id", handle.DelayKill)
	audit.POST("/reject", handle.RejectOrder)
	audit.POST("/execute", handle.ExecuteOrder)
	audit.PUT("/record", handle.FetchRecord)
	audit.PUT("/query/fetch", handle.FetchQueryOrder)
	audit.POST("/query/agreed", handle.AgreedQueryOrder)
	audit.POST("/query/disagreed", handle.DisAgreedQueryOrder)
	audit.POST("/query/undo", handle.SuperUndoQueryOrder)
	audit.PUT("/query/cancel", handle.QueryQuickCancel)
	audit.DELETE("/query/empty", handle.QueryDeleteEmptyRecord)
	audit.PUT("/query/fetch/record", handle.FetchQueryRecord)
	audit.PUT("/query/fetch/record/detail", handle.FetchQueryRecordDetail)
	audit.GET("/fetch_osc/:work_id", handle.OscPercent)
	audit.DELETE("/fetch_osc/:work_id", handle.OscKill)

	group := r.Group("/group", SuperManageGroup())
	group.GET("", handle.SuperGroup)
	group.POST("/update", handle.SuperGroupUpdate)
	group.DELETE("/del/:clear", handle.SuperClearUserRule)
	group.GET("/setting", handle.SuperFetchSetting)
	group.POST("/setting/add", handle.SuperSaveSetting)
	group.POST("/setting/roles", handle.SuperSaveRoles)
	group.PUT("/setting/test/:el", handle.SuperTestSetting)
	group.POST("/setting/del/order", handle.UndoAuditOrder)
	group.POST("/setting/del/query", handle.DelQueryOrder)
	group.POST("/board/post", handle.GeneralPostBoard)

	r.GET("/manage_user", handle.SuperFetchUser)
	r.POST("/manage_user", handle.SuperUserRegister)
	r.DELETE("/manage_user", handle.SuperDeleteUser)
	user := r.Group("/management_user", SuperManageUser())
	user.POST("/modify", handle.SuperModifyUser)
	user.POST("/password_reset", handle.SuperChangePassword)
	user.DELETE("/del/:user", handle.SuperDeleteUser)
	user.POST("/register", handle.SuperUserRegister)

	db := r.Group("/management_db", SuperManageDB())
	db.GET("", handle.SuperFetchDB)
	db.POST("", handle.SuperAddDB)
	db.PUT("", handle.SuperModifyDb)
	db.PUT("/test", handle.SuperTestDBConnect)
	db.DELETE("/del/:source", handle.SuperDeleteDb)
	db.DELETE("", handle.SuperDeleteDb)

	autoTask := r.Group("/auto", SuperManageGroup())
	autoTask.GET("", handle.SuperFetchAutoTaskSource)
	autoTask.POST("", handle.SuperReferAutoTask)
	autoTask.PUT("/fetch", handle.SuperFetchAutoTaskList)
	autoTask.POST("/edit", handle.SuperEditAutoTask)
	autoTask.DELETE("/:id", handle.SuperDeleteAutoTask)
	autoTask.POST("/active", handle.SuperAutoTaskActivation)
}
