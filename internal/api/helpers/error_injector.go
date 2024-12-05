package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rlapenok/effective_mobile_test/internal/domain/song_repository"
	"github.com/rlapenok/effective_mobile_test/internal/domain/swagger_client"
)

func SendError(c *gin.Context, err error) {
	switch e := err.(type) {
	case swagger_client.SwaggerClientError:
		{
			c.JSON(e.Code, gin.H{"desc": e.Error()})

		}
	case song_repository.SongRepositoryError:
		{
			switch e := e.Err.(type) {

			case *pq.Error:
				{
					switch e.Code {

					case "23505", "23503":
						{
							c.JSON(http.StatusBadRequest, ErrResponse{Desc: e.Error()})
						}
					default:
						{
							c.JSON(http.StatusInternalServerError, ErrResponse{Desc: e.Error()})

						}
					}

				}
			default:
				{
					c.JSON(http.StatusInternalServerError, ErrResponse{Desc: e.Error()})

				}
			}
		}
	default:
		if e.Error() == "not found" {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, ErrResponse{Desc: e.Error()})

		}
	}
}
