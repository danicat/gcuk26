package game

import (
	"fmt"
	"image/color"

	"london-eco-rider/internal/entity"
	"london-eco-rider/internal/render"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayState struct {
	Player         *entity.Player
	Items          []*entity.Item
	Obstacles      []*entity.Obstacle
	Background     *render.Background
	HUD            *render.HUD
	Score          int
	Bottles        int
	MaxBottles     int
	LevelWidth     float64
	CameraX        float64
	TimeRemaining  float64
	IsFever        bool
	FeverTimer     float64
	PromptMsg      string
	UnlockedBike   bool
	SelectingColor bool
	RecycleStation *entity.Obstacle
	BikeShopObj    *entity.Obstacle
}

func NewPlayState() *PlayState {
	return &PlayState{
		Background: render.NewBackground(640, 360),
		HUD:        render.NewHUD(640, 360),
		MaxBottles: 10,
	}
}

func (s *PlayState) Enter() {
	s.Player = entity.NewPlayer(50, 230)
	s.LevelWidth = 4500.0
	s.Score = 0
	s.Bottles = 0
	s.TimeRemaining = 180.0
	s.IsFever = false
	s.FeverTimer = 0.0
	s.UnlockedBike = false
	s.SelectingColor = false
	s.PromptMsg = "Collect plastic bottles on the pavement and find the Bike Shop & Recycling Bin!"

	s.generateLevel()
}

func (s *PlayState) generateLevel() {
	s.Items = nil
	s.Obstacles = nil

	// Phase 1 (0 to 600): Plastic Bottles and Ice Cream Cones on pavement
	s.Items = append(s.Items, entity.NewItem(entity.ItemPlasticBottle, 180, 230))
	s.Items = append(s.Items, entity.NewItem(entity.ItemIceCream, 260, 220))
	s.Items = append(s.Items, entity.NewItem(entity.ItemPlasticBottle, 340, 230))
	s.Items = append(s.Items, entity.NewItem(entity.ItemPlasticBottle, 460, 230))

	// Street Trees along initial pavement
	s.Obstacles = append(s.Obstacles, entity.NewObstacle(entity.ObstacleTree, 120, 268))
	s.Obstacles = append(s.Obstacles, entity.NewObstacle(entity.ObstacleTree, 400, 268))

	// London Bike Shop & First Recycling Station at x = 550 - 650
	s.BikeShopObj = entity.NewObstacle(entity.ObstacleBikeShop, 550, 268)
	s.RecycleStation = entity.NewObstacle(entity.ObstacleRecyclingBin, 650, 268)
	s.Obstacles = append(s.Obstacles, s.BikeShopObj)
	s.Obstacles = append(s.Obstacles, s.RecycleStation)
	s.Obstacles = append(s.Obstacles, entity.NewObstacle(entity.ObstacleIceCreamVan, 700, 268))

	// Phase 2 (700 to LevelEnd): Moving Buses, Phone Boxes, Ramps, Street Trees, and RARE Recycling Bins
	lastBinX := 650.0

	for x := 900.0; x < s.LevelWidth-300; x += 160.0 + rand.Float64()*120.0 {
		r := rand.Float64()
		var newObs *entity.Obstacle

		if rand.Float64() < 0.45 {
			s.Obstacles = append(s.Obstacles, entity.NewObstacle(entity.ObstacleTree, x-60.0, 268))
		}

		if r < 0.35 {
			newObs = entity.NewObstacle(entity.ObstacleRedBus, x, 268)
		} else if r < 0.60 {
			newObs = entity.NewObstacle(entity.ObstaclePhoneBox, x, 268)
		} else if r < 0.85 {
			newObs = entity.NewObstacle(entity.ObstacleRamp, x, 268)
		} else if x-lastBinX > 1200.0 {
			newObs = entity.NewObstacle(entity.ObstacleRecyclingBin, x, 268)
			s.Obstacles = append(s.Obstacles, entity.NewObstacle(entity.ObstacleIceCreamVan, x+50, 268))
			lastBinX = x
		}

		if newObs != nil {
			s.Obstacles = append(s.Obstacles, newObs)
		}

		itemX1 := x + 70.0
		if !s.isOverlapWithObstacles(itemX1) {
			s.Items = append(s.Items, entity.NewItem(entity.ItemIceCream, itemX1, 220-rand.Float64()*40.0))
		}

		itemX2 := x + 110.0
		if !s.isOverlapWithObstacles(itemX2) {
			s.Items = append(s.Items, entity.NewItem(entity.ItemPlasticBottle, itemX2, 220-rand.Float64()*30.0))
		}
	}
}

func (s *PlayState) isOverlapWithObstacles(x float64) bool {
	for _, obs := range s.Obstacles {
		if x >= obs.X-30.0 && x <= obs.X+obs.Width+30.0 {
			return true
		}
	}
	return false
}

func (s *PlayState) Update(g *Game) StateType {
	s.TimeRemaining -= 0.016
	if s.TimeRemaining <= 0 {
		return StateGameOver
	}

	if s.IsFever {
		s.FeverTimer -= 0.016
		if s.FeverTimer <= 0 {
			s.IsFever = false
		}
	}

	// Inputs
	pedalForward := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	pedalBack := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	isJumping := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	isInteractKey := isJumping || inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	triggerBoost := inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) || inpututil.IsKeyJustPressed(ebiten.KeyJ)

	if isJumping {
		g.Audio.PlayJump()
	}

	s.Player.Update(pedalForward, pedalBack, isJumping, triggerBoost)

	if s.Player.Physics.X < 20 {
		s.Player.Physics.X = 20
		s.Player.Physics.VX = 0
	}

	// Camera Follows Player
	s.CameraX = s.Player.Physics.X - 150
	if s.CameraX < 0 {
		s.CameraX = 0
	}
	if s.CameraX > s.LevelWidth-640 {
		s.CameraX = s.LevelWidth - 640
	}

	if s.Player.Physics.X >= s.LevelWidth-200 {
		g.FinalScore = s.Score
		return StateGameClear
	}

	playerAABB := s.Player.Physics.GetAABB()

	// Update Obstacles (Moving Buses)
	for _, obs := range s.Obstacles {
		obs.Update()
	}

	// Interactive Bike Shop & Color Selection
	s.handleBikeShopInteraction(isInteractKey, g)

	// Collectible items
	for _, item := range s.Items {
		item.Update()
		if !item.Collected && playerAABB.Intersects(item.GetAABB()) {
			item.Collected = true
			if item.Type == entity.ItemIceCream {
				g.Audio.PlayIceCream()
				s.Player.TriggerTurbo(3.0)
				pts := 100
				if s.IsFever {
					pts *= 2
				}
				s.Score += pts
				s.Player.Particles.EmitBurst(item.X, item.Y, 12, color.RGBA{255, 150, 200, 255})
			} else { // Plastic Bottle
				g.Audio.PlayBottle()
				s.Bottles++
				pts := 150
				if s.IsFever {
					pts *= 2
				}
				s.Score += pts
				s.Player.Particles.EmitBurst(item.X, item.Y, 12, color.RGBA{80, 220, 160, 255})

				if s.UnlockedBike && s.Bottles >= s.MaxBottles {
					s.Bottles = 0
					s.IsFever = true
					s.FeverTimer = 10.0
					s.Player.TriggerInvincibility(10.0)
					g.Audio.PlayBell()
				}
			}
		}
	}

	// Moving Platform Support & Fall Off Physics
	standingOnPlatform := false

	for _, obs := range s.Obstacles {
		if obs.Type == entity.ObstacleRecyclingBin || obs.Type == entity.ObstacleIceCreamVan || obs.Type == entity.ObstacleBikeShop || obs.Type == entity.ObstacleTree {
			continue // Non-blocking story decorations
		}

		obsAABB := obs.GetAABB()

		// Ramp Jump
		if obs.IsRamp && s.UnlockedBike && playerAABB.Intersects(obsAABB) {
			s.Player.Physics.VY = -11.5
			s.Player.Physics.Grounded = false
			continue
		}

		playerFootX := s.Player.Physics.X + s.Player.Physics.Width*0.5
		playerFootY := s.Player.Physics.Y + s.Player.Physics.Height

		// X-overlap over platform roof
		isOverPlatformX := playerFootX >= obs.X-2.0 && playerFootX <= obs.X+obs.Width+2.0
		obsTopY := obs.Y

		if isOverPlatformX && playerFootY >= obsTopY-8.0 && playerFootY <= obsTopY+16.0 {
			// STANDING ON TOP OF ROOF PLATFORM!
			s.Player.Physics.Y = obsTopY - s.Player.Physics.Height
			s.Player.Physics.VY = 0
			s.Player.Physics.Grounded = true
			standingOnPlatform = true

			// MOVING PLATFORM VELOCITY TRANSFER: Ride together with moving bus when idle!
			if obs.VX != 0 && !pedalForward && !pedalBack {
				s.Player.Physics.X += obs.VX
			}

			continue
		}

		// Ground-Level Side Collisions (evaluated when player is below top edge level)
		if playerFootY > obsTopY+16.0 && playerAABB.Intersects(obsAABB) {
			playerCenterX := s.Player.Physics.X + s.Player.Physics.Width*0.5
			obsCenterX := obs.X + obs.Width*0.5

			if obs.Type == entity.ObstacleRedBus {
				if playerCenterX > obsCenterX {
					s.Player.Physics.X = obs.X + obs.Width
					if s.Player.Physics.VX < 0 {
						s.Player.Physics.VX = 0
					}
				} else {
					obs.Stop()
					if !s.Player.Invincible {
						g.Audio.PlayCrash()
						s.Player.Physics.VX = 0
						s.Player.Physics.X = obs.X - s.Player.Physics.Width - 2.0
						s.Player.TriggerInvincibility(1.5)
						s.Player.Particles.EmitBurst(s.Player.Physics.X, s.Player.Physics.Y, 15, color.RGBA{255, 80, 80, 255})
					}
				}
			} else if obs.Type == entity.ObstaclePhoneBox {
				if playerCenterX < obsCenterX {
					s.Player.Physics.X = obs.X - s.Player.Physics.Width
					if s.Player.Physics.VX > 0 {
						s.Player.Physics.VX = 0
					}
				} else {
					s.Player.Physics.X = obs.X + obs.Width
					if s.Player.Physics.VX < 0 {
						s.Player.Physics.VX = 0
					}
				}
			}
		}
	}

	// Fall Off Platform Edge Check
	if !standingOnPlatform && s.Player.Physics.Y < s.Player.Physics.GroundY {
		s.Player.Physics.Grounded = false
	}

	return StatePlay
}

func (s *PlayState) handleBikeShopInteraction(isInteractKey bool, g *Game) {
	if s.RecycleStation == nil {
		return
	}

	playerCenterX := s.Player.Physics.X + s.Player.Physics.Width*0.5
	stationCenterX := s.RecycleStation.X + s.RecycleStation.Width*0.5
	dist := math.Abs(playerCenterX - stationCenterX)

	if dist <= 100.0 {
		if !s.UnlockedBike {
			if !s.SelectingColor {
				if s.Bottles >= 3 {
					s.PromptMsg = "PRESS SPACE / E / ENTER TO RECYCLE BOTTLES & SELECT BIKE! 🚴"
					if isInteractKey {
						s.SelectingColor = true
						g.Audio.PlayBell()
					}
				} else {
					s.PromptMsg = fmt.Sprintf("Collect %d more Plastic Bottle(s) on pavement for the Bike Shop!", 3-s.Bottles)
				}
			} else {
				s.PromptMsg = "SELECT BIKE COLOR: Press [1] RED | [2] BLUE | [3] GOLD"
				chosenColor := entity.BikeColorRed
				selected := false

				if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
					chosenColor = entity.BikeColorRed
					selected = true
				} else if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
					chosenColor = entity.BikeColorBlue
					selected = true
				} else if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
					chosenColor = entity.BikeColorGold
					selected = true
				}

				if selected {
					s.Bottles = 0
					s.SelectingColor = false
					s.UnlockedBike = true
					s.Player.SwitchToBike(chosenColor)
					s.Score += 500
					g.Audio.PlayBell()
					g.Audio.PlayIceCream()
					s.Player.TriggerTurbo(4.0)
					s.Player.Particles.EmitBurst(s.Player.Physics.X, s.Player.Physics.Y, 30, color.RGBA{255, 215, 0, 255})
					s.PromptMsg = "BICYCLE UNLOCKED! 🍦 Patrol London & recycle bottles!"
				}
			}
		} else {
			if s.Bottles > 0 {
				s.PromptMsg = "PRESS SPACE / E TO DUMP BOTTLES FOR BONUS SCORE! ♻️"
				if isInteractKey {
					s.Score += s.Bottles * 200
					s.Bottles = 0
					g.Audio.PlayBell()
					s.Player.Particles.EmitBurst(s.Player.Physics.X, s.Player.Physics.Y, 15, color.RGBA{80, 255, 120, 255})
				}
			}
		}
	} else {
		if !s.UnlockedBike {
			s.PromptMsg = "Walk forward and collect plastic bottles on the pavement!"
		} else if s.Player.TurboTimer > 0 {
			s.PromptMsg = "🍦 TURBO SPEED ACTIVE!"
		} else if s.IsFever {
			s.PromptMsg = "🌟 ECO FEVER! DOUBLE SCORE MULTIPLIER!"
		} else {
			s.PromptMsg = ""
		}
	}
}

func (s *PlayState) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen, s.CameraX)

	for _, obs := range s.Obstacles {
		obs.Draw(screen, s.CameraX)
	}
	for _, item := range s.Items {
		item.Draw(screen, s.CameraX)
	}

	s.Player.Draw(screen, s.CameraX)

	s.HUD.Draw(screen, s.Score, s.Bottles, s.MaxBottles, s.Player.TurboTimer, s.TimeRemaining, s.IsFever, s.Player.Mode, s.PromptMsg)
}

func (s *PlayState) Exit() {}
