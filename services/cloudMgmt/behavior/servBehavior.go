package behavior

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golf/cloudmgmt/services/cloudMgmt/model"
	"github.com/google/uuid"
)

type Server struct {
	serverData map[string]*model.ServerModel
}

func NewServerBehavior() *Server {
	return &Server{
		serverData: make(map[string]*model.ServerModel),
	}
}

func (server *Server) Create(serverModel *model.ServerModel) error {

	serverId := uuid.New().String()

	if _, ok := server.serverData[serverId]; ok {
		return fiber.NewError(fiber.StatusBadRequest, "Server with the same InfraId already exists")
	}

	serverModel.ServerId = serverId

	server.serverData[serverId] = serverModel

	return nil

}

func (server *Server) Get(serverId string) (*model.ServerModel, error) {
	if serverModel, ok := server.serverData[serverId]; ok {
		return serverModel, nil
	}
	return nil, fiber.NewError(fiber.StatusNotFound, "Server not found")
}

func (server *Server) GetByInfraId(infraId string) (*model.ServerModel, error) {
	for _, serverModel := range server.serverData {
		if serverModel.InfraId == infraId {
			return serverModel, nil
		}
	}
	return nil, fiber.NewError(fiber.StatusNotFound, "Server not found")
}

func (server *Server) GetAll() (map[string]*model.ServerModel, error) {
	if len(server.serverData) > 0 {
		return server.serverData, nil
	}

	return nil, fiber.NewError(fiber.StatusNotFound, "No any servers found")
}

func (server *Server) Update(serverId string, updatedModel *model.ServerModel) error {
	if _, ok := server.serverData[serverId]; !ok {
		return fiber.NewError(fiber.StatusNotFound, "Server not found")
	}

	server.serverData[serverId] = updatedModel
	return nil
}
