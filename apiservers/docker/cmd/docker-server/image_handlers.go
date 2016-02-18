package main

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/docker/restapi/operations"
	"github.com/vmware/vic/apiservers/docker/restapi/operations/image"
)

type ImageHandlersImpl struct{}

func (handler *ImageHandlersImpl) Configure(api *operations.DockerAPI) {
	api.ImageDeleteImagesNameHandler = image.DeleteImagesNameHandlerFunc(handler.DeleteImages)
	api.ImageGetImagesJSONHandler = image.GetImagesJSONHandlerFunc(handler.GetImages)
	api.ImageGetImagesNameJSONHandler = image.GetImagesNameJSONHandlerFunc(handler.GetImagesName)
	api.ImagePostImagesCreateHandler = image.PostImagesCreateHandlerFunc(handler.PostImagesCreate)
	api.ImageBuildHandler = image.BuildHandlerFunc(handler.Build)
	api.ImageCommitHandler = image.CommitHandlerFunc(handler.Commit)
	api.ImageHistoryHandler = image.HistoryHandlerFunc(handler.History)
	api.ImageLoadHandler = image.LoadHandlerFunc(handler.Load)
	api.ImagePushHandler = image.PushHandlerFunc(handler.Push)
	api.ImageSaveHandler = image.SaveHandlerFunc(handler.Save)
	api.ImageSaveAllHandler = image.SaveAllHandlerFunc(handler.SaveAll)
	api.ImageSearchHandler = image.SearchHandlerFunc(handler.Search)
	api.ImageTagHandler = image.TagHandlerFunc(handler.Tag)
}

func (handler *ImageHandlersImpl) DeleteImages(params image.DeleteImagesNameParams) middleware.Responder {
	return middleware.NotImplemented("operation image.DeleteImagesName has not yet been implemented")
}

func (handler *ImageHandlersImpl) GetImages(params image.GetImagesJSONParams) middleware.Responder {
	return middleware.NotImplemented("operation image.GetImagesJSON has not yet been implemented")
}

func (handler *ImageHandlersImpl) GetImagesName(params image.GetImagesNameJSONParams) middleware.Responder {
	return middleware.NotImplemented("operation image.GetImagesNameJSON has not yet been implemented")
}

func (handler *ImageHandlersImpl) PostImagesCreate(params image.PostImagesCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation image.PostImagesCreate has not yet been implemented")
}

func (handler *ImageHandlersImpl) Build(params image.BuildParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Build has not yet been implemented")
}

func (handler *ImageHandlersImpl) Commit(params image.CommitParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Commit has not yet been implemented")
}

func (handler *ImageHandlersImpl) History(params image.HistoryParams) middleware.Responder {
	return middleware.NotImplemented("operation image.History has not yet been implemented")
}

func (handler *ImageHandlersImpl) Load(params image.LoadParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Load has not yet been implemented")
}

func (handler *ImageHandlersImpl) Push(params image.PushParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Push has not yet been implemented")
}

func (handler *ImageHandlersImpl) Save(params image.SaveParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Save has not yet been implemented")
}

func (handler *ImageHandlersImpl) SaveAll(params image.SaveAllParams) middleware.Responder {
	return middleware.NotImplemented("operation image.SaveAll has not yet been implemented")
}

func (handler *ImageHandlersImpl) Search(params image.SearchParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Search has not yet been implemented")
}

func (handler *ImageHandlersImpl) Tag(params image.TagParams) middleware.Responder {
	return middleware.NotImplemented("operation image.Tag has not yet been implemented")
}
