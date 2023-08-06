package main

import (
	"fmt"
    "math/rand"
    "strings"
	"time"

    "github.com/deosjr/PokemonGo/src/model"
    "github.com/deosjr/PokemonGo/src/singleplayer"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var turn = 1

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pokemon",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		camPos = pixel.ZV
	)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

    pic, err := loadPicture("img/battlebgIndoorA.png")
    if err != nil {
        panic(err)
    }
	//batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
    sprite := pixel.NewSprite(pic, pic.Bounds())

    model.MustLoadConfig()
    p1 := singleplayer.GetRentalPokemon()
    p2 := singleplayer.GetRentalPokemon()
    // I dont know why I dont just expose currentHP?
    p1HP := p1.GetTotalHP()
    p2HP := p2.GetTotalHP()

    pic2, err := loadPicture(fmt.Sprintf("img/%s_back.png", strings.ToLower(p1.Name)))
    if err != nil {
        panic(err)
    }
    sprite2 := pixel.NewSprite(pic2, pic2.Bounds())

    pic3, err := loadPicture(fmt.Sprintf("img/%s.png", strings.ToLower(p2.Name)))
    if err != nil {
        panic(err)
    }
    sprite3 := pixel.NewSprite(pic3, pic3.Bounds())

    // pokemon game logic
    battle := model.NewSingleBattle(p1, p2)

    movebox1 := pixel.R(-200, -250, -30, -200)
    movebox2 := pixel.R(30, -250, 200, -200)
    movebox3 := pixel.R(-200, -350, -30, -300)
    movebox4 := pixel.R(30, -350, 200, -300)
    boxes := []pixel.Rect{movebox1, movebox2, movebox3, movebox4}

    atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
    move1 := text.New(movebox1.Center(), atlas)
    move2 := text.New(movebox2.Center(), atlas)
    move3 := text.New(movebox3.Center(), atlas)
    move4 := text.New(movebox4.Center(), atlas)
    moves := []*text.Text{move1, move2, move3, move4}

    imd := imdraw.New(nil)
	imd.Color = colornames.Aliceblue
    for i:=0; i<4; i++ {
        box := boxes[i]
	    imd.Push(box.Min, box.Max)
	    imd.Rectangle(0)
        if len(p1.Moves) <= i {
            continue
        }
        move := moves[i]
        move.Color = colornames.Black
        move.Dot.X -= move.BoundsOf(p1.Moves[i].Data.Name).W() / 2
        move.WriteString(p1.Moves[i].Data.Name)
    }

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

        var moveSelected uint8 = 0
        if win.JustPressed(pixelgl.MouseButtonLeft) {
            mousepos := cam.Unproject(win.MousePosition())
            switch {
            case movebox1.Contains(mousepos):
                moveSelected = 1
            case movebox2.Contains(mousepos):
                moveSelected = 2
            case movebox3.Contains(mousepos):
                moveSelected = 3
            case movebox4.Contains(mousepos):
                moveSelected = 4
            }
        }

        if moveSelected > 0 {
            sendMove(battle, moveSelected)
		    for _, l := range battle.Log().Logs()[turn-1] {
			    switch log := l.(type) {
			    case model.DamageLog:
				    if log.Index == 0 {
					    p1HP -= log.Damage
				    } else {
					    p2HP -= log.Damage
				    }
			    default:
			    }
			    fmt.Println(l)
		    }
        }
        if battle.IsOver() {
            fmt.Println("battle over")
            return
        }

		if win.Pressed(pixelgl.KeyH) || win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= dt
		}
		if win.Pressed(pixelgl.KeyL) || win.Pressed(pixelgl.KeyRight) {
			camPos.X += dt
		}
		if win.Pressed(pixelgl.KeyJ) || win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= dt
		}
		if win.Pressed(pixelgl.KeyK) || win.Pressed(pixelgl.KeyUp) {
			camPos.Y += dt
		}
		if win.Pressed(pixelgl.KeyQ) {
			return
		}

		// actually drawing things to screen
		win.Clear(colornames.Black)
        imd.Draw(win)
        for _, move := range moves {
            move.Draw(win, pixel.IM.Moved(camPos))
        }
        // we draw by center of sprite, and we start centered on (0,0)
        sprite.Draw(win, pixel.IM.Moved(camPos))
        sprite2.Draw(win, pixel.IM.Moved(camPos).Moved(pixel.V(-200,-150)))
        sprite3.Draw(win, pixel.IM.Moved(camPos).Moved(pixel.V(200,150)))
		//batch.Draw(win)
        /*
		for _, obj := range objects {
			drawTile(batch4, obj.(Unit).GetLoc(), unitSprite)
		}
		batch4.Draw(win)
        */
        imdHealth := imdraw.New(nil)
	    imdHealth.Color = colornames.Crimson
	    imdHealth.Push(pixel.V(200, -150), pixel.V(300, -160))
	    imdHealth.Rectangle(0)
	    imdHealth.Push(pixel.V(-200, 150), pixel.V(-300, 160))
	    imdHealth.Rectangle(0)
	    imdHealth.Color = colornames.Lime
	    imdHealth.Push(pixel.V(200, -150), pixel.V(200 + (float64(p1HP) / float64(p1.GetTotalHP()))*100 , -160))
	    imdHealth.Rectangle(0)
	    imdHealth.Push(pixel.V(-300, 150), pixel.V(-300 + (float64(p2HP) / float64(p2.GetTotalHP()))*100, 160))
	    imdHealth.Rectangle(0)
        imdHealth.Draw(win)

		win.Update()

		// frame counter
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func sendMove(battle model.Battle, choice uint8) {
		sourceCommand := model.MoveCommand{
			SourceIndex: 0,
			TargetIndex: 1,
			MoveIndex:   int(choice) - 1,
		}
		targetCommand := model.MoveCommand{
			SourceIndex: 1,
			TargetIndex: 0,
			MoveIndex:   random.Intn(4),
		}
		commands := []model.MoveCommand{sourceCommand, targetCommand}
		err := model.HandleTurn(battle, commands)
		if err != nil {
			panic(err)
		}
        turn++
}
