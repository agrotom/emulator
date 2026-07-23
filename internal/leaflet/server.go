package leaflet

import (
	"net/http"

	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/track"
	"github.com/gin-gonic/gin"
)

var StartCoords = mathutil.Vector2f{X: 51.546529, Y: 46.037097}
var EndCoords = mathutil.Vector2f{X: 51.536257, Y: 46.02405}

type LeafletService interface {
	StartLeaflet()
}

type baseLeafletService struct {
	r *gin.Engine
}

func CreateLeafletService() LeafletService {
	return &baseLeafletService{}
}

func (lf *baseLeafletService) InitLeaflet() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/route", func(c *gin.Context) {
		coords, _ := track.GetOSRMRoute(StartCoords, EndCoords)

		interpCoords := track.Interpolate(coords, track.StepM)

		records := track.Simulate(interpCoords, 500)

		c.JSON(200, records)
	})

	lf.r = r
}

func (lf *baseLeafletService) StartLeaflet() {
	lf.r.Run(":8080")
}
