package controller

import (
	"carrot-backyard/model"
	"carrot-backyard/param"
	"carrot-backyard/util"
	"carrot-backyard/util/context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllBallotRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != model.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	ballots, err := m.GetAllBallot()
	if err != nil {
		util.ErrorPrint(err, nil, "get all ballots failed")
		return context.Error(c, http.StatusInternalServerError, "", err)
	}

	return context.Success(c, ballots)
}

func CreateOneBallotByTitleRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != model.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}
	musterTitle := c.QueryParam("muster")
	ms, err := m.GetOneMusterByTitle(musterTitle)
	if musterTitle == "" || err != nil {
		return context.Error(c, http.StatusBadRequest, "muster cannot be empty", nil)
	}
	remark := c.QueryParam("remark")

	if err := m.CreateOneBallotByTitle(title, ms, remark); err != nil {
		util.ErrorPrint(err, nil, "create new muster failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, nil)
}

func DeleteOneBallotByTitleRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != model.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	if err := m.DeleteOneBallotByTitle(title); err != nil {
		util.ErrorPrint(err, nil, "delete muster failed")
		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
	}

	return context.Success(c, nil)
}

//func AddAnOptionToOneBallotRequestHandler(c echo.Context) error {
//	m := model.GetModel()
//	defer m.Close()
//
//	token := c.QueryParam("token")
//	if token != model.Token {
//		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
//	}
//
//	title := c.QueryParam("title")
//	if title == "" {
//		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
//	}
//	option := c.QueryParam("option")
//	if option == "" {
//		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
//	}
//
//	bt, err := m.AddAnOptionToOneBallot(title, option)
//	if err != nil {
//		util.ErrorPrint(err, nil, "add option failed")
//		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
//	}
//
//	return context.Success(c, bt)
//}
//
//func DeleteAnOptionOnOneBallotRequestHandler(c echo.Context) error {
//	m := model.GetModel()
//	defer m.Close()
//
//	token := c.QueryParam("token")
//	if token != model.Token {
//		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
//	}
//
//	title := c.QueryParam("title")
//	if title == "" {
//		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
//	}
//	option := c.QueryParam("option")
//	if option == "" {
//		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
//	}
//
//	bt, err := m.DeleteAnOptionOnOneBallot(title, option)
//	if err != nil {
//		util.ErrorPrint(err, nil, "delete option failed")
//		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
//	}
//
//	return context.Success(c, bt)
//}

func UpdateAnswerForOneMemberRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != model.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}
	answer := c.QueryParam("answer")
	if answer == "" {
		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
	}
	name := c.QueryParam("name")

	bt, err := m.UpdateAnswerForOneMember(title, answer, name)
	if err != nil {
		util.ErrorPrint(err, nil, "update options failed")
		return context.Error(c, http.StatusInternalServerError, "update in db failed", err)
	}

	return context.Success(c, bt)
}

func BroadCastMessageForNoAnswerer(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != model.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	message := c.QueryParam("message")

	bt, err := m.GetOneBallotByTitle(title)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "no ballot find", err)
	}

	var people []param.PersonWithQQ
	for _, member := range bt.TargetMember {
		if member.AnsweredFlag == false {
			people = append(people, member.People)
		}
	}

	failed := util.SendSameMessageToManyFriends(
		fmt.Sprintf("————————\n卡洛收到了，希望你能填写【%s】的祈愿！\n\"祈愿附文：%s\"\n————————", title, message), people)
	var failedName []string
	for _, person := range failed {
		failedName = append(failedName, person.Name)
	}
	return context.Success(c, failedName)
}
