package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombieroom/g"
)

type Zombie struct {
	Entity

	Direction g.V2
}

func NewZombie(bounds g.Rect) *Zombie {
	zombie := &Zombie{}

	zombie.Position = g.RandomV2(bounds)

	zombie.Elasticity = 0.2
	zombie.Mass = 1.0
	zombie.Dampening = 0.999
	zombie.Radius = 0.5

	zombie.Direction = g.V2{}

	zombie.CollisionLayer = ZombieLayer
	zombie.CollisionMask = HammerLayer

	return zombie
}

func (zombie *Zombie) Update(game *Game, dt float32) {
	var nearest *Entity
	mindist := float32(1000000.0)
	for _, player := range game.Players {
		dist := player.Survivor.Position.Sub(zombie.Position).Length()
		if dist < mindist {
			nearest = &player.Survivor
			mindist = dist
		}
	}

	if nearest == nil {
		return
	}

	target := nearest.Position.Sub(zombie.Position).Normalize().Scale(g.Pow(0.5, dt))
	zombie.Velocity = zombie.Velocity.Add(target).Normalize().Scale(0.5)
	zombie.Direction = zombie.Direction.Add(zombie.Velocity.Scale(g.Pow(0.5, dt))).Normalize()
}

func (zombie *Zombie) Respawn(bounds g.Rect) {
	if len(zombie.Collision) == 0 {
		return
	}

	zombie.Position = g.RandomV2(bounds)
}

func (zombie *Zombie) Render(game *Game) {
	gl.PushMatrix()
	{
		gl.Translatef(zombie.Position.X, zombie.Position.Y, 0)

		rotation := -(zombie.Velocity.Angle() + g.Tau/4)
		gl.Rotatef(g.RadToDeg(rotation), 0, 0, -1)

		tex := game.Assets.TextureRepeat("assets/zombie.png")
		tex.Draw(g.NewCircleRect(zombie.Radius))
	}
	gl.PopMatrix()
}
