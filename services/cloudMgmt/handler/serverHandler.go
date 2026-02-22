package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golf/cloudmgmt/appUtility/config"
	"github.com/golf/cloudmgmt/services/cloudMgmt/behavior"
	"github.com/golf/cloudmgmt/services/cloudMgmt/model"
	externalrepo "github.com/golf/cloudmgmt/services/cloudMgmt/repo/externalRepo"
	gatewayrepo "github.com/golf/cloudmgmt/services/cloudMgmt/repo/gatewayRepo"
)

type ServerHandler struct {
	serverBehavior *behavior.Server
}

func NewServerHandler(serverBehavior *behavior.Server) *ServerHandler {
	return &ServerHandler{
		serverBehavior: serverBehavior,
	}
}

func (handler *ServerHandler) GetAllServer(appConfig *config.AppConfig) fiber.Handler {

	return func(c *fiber.Ctx) error {

		// check db
		list, _ := handler.serverBehavior.GetAll()
		for _, server := range list {
			fmt.Printf("Infra Id in DB: %+v\n", server.InfraId)
			fmt.Printf("Server Id in DB: %+v\n", server.ServerId)
		}

		infraRegistRes, errGetRegist := getAllRegistedInfra(appConfig)
		if errGetRegist != nil {
			return errGetRegist
		}

		var serverListRes externalrepo.ResponseGetAllServer

		for _, serverFromInfra := range infraRegistRes.Resources {
			serverData, errGetInfra := handler.serverBehavior.GetByInfraId(serverFromInfra.ID)

			// the serverData will not appear bc the infra service is not return the actual created infra at the api

			if errGetInfra != nil {
				return errGetInfra
			}

			serverListRes.ServerList = append(serverListRes.ServerList, externalrepo.SkuByServerIdResponse{
				ServerID:  serverData.ServerId,
				SKU:       serverData.Sku,
				IsPowerOn: serverData.IsPowerOn,
			})
		}

		if len(serverListRes.ServerList) == 0 {
			serverListRes.IsFound = false
			serverListRes.ServerList = []externalrepo.SkuByServerIdResponse{}
		} else {
			serverListRes.IsFound = true
		}

		return c.Status(fiber.StatusOK).JSON(serverListRes)
	}

}

func (handler *ServerHandler) AddServer(appConfig *config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		urlInsertInfra := appConfig.InfraAPIBaseDomain + appConfig.InfraInsertPath

		var reqBody gatewayrepo.RequestAddInfra

		if errParseBody := c.BodyParser(&reqBody); errParseBody != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, errParseBody.Error())
		}

		infraListRes, errGetSku := getAllAvalibleSkuList(appConfig)
		if errGetSku != nil {
			return errGetSku
		}

		isValidSku := false

		for _, serverFromInfra := range infraListRes.SkuData {
			if reqBody.Sku == serverFromInfra.Sku {
				isValidSku = true
				break
			}
		}

		if !isValidSku {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid sku")
		}

		reqBodyBytes, errMarshal := json.Marshal(reqBody)
		if errMarshal != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errMarshal.Error())
		}

		resp, errConnect := http.Post(urlInsertInfra, "application/json", bytes.NewBuffer(reqBodyBytes))
		if errConnect != nil || resp.StatusCode != http.StatusOK {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to connect to infrastructure: %v", errConnect))
		}
		defer resp.Body.Close()

		bodyByte, errReadBody := io.ReadAll(resp.Body)
		if errReadBody != nil {
			return fiber.NewError(fiber.StatusConflict, errReadBody.Error())
		}

		var infraRes gatewayrepo.ResponseAddInfra

		if errUnmarshal := json.Unmarshal(bodyByte, &infraRes); errUnmarshal != nil {
			return fiber.NewError(fiber.StatusConflict, errUnmarshal.Error())
		}

		handler.serverBehavior.Create(&model.ServerModel{
			InfraId:   infraRes.Id,
			Sku:       reqBody.Sku,
			IsPowerOn: true,
		})

		return c.Status(fiber.StatusOK).JSON(externalrepo.ResponseAddServer{
			Success: true,
		})

	}
}

func (handler *ServerHandler) PowerControlServer(appConfig *config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {

		serverId := c.Params("serverId")
		if serverId == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Server ID is required")
		}

		serverData, errGetInfraId := handler.serverBehavior.Get(serverId)
		if errGetInfraId != nil {
			return fiber.NewError(fiber.StatusNotFound, "Server not found")
		}

		urlPowerControl := appConfig.InfraAPIBaseDomain + appConfig.InfraInsertPath + "/" + serverData.InfraId + "/power"

		var reqBody gatewayrepo.RequestControlInfraPower

		if errParseBody := c.BodyParser(&reqBody); errParseBody != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, errParseBody.Error())
		}

		reqBodyBytes, errMarshal := json.Marshal(reqBody)
		if errMarshal != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errMarshal.Error())
		}

		resp, errConnect := http.Post(urlPowerControl, "application/json", bytes.NewBuffer(reqBodyBytes))
		if errConnect != nil || resp.StatusCode != http.StatusOK {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to connect to infrastructure: %v", errConnect))
		}
		defer resp.Body.Close()

		bodyByte, errReadBody := io.ReadAll(resp.Body)
		if errReadBody != nil {
			return fiber.NewError(fiber.StatusConflict, errReadBody.Error())
		}

		var infraPowerRes gatewayrepo.ResponseControlInfraPower

		if errUnmarshal := json.Unmarshal(bodyByte, &infraPowerRes); errUnmarshal != nil {
			return fiber.NewError(fiber.StatusConflict, errUnmarshal.Error())
		}

		var isPowerOn bool

		switch reqBody.Action {
		case "on":
			isPowerOn = true
		case "off":
			isPowerOn = false
		default:
			return fiber.NewError(fiber.StatusBadRequest, "Invalid power action")
		}

		if errUpdate := handler.serverBehavior.Update(serverId, &model.ServerModel{
			ServerId:  serverData.ServerId,
			InfraId:   serverData.InfraId,
			Sku:       serverData.Sku,
			IsPowerOn: isPowerOn,
		}); errUpdate != nil {
			return errUpdate
		}

		return c.Status(fiber.StatusOK).JSON(externalrepo.ResponseControlPower{
			Success: true,
			State:   infraPowerRes.State,
		})

	}
}

func getAllAvalibleSkuList(appConfig *config.AppConfig) (gatewayrepo.ResponseGWGetAllInfra, error) {
	urlGetAllServer := appConfig.InfraAPIBaseDomain + appConfig.InfraGetAllPath

	resp, errConnect := http.Get(urlGetAllServer)
	if errConnect != nil || resp.StatusCode != http.StatusOK {
		return gatewayrepo.ResponseGWGetAllInfra{}, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to connect to infrastructure: %v", errConnect))
	}
	defer resp.Body.Close()

	bodyByte, errReadBody := io.ReadAll(resp.Body)
	if errReadBody != nil {
		return gatewayrepo.ResponseGWGetAllInfra{}, fiber.NewError(fiber.StatusConflict, errReadBody.Error())
	}

	var infraListRes gatewayrepo.ResponseGWGetAllInfra

	if errUnmarshal := json.Unmarshal(bodyByte, &infraListRes); errUnmarshal != nil {
		return gatewayrepo.ResponseGWGetAllInfra{}, fiber.NewError(fiber.StatusConflict, errUnmarshal.Error())
	}

	return infraListRes, nil
}

func getAllRegistedInfra(appConfig *config.AppConfig) (gatewayrepo.ResponseGetRegistedInfra, error) {
	urlGetRegistServer := appConfig.InfraAPIBaseDomain + appConfig.InfraInsertPath

	resp, errConnect := http.Get(urlGetRegistServer)
	if errConnect != nil || resp.StatusCode != http.StatusOK {
		return gatewayrepo.ResponseGetRegistedInfra{}, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to connect to infrastructure: %v", errConnect))
	}
	defer resp.Body.Close()

	bodyByte, errReadBody := io.ReadAll(resp.Body)
	if errReadBody != nil {
		return gatewayrepo.ResponseGetRegistedInfra{}, fiber.NewError(fiber.StatusConflict, errReadBody.Error())
	}

	var infraRegistRes gatewayrepo.ResponseGetRegistedInfra

	if errUnmarshal := json.Unmarshal(bodyByte, &infraRegistRes); errUnmarshal != nil {
		return gatewayrepo.ResponseGetRegistedInfra{}, fiber.NewError(fiber.StatusConflict, errUnmarshal.Error())
	}

	return infraRegistRes, nil
}
