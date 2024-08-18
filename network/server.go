package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"websocket_chatting/types"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	d := &api{server: server}

	server.engine.GET("/room-list", d.roomList)
	server.engine.POST("/make-room", d.makeRoom)
	server.engine.GET("/room", d.room)
	server.engine.GET("/enter-room", d.enterRoom)

}

func (a *api) roomList(c *gin.Context) {
	res, err := a.server.service.RoomList()
	if err != nil {
		response(c, http.StatusInternalServerError, err.Error())
		return
	}
	response(c, http.StatusOK, res)
}

func (a *api) makeRoom(c *gin.Context) {
	var req types.BodyRoomReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	err = a.server.service.MakeRoom(req.Name)
	if err != nil {
		response(c, http.StatusInternalServerError, err.Error())
		return
	}
	response(c, http.StatusOK, "Success")

}
func (a *api) room(c *gin.Context) {
	var req types.FormRoomReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := a.server.service.Room(req.Name)
	err = noResult(err)

	if err != nil {
		response(c, http.StatusInternalServerError, err.Error())
		return
	}
	response(c, http.StatusOK, res)
}
func (a *api) enterRoom(c *gin.Context) {

	var req types.FormRoomReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := a.server.service.EnterRoom(req.Name)
	if err != nil {
		response(c, http.StatusInternalServerError, err.Error())
		return
	}
	response(c, http.StatusOK, res)
}

func response(c *gin.Context, s int, res interface{}, data ...string) {
	c.JSON(s, types.NewRes(s, res, data...))
}
func noResult(err error) error {
	if strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil
	}
	return err
}
