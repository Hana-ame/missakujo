package backend

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/Hana-ame/missakujo/utils"

	"github.com/gofiber/fiber/v2"
)

type DelReqCtx struct {
	Host  string  `json:"host"`
	User  string  `json:"user"`
	Token string  `json:"token"`
	Since float64 `json:"since"`
	Until float64 `json:"until"`

	RenoteLessThan int `json:"renoteLessThan"`

	TimeOffset int `json:"timeOffset"`

	DeleteReply  bool `json:"deleteReply"`
	DeleteRenote bool `json:"deleteRenote"`
}

const timeForm = "2006-01-02 15:04:05"

var lm = utils.NewLockedMap()

func App() *fiber.App {
	app := fiber.New()

	app.Post("/delete", func(c *fiber.Ctx) error {
		req := new(DelReqCtx)
		// req := new(map[string]any)

		err := c.BodyParser(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// fmt.Println(req)
		// fmt.Printf("%T", (*req)["since"])
		// return nil

		sid := Wrapper(req)
		lm.Put(sid, time.Now().Unix()+999_999)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"sessionID": sid,
		})
	})

	app.Get("/log/:sessionID", func(c *fiber.Ctx) error {
		sid := c.Params("sessionID")
		// fmt.Println(sid)
		lm.Put(sid, time.Now().Unix()+30)

		file, err := os.OpenFile(sid+".txt", os.O_RDONLY, 0644)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		buf, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		n, err := c.Status(fiber.StatusOK).Write(buf)
		_ = n
		return err
	})

	app.Get("/webfinger/:acct", func(c *fiber.Ctx) error {
		acct := c.Params("acct")
		userId, err := utils.ResolveUser(acct)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"userId": userId,
		})
	})

	go func() {
		for {
			time.Sleep(60 * time.Second)
			lm.Iter(func(key string, v any) {
				vv, err := utils.GetWithType[int64](v)
				if err != nil {
					log.Println(err)
					return
				}

				if vv < time.Now().Unix() {
					err := os.Remove(key + ".txt")
					log.Println(err)
					lm.Remove(key)
				}
			})
		}
	}()

	return app
}
